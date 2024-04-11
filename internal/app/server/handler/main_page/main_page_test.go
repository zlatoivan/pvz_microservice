package main_page

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/delivery"
)

func getResp(t *testing.T, w *httptest.ResponseRecorder) (int, interface{}) {
	res := w.Result()
	defer func() {
		err := res.Body.Close()
		require.NoError(t, err)
	}()

	var respJSON interface{}
	switch res.StatusCode {
	case http.StatusOK:
		respID := delivery.ResponseComment{}
		err := json.NewDecoder(res.Body).Decode(&respID)
		require.NoError(t, err)
		respJSON = respID
	default:
		var respErr delivery.ResponseError
		err := json.NewDecoder(res.Body).Decode(&respErr)
		require.NoError(t, err)
		respJSON = respErr
	}

	return res.StatusCode, respJSON
}

func Test_MainPage(t *testing.T) {
	tests := []struct {
		name       string
		req        *http.Request
		wantStatus int
		wantJSON   interface{}
	}{
		{
			name: "success",
			req: httptest.NewRequest(
				http.MethodGet,
				"http://localhost:9000",
				nil,
			),
			wantStatus: http.StatusOK,
			wantJSON:   delivery.MakeRespComment("This is the main page"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			w := httptest.NewRecorder()

			// act
			MainPage(w, tt.req)

			// act
			respStatus, respJSON := getResp(t, w)

			// assert
			assert.Equal(t, tt.wantStatus, respStatus)
			assert.Equal(t, tt.wantJSON, respJSON)
		})
	}
}
