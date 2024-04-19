package kafka

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

type HandleFunc func(message *sarama.ConsumerMessage) error

type Consumer struct {
	consumer sarama.Consumer
}

// NewConsumer возвращает новый Consumer
func NewConsumer(brokers []string) (Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	newConsumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return Consumer{}, fmt.Errorf("sarama.NewConsumer: %w", err)
	}

	consumer := Consumer{
		consumer: newConsumer,
	}

	return consumer, nil
}

func (c Consumer) Subscribe(topic string, handler HandleFunc) error {
	partitionList, err := c.consumer.Partitions(topic)
	if err != nil {
		return fmt.Errorf("c.consumer.Partitions: %w", err)
	}

	initialOffset := sarama.OffsetNewest

	for _, partition := range partitionList {
		pc, err := c.consumer.ConsumePartition(topic, partition, initialOffset)
		if err != nil {
			return fmt.Errorf("c.consumer.ConsumePartition: %w", err)
		}

		errCh := make(chan string)
		go func(pc sarama.PartitionConsumer, partition int32, errCh chan string) {
			for {
				select {
				case message := <-pc.Messages():
					err = handler(message)
					if err != nil {
						errCh <- fmt.Sprintf("handler error: %v", err)
					}
				case errH := <-errCh:
					log.Printf("[Subscribe]: %v", errH)
					return
				default:
				}
			}
		}(pc, partition, errCh)
	}
	return nil
}

func GetLogHandler() HandleFunc {
	handler := func(message *sarama.ConsumerMessage) error {
		pm := CrudMessage{}
		err := json.Unmarshal(message.Value, &pm)
		if err != nil {
			return fmt.Errorf("[GetLogHandler] json.Unmarshal: %w", err)
		}
		fmt.Println("\n[MW]:")
		fmt.Printf("Time: %s\n", pm.TimeCreate)
		fmt.Printf("HTTP method: %s\n", pm.Type)
		fmt.Printf("Data body: %s\n", pm.Data)
		return nil
	}
	return handler
}
