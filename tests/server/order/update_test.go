//go:build integration
// +build integration

package order

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
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

func genHTTPUpdateOrderReq(t *testing.T, reqOrder delivery.RequestOrder) *http.Request {
	body, err := json.Marshal(reqOrder)
	require.NoError(t, err)
	req, err := http.NewRequest(
		http.MethodPut,
		"http://localhost:9000/api/v1/orders/id",
		bytes.NewReader(body),
	)
	require.NoError(t, err)
	username := "ivan"
	password := "order_best_pass"
	auth := username + ":" + password
	base64Auth := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", "Basic "+base64Auth)
	return req
}

func TestServer_UpdateOrder(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// arrange
		db, err := postgres.SetUp(ctx)
		require.NoError(t, err)
		id := dbCreateOrder(t, ctx, db, fixtures.ReqCreateOrderGood)
		reqUpdateOrderGood := fixtures.ReqCreateOrderGood
		reqUpdateOrderGood.ID = id
		reqUpdateOrderGood.Cost = 499
		req := genHTTPUpdateOrderReq(t, reqUpdateOrderGood)
		wantRespComment := delivery.MakeRespComment("Order updated")

		// act
		res, err := client.Do(req)
		require.NoError(t, err)
		respStatus, respJSON := getResp(t, res, "Comment")
		updatedOrder := delivery.GetOrderFromReqOrder(reqUpdateOrderGood)
		updatedOrderFromDB := dbGetByIDOrder(t, ctx, db, id)
		dbDeleteOrder(t, ctx, db, id)
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
		req := genHTTPUpdateOrderReq(t, reqUpdateOrderBadReq)
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
