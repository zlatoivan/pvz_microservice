package pvz

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/delivery"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/pvz/mock"
	"gitlab.ozon.dev/zlatoivan4/homework/tests/fixtures"
)

func genHTTPDeletePVZReq(t *testing.T, reqID interface{}) *http.Request {
	body, err := json.Marshal(reqID)
	require.NoError(t, err)
	req := httptest.NewRequest(
		http.MethodDelete,
		"http://localhost:9000/api/v1/pvzs/id",
		bytes.NewReader(body),
	)
	return req
}

func TestHandler_DeletePVZ(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mc := minimock.NewController(t)

	reqIDGood := delivery.RequestID{ID: fixtures.ID}
	reqIDBadReq := ""

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
				DeletePVZMock.
				Expect(ctx, fixtures.ID).
				Return(nil),
			req:        genHTTPDeletePVZReq(t, reqIDGood),
			wantStatus: http.StatusOK,
			wantJSON:   delivery.MakeRespComment("PVZ deleted"),
		},
		{
			name:       "bad request",
			service:    mock.NewServiceMock(mc),
			req:        genHTTPDeletePVZReq(t, reqIDBadReq),
			wantStatus: http.StatusBadRequest,
			wantJSON:   delivery.MakeRespErrInvalidData(errors.New("render.DecodeJSON: json: cannot unmarshal string into Go value of type delivery.RequestID")),
		},
		{
			name: "not found",
			service: mock.NewServiceMock(mc).
				DeletePVZMock.
				Expect(ctx, fixtures.ID).
				Return(ErrNotFound),
			req:        genHTTPDeletePVZReq(t, reqIDGood),
			wantStatus: http.StatusNotFound,
			wantJSON:   delivery.MakeRespErrNotFoundByID(errors.New("not found")),
		},
		{
			name: "server error",
			service: mock.NewServiceMock(mc).
				DeletePVZMock.
				Expect(ctx, fixtures.ID).
				Return(errors.New("")),
			req:        genHTTPDeletePVZReq(t, reqIDGood),
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
			s.DeletePVZ(w, tt.req)
			respStatus, respJSON := getResp(t, w, "Comment")

			// assert
			assert.Equal(t, tt.wantStatus, respStatus)
			assert.Equal(t, tt.wantJSON, respJSON)
		})
	}
}
