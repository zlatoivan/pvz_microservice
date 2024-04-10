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
	"gitlab.ozon.dev/zlatoivan4/homework/tests/fixtures"
	"gitlab.ozon.dev/zlatoivan4/homework/tests/postgres"
)

func TestHandler_DeleteOrder(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	method := http.MethodDelete
	endpoint := "/api/v1/orders/id"

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// arrange
		db, err := postgres.SetUp(ctx)
		require.NoError(t, err)
		id := dbCreateOrder(t, ctx, db, fixtures.ReqCreateOrderGood)
		reqIDGood := delivery.RequestID{ID: id}
		req := genHTTPReq(t, method, endpoint, reqIDGood)
		wantJSON := delivery.MakeRespComment("Order deleted")

		// act
		res, err := client.Do(req)
		require.NoError(t, err)
		respStatus, respJSON := getResp(t, res, "Comment")
		deletedOrderFromDB := dbGetByIDOrder(t, ctx, db, id)
		dbDeleteOrder(t, ctx, db, id)
		postgres.TearDown(ctx, db)

		// assert
		assert.Equal(t, respStatus, http.StatusOK)
		assert.Equal(t, wantJSON, respJSON)
		assert.True(t, deletedOrderFromDB.IsDeleted)
	})

	t.Run("bad request", func(t *testing.T) {
		t.Parallel()

		// arrange
		reqIDBadReq := delivery.RequestID{ID: uuid.Nil}
		req := genHTTPReq(t, method, endpoint, reqIDBadReq)
		wantJSON := delivery.MakeRespErrInvalidData(errors.New("id is nil"))

		// act
		res, err := client.Do(req)
		require.NoError(t, err)
		respStatus, respJSON := getResp(t, res, "Comment")

		// assert
		assert.Equal(t, respStatus, http.StatusBadRequest)
		assert.Equal(t, wantJSON, respJSON)
	})
}
