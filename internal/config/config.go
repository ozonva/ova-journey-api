package config

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	Host string `json:"host"`
	Port uint   `json:"port"`
}

func (c *Configuration) LoadConfigurationFromFile(path string) (conf Configuration, err error) {
	// вызов через функтор по условию задачи
	updateConfig := func(path string) (conf Configuration, err error) {
		file, err := os.Open(path)
		if err != nil {
			return
		}

		defer func() {
			if defErr := file.Close(); defErr != nil {
				err = defErr
			}
		}()

		decoder := json.NewDecoder(file)
		if err = decoder.Decode(&conf); err != nil {
			conf = Configuration{}
		}
		return
	}
	return updateConfig(path)
}
