package utils

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

var natsClient *nats.Conn

// InitNats initializes the NATS connection
func InitNats() error {
	natsUrl := GetEnv("NATS_URL", "nats://localhost:4222")

	opts := []nats.Option{
		nats.Name("cms-sys"),
		nats.MaxReconnects(10),
		nats.ReconnectWait(2 * time.Second),
		nats.Timeout(5 * time.Second),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			log.Printf("NATS disconnected: %v", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Println("NATS reconnected")
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			log.Println("NATS connection closed")
		}),
	}

	nc, err := nats.Connect(natsUrl, opts...)
	if err != nil {
		log.Printf("Failed to connect to NATS: %v", err)
		return err
	}

	natsClient = nc
	log.Println("Successfully connected to NATS")
	return nil
}

// GetNatsConnection returns the NATS connection
func GetNatsConnection() *nats.Conn {
	return natsClient
}

// PublishMessage publishes a message to NATS
func PublishMessage(subject string, data []byte) error {
	if natsClient == nil {
		return fmt.Errorf("NATS client is not initialized")
	}

	if !natsClient.IsConnected() {
		return fmt.Errorf("NATS is not connected")
	}

	return natsClient.Publish(subject, data)
}

// IsNatsConnected checks if NATS is connected
func IsNatsConnected() bool {
	if natsClient == nil {
		return false
	}
	return natsClient.IsConnected()
}

// CloseNats closes the NATS connection
func CloseNats() error {
	if natsClient != nil {
		natsClient.Close()
		log.Println("NATS connection closed")
	}
	return nil
}

// DrainNats drains the NATS connection
func DrainNats() error {
	if natsClient != nil {
		return natsClient.Drain()
	}
	return nil
}
