package kafka

import (
	"fmt"
	"time"

	"github.com/IBM/sarama"
)

type Consumer struct {
	brokers        []string
	SingleConsumer sarama.Consumer
}

// NewConsumer возвращает новый Consumer
func NewConsumer(brokers []string) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = false
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 5 * time.Second
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumer(brokers, config)

	if err != nil {
		return nil, fmt.Errorf("error NewConsumer: %w", err)
	}

	return &Consumer{
		brokers:        brokers,
		SingleConsumer: consumer,
	}, err
}
