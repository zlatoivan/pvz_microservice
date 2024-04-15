package kafka

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/config"
)

func New(cfg config.Config) (*Sender, error) {
	producer, err := NewProducer(cfg.Brokers)
	if err != nil {
		return nil, fmt.Errorf("kafka.NewProducer: %w", err)
	}
	consumer, err := NewConsumer(cfg.Brokers)
	if err != nil {
		return nil, fmt.Errorf("kafka.NewProducer: %w", err)
	}

	sender := NewSender(producer, cfg.Topic)

	handlers := map[string]HandleFunc{
		cfg.Topic: func(message *sarama.ConsumerMessage) {
			pm := CrudMessage{}
			err = json.Unmarshal(message.Value, &pm)
			if err != nil {
				log.Println("[error] json.Unmarshal", err)
			}

			fmt.Println("\n[MW]:")
			//fmt.Printf("Offset: %d\n", message.Offset)
			fmt.Printf("Time: %s\n", pm.TimeCreate)
			fmt.Printf("HTTP method: %s\n", pm.Type)
			fmt.Printf("Data body: %s\n", pm.Data)
		},
	}

	receiver := NewReceiver(consumer, handlers)
	err = receiver.Subscribe(cfg.Topic)
	if err != nil {
		return nil, fmt.Errorf("receiver.Subscribe: %w", err)
	}

	return sender, nil
}
