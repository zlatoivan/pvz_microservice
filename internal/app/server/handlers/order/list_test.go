package order

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handlers/delivery"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handlers/order/mock"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
	"gitlab.ozon.dev/zlatoivan4/homework/tests/fixtures"
)

func genHTTPListOrdersReq(t *testing.T) *http.Request {
	req := httptest.NewRequest(
		http.MethodGet,
		"http://localhost:9000/api/v1/orders",
		nil,
	)
	return req
}

func TestHandler_ListOrders(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mc := minimock.NewController(t)

	validOrders := []model.Order{
		fixtures.Order().Valid().V(),
		fixtures.Order().Valid().V(),
	}

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
				ListOrdersMock.
				Expect(ctx).
				Return(validOrders, nil),
			req:        genHTTPListOrdersReq(t),
			wantStatus: http.StatusOK,
			wantJSON:   delivery.MakeRespOrderList(validOrders),
		},
		{
			name: "server error",
			service: mock.NewServiceMock(mc).
				ListOrdersMock.
				Expect(ctx).
				Return([]model.Order{}, errors.New("")),
			req:        genHTTPListOrdersReq(t),
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
			s.ListOrders(w, tt.req)
			respStatus, respJSON := getResp(t, w, "ListOrders")

			// assert
			assert.Equal(t, tt.wantStatus, respStatus)
			assert.Equal(t, tt.wantJSON, respJSON)
		})
	}
}
