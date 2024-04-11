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
	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
	"gitlab.ozon.dev/zlatoivan4/homework/tests/fixtures"
)

func TestHandler_ListClientOrders(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	method := http.MethodGet
	endpoint := "/api/v1/orders/client/id"
	mc := minimock.NewController(t)

	reqIDGood := delivery.RequestID{ID: fixtures.ClientID}
	reqIDBadReq := delivery.RequestID{ID: uuid.Nil}
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
				ListClientOrdersMock.
				Expect(ctx, fixtures.ClientID).
				Return(validOrders, nil),
			req:        genHTTPReq(t, method, endpoint, reqIDGood),
			wantStatus: http.StatusOK,
			wantJSON:   delivery.MakeRespOrderList(validOrders),
		},
		{
			name:       "bad request",
			service:    mock.NewServiceMock(mc),
			req:        genHTTPReq(t, method, endpoint, reqIDBadReq),
			wantStatus: http.StatusBadRequest,
			wantJSON:   delivery.MakeRespErrInvalidData(errors.New("id is nil")),
		},
		{
			name: "server error",
			service: mock.NewServiceMock(mc).
				ListClientOrdersMock.
				Expect(ctx, fixtures.ClientID).
				Return(nil, errors.New("")),
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
			s.ListClientOrders(w, tt.req)
			respStatus, respJSON := getResp(t, w, "ListOrders")

			// assert
			assert.Equal(t, tt.wantStatus, respStatus)
			assert.Equal(t, tt.wantJSON, respJSON)
		})
	}
}