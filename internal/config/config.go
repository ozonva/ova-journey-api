package config

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	Host string `json:"host"`
	Port uint   `json:"port"`
}

func (c *Configuration) LoadConfiguration(path string) {
	updateConfig := func(path string) (Configuration, error) {
		var err error
		configuration := Configuration{}

		file, err := os.Open(path)
		if err != nil {
			return configuration, err
		}

		defer func() {
			defErr := file.Close()
			if defErr != nil {
				err = defErr
			}
		}()

		decoder := json.NewDecoder(file)
		err = decoder.Decode(&configuration)

		return configuration, err
	}

	newConfig, err := updateConfig(path)
	if err != nil {
		log.Panicf("Error occured while reading configuration file: %s", err)
	}
	if (newConfig == Configuration{}) {
		log.Panicf("Configration file is empty")
	}

	if *c != newConfig {
		*c = newConfig
		log.Printf("Configuration updated: %v", *c)
	}
}
