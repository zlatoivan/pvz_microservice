//go:build integration
// +build integration

package order

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/delivery"
	"gitlab.ozon.dev/zlatoivan4/homework/tests/postgres"
)

func TestServer_UpdateOrder(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	method := http.MethodPut
	endpoint := "/api/v1/orders/id"

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
		id := postgres.CreateOrder(t, ctx, db, reqOrder)
		reqUpdateOrderGood := postgres.GetByIDOrder(t, ctx, db, id)
		reqUpdateOrderGood.ID = id
		reqUpdateOrderGood.Cost = 499
		req := genHTTPReq(t, method, endpoint, reqUpdateOrderGood)
		wantComment := delivery.MakeRespComment("Order updated")

		// act
		res, err := client.Do(req)
		require.NoError(t, err)
		respComment := getCommentFromResp(t, res)
		updatedOrderFromDB := postgres.GetByIDOrder(t, ctx, db, id)
		postgres.DeleteOrder(t, ctx, db, id)
		postgres.TearDown(ctx, db)

		// assert
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, wantComment, respComment)
		assert.Equal(t, reqUpdateOrderGood, updatedOrderFromDB)
	})

	t.Run("bad request", func(t *testing.T) {
		t.Parallel()

		// arrange
		reqUpdateOrderBadReq := delivery.RequestOrder{ClientID: uuid.Nil}
		req := genHTTPReq(t, method, endpoint, reqUpdateOrderBadReq)
		wantErr := delivery.MakeRespErrInvalidData(errors.New("client id is nil"))

		// act
		res, err := client.Do(req)
		require.NoError(t, err)
		respErr := getErrorFromResp(t, res)

		// assert
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Equal(t, wantErr, respErr)
	})
}
