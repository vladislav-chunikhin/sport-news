package main

import (
	"log"
	"time"

	baseApp "github.com/vladislav-chunikhin/lib-go"
	"github.com/vladislav-chunikhin/lib-go/pkg/di"
	"github.com/vladislav-chunikhin/lib-go/pkg/logger"

	"github.com/vladislav-chunikhin/feed-fetcher/internal/app"
	"github.com/vladislav-chunikhin/feed-fetcher/internal/config"
	producerPkg "github.com/vladislav-chunikhin/feed-fetcher/internal/producer"
)

const defaultConfigFilePath = "./config/default.yaml"

func main() {
	a := baseApp.NewApp()

	// App initialization
	cfg := &config.Config{App: a.Config}
	err := a.LoadConfig(cfg, defaultConfigFilePath)
	if err != nil {
		log.Fatalf("config load error: %v", err)
	}

	if err = config.AppConfigure(cfg); err != nil {
		log.Fatalf("config init error: %v", err)
	}

	if err = a.Init(); err != nil {
		log.Fatalf("app init error: %v", err)
	}

	setupLocation(cfg, a.Logger)

	// RabbitMQ initialization
	var producer *producerPkg.Producer
	producer, err = producerPkg.NewProducer(&cfg.RabbitMq, a.Logger)
	if err != nil {
		a.Logger.Fatalf("failed to initialize rabbitmq producer: %v", err)
	}
	if err = producer.DeclareQueues(); err != nil {
		a.Logger.Fatalf("failed to declare queues: %v", err)
	}

	container := app.NewContainer(cfg, a.Logger, producer)
	container.RegisterClosers(a.Closer)

	if err = di.Build(a.Context, container); err != nil {
		defer a.Closer.CloseAll(a.Logger)
		a.Logger.Fatalf(err.Error())
	}

	if err = container.GetWorker(a.Context).Run(a.Context); err != nil {
		a.Logger.Fatalf("failed to run worker: %v", err)
	}

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
