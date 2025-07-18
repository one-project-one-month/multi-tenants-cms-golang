package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/wneessen/go-mail"
)

type EmailRequest struct {
	To           string                 `json:"to"`
	Subject      string                 `json:"subject"`
	TemplateName string                 `json:"templateName"`
	Data         map[string]interface{} `json:"data"`
	// Enhanced fields
	Priority    int               `json:"priority,omitempty"`    // 1=high, 2=normal, 3=low
	RetryCount  int               `json:"retryCount,omitempty"`  // Current retry attempt
	MaxRetries  int               `json:"maxRetries,omitempty"`  // Maximum retry attempts
	ScheduleAt  int64             `json:"scheduleAt,omitempty"`  // Unix timestamp for delayed sending
	CC          []string          `json:"cc,omitempty"`          // CC recipients
	BCC         []string          `json:"bcc,omitempty"`         // BCC recipients
	ReplyTo     string            `json:"replyTo,omitempty"`     // Reply-to address
	Attachments []Attachment      `json:"attachments,omitempty"` // File attachments
	Headers     map[string]string `json:"headers,omitempty"`     // Custom headers
	TrackingID  string            `json:"trackingId,omitempty"`  // For tracking purposes
}

type Attachment struct {
	Name    string `json:"name"`
	Content string `json:"content"` // Base64 encoded content
	Type    string `json:"type"`    // MIME type
}

type Config struct {
	SMTP struct {
		Host     string
		Port     int
		Username string
		Password string
		From     string
		// Enhanced SMTP config
		MaxConnections    int
		ConnectionTimeout time.Duration
		SendTimeout       time.Duration
		KeepAlive         bool
	}
	NATS struct {
		URL      string
		Stream   string
		Subjects []string
		// Enhanced NATS config
		MaxDeliver    int
		AckWait       time.Duration
		MaxAckPending int
	}
	Templates struct {
		Dir          string
		DefaultExt   string
		CacheEnabled bool
		// Enhanced template config
		WatchChanges   bool
		ReloadOnChange bool
	}
	// New sections
	Processing struct {
		Workers           int
		BatchSize         int
		RetryDelaySeconds int
		MaxRetries        int
	}
	Monitoring struct {
		EnableMetrics bool
		LogLevel      string
	}
}

type EmailMetrics struct {
	TotalSent     int64
	TotalFailed   int64
	TotalRetries  int64
	LastProcessed time.Time
	mu            sync.RWMutex
}

var metrics = &EmailMetrics{}

func main() {
	cfg := loadConfig()

	// Initialize SMTP client pool
	smtpPool := make(chan *mail.Client, cfg.Processing.Workers)
	for i := 0; i < cfg.Processing.Workers; i++ {
		client, err := createSMTPClient(cfg)
		if err != nil {
			log.Fatalf("Failed to create SMTP client: %v", err)
		}
		smtpPool <- client
	}

	nc, err := nats.Connect(cfg.NATS.URL)
	if err != nil {
		log.Fatalf("NATS connection failed: %v", err)
	}
	defer nc.Close()

	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatalf("JetStream init failed: %v", err)
	}

	streamConfig := jetstream.StreamConfig{
		Name:      cfg.NATS.Stream,
		Subjects:  cfg.NATS.Subjects,
		Retention: jetstream.WorkQueuePolicy,
		Storage:   jetstream.FileStorage,
		// Enhanced stream config
		MaxAge:       24 * time.Hour,
		MaxMsgs:      1000000,
		MaxBytes:     1024 * 1024 * 1024, // 1GB
		MaxConsumers: 10,
		Duplicates:   2 * time.Minute,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := js.CreateOrUpdateStream(ctx, streamConfig)
	if err != nil {
		log.Fatalf("Stream creation failed: %v", err)
	}

	consumer, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:       "email-processor",
		AckPolicy:     jetstream.AckExplicitPolicy,
		DeliverPolicy: jetstream.DeliverAllPolicy,
		// Enhanced consumer config
		MaxDeliver:    cfg.NATS.MaxDeliver,
		AckWait:       cfg.NATS.AckWait,
		MaxAckPending: cfg.NATS.MaxAckPending,
		FilterSubject: "",
		ReplayPolicy:  jetstream.ReplayInstantPolicy,
	})
	if err != nil {
		log.Fatalf("Consumer creation failed: %v", err)
	}

	log.Println("Email service started. Waiting for messages...")

	var templateCache *template.Template
	if cfg.Templates.CacheEnabled {
		templateCache, err = loadTemplates(cfg.Templates.Dir, cfg.Templates.DefaultExt)
		if err != nil {
			log.Fatalf("Failed to load templates: %v", err)
		}
	}

	// Start template watcher if enabled
	if cfg.Templates.WatchChanges {
		go watchTemplates(cfg.Templates.Dir, cfg.Templates.DefaultExt, &templateCache)
	}

	// Start metrics reporter
	if cfg.Monitoring.EnableMetrics {
		go reportMetrics()
	}

	// Start workers
	for i := 0; i < cfg.Processing.Workers; i++ {
		go func(workerID int) {
			_, err := consumer.Consume(func(msg jetstream.Msg) {
				processMessage(msg, smtpPool, cfg.SMTP.From, cfg.Templates, templateCache, cfg.Processing, workerID)
			})
			if err != nil {
				log.Printf("Worker %d consume failed: %v", workerID, err)
			}
		}(i)
	}

	select {}
}

