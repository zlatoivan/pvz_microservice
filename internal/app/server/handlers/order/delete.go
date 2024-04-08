package order

import (
	"errors"
	"log"
	"net/http"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handlers/delivery"
)

// DeleteOrder deletes Order
func (s Handler) DeleteOrder(w http.ResponseWriter, req *http.Request) {
	id, err := delivery.GetIDFromReq(req)
	if err != nil {
		log.Printf("[DeleteOrder] GetIDFromReq: %v", err)
		delivery.RenderResponse(w, req, http.StatusBadRequest, delivery.MakeRespErrInvalidData(err))
		return
	}

	err = s.service.DeleteOrder(req.Context(), id)
	if err != nil {
		log.Printf("[DeleteOrder] s.Service.DeleteOrder: %v\n", err)
		if errors.Is(err, ErrNotFound) {
			delivery.RenderResponse(w, req, http.StatusNotFound, delivery.MakeRespErrNotFoundByID(err))
			return
		}
		delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
		return
	}

	log.Println("[DeleteOrder] Order deleted")
	delivery.RenderResponse(w, req, http.StatusOK, delivery.MakeRespComment("Order deleted"))
}
