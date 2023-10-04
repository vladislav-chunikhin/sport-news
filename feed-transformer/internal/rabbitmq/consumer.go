package rabbitmq

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/vladislav-chunikhin/feed-transformer/internal/config"
	"github.com/vladislav-chunikhin/lib-go/pkg/logger"
)

const retryDelay = time.Second * 5

type Consumer struct {
	conn       *amqp.Connection
	cfg        *config.RabbitConfig
	logger     logger.Logger
	consumerCh *amqp.Channel
}

func NewConsumer(cfg *config.RabbitConfig, logger logger.Logger) (*Consumer, error) {
	if cfg == nil {
		return nil, fmt.Errorf("nil cfg")
	}

	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to create amqp connection: %w", err)
	}

	var consumerCh *amqp.Channel
	consumerCh, err = conn.Channel()
	if err != nil {
		if err = conn.Close(); err != nil {
			logger.Warnf("failed to close connection: %v", err)
		}
		return nil, fmt.Errorf("failed to create consumer channel: %w", err)
	}

	return &Consumer{
		conn:       conn,
		cfg:        cfg,
		logger:     logger,
		consumerCh: consumerCh,
	}, nil
}

func (c *Consumer) DeclareQueues() error {
	_, err := c.consumerCh.QueueDeclare(
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

func (c *Consumer) ConsumeHtafcFeed(ctx context.Context, handler func(amqp.Delivery) error) error {
	msgs, err := c.consumerCh.Consume(
		c.cfg.Queues.Htafc.Name,
		c.cfg.Consumers.Htafc.Name,
		c.cfg.Consumers.Htafc.AutoAck,
		c.cfg.Consumers.Htafc.Exclusive,
		c.cfg.Consumers.Htafc.NoLocal,
		c.cfg.Consumers.Htafc.NoWait,
		nil,
	)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			c.logger.Debugf("amqp consumer work stopped")
			return nil
		case msg := <-msgs:
			if err = handler(msg); err != nil {
				msg.Reject(true)
				c.logger.Debugf("message was handled with problems: %v", err)
				<-time.After(retryDelay)
			} else {
				msg.Ack(false)
				c.logger.Debugf("message was handled successfully")
			}
		}
	}
}

func (c *Consumer) Close() error {
	if err := c.consumerCh.Close(); err != nil {
		return err
	}

	if err := c.conn.Close(); err != nil {
		return err
	}

	c.logger.Debugf("rabbitmq disconnecting...")
	return nil
}
