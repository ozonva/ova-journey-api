package config

import "fmt"

// HealthCheckConfiguration type represents configuration for handler to HealthCheck
type HealthCheckConfiguration struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Path string `yaml:"path"`
}

// GetEndpointAddress - returns string in format "hostname:port" for endpoint
func (c *HealthCheckConfiguration) GetEndpointAddress() string {
	return fmt.Sprintf("%s:%v", c.Host, c.Port)
}
