package email

import (
	"bytes"
	"context"
	"html/template"
	"path/filepath"
	"sync"
	"time"

	"github.com/multi-tenants-cms-golang/email-service/config"
	"github.com/multi-tenants-cms-golang/email-service/internal/nats"
	"github.com/pkg/errors"
	"github.com/wneessen/go-mail"
	"go.uber.org/zap"
)

// Service handles email sending with connection pooling and retry mechanisms
type Service struct {
	logger    *zap.Logger
	config    *config.Config
	templates *templateCache
	pool      *smtpPool
	metrics   *emailMetrics
}

type templateCache struct {
	sync.RWMutex
	templates map[string]*template.Template
	dir       string
}

type smtpPool struct {
	clients chan *mail.Client
	factory func() (*mail.Client, error)
	logger  *zap.Logger
}

type emailMetrics struct {
	sentCounter     int
	failedCounter   int
	retryCounter    int
	poolWaitCounter int
	mu              sync.Mutex
}

// NewService creates a new email service with connection pooling
func NewService(logger *zap.Logger, cfg *config.Config) (*Service, error) {
	// Initialize template cache
	templateCache := &templateCache{
		templates: make(map[string]*template.Template),
		dir:       cfg.SMTP.TemplateDir,
	}

	// Preload templates
	if err := templateCache.preload(); err != nil {
		return nil, errors.Wrap(err, "failed to preload templates")
	}

	// Initialize metrics
	metrics := &emailMetrics{}

	// Initialize SMTP connection pool
	pool := &smtpPool{
		clients: make(chan *mail.Client, cfg.SMTP.PoolSize),
		factory: func() (*mail.Client, error) {
			return newSMTPClient(cfg)
		},
		logger: logger,
	}

	// Warm up the connection pool
	if err := pool.warm(cfg.SMTP.PoolSize); err != nil {
		return nil, errors.Wrap(err, "failed to warm up SMTP connection pool")
	}

	return &Service{
		logger:    logger,
		config:    cfg,
		templates: templateCache,
		pool:      pool,
		metrics:   metrics,
	}, nil
}

// newSMTPClient creates a new SMTP client with configured options
func newSMTPClient(cfg *config.Config) (*mail.Client, error) {
	opts := []mail.Option{
		mail.WithPort(cfg.SMTP.Port),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(cfg.SMTP.User),
		mail.WithPassword(cfg.SMTP.Password),
		mail.WithTimeout(cfg.SMTP.Timeout),
	}

	//if cfg.Security.TLSEnabled {
	//	opts = append(opts, mail.WithTLSPolicy(mail.TLSMandatory))
	//	if cfg.Security.ServerName != "" {
	//		opts = append(opts, mail.WWithSSLServerName(cfg.Security.ServerName))
	//	}
	//	if cfg.Security.Insecure {
	//		opts = append(opts, mail.WithSSLInsecureSkipVerify())
	//	}
	//}

	client, err := mail.NewClient(cfg.SMTP.Host, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create SMTP client")
	}

	return client, nil
}

// Send sends an email with retry logic and connection pooling
func (s *Service) Send(ctx context.Context, req nats.EmailRequest) error {
	startTime := time.Now()

	if err := validateEmailRequest(req); err != nil {
		s.metrics.recordFailure()
		return errors.Wrap(err, "invalid email request")
	}

	// Get template
	tpl, err := s.templates.get(req.Template)
	if err != nil {
		s.metrics.recordFailure()
		return errors.Wrap(err, "template error")
	}

	// Render template
	var body bytes.Buffer
	if err := tpl.Execute(&body, req.Data); err != nil {
		s.metrics.recordFailure()
		return errors.Wrap(err, "template execution error")
	}

	// Get subject
	subject := getSubject(req.Data)

	// Get SMTP client from pool
	client, err := s.pool.get(ctx)
	if err != nil {
		s.metrics.recordFailure()
		return errors.Wrap(err, "failed to get SMTP client")
	}
	defer s.pool.put(client)

	// Prepare message
	m := mail.NewMsg()
	if err := m.FromFormat(s.config.SMTP.FromName, s.config.SMTP.FromAddr); err != nil {
		s.metrics.recordFailure()
		return errors.Wrap(err, "failed to set from address")
	}
	if err := m.To(req.To); err != nil {
		s.metrics.recordFailure()
		return errors.Wrap(err, "failed to set to address")
	}
	m.Subject(subject)
	m.SetBodyString(mail.TypeTextHTML, body.String())

	// Send with retries
	var lastErr error
	for attempt := 1; attempt <= s.config.SMTP.MaxRetries; attempt++ {
		if err := ctx.Err(); err != nil {
			s.metrics.recordFailure()
			return errors.Wrap(err, "context canceled")
		}

		if err := client.DialAndSend(m); err != nil {
			lastErr = err
			s.metrics.recordRetry()
			s.logger.Warn("failed to send email, retrying",
				zap.String("to", req.To),
				zap.Int("attempt", attempt),
				zap.Error(err))

			if attempt < s.config.SMTP.MaxRetries {
				time.Sleep(s.config.SMTP.RetryDelay)
			}
			continue
		}

		s.metrics.recordSuccess()
		s.logger.Info("email sent successfully",
			zap.String("to", req.To),
			zap.String("template", req.Template),
			zap.String("subject", subject),
			zap.Duration("duration", time.Since(startTime)))
		return nil
	}

	s.metrics.recordFailure()
	return errors.Wrap(lastErr, "max retries exceeded")
}

