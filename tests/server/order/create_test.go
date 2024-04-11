//go:build integration
// +build integration

package order

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/delivery"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/service/order"
	"gitlab.ozon.dev/zlatoivan4/homework/tests/fixtures"
	"gitlab.ozon.dev/zlatoivan4/homework/tests/postgres"
)

func TestServer_CreateOrder(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	method := http.MethodPost
	endpoint := "/api/v1/orders"

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// arrange
		req := genHTTPReq(t, method, endpoint, fixtures.ReqCreateOrderGood)
		createdOrder, err := order.ApplyPackaging(delivery.GetOrderFromReqOrder(fixtures.ReqCreateOrderGood))
		require.NoError(t, err)

		// act
		res, err := client.Do(req)
		require.NoError(t, err)
		id := getOrderIDFromRespOrder(t, res)
		createdOrder.ID = id
		db, err := postgres.SetUp(ctx)
		require.NoError(t, err)
		createdOrderFromDB := postgres.GetByIDOrder(t, ctx, db, id)
		postgres.DeleteOrder(t, ctx, db, id)
		postgres.TearDown(ctx, db)

		// assert
		assert.Equal(t, http.StatusCreated, res.StatusCode)
		assert.Equal(t, createdOrder, createdOrderFromDB)
	})

	t.Run("bad request", func(t *testing.T) {
		t.Parallel()

		// arrange
		reqCreateOrderBadReq := delivery.RequestOrder{ClientID: uuid.Nil}
		req := genHTTPReq(t, method, endpoint, reqCreateOrderBadReq)
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
