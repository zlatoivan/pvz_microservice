package order

import (
	"log"
	"net/http"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/delivery"
)

// ListReturnedOrders returns list of returned orders
func (s Handler) ListReturnedOrders(w http.ResponseWriter, req *http.Request) {
	list, err := s.service.ListReturnedOrders(req.Context())
	if err != nil {
		log.Printf("[ListReturnedOrders] s.Service.ListReturnedOrders: %v\n", err)
		delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
		return
	}

	log.Printf("[ListReturnedOrders] Got list of returned orders! Length = %d.\n", len(list))
	delivery.RenderResponse(w, req, http.StatusOK, delivery.MakeRespOrderList(list))
}
