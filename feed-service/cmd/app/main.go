package main

import (
	"log"
	"time"

	app "github.com/vladislav-chunikhin/lib-go"
	"github.com/vladislav-chunikhin/lib-go/pkg/logger"
	"github.com/vladislav-chunikhin/lib-go/pkg/mongodb"

	feedRepoPkg "github.com/vladislav-chunikhin/feed-service/internal/repository/feed"
	feedServicePkg "github.com/vladislav-chunikhin/feed-service/internal/service/feed"
	"github.com/vladislav-chunikhin/feed-service/internal/transport"

	"github.com/vladislav-chunikhin/feed-service/internal/config"
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

	// Initializing application base layers
	var feedRepository *feedRepoPkg.Repository
	feedRepository, err = feedRepoPkg.NewRepository(a.Context, mongoClient, a.Logger)
	if err != nil {
		a.Logger.Fatalf("failed to initialize repository: %v", err)
	}

	feedService := feedServicePkg.NewService(feedRepository)
	api := transport.NewAPI(feedService, a.Logger)

	a.HttpServer.Handler = api.Router()

	// Closing resources
	a.Closer.Add(func() error {
		return mongoClient.Close(a.Context, a.Logger)
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
