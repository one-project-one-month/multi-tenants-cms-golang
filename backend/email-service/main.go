package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/resend/resend-go/v2"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/sirupsen/logrus"
)

// EmailProvider interface defines the contract for email providers
type EmailProvider interface {
	SendEmail(ctx context.Context, req SendEmailRequest) error
	GetProviderName() string
}

// EmailProviderFactory creates email providers
type EmailProviderFactory interface {
	CreateProvider(providerType string) (EmailProvider, error)
}

// Concrete factory implementation
type ConcreteEmailProviderFactory struct {
	logger       *logrus.Logger
	fromEmail    string
	fromName     string
	isDev        bool
	devEmailFile string
}

// SendGrid provider implementation
type SendGridProvider struct {
	client       *sendgrid.Client
	logger       *logrus.Logger
	fromEmail    string
	fromName     string
	isDev        bool
	devEmailFile string
}

// Resend provider implementation
type ResendProvider struct {
	client       *resend.Client
	logger       *logrus.Logger
	fromEmail    string
	fromName     string
	isDev        bool
	devEmailFile string
}

// EmailService struct
type EmailService struct {
	nats      *nats.Conn
	provider  EmailProvider
	log       *logrus.Logger
	templates map[string]EmailTemplate
	factory   EmailProviderFactory
}

// Email message and template structs remain the same
type EmailMessage struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	Type      string    `json:"type"`
	ExpiresAt time.Time `json:"expires_at"`
	Timestamp time.Time `json:"timestamp"`
}

type EmailTemplate struct {
	Subject     string
	HTMLContent string
	PlainText   string
}

type SendEmailRequest struct {
	To          string
	ToName      string
	Subject     string
	HTMLContent string
	PlainText   string
}

// Factory implementation
func NewEmailProviderFactory(logger *logrus.Logger, fromEmail, fromName string, isDev bool, devEmailFile string) EmailProviderFactory {
	return &ConcreteEmailProviderFactory{
		logger:       logger,
		fromEmail:    fromEmail,
		fromName:     fromName,
		isDev:        isDev,
		devEmailFile: devEmailFile,
	}
}

func (f *ConcreteEmailProviderFactory) CreateProvider(providerType string) (EmailProvider, error) {
	switch strings.ToLower(providerType) {
	case "sendgrid":
		return f.createSendGridProvider()
	case "resend":
		return f.createResendProvider()
	default:
		return nil, fmt.Errorf("unknown provider type: %s", providerType)
	}
}

func (f *ConcreteEmailProviderFactory) createSendGridProvider() (EmailProvider, error) {
	apiKey := getEnv("SENDGRID_API_KEY", "")
	if apiKey == "" {
		return nil, fmt.Errorf("SENDGRID_API_KEY environment variable is required")
	}

	client := sendgrid.NewSendClient(apiKey)

	return &SendGridProvider{
		client:       client,
		logger:       f.logger,
		fromEmail:    f.fromEmail,
		fromName:     f.fromName,
		isDev:        f.isDev,
		devEmailFile: f.devEmailFile,
	}, nil
}

func (f *ConcreteEmailProviderFactory) createResendProvider() (EmailProvider, error) {
	apiKey := getEnv("RESEND_API_KEY", "")
	if apiKey == "" {
		return nil, fmt.Errorf("RESEND_API_KEY environment variable is required")
	}

	client := resend.NewClient(apiKey)

	return &ResendProvider{
		client:       client,
		logger:       f.logger,
		fromEmail:    f.fromEmail,
		fromName:     f.fromName,
		isDev:        f.isDev,
		devEmailFile: f.devEmailFile,
	}, nil
}

