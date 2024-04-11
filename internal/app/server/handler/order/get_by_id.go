package order

import (
	"errors"
	"log"
	"net/http"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/delivery"
)

// GetOrderByID gets Order by ID
func (s Handler) GetOrderByID(w http.ResponseWriter, req *http.Request) {
	id, err := delivery.GetIDFromReq(req)
	if err != nil {
		log.Printf("[GetOrderByID] GetIDFromReq: %v", err)
		delivery.RenderResponse(w, req, http.StatusBadRequest, delivery.MakeRespErrInvalidData(err))
		return
	}

	order, err := s.service.GetOrderByID(req.Context(), id)
	if err != nil {
		log.Printf("[GetOrderByID] s.Service.GetOrderByID: %v\n", err)
		if errors.Is(err, ErrNotFound) {
			delivery.RenderResponse(w, req, http.StatusNotFound, delivery.MakeRespErrNotFoundByID(err))
			return
		}
		delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
		return
	}

	log.Println("[GetOrderByID] Got order by ID")
	delivery.RenderResponse(w, req, http.StatusOK, delivery.MakeRespOrder(order))
}
