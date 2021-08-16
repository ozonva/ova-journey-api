package config

import (
	"log"
	"sync"
	"time"
)

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
func (cu *ConfigurationUpdater) WatchConfigurationFile() {
	cu.enableConfigWatching = true
	go func() {
		for cu.enableConfigWatching {
			cu.loadConfiguration()
			time.Sleep(cu.updatePeriodicity)
		}
	}()
}

// UnWatchConfigurationFile - stop infinite task, checking updates in configuration file
func (cu *ConfigurationUpdater) UnWatchConfigurationFile() {
	cu.enableConfigWatching = false
}

func (cu *ConfigurationUpdater) loadConfiguration() {
	cu.fileMutex.Lock()
	defer cu.fileMutex.Unlock()

	newConfig, err := cu.configuration.LoadConfigurationFromFile(cu.filePath)
	if err != nil {
		log.Printf("Error occured while reading configuration file: %s", err)
	}
	if (newConfig == Configuration{}) {
		log.Panicf("Configration is empty")
	}
	if cu.configuration != newConfig {
		cu.configuration = newConfig
		log.Printf("Configuration updated: %v", cu.configuration)
	}
}
