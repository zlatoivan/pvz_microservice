package order

import (
	"log"
	"net/http"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/delivery"
)

// ListClientOrders gets list of client orders
func (s Handler) ListClientOrders(w http.ResponseWriter, req *http.Request) {
	id, err := delivery.GetIDFromReq(req)
	if err != nil {
		log.Printf("[ListClientOrders] GetIDFromReq: %v", err)
		delivery.RenderResponse(w, req, http.StatusBadRequest, delivery.MakeRespErrInvalidData(err))
		return
	}

	list, err := s.service.ListClientOrders(req.Context(), id)
	if err != nil {
		log.Printf("[ListClientOrders] s.Service.ListClientOrders: %v\n", err)
		delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
		return
	}

	log.Printf("[ListClientOrders] Got list of clients orders! Length = %d.\n", len(list))
	delivery.RenderResponse(w, req, http.StatusOK, delivery.MakeRespOrderList(list))
}
