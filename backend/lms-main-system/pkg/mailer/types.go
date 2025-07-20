package mailer

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"net/smtp"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type EmailJob struct {
	From      string
	To        []string // Multiple recipients
	Subject   string
	Body      map[string]interface{}
	Footer    map[string]interface{}
	Templates struct {
		BodyPath   string
		FooterPath string
	}
}

type Mailer struct {
	log       *logrus.Logger
	smtpHost  string
	smtpPort  int
	username  string
	password  string
	tlsConfig *tls.Config
	templates *template.Template
}

// New creates a new Mailer instance
func New(log *logrus.Logger, host string, port int, username, password string) *Mailer {
	return &Mailer{
		log:      log,
		smtpHost: host,
		smtpPort: port,
		username: username,
		password: password,
		tlsConfig: &tls.Config{
			ServerName: host, // Important for TLS verification
		},
		templates: template.New("email_templates"),
	}
}

// StartDispatcher initializes worker pool
func (m *Mailer) StartDispatcher(workerCount int, jobs <-chan EmailJob) {
	for i := 0; i < workerCount; i++ {
		go m.worker(i, jobs)
	}
}

func (m *Mailer) worker(id int, jobs <-chan EmailJob) {
	for job := range jobs {
		m.log.Debugf("Worker %d processing email to %v", id, job.To)
		if err := m.SendEmail(job); err != nil {
			m.log.Errorf("Worker %d failed: %v", id, err)
		}
	}
}

// SendEmail sends an email using SMTP
func (m *Mailer) SendEmail(job EmailJob) error {
	body, err := m.parseTemplate(job.Templates.BodyPath, job.Body)
	if err != nil {
		return fmt.Errorf("template parsing failed: %w", err)
	}

	footer, err := m.parseTemplate(job.Templates.FooterPath, job.Footer)
	if err != nil {
		return fmt.Errorf("footer parsing failed: %w", err)
	}

	// Construct message
	msg := "From: " + job.From + "\r\n" +
		"To: " + joinEmails(job.To) + "\r\n" +
		"Subject: " + job.Subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n\r\n" +
		body + "\r\n" + footer

	auth := smtp.PlainAuth("", m.username, m.password, m.smtpHost)

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", m.smtpHost, m.smtpPort), m.tlsConfig)
	if err != nil {
		return fmt.Errorf("TLS connection failed: %w", err)
	}
	defer func(conn *tls.Conn) {
		err := conn.Close()
		if err != nil {
			m.log.Errorf("TLS connection failed: %v", err)
		}
	}(conn)

	client, err := smtp.NewClient(conn, m.smtpHost)
	if err != nil {
		return fmt.Errorf("SMTP client creation failed: %w", err)
	}
	defer func(client *smtp.Client) {
		err := client.Close()
		if err != nil {
			m.log.Errorf("SMTP client close failed: %v", err)
		}
	}(client)

	// Auth
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP auth failed: %w", err)
	}

	if err = client.Mail(job.From); err != nil {
		return fmt.Errorf("MAIL command failed: %w", err)
	}

	for _, to := range job.To {
		if err = client.Rcpt(to); err != nil {
			return fmt.Errorf("RCPT command failed for %s: %w", to, err)
		}
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("DATA command failed: %w", err)
	}

	_, err = w.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("writing message failed: %w", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("closing writer failed: %w", err)
	}

	m.log.Infof("Email successfully sent to %v", job.To)
	return nil
}
func (m *Mailer) parseTemplate(path string, data interface{}) (string, error) {
	if path == "" {
		return "", nil
	}

	tpl := m.templates.Lookup(filepath.Base(path))
	if tpl == nil {
		var err error
		tpl, err = template.ParseFiles(path)
		if err != nil {
			return "", err
		}
		m.templates = tpl
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func joinEmails(emails []string) string {
	var buf bytes.Buffer
	for i, email := range emails {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(email)
	}
	return buf.String()
}
