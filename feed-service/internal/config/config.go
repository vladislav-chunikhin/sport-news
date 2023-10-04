package config

import (
	"fmt"

	app "github.com/vladislav-chunikhin/lib-go"
	"github.com/vladislav-chunikhin/lib-go/pkg/mongodb"
)

type Config struct {
	App                    *app.Config
	LogLevel               string         `yaml:"logLevel"`
	TimeZone               string         `yaml:"timeZone"`
	HTTPServerReadTimeout  int64          `yaml:"httpServerReadTimeout"`
	HTTPServerWriteTimeout int64          `yaml:"httpServerWriteTimeout"`
	HTTPServerPort         int            `yaml:"httpServerPort"`
	HTTPDebugPort          int            `yaml:"httpDebugPort"`
	Mongodb                mongodb.Config `yaml:"mongodb"`
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
