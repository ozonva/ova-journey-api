package config

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
