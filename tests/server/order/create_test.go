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

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handlers/delivery"
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

		// act
		res, err := client.Do(req)
		require.NoError(t, err)
		respID := getOrderIDFromRespOrder(t, res)
		createdOrderRaw := delivery.GetOrderFromReqOrder(fixtures.ReqCreateOrderGood)
		createdOrderRaw.ID = respID
		createdOrder, err := order.ApplyPackaging(createdOrderRaw)
		require.NoError(t, err)
		db, err := postgres.SetUp(ctx)
		require.NoError(t, err)
		createdOrderFromDB := postgres.GetByIDOrder(t, ctx, db, respID)
		postgres.DeleteOrder(t, ctx, db, respID)
		postgres.TearDown(ctx, db)

		// assert
		assert.Equal(t, http.StatusCreated, res.StatusCode)
		assert.Equal(t, createdOrder, createdOrderFromDB)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// arrange
		reqCreateOrderBadReq := delivery.RequestOrder{
			ClientID: uuid.Nil,
		}
		req := genHTTPReq(t, method, endpoint, reqCreateOrderBadReq)
		wantJSON := delivery.MakeRespErrInvalidData(errors.New("client id is nil"))

		// act
		res, err := client.Do(req)
		require.NoError(t, err)
		_, respErr := getResp(t, res, "")

		// assert
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Equal(t, wantJSON, respErr)
	})
}
