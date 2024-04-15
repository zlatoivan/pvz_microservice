package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
)

type HandleFunc func(message *sarama.ConsumerMessage)

type Receiver struct {
	consumer *Consumer
	handlers map[string]HandleFunc
}

func NewReceiver(consumer *Consumer, handlers map[string]HandleFunc) *Receiver {
	return &Receiver{
		consumer: consumer,
		handlers: handlers,
	}
}

func (r *Receiver) Subscribe(topic string) error {
	handler, ok := r.handlers[topic]

	if !ok {
		return fmt.Errorf("can not find handler")
	}

	partitionList, err := r.consumer.SingleConsumer.Partitions(topic)
	if err != nil {
		return fmt.Errorf("r.consumer.SingleConsumer.Partitions: %w", err)
	}

	initialOffset := sarama.OffsetNewest

	for _, partition := range partitionList {
		pc, err := r.consumer.SingleConsumer.ConsumePartition(topic, partition, initialOffset)

		if err != nil {
			return fmt.Errorf("r.consumer.SingleConsumer.ConsumePartition: %w", err)
		}

		go func(pc sarama.PartitionConsumer, partition int32) {
			for message := range pc.Messages() {
				handler(message)
			}
		}(pc, partition)
	}

	return nil
}

//func (r *Receiver) StartConsume(topic string) error {
//	err := r.Subscribe(topic)
//
//	if err != nil {
//		return fmt.Errorf("r.Subscribe: %w", err)
//	}
//
//	return nil
//}