// SendGrid provider methods
func (s *SendGridProvider) SendEmail(ctx context.Context, req SendEmailRequest) error {
	from := mail.NewEmail(s.fromName, s.fromEmail)
	to := mail.NewEmail(req.ToName, req.To)

	message := mail.NewSingleEmail(from, req.Subject, to, req.PlainText, req.HTMLContent)

	if s.isDev {
		condition := true
		message.SetMailSettings(&mail.MailSettings{
			SandboxMode: &mail.Setting{
				Enable: &condition,
			},
		})
		s.logEmailToFile(req)
		s.logger.WithFields(logrus.Fields{
			"to":       req.To,
			"subject":  req.Subject,
			"mode":     "sandbox",
			"provider": "sendgrid",
		}).Info("Email sent in sandbox mode")
	}

	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		response, err := s.client.Send(message)
		if err != nil {
			lastErr = fmt.Errorf("sendgrid request failed: %w", err)
			s.logger.WithError(err).WithField("attempt", attempt+1).Warn("Failed to send email, retrying...")
			time.Sleep(time.Duration(attempt+1) * time.Second)
			continue
		}

		if response.StatusCode >= 200 && response.StatusCode < 300 {
			s.logger.WithFields(logrus.Fields{
				"to":          req.To,
				"status_code": response.StatusCode,
				"provider":    "sendgrid",
				"mode":        s.getMode(),
			}).Info("Email sent successfully")
			return nil
		}

		lastErr = fmt.Errorf("sendgrid returned status %d: %s", response.StatusCode, response.Body)
		s.logger.WithFields(logrus.Fields{
			"status_code": response.StatusCode,
			"body":        response.Body,
			"attempt":     attempt + 1,
			"provider":    "sendgrid",
		}).Warn("SendGrid returned error status, retrying...")

		time.Sleep(time.Duration(attempt+1) * time.Second)
	}

	return fmt.Errorf("failed to send email after 3 attempts: %w", lastErr)
}

func (s *SendGridProvider) GetProviderName() string {
	return "sendgrid"
}

func (s *SendGridProvider) logEmailToFile(req SendEmailRequest) {
	emailLog := map[string]interface{}{
		"timestamp":    time.Now().Format(time.RFC3339),
		"provider":     "sendgrid",
		"to":           req.To,
		"to_name":      req.ToName,
		"subject":      req.Subject,
		"html_content": req.HTMLContent,
		"plain_text":   req.PlainText,
	}

	logData, _ := json.MarshalIndent(emailLog, "", "  ")

	file, err := os.OpenFile(s.devEmailFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		s.logger.WithError(err).Error("Failed to open email log file")
		return
	}
	defer file.Close()

	file.WriteString(string(logData) + "\n---\n")
}

func (s *SendGridProvider) getMode() string {
	if s.isDev {
		return "sandbox"
	}
	return "production"
}

// Resend provider methods
func (r *ResendProvider) SendEmail(ctx context.Context, req SendEmailRequest) error {
	from := fmt.Sprintf("%s <%s>", r.fromName, r.fromEmail)

	params := &resend.SendEmailRequest{
		From:    from,
		To:      []string{req.To},
		Subject: req.Subject,
		Html:    req.HTMLContent,
		Text:    req.PlainText,
	}

	if r.isDev {
		r.logEmailToFile(req)
		r.logger.WithFields(logrus.Fields{
			"to":       req.To,
			"subject":  req.Subject,
			"mode":     "development",
			"provider": "resend",
		}).Info("Email logged in development mode")
	}

	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		sent, err := r.client.Emails.Send(params)
		if err != nil {
			lastErr = fmt.Errorf("resend request failed: %w", err)
			r.logger.WithError(err).WithField("attempt", attempt+1).Warn("Failed to send email, retrying...")
			time.Sleep(time.Duration(attempt+1) * time.Second)
			continue
		}

		r.logger.WithFields(logrus.Fields{
			"to":       req.To,
			"email_id": sent.Id,
			"provider": "resend",
			"mode":     r.getMode(),
		}).Info("Email sent successfully")
		return nil
	}

	return fmt.Errorf("failed to send email after 3 attempts: %w", lastErr)
}

func (r *ResendProvider) GetProviderName() string {
	return "resend"
}

func (r *ResendProvider) logEmailToFile(req SendEmailRequest) {
	emailLog := map[string]interface{}{
		"timestamp":    time.Now().Format(time.RFC3339),
		"provider":     "resend",
		"to":           req.To,
		"to_name":      req.ToName,
		"subject":      req.Subject,
		"html_content": req.HTMLContent,
		"plain_text":   req.PlainText,
	}

	logData, _ := json.MarshalIndent(emailLog, "", "  ")

	file, err := os.OpenFile(r.devEmailFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		r.logger.WithError(err).Error("Failed to open email log file")
		return
	}
	defer file.Close()

	file.WriteString(string(logData) + "\n---\n")
}

func (r *ResendProvider) getMode() string {
	if r.isDev {
		return "development"
	}
	return "production"
}

