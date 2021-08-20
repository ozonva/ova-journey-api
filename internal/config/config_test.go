package config

import (
	"github.com/stretchr/testify/assert"
	"io/fs"
	"testing"
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
			configPath: "test_configs/config_test.json",
			result: testLoadConfigurationFromFileResult{
				conf: Configuration{
					Host: "localhost",
					Port: 777,
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
