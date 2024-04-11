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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/delivery"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/pvz/mock"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
	"gitlab.ozon.dev/zlatoivan4/homework/tests/fixtures"
)

func genHTTPCreatePVZReq(t *testing.T, reqOrder interface{}) *http.Request {
	body, err := json.Marshal(reqOrder)
	require.NoError(t, err)
	req := httptest.NewRequest(
		http.MethodPost,
		"http://localhost:9000/api/v1/pvzs",
		bytes.NewReader(body),
	)
	return req
}

func TestHandler_CreatePVZ(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mc := minimock.NewController(t)

	reqPVZGood := delivery.RequestPVZ{
		Name:     fixtures.Name,
		Address:  fixtures.Address,
		Contacts: fixtures.Contacts,
	}

	reqPVZBadReq := ""

	validPVZ := model.PVZ{
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
				CreatePVZMock.
				Expect(ctx, validPVZ).
				Return(fixtures.ID, nil),
			req:        genHTTPCreatePVZReq(t, reqPVZGood),
			wantStatus: http.StatusCreated,
			wantJSON:   delivery.MakeRespId(fixtures.ID),
		},
		{
			name:       "bad request",
			service:    mock.NewServiceMock(mc),
			req:        genHTTPCreatePVZReq(t, reqPVZBadReq),
			wantStatus: http.StatusBadRequest,
			wantJSON:   delivery.MakeRespErrInvalidData(errors.New("json.Unmarshal: json: cannot unmarshal string into Go value of type delivery.RequestPVZ")),
		},
		{
			name: "already exists",
			service: mock.NewServiceMock(mc).
				CreatePVZMock.
				Expect(ctx, validPVZ).
				Return(uuid.UUID{}, ErrAlreadyExists),
			req:        genHTTPCreatePVZReq(t, reqPVZGood),
			wantStatus: http.StatusConflict,
			wantJSON:   delivery.MakeRespErrAlreadyExists(ErrAlreadyExists),
		},
		{
			name: "server error",
			service: mock.NewServiceMock(mc).
				CreatePVZMock.
				Expect(ctx, validPVZ).
				Return(uuid.UUID{}, errors.New("")),
			req:        genHTTPCreatePVZReq(t, reqPVZGood),
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
			s.CreatePVZ(w, tt.req)
			respStatus, respJSON := getResp(t, w, "ID")

			// assert
			assert.Equal(t, tt.wantStatus, respStatus)
			assert.Equal(t, tt.wantJSON, respJSON)
		})
	}
}
