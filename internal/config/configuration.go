package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

// Configuration type represents application configuration
type Configuration struct {
	Project   *ProjectConfiguration  `yaml:"project"`
	GRPC      *EndpointConfiguration `yaml:"grpc"`
	Gateway   *EndpointConfiguration `yaml:"gateway"`
	Database  *DatabaseConfiguration `yaml:"database"`
	ChunkSize int                    `yaml:"chunkSize"`
	Jaeger    *EndpointConfiguration `yaml:"jaeger"`
}

// LoadConfigurationFromFile - method for load Configuration from JSON file.
// The method return empty Configuration if any error occur in loading process.
func (c *Configuration) LoadConfigurationFromFile(path string) (conf Configuration, err error) {
	updateConfig := func(path string) (conf Configuration, err error) {
		var file *os.File
		file, err = os.Open(path)
		if err != nil {
			return
		}

		defer func() {
			if defErr := file.Close(); defErr != nil {
				err = defErr
			}
		}()

		decoder := yaml.NewDecoder(file)
		if err = decoder.Decode(&conf); err != nil {
			conf = Configuration{}
		}
		return
	}
	return updateConfig(path)
}
