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
	"gitlab.ozon.dev/zlatoivan4/homework/tests/fixtures"
)

func genHTTPGiveOutReq(t *testing.T, reqID interface{}) *http.Request {
	body, err := json.Marshal(reqID)
	require.NoError(t, err)
	req := httptest.NewRequest(
		http.MethodGet,
		"http://localhost:9000/api/v1/orders/id",
		bytes.NewReader(body),
	)
	return req
}

func TestHandler_GiveOutOrders(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
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
			req:        genHTTPGiveOutReq(t, reqGood),
			wantStatus: http.StatusOK,
			wantJSON:   delivery.MakeRespComment("Orders are given out"),
		},
		{
			name:       "bad request",
			service:    mock.NewServiceMock(mc),
			req:        genHTTPGiveOutReq(t, reqBadReq),
			wantStatus: http.StatusBadRequest,
			wantJSON:   delivery.MakeRespErrInvalidData(errors.New("json.Unmarshal: json: cannot unmarshal string into Go value of type delivery.RequestGiveOut")),
		},
		{
			name: "not found",
			service: mock.NewServiceMock(mc).
				GiveOutOrdersMock.
				Expect(ctx, reqGood.ClientID, reqGood.IDs).
				Return(ErrNotFound),
			req:        genHTTPGiveOutReq(t, reqGood),
			wantStatus: http.StatusNotFound,
			wantJSON:   delivery.MakeRespErrNotFoundByID(ErrNotFound),
		},
		{
			name: "server error",
			service: mock.NewServiceMock(mc).
				GiveOutOrdersMock.
				Expect(ctx, reqGood.ClientID, reqGood.IDs).
				Return(errors.New("")),
			req:        genHTTPGiveOutReq(t, reqGood),
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
			s.GiveOutOrders(w, tt.req)
			respStatus, respJSON := getResp(t, w, "Comment")

			// assert
			assert.Equal(t, tt.wantStatus, respStatus)
			assert.Equal(t, tt.wantJSON, respJSON)
		})
	}
}
