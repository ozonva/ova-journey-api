package config

import "fmt"

// EndpointConfiguration type represents configuration for network endpoint (host and port)
type EndpointConfiguration struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// GetEndpointAddress - returns string in format "hostname:port" for endpoint
func (c *EndpointConfiguration) GetEndpointAddress() string {
	return fmt.Sprintf("%s:%v", c.Host, c.Port)
}
