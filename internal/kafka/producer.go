package kafka

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/ozonva/ova-journey-api/internal/config"
	"github.com/rs/zerolog/log"
)

// Producer - interface for work with Kafka
type Producer interface {
	Send(message Message) error
	Close() error
}

type producer struct {
	syncProducer sarama.SyncProducer
	topic        string
}

// NewProducer - creates new Producer for work with Kafka
func NewProducer(configuration *config.KafkaConfiguration) (Producer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Partitioner = sarama.NewRandomPartitioner
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Return.Successes = true

	syncProducer, err := sarama.NewSyncProducer(configuration.Brokers, saramaConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("Kafka producer: failed to create")
		return nil, err
	}

	return &producer{
		topic:        configuration.Topic,
		syncProducer: syncProducer,
	}, nil
}

// Send - Send new message to Kafka
func (p *producer) Send(message Message) error {
	jsonMes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, _, err = p.syncProducer.SendMessage(
		&sarama.ProducerMessage{
			Topic:     p.topic,
			Partition: -1,
			Key:       sarama.StringEncoder(p.topic),
			Value:     sarama.StringEncoder(jsonMes),
		})
	return err
}

func (p *producer) Close() error {
	return p.syncProducer.Close()
}
