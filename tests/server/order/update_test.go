//go:build integration
// +build integration

package order

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handlers/delivery"
	"gitlab.ozon.dev/zlatoivan4/homework/tests/fixtures"
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
		db, err := postgres.SetUp(ctx)
		require.NoError(t, err)
		id := postgres.CreateOrder(t, ctx, db, fixtures.ReqCreateOrderGood)
		reqUpdateOrderGood := fixtures.ReqCreateOrderGood
		reqUpdateOrderGood.ID = id
		reqUpdateOrderGood.Cost = 499
		req := genHTTPReq(t, method, endpoint, reqUpdateOrderGood)
		wantRespComment := delivery.MakeRespComment("Order updated")

		// act
		res, err := client.Do(req)
		require.NoError(t, err)
		respStatus, respJSON := getResp(t, res, "Comment")
		updatedOrder := delivery.GetOrderFromReqOrder(reqUpdateOrderGood)
		updatedOrderFromDB := postgres.GetByIDOrder(t, ctx, db, id)
		postgres.DeleteOrder(t, ctx, db, id)
		postgres.TearDown(ctx, db)

		// assert
		assert.Equal(t, http.StatusOK, respStatus)
		assert.Equal(t, wantRespComment, respJSON)
		assert.Equal(t, updatedOrder, updatedOrderFromDB)
	})

	t.Run("bad request", func(t *testing.T) {
		t.Parallel()

		// arrange
		reqUpdateOrderBadReq := fixtures.ReqCreateOrderGood
		reqUpdateOrderBadReq.ClientID = uuid.Nil
		req := genHTTPReq(t, method, endpoint, reqUpdateOrderBadReq)
		wantRespComment := delivery.MakeRespErrInvalidData(errors.New("client id is nil"))

		// act
		res, err := client.Do(req)
		require.NoError(t, err)
		respStatus, respJSON := getResp(t, res, "")

		// assert
		assert.Equal(t, http.StatusBadRequest, respStatus)
		assert.Equal(t, wantRespComment, respJSON)
	})
}
