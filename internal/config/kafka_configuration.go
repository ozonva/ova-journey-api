package config

// KafkaConfiguration type represents configuration for Kafka
type KafkaConfiguration struct {
	Topic   string   `yaml:"topic"`
	Brokers []string `yaml:"brokers"`
}
