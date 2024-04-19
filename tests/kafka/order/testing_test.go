package order

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IBM/sarama"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/config"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/kafka"
)

var (
	client http.Client
	url    = "http://localhost:9000"
)

func addAuthHeaders(req *http.Request) {
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

	handler := func(message *sarama.ConsumerMessage) error {
		pm := kafka.CrudMessage{}
		err = json.Unmarshal(message.Value, &pm)
		if err != nil {
			return fmt.Errorf("[consumerInit] json.Unmarshal: %w", err)
		}
		channelKafka <- pm
		return nil
	}

	err = consumer.Subscribe(cfg.Topic, handler)
	if err != nil {
		return fmt.Errorf("consumer.Subscribe: %w", err)
	}

	return nil
}
