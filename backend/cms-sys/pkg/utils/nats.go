package utils

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

var (
	natsClient  *nats.Conn
	js          nats.JetStreamContext
	initOnce    sync.Once
	initialized bool
)

// InitNats initializes NATS connection and JetStream with proper stream handling
func InitNats() error {
	var initErr error
	initOnce.Do(func() {
		natsUrl := GetEnv("NATS_URL", "nats://localhost:4222")

		opts := []nats.Option{
			nats.Name("EMAIL_SERVICE"),
			nats.MaxReconnects(-1),
			nats.ReconnectWait(2 * time.Second),
			nats.Timeout(10 * time.Second),
			nats.PingInterval(30 * time.Second),
			nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
				log.Printf("NATS disconnected: %v", err)
			}),
			nats.ReconnectHandler(func(nc *nats.Conn) {
				log.Println("NATS reconnected")
				if js, err := nc.JetStream(); err == nil {
					err := initJetStream(js)
					if err != nil {
						log.Printf("NATS stream initialization failed: %v", err)
						return
					}
				}
			}),
		}

		// Connect to NATS
		nc, err := nats.Connect(natsUrl, opts...)
		if err != nil {
			initErr = fmt.Errorf("failed to connect to NATS: %w", err)
			return
		}

		// Initialize JetStream
		js, err = nc.JetStream(nats.PublishAsyncMaxPending(256))
		if err != nil {
			initErr = fmt.Errorf("failed to initialize JetStream: %w", err)
			nc.Close()
			return
		}

		natsClient = nc
		initialized = true

		// Initialize streams with conflict resolution
		if err := initJetStream(js); err != nil {
			initErr = fmt.Errorf("failed to initialize streams: %w", err)
			return
		}

		log.Println("NATS and JetStream successfully initialized")
	})

	return initErr
}

func initJetStream(js nats.JetStreamContext) error {
	streamConfig := &nats.StreamConfig{
		Name:         "EMAILS",
		Subjects:     []string{"email.verification", "email.notification", "page.approval"},
		Retention:    nats.WorkQueuePolicy,
		Storage:      nats.FileStorage,
		MaxAge:       24 * time.Hour,
		MaxMsgs:      1000000,
		MaxBytes:     1024 * 1024 * 1024, // 1GB
		MaxConsumers: 10,                 // Add this line
		Duplicates:   2 * time.Minute,
	}

	existingInfo, err := js.StreamInfo("EMAILS")
	if err == nil {
		if existingInfo.Config.MaxConsumers != streamConfig.MaxConsumers {
			log.Printf("Warning: Existing stream has MaxConsumers=%d, but config wants %d. Cannot change this field.",
				existingInfo.Config.MaxConsumers, streamConfig.MaxConsumers)
			return nil
		}

		_, err = js.UpdateStream(streamConfig)
		if err != nil {
			return fmt.Errorf("failed to update existing stream: %w", err)
		}
		log.Println("Updated existing EMAILS stream configuration")
		return nil
	}

	_, err = js.AddStream(streamConfig)
	if err != nil {
		return fmt.Errorf("failed to create new stream: %w", err)
	}

	log.Println("Created new EMAILS stream")
	return nil
}

// GetJetStream safely returns the JetStream context
func GetJetStream() (nats.JetStreamContext, error) {
	if !initialized {
		return nil, fmt.Errorf("JetStream not initialized")
	}
	return js, nil
}

// GetNatsConnection returns the core NATS connection
func GetNatsConnection() (*nats.Conn, error) {
	if !initialized {
		return nil, fmt.Errorf("NATS not initialized")
	}
	return natsClient, nil
}

// PublishMessage publishes a message to NATS with JetStream
func PublishMessage(subject string, data []byte) error {
	if !initialized {
		return fmt.Errorf("NATS not initialized")
	}

	js, err := GetJetStream()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = js.Publish(subject, data, nats.Context(ctx))
	return err
}

// PublishMessageAsync publishes a message asynchronously with confirmation
//func PublishMessageAsync(subject string, data []byte, ackHandler func(*nats.PubAck, error)) error {
//	if !initialized {
//		return fmt.Errorf("NATS not initialized")
//	}
//
//	js, err := GetJetStream()
//	if err != nil {
//		return err
//	}
//
//	_, err = js.PublishAsync(subject, data, ackHandler)
//	return err
//}

// Subscribe creates a JetStream subscription
func Subscribe(subject string, handler func(*nats.Msg)) (*nats.Subscription, error) {
	if !initialized {
		return nil, fmt.Errorf("NATS not initialized")
	}

	js, err := GetJetStream()
	if err != nil {
		return nil, err
	}

	return js.Subscribe(subject, handler,
		nats.Durable("email-service"),
		nats.ManualAck(),
		nats.AckExplicit(),
		nats.DeliverAll(),
		nats.MaxDeliver(5),
	)
}

// IsConnected checks if NATS is connected
func IsConnected() bool {
	return initialized && natsClient != nil && natsClient.IsConnected()
}

func Close() error {
	if natsClient != nil {
		if err := natsClient.Drain(); err != nil {
			return fmt.Errorf("failed to drain NATS connection: %w", err)
		}
		natsClient.Close()
		initialized = false
	}
	return nil
}
