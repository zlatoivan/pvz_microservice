//go:build integration
// +build integration

package order

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handlers/delivery"
	"gitlab.ozon.dev/zlatoivan4/homework/tests/fixtures"
	"gitlab.ozon.dev/zlatoivan4/homework/tests/postgres"
)

func TestServer_GetOrderByID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	method := http.MethodGet
	endpoint := "/api/v1/orders/id"

	t.Run("success", func(t *testing.T) {
		//t.Parallel()

		// arrange
		db, err := postgres.SetUp(ctx)
		require.NoError(t, err)
		id := dbCreateOrder(t, ctx, db, fixtures.ReqCreateOrderGood)
		reqIDGood := delivery.RequestID{ID: id}
		req := genHTTPReq(t, method, endpoint, reqIDGood)
		respGetByIDOrder := delivery.GetOrderFromReqOrder(fixtures.ReqCreateOrderGood)
		respGetByIDOrder.ID = id
		respGetByIDOrder.StoresTill = fixtures.StoresTill
		wantResp := delivery.MakeRespOrder(respGetByIDOrder)

		// act
		res, err := client.Do(req)
		require.NoError(t, err)
		respStatus, respOrder := getResp(t, res, "Order")
		dbDeleteOrder(t, ctx, db, id)
		postgres.TearDown(ctx, db)

		// assert
		assert.Equal(t, http.StatusOK, respStatus)
		assert.Equal(t, wantResp, respOrder)
	})

	t.Run("bad request", func(t *testing.T) {
		t.Parallel()

		// arrange
		db, err := postgres.SetUp(ctx)
		require.NoError(t, err)
		id := dbCreateOrder(t, ctx, db, fixtures.ReqCreateOrderGood)
		reqIDBadReq := ""
		req := genHTTPReq(t, method, endpoint, reqIDBadReq)
		wantResp := delivery.MakeRespErrInvalidData(errors.New("json.Unmarshal: json: cannot unmarshal string into Go value of type delivery.RequestID"))

		// act
		res, err := client.Do(req)
		require.NoError(t, err)
		respStatus, respOrder := getResp(t, res, "")
		dbDeleteOrder(t, ctx, db, id)
		postgres.TearDown(ctx, db)

		// assert
		assert.Equal(t, http.StatusBadRequest, respStatus)
		assert.Equal(t, wantResp, respOrder)
	})
}
