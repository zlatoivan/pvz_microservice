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

func genHTTPCreateOrderReq(t *testing.T, reqOrder delivery.RequestOrder) *http.Request {
	body, err := json.Marshal(reqOrder)
	require.NoError(t, err)
	req := httptest.NewRequest(
		http.MethodPost,
		"http://localhost:9000/api/v1/orders",
		bytes.NewReader(body),
	)
	return req
}

func TestHandler_CreateOrder(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mc := minimock.NewController(t)

	reqOrderGood := delivery.RequestOrder{
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
		ClientID:      fixtures.ClientID,
		StoresTill:    fixtures.StoresTill,
		Weight:        fixtures.Weight,
		Cost:          fixtures.Cost,
		PackagingType: fixtures.PackagingType,
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
				CreateOrderMock.
				Expect(ctx, validOrder).
				Return(fixtures.ID, nil),
			req:        genHTTPCreateOrderReq(t, reqOrderGood),
			wantStatus: http.StatusCreated,
			wantJSON:   delivery.MakeRespId(fixtures.ID),
		},
		{
			name:       "bad request",
			service:    mock.NewServiceMock(mc),
			req:        genHTTPCreateOrderReq(t, reqOrderBadReq),
			wantStatus: http.StatusBadRequest,
			wantJSON:   delivery.MakeRespErrInvalidData(errors.New("client id is nil")),
		},
		{
			name: "already exists",
			service: mock.NewServiceMock(mc).
				CreateOrderMock.
				Expect(ctx, validOrder).
				Return(uuid.UUID{}, ErrAlreadyExists),
			req:        genHTTPCreateOrderReq(t, reqOrderGood),
			wantStatus: http.StatusConflict,
			wantJSON:   delivery.MakeRespErrAlreadyExists(ErrAlreadyExists),
		},
		{
			name: "server error",
			service: mock.NewServiceMock(mc).
				CreateOrderMock.
				Expect(ctx, validOrder).
				Return(uuid.UUID{}, errors.New("")),
			req:        genHTTPCreateOrderReq(t, reqOrderGood),
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
			s.CreateOrder(w, tt.req)
			respStatus, respJSON := getResp(t, w, "ID")

			// assert
			assert.Equal(t, tt.wantStatus, respStatus)
			assert.Equal(t, tt.wantJSON, respJSON)
		})
	}
}
