package email

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"path/filepath"
	"time"

	"github.com/multi-tenants-cms-golang/email-service/internal/nats"
	"github.com/wneessen/go-mail"
	"go.uber.org/zap"
)

type Service struct {
	client    *mail.Client
	logger    *zap.Logger
	tpDir     string
	fromAddr  string
	templates map[string]*template.Template
}

func NewService(
	logger *zap.Logger,
	smtpHost string,
	smtpPort int,
	smtpUser string,
	smtpPassword string,
	fromAddr string,
	tpDir string,
) (*Service, error) {
	client, err := mail.NewClient(
		smtpHost,
		mail.WithPort(smtpPort),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(smtpUser),
		mail.WithPassword(smtpPassword),
		mail.WithTLSPolicy(mail.TLSMandatory),
		mail.WithTimeout(30*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create mail client: %w", err)
	}

	service := &Service{
		client:    client,
		logger:    logger,
		tpDir:     tpDir,
		fromAddr:  fromAddr,
		templates: make(map[string]*template.Template),
	}

	return service, nil
}

func (s *Service) Send(ctx context.Context, req nats.EmailRequest) error {
	if err := s.validateRequest(req); err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}

	tpl, err := s.getTemplate(req.Template)
	if err != nil {
		return fmt.Errorf("template error: %w", err)
	}

	var body bytes.Buffer
	if err := tpl.Execute(&body, req); err != nil {
		return fmt.Errorf("template execution error: %w", err)
	}

	m := mail.NewMsg()
	if err := m.From(s.fromAddr); err != nil {
		return fmt.Errorf("failed to set from address: %w", err)
	}
	if err := m.To(req.To); err != nil {
		return fmt.Errorf("failed to set to address: %w", err)
	}

	subject, ok := req.Data["subject"].(string)
	if !ok || subject == "" {
		subject = "No Subject"
	}
	m.Subject(subject)
	m.SetBodyString(mail.TypeTextHTML, body.String())

	if err := s.client.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	s.logger.Info("email sent successfully",
		zap.String("to", req.To),
		zap.String("template", req.Template),
		zap.String("subject", subject))
	return nil
}

func (s *Service) validateRequest(req nats.EmailRequest) error {
	if req.To == "" {
		return fmt.Errorf("recipient email is required")
	}
	if req.Template == "" {
		return fmt.Errorf("template is required")
	}
	return nil
}

func (s *Service) getTemplate(templateName string) (*template.Template, error) {
	if tpl, exists := s.templates[templateName]; exists {
		return tpl, nil
	}

	tplPath := filepath.Join(s.tpDir, templateName+".html")
	tpl, err := template.ParseFiles(tplPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template %s: %w", templateName, err)
	}

	s.templates[templateName] = tpl
	return tpl, nil
}
