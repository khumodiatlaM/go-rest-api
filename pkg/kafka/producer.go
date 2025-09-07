package kafka

import (
	"go-rest-api/pkg/logger"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Producer struct {
	producer *kafka.Producer
	log      logger.CustomLogger
}

func NewProducer(brokers string, logger logger.CustomLogger) (Producer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
	})
	if err != nil {
		logger.Error("failed to create kafka producer", err)
		return Producer{}, err
	}
	return Producer{
		producer: producer,
		log:      logger,
	}, nil
}

func (p *Producer) Produce(topic string, key string, value []byte) error {
	deliveryChan := make(chan kafka.Event)
	err := p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          value,
	}, deliveryChan)
	if err != nil {
		return err
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		return m.TopicPartition.Error
	}
	close(deliveryChan)
	return nil
}

func (p *Producer) Close() {
	p.producer.Close()
}
