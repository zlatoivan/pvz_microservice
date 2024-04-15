//go:build integration
// +build integration

package order

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/delivery"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/kafka"
	"gitlab.ozon.dev/zlatoivan4/homework/tests/fixtures"
	"gitlab.ozon.dev/zlatoivan4/homework/tests/postgres"
)

func TestKafka_UpdateOrder(t *testing.T) {
	//t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// arrange
		ctx := context.Background()
		db, err := postgres.SetUp(ctx)
		require.NoError(t, err)
		id := postgres.CreateOrder(t, ctx, db, fixtures.ReqCreateOrderGood)
		reqOrder := fixtures.ReqCreateOrderGood
		reqOrder.ID = id
		body, err := json.Marshal(reqOrder)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPut, url+"/api/v1/orders/id", bytes.NewReader(body))
		require.NoError(t, err)
		addAuthHeaders(t, req)

		channelKafka := make(chan kafka.CrudMessage)
		err = consumerInit(channelKafka)
		require.NoError(t, err)

		// act
		res, err := client.Do(req)
		require.NoError(t, err)

		var reqOrderFromKafka delivery.RequestOrder
		crudMsg := <-channelKafka
		err = json.Unmarshal([]byte(crudMsg.Data), &reqOrderFromKafka)
		require.NoError(t, err)

		postgres.DeleteOrder(t, ctx, db, id)
		postgres.TearDown(ctx, db)

		// assert
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, reqOrder, reqOrderFromKafka)
	})
}
