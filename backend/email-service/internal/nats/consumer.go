package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

type EmailRequest struct {
	To       string                 `json:"to"`
	Template string                 `json:"template"`
	Data     map[string]interface{} `json:"data"`
}

type Consumer struct {
	conn         *nats.Conn
	js           jetstream.JetStream
	stream       jetstream.Stream
	cons         jetstream.Consumer
	logger       *zap.Logger
	emailCh      chan<- EmailRequest
	streamName   string
	subject      string
	consumerName string
}

func NewConsumer(
	conn *nats.Conn,
	logger *zap.Logger,
	emailCh chan<- EmailRequest,
	streamName, subject, consumerName string,
) (*Consumer, error) {
	js, err := jetstream.New(conn)
	if err != nil {
		return nil, fmt.Errorf("jetstream init error: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stream, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:      streamName,
		Subjects:  []string{subject},
		Storage:   jetstream.FileStorage,
		Retention: jetstream.WorkQueuePolicy,
		MaxAge:    24 * time.Hour,
	})
	if err != nil {
		return nil, fmt.Errorf("stream creation error: %w", err)
	}

	cons, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:       consumerName,
		AckPolicy:     jetstream.AckExplicitPolicy,
		DeliverPolicy: jetstream.DeliverAllPolicy,
		MaxDeliver:    3,
		AckWait:       30 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("consumer creation error: %w", err)
	}

	return &Consumer{
		conn:         conn,
		js:           js,
		stream:       stream,
		cons:         cons,
		logger:       logger,
		emailCh:      emailCh,
		streamName:   streamName,
		subject:      subject,
		consumerName: consumerName,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	c.logger.Info("starting NATS consumer",
		zap.String("stream", c.streamName),
		zap.String("consumer", c.consumerName))

	consCtx, err := c.cons.Consume(func(msg jetstream.Msg) {
		c.logger.Debug("Received NATS message",
			zap.String("data", string(msg.Data())))
		var req EmailRequest
		if err := json.Unmarshal(msg.Data(), &req); err != nil {
			c.logger.Error("failed to unmarshal email request",
				zap.Error(err),
				zap.String("data", string(msg.Data())))
			if err := msg.Ack(); err != nil {
				c.logger.Error("failed to ack invalid message", zap.Error(err))
			}
			return
		}

		select {
		case c.emailCh <- req:
			if err := msg.Ack(); err != nil {
				c.logger.Error("failed to ack processed message", zap.Error(err))
			}
		case <-ctx.Done():
			if err := msg.Nak(); err != nil {
				c.logger.Error("failed to nak message on shutdown", zap.Error(err))
			}
			return
		default:
			if err := msg.Term(); err != nil {
				c.logger.Error("failed to terminate message", zap.Error(err))
			}
			c.logger.Warn("email channel full, message terminated",
				zap.String("email", req.To))
		}
	})
	if err != nil {
		return fmt.Errorf("jetstream consume error: %w", err)
	}

	<-ctx.Done()
	c.logger.Info("stopping NATS consumer")
	consCtx.Stop()
	return nil
}