func createSMTPClient(cfg Config) (*mail.Client, error) {
	options := []mail.Option{
		mail.WithPort(cfg.SMTP.Port),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(cfg.SMTP.Username),
		mail.WithPassword(cfg.SMTP.Password),
		mail.WithTLSPolicy(mail.TLSMandatory),
	}

	if cfg.SMTP.ConnectionTimeout > 0 {
		options = append(options, mail.WithTimeout(cfg.SMTP.ConnectionTimeout))
	}

	return mail.NewClient(cfg.SMTP.Host, options...)
}

func processMessage(msg jetstream.Msg, smtpPool chan *mail.Client, from string, templateCfg struct {
	Dir            string
	DefaultExt     string
	CacheEnabled   bool
	WatchChanges   bool
	ReloadOnChange bool
}, templateCache *template.Template, processingCfg struct {
	Workers           int
	BatchSize         int
	RetryDelaySeconds int
	MaxRetries        int
}, workerID int) {

	log.Printf("Worker %d received message on subject: %s", workerID, msg.Subject())

	var emailReq EmailRequest
	if err := json.Unmarshal(msg.Data(), &emailReq); err != nil {
		log.Printf("Worker %d failed to parse message: %v", workerID, err)
		if err := msg.Nak(); err != nil {
			log.Printf("Worker %d failed to NAK message: %v", workerID, err)
		}
		return
	}

	// Set defaults
	if emailReq.TemplateName == "" {
		emailReq.TemplateName = "default"
	}
	if emailReq.MaxRetries == 0 {
		emailReq.MaxRetries = processingCfg.MaxRetries
	}

	// Check if message should be delayed
	if emailReq.ScheduleAt > 0 && time.Now().Unix() < emailReq.ScheduleAt {
		log.Printf("Worker %d delaying message until %v", workerID, time.Unix(emailReq.ScheduleAt, 0))
		// Requeue with delay
		if err := msg.Nak(); err != nil {
			log.Printf("Worker %d failed to NAK delayed message: %v", workerID, err)
		}
		return
	}

	// Get SMTP client from pool
	client := <-smtpPool
	defer func() {
		smtpPool <- client
	}()

	if err := sendEnhancedEmail(client, from, emailReq, templateCfg, templateCache, workerID); err != nil {
		log.Printf("Worker %d failed to send email: %v", workerID, err)

		// Retry logic
		if emailReq.RetryCount < emailReq.MaxRetries {
			emailReq.RetryCount++

			// Exponential backoff
			delay := time.Duration(processingCfg.RetryDelaySeconds) * time.Second * time.Duration(emailReq.RetryCount)
			log.Printf("Worker %d retrying email (attempt %d/%d) after %v", workerID, emailReq.RetryCount, emailReq.MaxRetries, delay)

			//retryData, _ := json.Marshal(emailReq)
			if err := msg.Nak(); err != nil {
				log.Printf("Worker %d failed to NAK for retry: %v", workerID, err)
			}

			metrics.mu.Lock()
			metrics.TotalRetries++
			metrics.mu.Unlock()

			time.Sleep(delay)
			return
		}

		// Max retries exceeded
		log.Printf("Worker %d max retries exceeded for email to %s", workerID, emailReq.To)
		metrics.mu.Lock()
		metrics.TotalFailed++
		metrics.mu.Unlock()

		if err := msg.Ack(); err != nil {
			log.Printf("Worker %d failed to ACK failed message: %v", workerID, err)
		}
		return
	}

	log.Printf("Worker %d email sent to %s (tracking: %s)", workerID, emailReq.To, emailReq.TrackingID)

	metrics.mu.Lock()
	metrics.TotalSent++
	metrics.LastProcessed = time.Now()
	metrics.mu.Unlock()

	if err := msg.Ack(); err != nil {
		log.Printf("Worker %d failed to ACK message: %v", workerID, err)
	}
}

