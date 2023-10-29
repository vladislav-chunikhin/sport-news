package producer

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/vladislav-chunikhin/feed-fetcher/internal/config"
	"github.com/vladislav-chunikhin/lib-go/pkg/logger"
)

const (
	msgContentType  = "application/json"
	expirationValue = "3600000" // 1 hour

	typeKey   = "type"
	htafcType = "htafc"
	appID     = "feed-fetcher"
)

type Producer struct {
	conn       *amqp.Connection
	cfg        *config.RabbitConfig
	logger     logger.Logger
	producerCh *amqp.Channel
}

func NewProducer(cfg *config.RabbitConfig, logger logger.Logger) (*Producer, error) {
	if cfg == nil {
		return nil, fmt.Errorf("nil cfg")
	}

	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to create amqp connection: %w", err)
	}

	var producerCh *amqp.Channel
	producerCh, err = conn.Channel()
	if err != nil {
		if err = conn.Close(); err != nil {
			logger.Warnf("failed to close connection: %v", err)
		}
		return nil, fmt.Errorf("failed to create producer channel: %w", err)
	}

	return &Producer{
		conn:       conn,
		cfg:        cfg,
		logger:     logger,
		producerCh: producerCh,
	}, nil
}

func (c *Producer) DeclareQueues() error {
	_, err := c.producerCh.QueueDeclare(
		c.cfg.Queues.Htafc.Name,
		c.cfg.Queues.Htafc.Durable,
		c.cfg.Queues.Htafc.AutoDelete,
		c.cfg.Queues.Htafc.Exclusive,
		c.cfg.Queues.Htafc.NoWait,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *Producer) PublishHtafcFeed(ctx context.Context, message []byte) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, c.cfg.Timeout)
	defer cancel()

	err := c.producerCh.PublishWithContext(
		timeoutCtx,
		"",
		c.cfg.Queues.Htafc.Name,
		false, // Mandatory flag (false means the message can be silently dropped)
		false, // Immediate flag (false means the message will be queued immediately)
		amqp.Publishing{
			ContentType:  msgContentType,
			Body:         message,
			DeliveryMode: amqp.Persistent, // Message delivery mode (Persistent makes the message durable)
			MessageId:    uuid.New().String(),
			Headers:      amqp.Table{typeKey: htafcType},
			Expiration:   expirationValue,
			AppId:        appID,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *Producer) Close() error {
	if err := c.producerCh.Close(); err != nil {
		return err
	}

	if err := c.conn.Close(); err != nil {
		return err
	}

	c.logger.Debugf("rabbitmq disconnecting...")
	return nil
}
