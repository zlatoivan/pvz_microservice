package pvz

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/delivery"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/pvz/mock"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
	"gitlab.ozon.dev/zlatoivan4/homework/tests/fixtures"
)

func genHTTPListPVZsReq(t *testing.T) *http.Request {
	req := httptest.NewRequest(
		http.MethodGet,
		"http://localhost:9000/api/v1/pvzs",
		nil,
	)
	return req
}

func TestHandler_ListPVZs(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mc := minimock.NewController(t)

	validPVZs := []model.PVZ{
		fixtures.PVZ().Valid().V(),
		fixtures.PVZ().Valid().V(),
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
				ListPVZsMock.
				Expect(ctx).
				Return(validPVZs, nil),
			req:        genHTTPListPVZsReq(t),
			wantStatus: http.StatusOK,
			wantJSON:   delivery.MakeRespPVZList(validPVZs),
		},
		{
			name: "server error",
			service: mock.NewServiceMock(mc).
				ListPVZsMock.
				Expect(ctx).
				Return([]model.PVZ{}, errors.New("")),
			req:        genHTTPListPVZsReq(t),
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
			s.ListPVZs(w, tt.req)
			respStatus, respJSON := getResp(t, w, "ListPVZs")

			// assert
			assert.Equal(t, tt.wantStatus, respStatus)
			assert.Equal(t, tt.wantJSON, respJSON)
		})
	}
}
