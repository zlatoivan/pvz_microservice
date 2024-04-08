package order

import (
	"errors"
	"log"
	"net/http"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handlers/delivery"
)

// CreateOrder creates order
func (s Handler) CreateOrder(w http.ResponseWriter, req *http.Request) {
	newOrder, err := delivery.GetOrderFromReq(req)
	if err != nil {
		log.Printf("[CreateOrder] GetOrderFromReq: %v\n", err)
		delivery.RenderResponse(w, req, http.StatusBadRequest, delivery.MakeRespErrInvalidData(err))
		return
	}

	id, err := s.service.CreateOrder(req.Context(), newOrder)
	if err != nil {
		log.Printf("[CreateOrder] s.Service.CreateOrder: %v\n", err)
		if errors.Is(err, ErrAlreadyExists) {
			delivery.RenderResponse(w, req, http.StatusConflict, delivery.MakeRespErrAlreadyExists(err))
			return
		}
		delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
		return
	}

	log.Printf("[CreateOrder] Order created. id = %s\n", id)
	delivery.RenderResponse(w, req, http.StatusCreated, delivery.MakeRespId(id))
}
