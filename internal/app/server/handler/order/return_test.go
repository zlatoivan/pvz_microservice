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

func TestHandler_ReturnOrder(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	method := http.MethodDelete
	endpoint := "/api/v1/orders/id/client/id/return"
	mc := minimock.NewController(t)

	reqIDGood := delivery.RequestReturn{ClientID: fixtures.ClientID, ID: fixtures.ID}
	reqIDBadReq := delivery.RequestReturn{ClientID: uuid.Nil, ID: fixtures.ID}

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
				ReturnOrderMock.
				Expect(ctx, fixtures.ClientID, fixtures.ID).
				Return(nil),
			req:        genHTTPReq(t, method, endpoint, reqIDGood),
			wantStatus: http.StatusOK,
			wantJSON:   delivery.MakeRespComment("Order returned"),
		},
		{
			name:       "bad request",
			service:    mock.NewServiceMock(mc),
			req:        genHTTPReq(t, method, endpoint, reqIDBadReq),
			wantStatus: http.StatusBadRequest,
			wantJSON:   delivery.MakeRespErrInvalidData(errors.New("client id is nil")),
		},
		{
			name: "not found",
			service: mock.NewServiceMock(mc).
				ReturnOrderMock.
				Expect(ctx, fixtures.ClientID, fixtures.ID).
				Return(ErrNotFound),
			req:        genHTTPReq(t, method, endpoint, reqIDGood),
			wantStatus: http.StatusNotFound,
			wantJSON:   delivery.MakeRespErrNotFoundByID(errors.New("not found")),
		},
		{
			name: "server error",
			service: mock.NewServiceMock(mc).
				ReturnOrderMock.
				Expect(ctx, fixtures.ClientID, fixtures.ID).
				Return(errors.New("")),
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
			s.ReturnOrder(w, tt.req)
			respStatus, respJSON := getResp(t, w, "Comment")

			// assert
			assert.Equal(t, tt.wantStatus, respStatus)
			assert.Equal(t, tt.wantJSON, respJSON)
		})
	}
}