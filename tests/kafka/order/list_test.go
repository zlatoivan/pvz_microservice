//go:build integration
// +build integration

package order

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/kafka"
)

func TestKafka_ListOrders(t *testing.T) {
	//t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// arrange
		req, err := http.NewRequest(http.MethodGet, url+"/api/v1/orders", nil)
		require.NoError(t, err)
		addAuthHeaders(t, req)

		channelKafka := make(chan kafka.CrudMessage)
		err = consumerInit(channelKafka)
		require.NoError(t, err)

		// act
		res, err := client.Do(req)
		require.NoError(t, err)
		crudMsg := <-channelKafka

		// assert
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "", crudMsg.Data)
	})
}
