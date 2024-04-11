package order

import (
	"log"
	"net/http"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/delivery"
)

// ListOrders gets list of orders
func (s Handler) ListOrders(w http.ResponseWriter, req *http.Request) {
	list, err := s.service.ListOrders(req.Context())
	if err != nil {
		log.Printf("[ListOrders] s.Service.ListOrders: %v\n", err)
		delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
		return
	}

	log.Printf("[ListOrders] Got list of orders. Length = %d.\n", len(list))
	delivery.RenderResponse(w, req, http.StatusOK, delivery.MakeRespOrderList(list))
}
