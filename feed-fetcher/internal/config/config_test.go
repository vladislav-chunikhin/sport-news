package config

import (
	"testing"

	"github.com/stretchr/testify/require"
	app "github.com/vladislav-chunikhin/lib-go"
)

func TestAppConfigure(t *testing.T) {
	testCases := []struct {
		name           string
		cfg            *Config
		expectedAppCfg *app.Config
		expectedErrMsg string
	}{
		{
			name: "normal case",
			cfg: &Config{
				App:                    &app.Config{},
				LogLevel:               "debug",
				HTTPServerReadTimeout:  60,
				HTTPServerWriteTimeout: 120,
				HTTPServerPort:         443,
				HTTPDebugPort:          80,
			},
			expectedAppCfg: &app.Config{
				LoggerLevel:            "debug",
				HTTPServerReadTimeout:  60,
				HTTPServerWriteTimeout: 120,
				Port:                   "443",
				MtPort:                 "80",
			},
		},
		{
			name:           "nil config",
			cfg:            nil,
			expectedErrMsg: "nil config",
		},
		{
			name:           "nil app config",
			cfg:            &Config{},
			expectedErrMsg: "nil config",
		},
	}

	for _, tc := range testCases {
		var err error
		if err = AppConfigure(tc.cfg); err != nil {
			require.Equal(t, tc.expectedErrMsg, err.Error())
			continue
		}
		require.NoError(t, err)
		require.Equal(t, tc.expectedAppCfg, tc.cfg.App)
	}
}