// validateEmailRequest validates the email request
func validateEmailRequest(req nats.EmailRequest) error {
	if req.To == "" {
		return errors.New("recipient email is required")
	}
	if req.Template == "" {
		return errors.New("template is required")
	}
	return nil
}

// getSubject extracts subject from email data
func getSubject(data map[string]interface{}) string {
	if subject, ok := data["subject"].(string); ok && subject != "" {
		return subject
	}
	return "No Subject"
}

// templateCache methods
func (tc *templateCache) preload() error {
	// Implement template preloading if needed
	return nil
}

func (tc *templateCache) get(templateName string) (*template.Template, error) {
	tc.RLock()
	tpl, exists := tc.templates[templateName]
	tc.RUnlock()

	if exists {
		return tpl, nil
	}

	// Load template if not in cache
	tc.Lock()
	defer tc.Unlock()

	// Double-check in case another goroutine loaded it
	if tpl, exists := tc.templates[templateName]; exists {
		return tpl, nil
	}

	tplPath := filepath.Join(tc.dir, templateName+".html")
	tpl, err := template.ParseFiles(tplPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse template %s", templateName)
	}

	tc.templates[templateName] = tpl
	return tpl, nil
}

// smtpPool methods
func (p *smtpPool) warm(count int) error {
	for i := 0; i < count; i++ {
		client, err := p.factory()
		if err != nil {
			return errors.Wrap(err, "failed to create SMTP client during warmup")
		}
		p.clients <- client
	}
	return nil
}

func (p *smtpPool) get(ctx context.Context) (*mail.Client, error) {
	select {
	case client := <-p.clients:
		return client, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		p.logger.Warn("SMTP connection pool exhausted, creating new connection")
		return p.factory()
	}
}

func (p *smtpPool) put(client *mail.Client) {
	select {
	case p.clients <- client:
		// Returned to pool
	default:
		// Pool is full, close the connection
		if err := client.Close(); err != nil {
			p.logger.Error("failed to close SMTP connection", zap.Error(err))
		}
	}
}

// emailMetrics methods
func (m *emailMetrics) recordSuccess() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sentCounter++
}

func (m *emailMetrics) recordFailure() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.failedCounter++
}

func (m *emailMetrics) recordRetry() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.retryCounter++
}

func (m *emailMetrics) recordPoolWait() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.poolWaitCounter++
}

// GetMetrics returns current metrics
func (s *Service) GetMetrics() map[string]int {
	s.metrics.mu.Lock()
	defer s.metrics.mu.Unlock()

	return map[string]int{
		"sent":      s.metrics.sentCounter,
		"failed":    s.metrics.failedCounter,
		"retries":   s.metrics.retryCounter,
		"poolWaits": s.metrics.poolWaitCounter,
	}
}

// HealthCheck verifies SMTP connectivity
func (s *Service) HealthCheck(ctx context.Context) error {
	client, err := s.pool.get(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get SMTP client for health check")
	}
	defer s.pool.put(client)

	// Test connection
	if _, err := client.DialToSMTPClientWithContext(ctx); err != nil {
		return errors.Wrap(err, "SMTP health check failed")
	}
	defer func(client *mail.Client) {
		err := client.Close()
		if err != nil {
			s.logger.Error("failed to close SMTP client", zap.Error(err))
		}
	}(client)

	return nil
}

// Close cleans up resources
func (s *Service) Close() error {
	close(s.pool.clients)
	for client := range s.pool.clients {
		if err := client.Close(); err != nil {
			s.logger.Error("failed to close SMTP client", zap.Error(err))
		}
	}
	return nil
}
