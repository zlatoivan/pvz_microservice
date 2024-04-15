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

func TestHandler_GiveOutOrders(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	method := http.MethodGet
	endpoint := "/api/v1/orders/id"
	mc := minimock.NewController(t)

	reqGood := delivery.RequestGiveOut{
		ClientID: fixtures.ClientID,
		IDs:      []uuid.UUID{fixtures.ID, fixtures.ID2},
	}
	reqBadReq := ""

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
				GiveOutOrdersMock.
				Expect(ctx, reqGood.ClientID, reqGood.IDs).
				Return(nil),
			req:        genHTTPReq(t, method, endpoint, reqGood),
			wantStatus: http.StatusOK,
			wantJSON:   delivery.MakeRespComment("Orders are given out"),
		},
		{
			name:       "bad request",
			service:    mock.NewServiceMock(mc),
			req:        genHTTPReq(t, method, endpoint, reqBadReq),
			wantStatus: http.StatusBadRequest,
			wantJSON:   delivery.MakeRespErrInvalidData(errors.New("render.DecodeJSON: json: cannot unmarshal string into Go value of type delivery.RequestGiveOut")),
		},
		{
			name: "not found",
			service: mock.NewServiceMock(mc).
				GiveOutOrdersMock.
				Expect(ctx, reqGood.ClientID, reqGood.IDs).
				Return(ErrNotFound),
			req:        genHTTPReq(t, method, endpoint, reqGood),
			wantStatus: http.StatusNotFound,
			wantJSON:   delivery.MakeRespErrNotFoundByID(ErrNotFound),
		},
		{
			name: "server error",
			service: mock.NewServiceMock(mc).
				GiveOutOrdersMock.
				Expect(ctx, reqGood.ClientID, reqGood.IDs).
				Return(errors.New("")),
			req:        genHTTPReq(t, method, endpoint, reqGood),
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
			s.GiveOutOrders(w, tt.req)
			respStatus, respJSON := getResp(t, w, "Comment")

			// assert
			assert.Equal(t, tt.wantStatus, respStatus)
			assert.Equal(t, tt.wantJSON, respJSON)
		})
	}
}