// EmailService methods
func NewEmailService() (*EmailService, error) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	natsURL := getEnv("NATS_URL", "nats://localhost:4222")
	fromEmail := getEnv("FROM_EMAIL", "noreply@yourapp.com")
	fromName := getEnv("FROM_NAME", "Your App")
	isDev := strings.ToLower(getEnv("ENVIRONMENT", "development")) == "development"
	devEmailFile := getEnv("DEV_EMAIL_FILE", "sent_emails.log")
	providerType := getEnv("EMAIL_PROVIDER", "resend") // Default to resend

	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	factory := NewEmailProviderFactory(logger, fromEmail, fromName, isDev, devEmailFile)
	provider, err := factory.CreateProvider(providerType)
	if err != nil {
		return nil, fmt.Errorf("failed to create email provider: %w", err)
	}

	service := &EmailService{
		nats:      nc,
		provider:  provider,
		log:       logger,
		templates: initializeTemplates(),
		factory:   factory,
	}

	return service, nil
}

func (s *EmailService) Start() error {
	s.log.WithField("provider", s.provider.GetProviderName()).Info("Starting email service...")

	_, err := s.nats.Subscribe("email.verification", s.handleEmailVerification)
	if err != nil {
		return fmt.Errorf("failed to subscribe to email.verification: %w", err)
	}

	_, err = s.nats.Subscribe("email.password_reset", s.handlePasswordReset)
	if err != nil {
		return fmt.Errorf("failed to subscribe to email.password_reset: %w", err)
	}

	s.log.Info("Email service started successfully")
	return nil
}

func (s *EmailService) Shutdown() {
	s.log.Info("Shutting down email service...")
	if s.nats != nil {
		s.nats.Close()
	}
}

func (s *EmailService) SwitchProvider(providerType string) error {
	newProvider, err := s.factory.CreateProvider(providerType)
	if err != nil {
		return fmt.Errorf("failed to switch to provider %s: %w", providerType, err)
	}

	s.provider = newProvider
	s.log.WithField("provider", providerType).Info("Switched email provider")
	return nil
}

func (s *EmailService) handleEmailVerification(msg *nats.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var emailMsg EmailMessage
	if err := json.Unmarshal(msg.Data, &emailMsg); err != nil {
		s.log.WithError(err).Error("Failed to unmarshal email message")
		return
	}

	s.log.WithFields(logrus.Fields{
		"user_id":  emailMsg.UserID,
		"email":    emailMsg.Email,
		"type":     emailMsg.Type,
		"provider": s.provider.GetProviderName(),
	}).Info("Processing email verification request")

	if err := s.sendVerificationEmail(ctx, emailMsg); err != nil {
		s.log.WithError(err).Error("Failed to send verification email")
		return
	}

	s.log.WithFields(logrus.Fields{
		"user_id":  emailMsg.UserID,
		"email":    emailMsg.Email,
		"provider": s.provider.GetProviderName(),
	}).Info("Verification email sent successfully")
}

func (s *EmailService) handlePasswordReset(msg *nats.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var emailMsg EmailMessage
	if err := json.Unmarshal(msg.Data, &emailMsg); err != nil {
		s.log.WithError(err).Error("Failed to unmarshal password reset message")
		return
	}

	s.log.WithFields(logrus.Fields{
		"user_id":  emailMsg.UserID,
		"email":    emailMsg.Email,
		"type":     emailMsg.Type,
		"provider": s.provider.GetProviderName(),
	}).Info("Processing password reset request")

	if err := s.sendPasswordResetEmail(ctx, emailMsg); err != nil {
		s.log.WithError(err).Error("Failed to send password reset email")
		return
	}

	s.log.WithFields(logrus.Fields{
		"user_id":  emailMsg.UserID,
		"email":    emailMsg.Email,
		"provider": s.provider.GetProviderName(),
	}).Info("Password reset email sent successfully")
}

func (s *EmailService) sendVerificationEmail(ctx context.Context, emailMsg EmailMessage) error {
	template, exists := s.templates["email_verification"]
	if !exists {
		return fmt.Errorf("email verification template not found")
	}

	htmlContent := s.replacePlaceholders(template.HTMLContent, emailMsg)
	plainText := s.replacePlaceholders(template.PlainText, emailMsg)
	subject := s.replacePlaceholders(template.Subject, emailMsg)

	return s.provider.SendEmail(ctx, SendEmailRequest{
		To:          emailMsg.Email,
		ToName:      emailMsg.Name,
		Subject:     subject,
		HTMLContent: htmlContent,
		PlainText:   plainText,
	})
}

