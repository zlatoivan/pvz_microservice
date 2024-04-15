package order

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/IBM/sarama"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/kafka"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/config"
)

var (
	client http.Client
	url    = "http://localhost:9000"
)

func addAuthHeaders(t *testing.T, req *http.Request) {
	username := "ivan"
	password := "order_best_pass"
	auth := username + ":" + password
	base64Auth := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", "Basic "+base64Auth)
}

func consumerInit(channelKafka chan<- kafka.CrudMessage) error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("config.New: %w", err)
	}

	consumer, err := kafka.NewConsumer(cfg.Brokers)
	if err != nil {
		return fmt.Errorf("kafka.NewConsumer: %w", err)
	}

	handlers := map[string]kafka.HandleFunc{
		cfg.Topic: func(message *sarama.ConsumerMessage) {
			pm := kafka.CrudMessage{}
			err = json.Unmarshal(message.Value, &pm)
			if err != nil {
				log.Println("Consumer error", err)
			}
			channelKafka <- pm
		},
	}

	receiver := kafka.NewReceiver(consumer, handlers)
	err = receiver.Subscribe(cfg.Topic)
	if err != nil {
		return fmt.Errorf("receiver.Subscribe: %w", err)
	}

	return nil
}
