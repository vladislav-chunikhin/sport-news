package main

import (
	"log"
	"time"

	"github.com/vladislav-chunikhin/feed-transformer/internal/provider/htafc"
	"github.com/vladislav-chunikhin/feed-transformer/internal/rabbitmq"
	feedRepoPkg "github.com/vladislav-chunikhin/feed-transformer/internal/repository/feed"
	"github.com/vladislav-chunikhin/feed-transformer/internal/transformer"
	app "github.com/vladislav-chunikhin/lib-go"
	"github.com/vladislav-chunikhin/lib-go/pkg/logger"
	"github.com/vladislav-chunikhin/lib-go/pkg/mongodb"
	"github.com/vladislav-chunikhin/lib-go/pkg/redis"

	"github.com/vladislav-chunikhin/feed-transformer/internal/config"
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

	// MongoDB initialization
	var mongoClient *mongodb.Client
	mongoClient, err = mongodb.NewClient(a.Context, &cfg.Mongodb)
	if err != nil {
		a.Logger.Fatalf("failed to initialize the mongo client: %v", err)
	}
	if err = mongoClient.Ping(a.Context); err != nil {
		a.Logger.Fatalf("failed to ping the mongo db: %v", err)
	}

	// RabbitMQ initialization
	var rabbitConsumer *rabbitmq.Consumer
	rabbitConsumer, err = rabbitmq.NewConsumer(&cfg.RabbitMq, a.Logger)
	if err != nil {
		a.Logger.Fatalf("failed to initialize rabbitmq consumer: %v", err)
	}
	if err = rabbitConsumer.DeclareQueues(); err != nil {
		a.Logger.Fatalf("failed to declare queues: %v", err)
	}

	// Feed provider initialization
	var htafcClient *htafc.Htafc
	htafcClient, err = htafc.New(&cfg.FeedProviders.Htafc, cfg.RateLimiter.Limit, a.Logger)
	if err != nil {
		a.Logger.Fatalf("failed to initialize htafc client: %v", err)
	}

	// Redis initialization
	var cacheClient *redis.Client
	cacheClient, err = redis.NewClient(&cfg.Redis)
	if err != nil {
		a.Logger.Fatalf("failed to initialize redis client: %v", err)
	}

	// Initializing application base layers
	var feedRepository *feedRepoPkg.Repository
	feedRepository, err = feedRepoPkg.NewRepository(a.Context, mongoClient, a.Logger)
	if err != nil {
		a.Logger.Fatalf("failed to initialize repository: %v", err)
	}

	transformerService := transformer.NewService(a.Context, rabbitConsumer, htafcClient, feedRepository, cacheClient, a.Logger)

	// Running consumer
	go func() {
		if err = transformerService.Run(a.Context); err != nil {
			a.Logger.Fatalf("failed to run consumer: %v", err)
		}
	}()

	// Closing resources
	a.Closer.Add(func() error {
		return mongoClient.Close(a.Context, a.Logger)
	})
	a.Closer.Add(func() error {
		return rabbitConsumer.Close()
	})
	a.Closer.Add(func() error {
		return cacheClient.Close(a.Logger)
	})

	// Run app
	a.Run()
}

func setupLocation(cfg *config.Config, l logger.Logger) {
	if len(cfg.TimeZone) != 0 {
		loc, err := time.LoadLocation(cfg.TimeZone)
		if err == nil {
			l.Debugf("init location from config: %s", cfg.TimeZone)
			time.Local = loc
		}
	}
}
