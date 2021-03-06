package config

import (
	"reflect"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// ConfigurationUpdater type represents object for check Configuration updates in JSON file.
type ConfigurationUpdater struct {
	updatePeriodicity    time.Duration
	filePath             string
	configuration        Configuration
	fileMutex            sync.Mutex
	enableConfigWatching bool
}

// NewConfigurationUpdater - create ConfigurationUpdater object for check Configuration updates in JSON file.
// For run periodic checking of JSON configuration file use ConfigurationUpdater.WatchConfigurationFile.
// UpdatePeriodicity - time period between repeatable reads of configuration file.
// ConfigurationFilePath - path to JSON file with app configuration.
func NewConfigurationUpdater(updatePeriodicity time.Duration, configurationFilePath string) *ConfigurationUpdater {
	return &ConfigurationUpdater{
		updatePeriodicity: updatePeriodicity,
		filePath:          configurationFilePath,
	}
}

// GetConfiguration - get actual configuration
func (cu *ConfigurationUpdater) GetConfiguration() *Configuration {
	if (cu.configuration == Configuration{}) {
		cu.loadConfiguration()
	}
	return &cu.configuration
}

// WatchConfigurationFile - run infinite task, checking updates in configuration file
func (cu *ConfigurationUpdater) WatchConfigurationFile(configurationUpdateHandler func(conf Configuration)) {
	cu.enableConfigWatching = true
	go func() {
		for cu.enableConfigWatching {
			if cu.loadConfiguration() {
				configurationUpdateHandler(cu.configuration)
			}
			time.Sleep(cu.updatePeriodicity)
		}
	}()
}

// UnWatchConfigurationFile - stop infinite task, checking updates in configuration file
func (cu *ConfigurationUpdater) UnWatchConfigurationFile() {
	cu.enableConfigWatching = false
}

// loadConfiguration - update configuration from file,
// returns true, if configuration was updated
func (cu *ConfigurationUpdater) loadConfiguration() bool {
	cu.fileMutex.Lock()
	defer cu.fileMutex.Unlock()

	newConfig, err := cu.configuration.LoadConfigurationFromFile(cu.filePath)
	if err != nil {
		log.Err(err).Msg("Error occurred while reading configuration file")
	}
	if (newConfig == Configuration{}) {
		log.Fatal().Msg("Configuration is empty")
	}
	if !reflect.DeepEqual(&cu.configuration, &newConfig) {
		cu.configuration = newConfig
		log.Info().Msg("Configuration updated")
		return true
	}

	return false
}
