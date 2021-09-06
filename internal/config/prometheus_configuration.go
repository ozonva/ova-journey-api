package config

import "fmt"

// PrometheusConfiguration type represents configuration for handler to prometheus
type PrometheusConfiguration struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Path string `yaml:"path"`
}

// GetEndpointAddress - returns string in format "hostname:port" for endpoint
func (c *PrometheusConfiguration) GetEndpointAddress() string {
	return fmt.Sprintf("%s:%v", c.Host, c.Port)
}
