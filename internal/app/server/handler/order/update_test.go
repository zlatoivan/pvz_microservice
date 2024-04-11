package order

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/delivery"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/order/mock"
	"gitlab.ozon.dev/zlatoivan4/homework/tests/fixtures"
)

func TestHandler_UpdateOrder(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	method := http.MethodPut
	endpoint := "/api/v1/orders/id"
	mc := minimock.NewController(t)

	reqOrderGood := fixtures.ReqCreateOrderGood
	reqOrderGood.ID = fixtures.ID
	reqOrderBadReq := delivery.RequestOrder{ClientID: uuid.Nil}
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
				UpdateOrderMock.
				Expect(ctx, validOrder).
				Return(nil),
			req:        genHTTPReq(t, method, endpoint, reqOrderGood),
			wantStatus: http.StatusOK,
			wantJSON:   delivery.MakeRespComment("Order updated"),
		},
		{
			name:       "bad request",
			service:    mock.NewServiceMock(mc),
			req:        genHTTPReq(t, method, endpoint, reqOrderBadReq),
			wantStatus: http.StatusBadRequest,
			wantJSON:   delivery.MakeRespErrInvalidData(errors.New("client id is nil")),
		},
		{
			name: "not found",
			service: mock.NewServiceMock(mc).
				UpdateOrderMock.
				Expect(ctx, validOrder).
				Return(ErrNotFound),
			req:        genHTTPReq(t, method, endpoint, reqOrderGood),
			wantStatus: http.StatusNotFound,
			wantJSON:   delivery.MakeRespErrNotFoundByID(errors.New("not found")),
		},
		{
			name: "server error",
			service: mock.NewServiceMock(mc).
				UpdateOrderMock.
				Expect(ctx, validOrder).
				Return(errors.New("")),
			req:        genHTTPReq(t, method, endpoint, reqOrderGood),
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
			s.UpdateOrder(w, tt.req)
			respStatus, respJSON := getResp(t, w, "Comment")

			// assert
			assert.Equal(t, tt.wantStatus, respStatus)
			assert.Equal(t, tt.wantJSON, respJSON)
		})
	}
}
