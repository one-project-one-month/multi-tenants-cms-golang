package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/wneessen/go-mail"
)

type EmailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func main() {
	// Load configuration from environment
	smtpHost := getEnv("SMTP_HOST", "smtp.gmail.com")
	smtpPort := getEnvInt("SMTP_PORT", 587)
	smtpUser := getEnv("SMTP_USER", "")
	smtpPass := getEnv("SMTP_PASSWORD", "")
	fromAddr := getEnv("FROM_ADDR", smtpUser)
	natsUrl := getEnv("NATS_URL", "nats://localhost:4222")

	// Initialize SMTP client
	smtpClient, err := mail.NewClient(
		smtpHost,
		mail.WithPort(smtpPort),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(smtpUser),
		mail.WithPassword(smtpPass),
		mail.WithTLSPolicy(mail.TLSMandatory),
	)
	if err != nil {
		log.Fatalf("Failed to create SMTP client: %v", err)
	}

	// Connect to NATS
	nc, err := nats.Connect(natsUrl)
	if err != nil {
		log.Fatalf("NATS connection failed: %v", err)
	}
	defer nc.Close()

	// Create JetStream context
	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatalf("JetStream init failed: %v", err)
	}

	// Create stream configuration
	streamConfig := jetstream.StreamConfig{
		Name:      "EMAILS",
		Subjects:  []string{"email.verification", "email.notification"},
		Retention: jetstream.WorkQueuePolicy,
		Storage:   jetstream.FileStorage,
	}

	// Create or update stream
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := js.CreateOrUpdateStream(ctx, streamConfig)
	if err != nil {
		log.Fatalf("Stream creation failed: %v", err)
	}

	// Create consumer
	consumerConfig := jetstream.ConsumerConfig{
		Durable:       "email-processor",
		AckPolicy:     jetstream.AckExplicitPolicy,
		DeliverPolicy: jetstream.DeliverAllPolicy,
	}

	consumer, err := stream.CreateOrUpdateConsumer(ctx, consumerConfig)
	if err != nil {
		log.Fatalf("Consumer creation failed: %v", err)
	}

	log.Println("Email service started. Waiting for messages...")

	_, err = consumer.Consume(func(msg jetstream.Msg) {
		log.Printf("Received message on subject: %s", msg.Subject())

		var emailReq EmailRequest
		if err := json.Unmarshal(msg.Data(), &emailReq); err != nil {
			log.Printf("Failed to parse message: %v", err)
			msg.Nak()
			return
		}

		if err := sendEmail(smtpClient, fromAddr, emailReq); err != nil {
			log.Printf("Failed to send email: %v", err)
			msg.Nak()
			return
		}

		log.Printf("Email sent to %s", emailReq.To)
		msg.Ack()
	})
	if err != nil {
		log.Fatalf("Consume failed: %v", err)
	}

	select {}
}

func sendEmail(client *mail.Client, from string, req EmailRequest) error {
	m := mail.NewMsg()
	if err := m.From(from); err != nil {
		return fmt.Errorf("failed to set From address: %w", err)
	}
	if err := m.To(req.To); err != nil {
		return fmt.Errorf("failed to set To address: %w", err)
	}

	m.Subject(req.Subject)
	m.SetBodyString(mail.TypeTextHTML, req.Body)

	if err := client.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
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
