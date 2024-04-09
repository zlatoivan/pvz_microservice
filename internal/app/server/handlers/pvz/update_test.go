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

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handlers/delivery"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handlers/pvz/mock"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
	"gitlab.ozon.dev/zlatoivan4/homework/tests/fixtures"
)

func genHTTPUpdatePVZReq(t *testing.T, reqOrder interface{}) *http.Request {
	body, err := json.Marshal(reqOrder)
	require.NoError(t, err)
	req := httptest.NewRequest(
		http.MethodPut,
		"http://localhost:9000/api/v1/pvzs/id",
		bytes.NewReader(body),
	)
	return req
}

func TestHandler_UpdatePVZ(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mc := minimock.NewController(t)

	reqPVZGood := delivery.RequestPVZ{
		ID:       fixtures.ID,
		Name:     fixtures.Name,
		Address:  fixtures.Address,
		Contacts: fixtures.Contacts,
	}

	reqPVZBadReq := ""

	validPVZ := model.PVZ{
		ID:       fixtures.ID,
		Name:     fixtures.Name,
		Address:  fixtures.Address,
		Contacts: fixtures.Contacts,
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
				UpdatePVZMock.
				Expect(ctx, validPVZ).
				Return(nil),
			req:        genHTTPUpdatePVZReq(t, reqPVZGood),
			wantStatus: http.StatusOK,
			wantJSON:   delivery.MakeRespComment("PVZ updated"),
		},
		{
			name:       "bad request",
			service:    mock.NewServiceMock(mc),
			req:        genHTTPUpdatePVZReq(t, reqPVZBadReq),
			wantStatus: http.StatusBadRequest,
			wantJSON:   delivery.MakeRespErrInvalidData(errors.New("json.Unmarshal: json: cannot unmarshal string into Go value of type delivery.RequestPVZ")),
		},
		{
			name: "not found",
			service: mock.NewServiceMock(mc).
				UpdatePVZMock.
				Expect(ctx, validPVZ).
				Return(ErrNotFound),
			req:        genHTTPUpdatePVZReq(t, reqPVZGood),
			wantStatus: http.StatusNotFound,
			wantJSON:   delivery.MakeRespErrNotFoundByID(errors.New("not found")),
		},
		{
			name: "server error",
			service: mock.NewServiceMock(mc).
				UpdatePVZMock.
				Expect(ctx, validPVZ).
				Return(errors.New("")),
			req:        genHTTPUpdatePVZReq(t, reqPVZGood),
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
			s.UpdatePVZ(w, tt.req)
			respStatus, respJSON := getResp(t, w, "Comment")

			// assert
			assert.Equal(t, tt.wantStatus, respStatus)
			assert.Equal(t, tt.wantJSON, respJSON)
		})
	}
}