func sendEnhancedEmail(client *mail.Client, from string, req EmailRequest, templateCfg struct {
	Dir            string
	DefaultExt     string
	CacheEnabled   bool
	WatchChanges   bool
	ReloadOnChange bool
}, templateCache *template.Template, workerID int) error {

	m := mail.NewMsg()

	if err := m.From(from); err != nil {
		return fmt.Errorf("failed to set From address: %w", err)
	}
	if err := m.To(req.To); err != nil {
		return fmt.Errorf("failed to set To address: %w", err)
	}

	// Enhanced recipients
	if len(req.CC) > 0 {
		if err := m.Cc(req.CC...); err != nil {
			return fmt.Errorf("failed to set CC addresses: %w", err)
		}
	}
	if len(req.BCC) > 0 {
		if err := m.Bcc(req.BCC...); err != nil {
			return fmt.Errorf("failed to set BCC addresses: %w", err)
		}
	}
	if req.ReplyTo != "" {
		if err := m.ReplyTo(req.ReplyTo); err != nil {
			return fmt.Errorf("failed to set Reply-To address: %w", err)
		}
	}

	m.Subject(req.Subject)

	// Set priority
	switch req.Priority {
	case 1:
		m.SetImportance(mail.ImportanceHigh)
	case 3:
		m.SetImportance(mail.ImportanceLow)
	default:
		m.SetImportance(mail.ImportanceNormal)
	}

	for key, value := range req.Headers {
		m.SetGenHeader(mail.Header(key), value)
	}

	// Add tracking header
	if req.TrackingID != "" {
		m.SetGenHeader("X-Tracking-ID", req.TrackingID)
	}

	// Render template
	htmlContent, err := renderTemplate(req, templateCfg, templateCache)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	m.SetBodyString(mail.TypeTextHTML, htmlContent)

	// Add attachments
	for _, attachment := range req.Attachments {
		if err := m.AttachReader(attachment.Name, strings.NewReader(attachment.Content)); err != nil {
			return fmt.Errorf("failed to add attachment %s: %w", attachment.Name, err)
		}
	}

	if err := client.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func renderTemplate(req EmailRequest, templateCfg struct {
	Dir            string
	DefaultExt     string
	CacheEnabled   bool
	WatchChanges   bool
	ReloadOnChange bool
}, templateCache *template.Template) (string, error) {

	templatePath := filepath.Join(templateCfg.Dir, req.TemplateName+templateCfg.DefaultExt)

	if templateCache != nil {
		tmpl := templateCache.Lookup(filepath.Base(templatePath))
		if tmpl == nil {
			return "", fmt.Errorf("template %s not found in cache", templatePath)
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, req.Data); err != nil {
			return "", fmt.Errorf("failed to execute template: %w", err)
		}
		return buf.String(), nil
	}

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, req.Data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	return buf.String(), nil
}

func loadTemplates(dir, ext string) (*template.Template, error) {
	if ext == "" {
		ext = ".html"
	}

	pattern := filepath.Join(dir, "*"+ext)
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to find template files: %w", err)
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no templates found in %s with extension %s", dir, ext)
	}

	tmpl := template.New("").Funcs(template.FuncMap{
		"default": func(defVal interface{}, val interface{}) interface{} {
			if val == nil || val == "" {
				return defVal
			}
			return val
		},
		// Enhanced template functions
		"formatDate": func(timestamp int64) string {
			return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
		},
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
		"title": strings.Title,
	})

	for _, file := range files {
		_, err := tmpl.ParseFiles(file)
		if err != nil {
			return nil, fmt.Errorf("failed to parse template %s: %w", file, err)
		}
	}

	return tmpl, nil
}

