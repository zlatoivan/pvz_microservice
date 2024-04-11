package order

import (
	"errors"
	"log"
	"net/http"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/delivery"
)

// UpdateOrder updates Order
func (s Handler) UpdateOrder(w http.ResponseWriter, req *http.Request) {
	updOrder, err := delivery.GetOrderFromReq(req)
	if err != nil {
		log.Printf("[UpdateOrder] GetOrderFromReq: %v", err)
		delivery.RenderResponse(w, req, http.StatusBadRequest, delivery.MakeRespErrInvalidData(err))
		return
	}

	err = s.service.UpdateOrder(req.Context(), updOrder)
	if err != nil {
		log.Printf("[UpdateOrder] s.Service.UpdateOrder: %v\n", err)
		if errors.Is(err, ErrNotFound) {
			delivery.RenderResponse(w, req, http.StatusNotFound, delivery.MakeRespErrNotFoundByID(err))
			return
		}
		delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
		return
	}

	log.Println("[UpdateOrder] Order updated")
	delivery.RenderResponse(w, req, http.StatusOK, delivery.MakeRespComment("Order updated"))
}
