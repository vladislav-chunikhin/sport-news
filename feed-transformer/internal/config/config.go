package config

import (
	"fmt"
	"time"

	app "github.com/vladislav-chunikhin/lib-go"
	"github.com/vladislav-chunikhin/lib-go/pkg/circuitbreaker"
	"github.com/vladislav-chunikhin/lib-go/pkg/limiter"
	"github.com/vladislav-chunikhin/lib-go/pkg/mongodb"
	"github.com/vladislav-chunikhin/lib-go/pkg/redis"
)

type Config struct {
	App                    *app.Config
	LogLevel               string              `yaml:"logLevel"`
	TimeZone               string              `yaml:"timeZone"`
	HTTPServerReadTimeout  int64               `yaml:"httpServerReadTimeout"`
	HTTPServerWriteTimeout int64               `yaml:"httpServerWriteTimeout"`
	HTTPServerPort         int                 `yaml:"httpServerPort"`
	HTTPDebugPort          int                 `yaml:"httpDebugPort"`
	Mongodb                mongodb.Config      `yaml:"mongodb"`
	RabbitMq               RabbitConfig        `yaml:"rabbitmq"`
	FeedProviders          FeedProvidersConfig `yaml:"feedProviders"`
	Redis                  redis.Config        `yaml:"redis"`
	RateLimiter            limiter.Config      `yaml:"rateLimiter"`
}

type RabbitConfig struct {
	URL       string          `yaml:"url"`
	Timeout   time.Duration   `yaml:"timeout"`
	Queues    QueuesConfig    `yaml:"queues"`
	Consumers ConsumersConfig `yaml:"consumers"`
}

type QueuesConfig struct {
	Htafc QueueConfig `yaml:"htafc"`
}

type ConsumersConfig struct {
	Htafc ConsumerConfig `yaml:"htafc"`
}

type QueueConfig struct {
	Name       string `yaml:"name"`
	Durable    bool   `yaml:"durable"`
	AutoDelete bool   `yaml:"autoDelete"`
	Exclusive  bool   `yaml:"exclusive"`
	NoWait     bool   `yaml:"noWait"`
}

type ConsumerConfig struct {
	Name      string `yaml:"name"`
	AutoAck   bool   `yaml:"autoAck"`
	Exclusive bool   `yaml:"exclusive"`
	NoLocal   bool   `yaml:"noLocal"`
	NoWait    bool   `yaml:"noWait"`
}

type FeedProvidersConfig struct {
	Htafc FeedProviderConfig `yaml:"htafc"`
}

type FeedProviderConfig struct {
	Address        string                `yaml:"address"`
	Timeout        time.Duration         `yaml:"timeout"`
	CircuitBreaker circuitbreaker.Config `yaml:"circuitBreaker"`
}

func AppConfigure(cfg *Config) {
	cfg.App.LoggerLevel = cfg.LogLevel
	cfg.App.HTTPServerReadTimeout = cfg.HTTPServerReadTimeout
	cfg.App.HTTPServerWriteTimeout = cfg.HTTPServerWriteTimeout

	if cfg.HTTPServerPort != 0 {
		cfg.App.Port = fmt.Sprintf("%d", cfg.HTTPServerPort)
	}
	if cfg.HTTPDebugPort != 0 {
		cfg.App.MtPort = fmt.Sprintf("%d", cfg.HTTPDebugPort)
	}
}