func watchTemplates(dir, ext string, templateCache **template.Template) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		newCache, err := loadTemplates(dir, ext)
		if err != nil {
			log.Printf("Failed to reload templates: %v", err)
			continue
		}
		*templateCache = newCache
		log.Println("Templates reloaded")
	}
}

func reportMetrics() {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		metrics.mu.RLock()
		log.Printf("Metrics - Sent: %d, Failed: %d, Retries: %d, Last: %v, Goroutines: %d",
			metrics.TotalSent, metrics.TotalFailed, metrics.TotalRetries,
			metrics.LastProcessed.Format("15:04:05"), runtime.NumGoroutine())
		metrics.mu.RUnlock()
	}
}

func loadConfig() Config {
	var cfg Config

	cfg.SMTP.Host = getEnv("SMTP_HOST", "smtp.gmail.com")
	cfg.SMTP.Port = getEnvInt("SMTP_PORT", 587)
	cfg.SMTP.Username = getEnv("SMTP_USER", "")
	cfg.SMTP.Password = getEnv("SMTP_PASSWORD", "")
	cfg.SMTP.From = getEnv("FROM_ADDR", cfg.SMTP.Username)
	cfg.SMTP.MaxConnections = getEnvInt("SMTP_MAX_CONNECTIONS", 10)
	cfg.SMTP.ConnectionTimeout = time.Duration(getEnvInt("SMTP_CONNECTION_TIMEOUT", 30)) * time.Second
	cfg.SMTP.SendTimeout = time.Duration(getEnvInt("SMTP_SEND_TIMEOUT", 60)) * time.Second
	cfg.SMTP.KeepAlive = getEnvBool("SMTP_KEEP_ALIVE", true)

	cfg.NATS.URL = getEnv("NATS_URL", "nats://localhost:4222")
	cfg.NATS.Stream = getEnv("NATS_STREAM", "EMAILS")
	subjects := getEnv("NATS_SUBJECTS", "email.verification,email.notification,page.approval")
	cfg.NATS.Subjects = strings.Split(subjects, ",")
	cfg.NATS.MaxDeliver = getEnvInt("NATS_MAX_DELIVER", 3)
	cfg.NATS.AckWait = time.Duration(getEnvInt("NATS_ACK_WAIT", 300)) * time.Second
	cfg.NATS.MaxAckPending = getEnvInt("NATS_MAX_ACK_PENDING", 100)

	cfg.Templates.Dir = getEnv("TEMPLATES_DIR", "templates")
	cfg.Templates.DefaultExt = getEnv("TEMPLATES_EXT", ".html")
	cfg.Templates.CacheEnabled = getEnvBool("TEMPLATES_CACHE", true)
	cfg.Templates.WatchChanges = getEnvBool("TEMPLATES_WATCH", false)
	cfg.Templates.ReloadOnChange = getEnvBool("TEMPLATES_RELOAD", false)

	cfg.Processing.Workers = getEnvInt("PROCESSING_WORKERS", runtime.NumCPU())
	cfg.Processing.BatchSize = getEnvInt("PROCESSING_BATCH_SIZE", 10)
	cfg.Processing.RetryDelaySeconds = getEnvInt("PROCESSING_RETRY_DELAY", 60)
	cfg.Processing.MaxRetries = getEnvInt("PROCESSING_MAX_RETRIES", 3)

	cfg.Monitoring.EnableMetrics = getEnvBool("MONITORING_METRICS", true)
	cfg.Monitoring.LogLevel = getEnv("MONITORING_LOG_LEVEL", "INFO")

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		var result int
		if _, err := fmt.Sscanf(value, "%d", &result); err == nil {
			return result
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		return strings.ToLower(value) == "true"
	}
	return defaultValue
}