func (s *EmailService) sendPasswordResetEmail(ctx context.Context, emailMsg EmailMessage) error {
	template, exists := s.templates["password_reset"]
	if !exists {
		return fmt.Errorf("password reset template not found")
	}

	htmlContent := s.replacePlaceholders(template.HTMLContent, emailMsg)
	plainText := s.replacePlaceholders(template.PlainText, emailMsg)
	subject := s.replacePlaceholders(template.Subject, emailMsg)

	return s.provider.SendEmail(ctx, SendEmailRequest{
		To:          emailMsg.Email,
		ToName:      emailMsg.Name,
		Subject:     subject,
		HTMLContent: htmlContent,
		PlainText:   plainText,
	})
}

func (s *EmailService) replacePlaceholders(template string, emailMsg EmailMessage) string {
	replacements := map[string]string{
		"{{.Name}}":  emailMsg.Name,
		"{{.Code}}":  emailMsg.Code,
		"{{.Email}}": emailMsg.Email,
	}

	result := template
	for placeholder, value := range replacements {
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}

func initializeTemplates() map[string]EmailTemplate {
	return map[string]EmailTemplate{
		"email_verification": {
			Subject:     "Verify Your Email Address",
			HTMLContent: `<!DOCTYPE html><html><head><meta charset="UTF-8"><title>Email Verification</title><style>body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }.container { max-width: 600px; margin: 0 auto; padding: 20px; }.header { background: #007bff; color: white; padding: 20px; text-align: center; }.content { padding: 20px; background: #f9f9f9; }.code { font-size: 24px; font-weight: bold; color: #007bff; text-align: center; padding: 10px; }.footer { padding: 20px; text-align: center; font-size: 12px; color: #666; }</style></head><body><div class="container"><div class="header"><h1>Email Verification</h1></div><div class="content"><p>Hello {{.Name}},</p><p>Thank you for signing up! Please use the following verification code to complete your registration:</p><div class="code">{{.Code}}</div><p>This code will expire in 10 minutes.</p><p>If you didn't request this verification, please ignore this email.</p></div><div class="footer"><p>This email was sent automatically. Please do not reply.</p></div></div></body></html>`,
			PlainText: `Hello {{.Name}},

Thank you for signing up! Please use the following verification code to complete your registration:

{{.Code}}

This code will expire in 10 minutes.

If you didn't request this verification, please ignore this email.

This email was sent automatically. Please do not reply.`,
		},
		"password_reset": {
			Subject:     "Password Reset Request",
			HTMLContent: `<!DOCTYPE html><html><head><meta charset="UTF-8"><title>Password Reset</title><style>body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }.container { max-width: 600px; margin: 0 auto; padding: 20px; }.header { background: #dc3545; color: white; padding: 20px; text-align: center; }.content { padding: 20px; background: #f9f9f9; }.code { font-size: 24px; font-weight: bold; color: #dc3545; text-align: center; padding: 10px; }.footer { padding: 20px; text-align: center; font-size: 12px; color: #666; }</style></head><body><div class="container"><div class="header"><h1>Password Reset</h1></div><div class="content"><p>Hello {{.Name}},</p><p>We received a request to reset your password. Please use the following code:</p><div class="code">{{.Code}}</div><p>This code will expire in 10 minutes.</p><p>If you didn't request this password reset, please ignore this email and your password will remain unchanged.</p></div><div class="footer"><p>This email was sent automatically. Please do not reply.</p></div></div></body></html>`,
			PlainText: `Hello {{.Name}},

We received a request to reset your password. Please use the following code:

{{.Code}}

This code will expire in 10 minutes.

If you didn't request this password reset, please ignore this email and your password will remain unchanged.

This email was sent automatically. Please do not reply.`,
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	service, err := NewEmailService()
	if err != nil {
		log.Fatal("Failed to initialize email service:", err)
	}
	defer service.nats.Close()

	if err := service.Start(); err != nil {
		log.Fatal("Failed to start email service:", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	service.log.Info("Email service shutting down...")
	service.Shutdown()
}
