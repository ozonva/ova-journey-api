package config

import "fmt"

// DatabaseConfiguration type represents configuration with information connection parameters to database
type DatabaseConfiguration struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	SslMode  string `yaml:"sslMode"`
	Driver   string `yaml:"driver"`
}

// GetDataSourceName - return dataSourceName in format "host=%v port=%v dbname=%v user=%v password=%v sslmode=%v"
func (configuration DatabaseConfiguration) GetDataSourceName() string {
	return fmt.Sprintf("host=%v port=%v dbname=%v user=%v password=%v sslmode=%v",
		configuration.Host,
		configuration.Port,
		configuration.Name,
		configuration.User,
		configuration.Password,
		configuration.SslMode,
	)
}
