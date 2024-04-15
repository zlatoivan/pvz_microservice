package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
)

// Producer ...
type Producer struct {
	brokers      []string
	syncProducer sarama.SyncProducer
}

func newSyncProducer(brokers []string) (sarama.SyncProducer, error) {
	syncProducerConfig := sarama.NewConfig()
	syncProducerConfig.Producer.Partitioner = sarama.NewRandomPartitioner
	syncProducerConfig.Producer.RequiredAcks = sarama.WaitForAll
	syncProducerConfig.Producer.Idempotent = true
	syncProducerConfig.Net.MaxOpenRequests = 1
	syncProducerConfig.Producer.CompressionLevel = sarama.CompressionLevelDefault
	syncProducerConfig.Producer.Return.Successes = true
	syncProducerConfig.Producer.Return.Errors = true
	syncProducerConfig.Producer.Compression = sarama.CompressionGZIP

	syncProducer, err := sarama.NewSyncProducer(brokers, syncProducerConfig)
	if err != nil {
		return nil, fmt.Errorf("sarama.NewSyncProducer: %w", err)
	}

	return syncProducer, nil
}

// NewProducer возвращает новую Producer
func NewProducer(brokers []string) (*Producer, error) {
	syncProducer, err := newSyncProducer(brokers)
	if err != nil {
		return nil, fmt.Errorf("newSyncProducer: %w", err)
	}

	producer := &Producer{
		brokers:      brokers,
		syncProducer: syncProducer,
	}

	return producer, nil
}

// SendSyncMessage отправляет одно сообщение
func (p *Producer) SendSyncMessage(message *sarama.ProducerMessage) error {
	_, _, err := p.syncProducer.SendMessage(message)
	if err != nil {
		return fmt.Errorf("p.syncProducer.SendMessage: %w", err)
	}
	return nil
}

// SendSyncMessages отправляет список сообщений
func (p *Producer) SendSyncMessages(messages []*sarama.ProducerMessage) error {
	err := p.syncProducer.SendMessages(messages)
	if err != nil {
		return fmt.Errorf("p.syncProducer.SendMessages: %w", err)
	}

	return nil
}

// Close ...
func (p *Producer) Close() error {
	err := p.syncProducer.Close()
	if err != nil {
		return fmt.Errorf("p.syncProducer.Close %w", err)
	}

	return nil
}
