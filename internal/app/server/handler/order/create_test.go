package order

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/delivery"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/order/mock"
	"gitlab.ozon.dev/zlatoivan4/homework/tests/fixtures"
)

func TestHandler_CreateOrder(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	method := http.MethodPost
	endpoint := "/api/v1/orders"
	mc := minimock.NewController(t)

	ClientID, _ := uuid.Parse("88cda6c0-36fc-4be4-b976-e11a8a7a8f7e")
	StoresTill, _ := time.Parse(time.RFC3339, "2024-04-22T12:12:00Z")
	reqOrderGood := delivery.RequestOrder{
		ClientID:      ClientID,
		StoresTill:    StoresTill,
		Weight:        29,
		Cost:          1100,
		PackagingType: "box",
	}
	reqOrderBadReq := delivery.RequestOrder{ClientID: uuid.Nil}
	validOrder := delivery.GetOrderFromReqOrder(reqOrderGood)

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
			req:        genHTTPReq(t, method, endpoint, reqOrderGood),
			wantStatus: http.StatusCreated,
			wantJSON:   delivery.MakeRespId(fixtures.ID),
		},
		{
			name:       "bad request",
			service:    mock.NewServiceMock(mc),
			req:        genHTTPReq(t, method, endpoint, reqOrderBadReq),
			wantStatus: http.StatusBadRequest,
			wantJSON:   delivery.MakeRespErrInvalidData(errors.New("client id is nil")),
		},
		{
			name: "already exists",
			service: mock.NewServiceMock(mc).
				CreateOrderMock.
				Expect(ctx, validOrder).
				Return(uuid.UUID{}, ErrAlreadyExists),
			req:        genHTTPReq(t, method, endpoint, reqOrderGood),
			wantStatus: http.StatusConflict,
			wantJSON:   delivery.MakeRespErrAlreadyExists(ErrAlreadyExists),
		},
		{
			name: "server error",
			service: mock.NewServiceMock(mc).
				CreateOrderMock.
				Expect(ctx, validOrder).
				Return(uuid.UUID{}, errors.New("")),
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
			s.CreateOrder(w, tt.req)
			respStatus, respJSON := getResp(t, w, "ID")

			// assert
			assert.Equal(t, tt.wantStatus, respStatus)
			assert.Equal(t, tt.wantJSON, respJSON)
		})
	}
}
