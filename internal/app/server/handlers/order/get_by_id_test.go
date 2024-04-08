package order

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handlers/delivery"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handlers/order/mock"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
	"gitlab.ozon.dev/zlatoivan4/homework/tests/fixtures"
)

func genHTTPGetOrderByIDReq(t *testing.T, reqID interface{}) *http.Request {
	body, err := json.Marshal(reqID)
	require.NoError(t, err)
	req := httptest.NewRequest(
		http.MethodGet,
		"http://localhost:9000/api/v1/orders/id",
		bytes.NewReader(body),
	)
	return req
}

func TestHandler_GetOrderByID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mc := minimock.NewController(t)

	reqIDGood := delivery.RequestID{ID: fixtures.ID}
	reqIDBadReq := ""
	validOrder := fixtures.Order().Valid().V()

	tests := []struct {
		name       string
		service    Service
		req        *http.Request
		wantStatus int
		wantJSON   interface{}
	}{
		{
			name: "success",
			service: mock.NewServiceMock(mc).
				GetOrderByIDMock.
				Expect(ctx, fixtures.ID).
				Return(validOrder, nil),
			req:        genHTTPGetOrderByIDReq(t, reqIDGood),
			wantStatus: http.StatusOK,
			wantJSON:   delivery.MakeRespOrder(validOrder),
		},
		{
			name:       "bad request",
			service:    mock.NewServiceMock(mc),
			req:        genHTTPGetOrderByIDReq(t, reqIDBadReq),
			wantStatus: http.StatusBadRequest,
			wantJSON:   delivery.MakeRespErrInvalidData(errors.New("json.Unmarshal: json: cannot unmarshal string into Go value of type delivery.RequestID")),
		},
		{
			name: "not found",
			service: mock.NewServiceMock(mc).
				GetOrderByIDMock.
				Return(model.Order{}, ErrNotFound),
			req:        genHTTPGetOrderByIDReq(t, reqIDGood),
			wantStatus: http.StatusNotFound,
			wantJSON:   delivery.MakeRespErrNotFoundByID(errors.New("not found")),
		},
		{
			name: "server error",
			service: mock.NewServiceMock(mc).
				GetOrderByIDMock.
				Return(model.Order{}, errors.New("")),
			req:        genHTTPGetOrderByIDReq(t, reqIDGood),
			wantStatus: http.StatusInternalServerError,
			wantJSON:   delivery.MakeRespErrInternalServer(errors.New("")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// arrange
			s := Handler{
				service: tt.service,
			}
			w := httptest.NewRecorder()

			// act
			s.GetOrderByID(w, tt.req)
			respStatus, respJSON := getResp(t, w, "Order")

			// assert
			assert.Equal(t, tt.wantStatus, respStatus)
			assert.Equal(t, tt.wantJSON, respJSON)
		})
	}
}
