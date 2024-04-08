package order

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handlers/delivery"
)

func getResp(t *testing.T, w *httptest.ResponseRecorder, respType string) (int, interface{}) {
	res := w.Result()
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
