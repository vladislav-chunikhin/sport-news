package main

import (
	"log"
	"time"

	app "github.com/vladislav-chunikhin/lib-go"
	"github.com/vladislav-chunikhin/lib-go/pkg/logger"

	"github.com/vladislav-chunikhin/feed-fetcher/internal/fetcher"
	"github.com/vladislav-chunikhin/feed-fetcher/internal/provider/htafc"
	"github.com/vladislav-chunikhin/feed-fetcher/internal/rabbitmq"

	"github.com/vladislav-chunikhin/feed-fetcher/internal/config"
	workerPkg "github.com/vladislav-chunikhin/feed-fetcher/internal/worker"
)

const defaultConfigFilePath = "./config/default.yaml"

func main() {
	a := app.NewApp()

	cfg := &config.Config{App: a.Config}
	err := a.LoadConfig(cfg, defaultConfigFilePath)
	if err != nil {
		log.Fatalf("config load error: %v", err)
	}

	config.AppConfigure(cfg)

	if err = a.Init(); err != nil {
		log.Fatalf("app init error: %v", err)
	}

	setupLocation(cfg, a.Logger)

	// RabbitMQ initialization
	var rabbitProducer *rabbitmq.Producer
	rabbitProducer, err = rabbitmq.NewProducer(&cfg.RabbitMq, a.Logger)
	if err != nil {
		a.Logger.Fatalf("failed to initialize rabbitmq producer: %v", err)
	}
	if err = rabbitProducer.DeclareQueues(); err != nil {
		a.Logger.Fatalf("failed to declare queues: %v", err)
	}

	// Feed provider initialization
	var feedProvider *htafc.Client
	feedProvider, err = htafc.New(&cfg.FeedProviders.Htafc, a.Logger)
	if err != nil {
		a.Logger.Fatalf("failed to initialize feed provider: %v", err)
	}

	// Fetcher initialization
	var feedFetcher *fetcher.Fetcher
	feedFetcher, err = fetcher.NewFetcher(rabbitProducer, feedProvider, a.Logger)
	if err != nil {
		a.Logger.Fatalf("failed to initialize feed fetcher: %v", err)
	}

	// Worker initialization
	var worker *workerPkg.Worker
	if worker, err = workerPkg.NewWorker(feedFetcher, cfg.Worker.Interval, a.Logger); err != nil {
		a.Logger.Fatalf("failed to initialize the worker: %v", err)
	}

	if err = worker.Run(a.Context); err != nil {
		a.Logger.Fatalf("failed to run the worker: %v", err)
	}

	// Closing resources
	a.Closer.Add(func() error {
		worker.Stop()
		return nil
	})

	a.Closer.Add(func() error {
		return rabbitProducer.Close()
	})

	// Run app
	a.Run()
}

func setupLocation(cfg *config.Config, logger logger.Logger) {
	if len(cfg.TimeZone) != 0 {
		loc, err := time.LoadLocation(cfg.TimeZone)
		if err == nil {
			logger.Debugf("init location from config: %s", cfg.TimeZone)
			time.Local = loc
		}
	}
}
