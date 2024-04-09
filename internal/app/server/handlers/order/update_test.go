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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handlers/delivery"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handlers/order/mock"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
	"gitlab.ozon.dev/zlatoivan4/homework/tests/fixtures"
)

func genHTTPUpdateOrderReq(t *testing.T, reqOrder delivery.RequestOrder) *http.Request {
	body, err := json.Marshal(reqOrder)
	require.NoError(t, err)
	req := httptest.NewRequest(
		http.MethodPut,
		"http://localhost:9000/api/v1/orders/id",
		bytes.NewReader(body),
	)
	return req
}

func TestHandler_UpdateOrder(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mc := minimock.NewController(t)

	reqOrderGood := delivery.RequestOrder{
		ID:            fixtures.ID,
		ClientID:      fixtures.ClientID,
		StoresTill:    "2024-04-22T13:14:00Z",
		Weight:        fixtures.Weight,
		Cost:          fixtures.Cost,
		PackagingType: fixtures.PackagingType,
	}

	reqOrderBadReq := delivery.RequestOrder{
		ClientID: uuid.Nil,
	}

	validOrder := model.Order{
		ID:            fixtures.ID,
		ClientID:      fixtures.ClientID,
		StoresTill:    fixtures.StoresTill,
		Weight:        fixtures.Weight,
		Cost:          fixtures.Cost,
		PackagingType: fixtures.PackagingType,
	}

	tests := []struct {
		name            string
		service         Service
		req             *http.Request
		wantStatus int
		wantJSON   interface{}
	}{
		{
			name: "success",
			service: mock.NewServiceMock(mc).
				UpdateOrderMock.
				Expect(ctx, validOrder).
				Return(nil),
			req:        genHTTPUpdateOrderReq(t, reqOrderGood),
			wantStatus: http.StatusOK,
			wantJSON:   delivery.MakeRespComment("Order updated"),
		},
		{
			name:       "bad request",
			service:    mock.NewServiceMock(mc),
			req:        genHTTPUpdateOrderReq(t, reqOrderBadReq),
			wantStatus: http.StatusBadRequest,
			wantJSON:   delivery.MakeRespErrInvalidData(errors.New("client id is nil")),
		},
		{
			name: "not found",
			service: mock.NewServiceMock(mc).
				UpdateOrderMock.
				Expect(ctx, validOrder).
				Return(ErrNotFound),
			req:        genHTTPUpdateOrderReq(t, reqOrderGood),
			wantStatus: http.StatusNotFound,
			wantJSON:   delivery.MakeRespErrNotFoundByID(errors.New("not found")),
		},
		{
			name: "server error",
			service: mock.NewServiceMock(mc).
				UpdateOrderMock.
				Expect(ctx, validOrder).
				Return(errors.New("")),
			req:        genHTTPUpdateOrderReq(t, reqOrderGood),
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
