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
		id := postgres.CreateOrder(t, ctx, db, fixtures.ReqCreateOrderGood)
		reqIDGood := delivery.RequestID{ID: id}
		req := genHTTPReq(t, method, endpoint, reqIDGood)
		orderFromDB := postgres.GetByIDOrder(t, ctx, db, id)

		// act
		res, err := client.Do(req)
		require.NoError(t, err)
		respOrder := getOrderFromResp(t, res)
		postgres.DeleteOrder(t, ctx, db, id)
		postgres.TearDown(ctx, db)

		// assert
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, orderFromDB, respOrder)
	})

	t.Run("bad request", func(t *testing.T) {
		t.Parallel()

		// arrange
		reqIDBadReq := ""
		req := genHTTPReq(t, method, endpoint, reqIDBadReq)
		wantErr := delivery.MakeRespErrInvalidData(errors.New("json.Unmarshal: json: cannot unmarshal string into Go value of type delivery.RequestID"))

		// act
		res, err := client.Do(req)
		require.NoError(t, err)
		respErr := getErrorFromResp(t, res)

		// assert
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Equal(t, wantErr, respErr)
	})
}
