package middleware

import (
	"fmt"
	"net/http"
	"time"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/delivery"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/kafka"
)

type MW struct {
	Sender *kafka.Sender
}

func New(sender *kafka.Sender) MW {
	return MW{Sender: sender}
}

func (m *MW) DataLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		data, err := delivery.GetRawDataFromReq(req)
		if err != nil {
			err = m.Sender.SendMessage(kafka.CrudMessage{
				TimeCreate: time.Now(),
				Type:       req.Method,
				Data:       fmt.Sprintf("[MW]: delivery.GetRawDataFromReq: %v", err),
			})
			return
		}
		err = m.Sender.SendMessage(kafka.CrudMessage{
			TimeCreate: time.Now(),
			Type:       req.Method,
			Data:       string(data),
		})

		next.ServeHTTP(w, req)
	})
}
