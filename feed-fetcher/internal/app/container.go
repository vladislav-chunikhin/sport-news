package app

import (
	"context"
	"fmt"

	"github.com/vladislav-chunikhin/lib-go/pkg/logger"
	"github.com/vladislav-chunikhin/lib-go/pkg/shutdown"

	"github.com/vladislav-chunikhin/feed-fetcher/internal/fetcher"
	"github.com/vladislav-chunikhin/feed-fetcher/internal/producer"
	workerPkg "github.com/vladislav-chunikhin/feed-fetcher/internal/worker"

	"github.com/vladislav-chunikhin/feed-fetcher/internal/config"
)

type Container struct {
	config   *config.Config
	logger   logger.Logger
	producer *producer.Producer

	feedProvider fetcher.FeedProvider

	feedFetcher *fetcher.Fetcher
	worker      *workerPkg.Worker

	closers []func() error
}

func NewContainer(
	cfg *config.Config,
	log logger.Logger,
	producer *producer.Producer,
) *Container {
	return &Container{
		config:   cfg,
		logger:   log,
		producer: producer,
	}
}

func (c *Container) Config() *config.Config {
	return c.config
}

func (c *Container) Logger() logger.Logger {
	return c.logger
}

func (c *Container) Producer() *producer.Producer {
	return c.producer
}

func (c *Container) GetWorker(ctx context.Context) *workerPkg.Worker {
	if c.worker == nil {
		var err error
		c.worker, err = workerPkg.NewWorker(
			c.GetFetcher(ctx),
			c.config.Worker,
			c.logger,
		)

		if err != nil {
			panic(fmt.Errorf("failed to create worker: %v", err))
		}
	}

	return c.worker
}

func (c *Container) RegisterClosers(reg shutdown.Closer) {
	c.addCloser(func() error {
		c.worker.Stop()
		return nil
	})

	c.addCloser(func() error {
		return c.Producer().Close()
	})

	if len(c.closers) > 0 {
		for _, fn := range c.closers {
			reg.Add(fn)
		}
	}
}

func (c *Container) addCloser(fn func() error) {
	c.closers = append(c.closers, fn)
}
