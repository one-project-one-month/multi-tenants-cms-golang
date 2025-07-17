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
	"strings"
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
}

type Config struct {
	SMTP struct {
		Host     string
		Port     int
		Username string
		Password string
		From     string
	}
	NATS struct {
		URL      string
		Stream   string
		Subjects []string
	}
	Templates struct {
		Dir          string
		DefaultExt   string
		CacheEnabled bool
	}
}

func main() {
	// Load configuration
	cfg := loadConfig()

	// Initialize SMTP client
	smtpClient, err := mail.NewClient(
		cfg.SMTP.Host,
		mail.WithPort(cfg.SMTP.Port),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(cfg.SMTP.Username),
		mail.WithPassword(cfg.SMTP.Password),
		mail.WithTLSPolicy(mail.TLSMandatory),
	)
	if err != nil {
		log.Fatalf("Failed to create SMTP client: %v", err)
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

	_, err = consumer.Consume(func(msg jetstream.Msg) {
		processMessage(msg, smtpClient, cfg.SMTP.From, cfg.Templates, templateCache)
	})
	if err != nil {
		log.Fatalf("Consume failed: %v", err)
	}

	select {}
}

func processMessage(msg jetstream.Msg, client *mail.Client, from string, templateCfg struct {
	Dir          string
	DefaultExt   string
	CacheEnabled bool
}, templateCache *template.Template) {
	log.Printf("Received message on subject: %s", msg.Subject())

	var emailReq EmailRequest
	if err := json.Unmarshal(msg.Data(), &emailReq); err != nil {
		log.Printf("Failed to parse message: %v", err)
		if err := msg.Nak(); err != nil {
			log.Printf("Failed to NAK message: %v", err)
		}
		return
	}

	if emailReq.TemplateName == "" {
		emailReq.TemplateName = "default"
	}

	if err := sendEmail(client, from, emailReq, templateCfg, templateCache); err != nil {
		log.Printf("Failed to send email: %v", err)
		if err := msg.Nak(); err != nil {
			log.Printf("Failed to NAK message: %v", err)
		}
		return
	}

	log.Printf("Email sent to %s", emailReq.To)
	if err := msg.Ack(); err != nil {
		log.Printf("Failed to ACK message: %v", err)
	}
}

func sendEmail(client *mail.Client, from string, req EmailRequest, templateCfg struct {
	Dir          string
	DefaultExt   string
	CacheEnabled bool
}, templateCache *template.Template) error {
	m := mail.NewMsg()
	if err := m.From(from); err != nil {
		return fmt.Errorf("failed to set From address: %w", err)
	}
	if err := m.To(req.To); err != nil {
		return fmt.Errorf("failed to set To address: %w", err)
	}

	m.Subject(req.Subject)

	htmlContent, err := renderTemplate(req, templateCfg, templateCache)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	m.SetBodyString(mail.TypeTextHTML, htmlContent)

	if err := client.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func renderTemplate(req EmailRequest, templateCfg struct {
	Dir          string
	DefaultExt   string
	CacheEnabled bool
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
	})

	// Parse each file individually to get better error messages
	for _, file := range files {
		_, err := tmpl.ParseFiles(file)
		if err != nil {
			return nil, fmt.Errorf("failed to parse template %s: %w", file, err)
		}
	}

	return tmpl, nil
}
func loadConfig() Config {
	var cfg Config

	cfg.SMTP.Host = getEnv("SMTP_HOST", "smtp.gmail.com")
	cfg.SMTP.Port = getEnvInt("SMTP_PORT", 587)
	cfg.SMTP.Username = getEnv("SMTP_USER", "")
	cfg.SMTP.Password = getEnv("SMTP_PASSWORD", "")
	cfg.SMTP.From = getEnv("FROM_ADDR", cfg.SMTP.Username)

	cfg.NATS.URL = getEnv("NATS_URL", "nats://localhost:4222")
	cfg.NATS.Stream = getEnv("NATS_STREAM", "EMAILS")
	subjects := getEnv("NATS_SUBJECTS", "email.verification,email.notification,page.approval")
	cfg.NATS.Subjects = strings.Split(subjects, ",")

	cfg.Templates.Dir = getEnv("TEMPLATES_DIR", "templates")
	cfg.Templates.DefaultExt = getEnv("TEMPLATES_EXT", ".html")
	cfg.Templates.CacheEnabled = getEnvBool("TEMPLATES_CACHE", true)

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
