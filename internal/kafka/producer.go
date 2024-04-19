package kafka

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
)

type CrudMessage struct {
	TimeCreate time.Time
	Type       string
	Data       string
}

// Producer ...
type Producer struct {
	producer sarama.SyncProducer
	topic    string
}

// NewProducer возвращает новую Producer
func NewProducer(brokers []string, topic string) (*Producer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Idempotent = true
	cfg.Producer.Return.Successes = true
	cfg.Net.MaxOpenRequests = 1

	syncProducer, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		return nil, fmt.Errorf("sarama.NewSyncProducer: %w", err)
	}

	producer := &Producer{
		producer: syncProducer,
		topic:    topic,
	}

	return producer, nil
}

func (p Producer) SendMessage(message CrudMessage) error {
	msg, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	kafkaMsg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(msg),
	}

	_, _, err = p.producer.SendMessage(kafkaMsg)
	if err != nil {
		return fmt.Errorf("p.producer.SendMessage: %w", err)
	}

	return nil
}

// Close ...
func (p Producer) Close() error {
	err := p.producer.Close()
	if err != nil {
		return fmt.Errorf("p.producer.Close %w", err)
	}

	return nil
}
