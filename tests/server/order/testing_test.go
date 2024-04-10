//go:build integration
// +build integration

package order

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handlers/delivery"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

var (
	client http.Client
	url    = "http://localhost:9000"
)

func addAuthHeaders(t *testing.T, req *http.Request) {
	username := "ivan"
	password := "order_best_pass"
	auth := username + ":" + password
	base64Auth := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", "Basic "+base64Auth)
}

func genHTTPReq(t *testing.T, method string, endpoint string, reqData interface{}) *http.Request {
	body, err := json.Marshal(reqData)
	require.NoError(t, err)
	req, err := http.NewRequest(method, url+endpoint, bytes.NewReader(body))
	require.NoError(t, err)
	addAuthHeaders(t, req)
	return req
}

func getOrderFromResp(t *testing.T, res *http.Response) model.Order {
	defer func() {
		err := res.Body.Close()
		require.NoError(t, err)
	}()
	var respOrder delivery.ResponseOrder
	err := json.NewDecoder(res.Body).Decode(&respOrder)
	require.NoError(t, err)
	order := delivery.GetOrderFromRespOrder(respOrder)
	return order
}

func getOrderIDFromRespOrder(t *testing.T, res *http.Response) uuid.UUID {
	defer func() {
		err := res.Body.Close()
		require.NoError(t, err)
	}()
	var respID delivery.ResponseID
	err := json.NewDecoder(res.Body).Decode(&respID)
	require.NoError(t, err)
	return respID.ID
}

func getOrdersFromRespListOrders(t *testing.T, res *http.Response) []model.Order {
	defer func() {
		err := res.Body.Close()
		require.NoError(t, err)
	}()
	var respOrders []delivery.ResponseOrder
	err := json.NewDecoder(res.Body).Decode(&respOrders)
	require.NoError(t, err)

	orders := make([]model.Order, 0)
	for _, r := range respOrders {
		orders = append(orders, delivery.GetOrderFromRespOrder(r))
	}

	return orders
}

func getCommentFromResp(t *testing.T, res *http.Response) delivery.ResponseComment {
	defer func() {
		err := res.Body.Close()
		require.NoError(t, err)
	}()
	var respComment delivery.ResponseComment
	err := json.NewDecoder(res.Body).Decode(&respComment)
	require.NoError(t, err)
	return respComment
}

func getErrorFromResp(t *testing.T, res *http.Response) delivery.ResponseError {
	defer func() {
		err := res.Body.Close()
		require.NoError(t, err)
	}()
	var respErr delivery.ResponseError
	err := json.NewDecoder(res.Body).Decode(&respErr)
	require.NoError(t, err)
	return respErr
}

func checkIn(t *testing.T, order model.Order, orders []model.Order) bool {
	for _, r := range orders {
		if order == r {
			return true
		}
	}
	return false
}
