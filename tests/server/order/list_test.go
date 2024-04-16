package order

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/zlatoivan4/homework/tests/fixtures"
	"gitlab.ozon.dev/zlatoivan4/homework/tests/postgres"
)

func TestServer_ListOrders(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	method := http.MethodGet
	endpoint := "/api/v1/orders"

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// arrange
		db, err := postgres.SetUp(ctx)
		require.NoError(t, err)
		id1 := postgres.CreateOrder(t, ctx, db, fixtures.ReqCreateOrderGood)
		id2 := postgres.CreateOrder(t, ctx, db, fixtures.ReqCreateOrderGood)
		orderFromDB1 := postgres.GetByIDOrder(t, ctx, db, id1)
		orderFromDB2 := postgres.GetByIDOrder(t, ctx, db, id2)
		req := genHTTPReq(t, method, endpoint, "")

		// act
		res, err := client.Do(req)
		require.NoError(t, err)
		orders := getOrdersFromRespListOrders(t, res)
		wantTrue := checkIn(orderFromDB1, orders) && checkIn(orderFromDB2, orders)
		postgres.DeleteOrder(t, ctx, db, id1)
		postgres.DeleteOrder(t, ctx, db, id2)
		postgres.TearDown(ctx, db)

		// assert
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.True(t, wantTrue)
	})
}
