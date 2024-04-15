//go:build integration
// +build integration

package order

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/delivery"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/kafka"
	"gitlab.ozon.dev/zlatoivan4/homework/tests/fixtures"
)

func TestKafka_CreateOrder(t *testing.T) {
	//t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// arrange
		reqOrder := fixtures.ReqCreateOrderGood
		body, err := json.Marshal(reqOrder)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, url+"/api/v1/orders", bytes.NewReader(body))
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

		// assert
		assert.Equal(t, http.StatusCreated, res.StatusCode)
		assert.Equal(t, reqOrder, reqOrderFromKafka)
	})
}
