//go:build integration
// +build integration

package order

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/delivery"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/kafka"
	"gitlab.ozon.dev/zlatoivan4/homework/tests/postgres"
)

func TestKafka_GetByIdOrder(t *testing.T) {
	//t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// arrange
		ClientID, _ := uuid.Parse("88cda6c0-36fc-4be4-b976-e11a8a7a8f7e")
		StoresTill, _ := time.Parse(time.RFC3339, "2024-04-22T12:12:00Z")
		reqOrder := delivery.RequestOrder{
			ClientID:      ClientID,
			StoresTill:    StoresTill,
			Weight:        29,
			Cost:          1100,
			PackagingType: "box",
		}
		ctx := context.Background()
		db, err := postgres.SetUp(ctx)
		require.NoError(t, err)
		id := postgres.CreateOrder(t, ctx, db, reqOrder)
		reqID := delivery.RequestID{ID: id}
		body, err := json.Marshal(reqID)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodGet, url+"/api/v1/orders/id", bytes.NewReader(body))
		require.NoError(t, err)
		addAuthHeaders(req)

		channelKafka := make(chan kafka.CrudMessage)
		err = consumerInit(channelKafka)
		require.NoError(t, err)

		// act
		res, err := client.Do(req)
		require.NoError(t, err)

		var reqIDFromKafka delivery.RequestID
		crudMsg := <-channelKafka
		err = json.Unmarshal([]byte(crudMsg.Data), &reqIDFromKafka)
		require.NoError(t, err)

		postgres.DeleteOrder(t, ctx, db, id)
		postgres.TearDown(ctx, db)

		// assert
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, reqID, reqIDFromKafka)
	})
}
