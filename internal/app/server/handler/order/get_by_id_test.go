package order

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/delivery"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/order/mock"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
	"gitlab.ozon.dev/zlatoivan4/homework/tests/fixtures"
)

func TestHandler_GetOrderByID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	method := http.MethodGet
	endpoint := "/api/v1/orders/id"
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
			req:        genHTTPReq(t, method, endpoint, reqIDGood),
			wantStatus: http.StatusOK,
			wantJSON:   delivery.MakeRespOrder(validOrder),
		},
		{
			name:       "bad request",
			service:    mock.NewServiceMock(mc),
			req:        genHTTPReq(t, method, endpoint, reqIDBadReq),
			wantStatus: http.StatusBadRequest,
			wantJSON:   delivery.MakeRespErrInvalidData(errors.New("json.Unmarshal: json: cannot unmarshal string into Go value of type delivery.RequestID")),
		},
		{
			name: "not found",
			service: mock.NewServiceMock(mc).
				GetOrderByIDMock.
				Return(model.Order{}, ErrNotFound),
			req:        genHTTPReq(t, method, endpoint, reqIDGood),
			wantStatus: http.StatusNotFound,
			wantJSON:   delivery.MakeRespErrNotFoundByID(errors.New("not found")),
		},
		{
			name: "server error",
			service: mock.NewServiceMock(mc).
				GetOrderByIDMock.
				Return(model.Order{}, errors.New("")),
			req:        genHTTPReq(t, method, endpoint, reqIDGood),
			wantStatus: http.StatusInternalServerError,
			wantJSON:   delivery.MakeRespErrInternalServer(errors.New("")),
		},
	}

	for _, tt := range tests {
		tt := tt
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
