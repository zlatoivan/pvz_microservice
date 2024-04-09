//go:build integration
// +build integration

package order

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
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

func genHTTPCreateOrderReq(t *testing.T, reqOrder delivery.RequestOrder) *http.Request {
	body, err := json.Marshal(reqOrder)
	require.NoError(t, err)
	req, err := http.NewRequest(
		http.MethodPost,
		"http://localhost:9000/api/v1/orders",
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

func TestServer_CreateOrder(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// arrange
		req := genHTTPCreateOrderReq(t, fixtures.ReqCreateOrderGood)

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
		createdOrderFromDB := dbGetByIDOrder(t, ctx, db, respID)
		dbDeleteOrder(t, ctx, db, respID)
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
		req := genHTTPCreateOrderReq(t, reqCreateOrderBadReq)
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
