package config

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testLoadConfigurationFromFileResult struct {
	conf Configuration
	err  error
}

func TestConfiguration_LoadConfigurationFromFile(t *testing.T) {
	testTable := []struct {
		name       string
		config     Configuration
		configPath string
		result     testLoadConfigurationFromFileResult
	}{
		{
			name:       "correct configuration",
			configPath: "test_configs/config_test.yaml",
			result: testLoadConfigurationFromFileResult{
				conf: Configuration{
					Project: &ProjectConfiguration{
						Name:    "Journey API for Amazon Voice Assistant Project",
						Version: "v0.0.0-test",
					},
					GRPC: &EndpointConfiguration{
						Host: "127.0.0.1",
						Port: 9090,
					},
					Gateway: &EndpointConfiguration{
						Host: "127.0.0.1",
						Port: 8080,
					},
					Database: &DatabaseConfiguration{
						Host:     "database",
						Port:     5432,
						User:     "postgres",
						Password: "postgres",
						Name:     "ova_journey_api",
						SslMode:  "disable",
						Driver:   "pgx",
					},
				},
				err: nil,
			},
		},
		{
			name:       "non existing file",
			configPath: "test_configs/config_non_exist_test.json",
			result: testLoadConfigurationFromFileResult{
				conf: Configuration{},
				err:  fs.ErrNotExist,
			},
		},
	}

	conf := Configuration{}
	for _, testCase := range testTable {
		result, err := conf.LoadConfigurationFromFile(testCase.configPath)
		assert.ErrorIs(t, err, testCase.result.err, testCase.name)
		assert.Equal(t, testCase.result.conf, result, testCase.name)
	}
}
