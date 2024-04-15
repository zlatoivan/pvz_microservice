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

type Sender struct {
	producer *Producer
	topic    string
}

func NewSender(producer *Producer, topic string) *Sender {
	return &Sender{
		producer,
		topic,
	}
}

func (s *Sender) SendMessage(message CrudMessage) error {
	kafkaMsg, err := s.buildMessage(message)
	if err != nil {
		return fmt.Errorf("s.buildMessage: %w", err)
	}

	err = s.producer.SendSyncMessage(kafkaMsg)
	if err != nil {
		return fmt.Errorf("s.producer.SendSyncMessage: %w", err)
	}

	return nil
}

func (s *Sender) sendMessages(messages []CrudMessage) error {
	var kafkaMsg []*sarama.ProducerMessage
	for _, m := range messages {
		message, err := s.buildMessage(m)
		if err != nil {
			return fmt.Errorf("s.buildMessage: %w", err)
		}
		kafkaMsg = append(kafkaMsg, message)
	}

	err := s.producer.SendSyncMessages(kafkaMsg)
	if err != nil {
		return fmt.Errorf("s.producer.SendSyncMessages: %w", err)
	}

	return nil
}

func (s *Sender) buildMessage(message CrudMessage) (*sarama.ProducerMessage, error) {
	msg, err := json.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}

	return &sarama.ProducerMessage{
		Topic:     s.topic,
		Value:     sarama.ByteEncoder(msg),
		Partition: -1,
	}, nil
}
