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

func getResp(t *testing.T, res *http.Response, respType string) (int, interface{}) {
	defer func() {
		err := res.Body.Close()
		require.NoError(t, err)
	}()

	var respJSON interface{}

	switch res.StatusCode {
	case http.StatusCreated:
		respID := delivery.ResponseID{}
		err := json.NewDecoder(res.Body).Decode(&respID)
		require.NoError(t, err)
		respJSON = respID

	case http.StatusOK:
		switch respType {
		case "ID":
			respID := delivery.ResponseID{}
			err := json.NewDecoder(res.Body).Decode(&respID)
			require.NoError(t, err)
			respJSON = respID
		case "Comment":
			var respComment delivery.ResponseComment
			err := json.NewDecoder(res.Body).Decode(&respComment)
			require.NoError(t, err)
			respJSON = respComment
		case "Order":
			var respOrder delivery.ResponseOrder
			err := json.NewDecoder(res.Body).Decode(&respOrder)
			require.NoError(t, err)
			respJSON = respOrder
		case "ListOrders":
			var respOrders []delivery.ResponseOrder
			err := json.NewDecoder(res.Body).Decode(&respOrders)
			require.NoError(t, err)
			respJSON = respOrders
		}

	default:
		var respErr delivery.ResponseError
		err := json.NewDecoder(res.Body).Decode(&respErr)
		require.NoError(t, err)
		respJSON = respErr
	}

	return res.StatusCode, respJSON
}
