package config

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

var configurationTestFile = "test_configs/config_test.yaml"

func TestNewConfigurationUpdater(t *testing.T) {
	testExample := &ConfigurationUpdater{
		updatePeriodicity: time.Second,
		filePath:          configurationTestFile,
	}

	result := NewConfigurationUpdater(testExample.updatePeriodicity, testExample.filePath)

	assert.Equal(t, testExample, result, "Should create correct pointer to ConfigurationUpdater object")
}

func TestConfigurationUpdater_GetConfiguration(t *testing.T) {
	config := Configuration{}
	config, _ = config.LoadConfigurationFromFile(configurationTestFile)
	cu := NewConfigurationUpdater(time.Second, configurationTestFile)

	result := cu.GetConfiguration()

	assert.Equal(t, &config, result, "Should return equal config file")
}

func TestConfigurationUpdater_WatchConfigurationFile(t *testing.T) {
	config := Configuration{}
	config, _ = config.LoadConfigurationFromFile(configurationTestFile)

	wg := sync.WaitGroup{}
	wg.Add(1)

	var result Configuration
	mockFunc := func(conf Configuration) {
		result = conf
		wg.Done()
	}

	cu := NewConfigurationUpdater(time.Millisecond*250, configurationTestFile)
	cu.WatchConfigurationFile(mockFunc)
	wg.Wait()

	assert.Equal(t, true, cu.enableConfigWatching, "Config watching should be enabled")
	assert.Equal(t, config, result, "Should process equal config file")
}

func TestConfigurationUpdater_UnWatchConfigurationFile(t *testing.T) {

	cu := NewConfigurationUpdater(time.Millisecond*250, configurationTestFile)
	cu.UnWatchConfigurationFile()

	assert.Equal(t, false, cu.enableConfigWatching, "Config watching should be disabled")

}
