package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

// Configuration type represents application configuration
type Configuration struct {
	Project  *ProjectConfiguration  `yaml:"project"`
	GRPC     *EndpointConfiguration `yaml:"grpc"`
	Gateway  *EndpointConfiguration `yaml:"gateway"`
	Database *DatabaseConfiguration `yaml:"database"`
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

// CompareConfigurations - compare two configuration structs.
// Returns true if all configuration fields in struct 'a' are equal configuration fields in struct 'b'.
// Return false if any configuration field is not equal or any of structs is empty.
func CompareConfigurations(a, b *Configuration) bool {
	if &a == &b {
		return true
	}
	if (Configuration{}) == *a || (Configuration{}) == *b {
		return false
	}
	if a.Project.Name != b.Project.Name || a.Project.Version != b.Project.Version {
		return false
	}
	if a.GRPC.Host != b.GRPC.Host || a.GRPC.Port != b.GRPC.Port {
		return false
	}
	if a.Gateway.Host != b.Gateway.Host || a.Gateway.Port != b.Gateway.Port {
		return false
	}
	return true
}
