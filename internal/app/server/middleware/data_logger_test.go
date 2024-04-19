package middleware

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/middleware/mock"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/kafka"
)

const body = `{
    "client_id": "9967bb48-bd6f-4ad0-924d-8c9094c4d8c2",
    "stores_till": "2024-04-22T13:14:01Z",
    "weight": 29,
    "cost": 1100,
    "packaging_type": "box"
}`

func TestMW_DataLogger(t *testing.T) {
	t.Parallel()

	mc := minimock.NewController(t)
	message := kafka.CrudMessage{
		TimeCreate: time.Now(),
		Type:       http.MethodPost,
		Data:       body,
	}

	tests := []struct {
		name     string
		producer Producer
		wantErr  error
	}{
		{
			name: "success",
			producer: mock.NewProducerMock(mc).
				SendMessageMock.
				Expect(message).
				Return(nil),
			wantErr: nil,
		},
		{
			name: "json marshall error",
			producer: mock.NewProducerMock(mc).
				SendMessageMock.
				Expect(message).
				Return(fmt.Errorf("json.Marshal")),
			wantErr: fmt.Errorf("json.Marshal"),
		},
		{
			name: "send message error",
			producer: mock.NewProducerMock(mc).
				SendMessageMock.
				Expect(message).
				Return(fmt.Errorf("p.producer.SendMessage")),
			wantErr: fmt.Errorf("p.producer.SendMessage"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// arrange

			// act
			err := tt.producer.SendMessage(message)

			// assert
			require.Equal(t, tt.wantErr, err)
		})
	}
}
