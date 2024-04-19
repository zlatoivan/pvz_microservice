//go:generate minimock -i Producer -o mock/data_logger_mock.go -p mock -g

package middleware

import (
	"fmt"
	"net/http"
	"time"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/delivery"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/kafka"
)

type Producer interface {
	SendMessage(message kafka.CrudMessage) error
}

type MW struct {
	producer Producer
}

func New(producer Producer) MW {
	return MW{producer: producer}
}

func (m *MW) DataLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		data, err := delivery.GetRawDataFromReq(req)
		if err != nil {
			data = fmt.Sprintf("[MW]: delivery.GetRawDataFromReq: %v", err)
		}
		err = m.producer.SendMessage(kafka.CrudMessage{
			TimeCreate: time.Now(),
			Type:       req.Method,
			Data:       data,
		})

		next.ServeHTTP(w, req)
	})
}
