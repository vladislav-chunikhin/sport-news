package config

import (
	"testing"

	"github.com/stretchr/testify/require"
	app "github.com/vladislav-chunikhin/lib-go"
)

func TestAppConfigure(t *testing.T) {
	appConfig := &app.Config{
		AppName:                "testApp",
		AppEnvironment:         "testEnv",
		AppVersion:             "testVersion",
		LoggerLevel:            "info",
		HTTPServerReadTimeout:  100,
		HTTPServerWriteTimeout: 100,
	}
	AppConfigure(&Config{
		App:                    appConfig,
		LogLevel:               "debug",
		HTTPServerWriteTimeout: 200,
		HTTPServerReadTimeout:  200,
	})

	require.Equal(t, "debug", appConfig.LoggerLevel)
	require.Equal(t, int64(200), appConfig.HTTPServerWriteTimeout)
	require.Equal(t, int64(200), appConfig.HTTPServerReadTimeout)
	require.Equal(t, "testApp", appConfig.AppName)
	require.Equal(t, "testEnv", appConfig.AppEnvironment)
	require.Equal(t, "testVersion", appConfig.AppVersion)
}
