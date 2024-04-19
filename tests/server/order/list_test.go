package order

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/delivery"
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
		ClientID, _ := uuid.Parse("88cda6c0-36fc-4be4-b976-e11a8a7a8f7e")
		StoresTill, _ := time.Parse(time.RFC3339, "2024-04-22T12:12:00Z")
		reqOrder := delivery.RequestOrder{
			ClientID:      ClientID,
			StoresTill:    StoresTill,
			Weight:        29,
			Cost:          1100,
			PackagingType: "box",
		}
		db, err := postgres.SetUp(ctx)
		require.NoError(t, err)
		id1 := postgres.CreateOrder(t, ctx, db, reqOrder)
		id2 := postgres.CreateOrder(t, ctx, db, reqOrder)
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
