package order

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/delivery"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

// GetOrderByID gets Order by ID
func (s Handler) GetOrderByID(w http.ResponseWriter, req *http.Request) {
	id, err := delivery.GetIDFromReq(req)
	if err != nil {
		log.Printf("[GetOrderByID] GetIDFromReq: %v\n", err)
		delivery.RenderResponse(w, req, http.StatusBadRequest, delivery.MakeRespErrInvalidData(err))
		return
	}

	var order model.Order
	orderRaw, err := s.cache.Get(req.Context(), id.String())
	if err == nil {
		err = json.Unmarshal(orderRaw, &order)
		if err != nil {
			log.Printf("[GetOrderByID] json.Unmarshal: %v\n", err)
			delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
			return
		}
		err = s.cache.Set(req.Context(), id.String(), orderRaw)
		if err != nil {
			log.Printf("[GetOrderByID] s.cache.Set: %v\n", err)
			delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
			return
		}
		log.Printf("[GetOrderByID] Got order by ID\n")
		delivery.RenderResponse(w, req, http.StatusOK, delivery.MakeRespOrder(order))
		return
	}

	order, err = s.service.GetOrderByID(req.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			log.Printf("[GetOrderByID] s.service.GetOrderByID: %v\n", err)
			delivery.RenderResponse(w, req, http.StatusNotFound, delivery.MakeRespErrNotFoundByID(err))
			return
		}
		log.Printf("[GetOrderByID] s.service.GetOrderByID: %v\n", err)
		delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
		return
	}

	orderRaw, err = json.Marshal(order)
	if err != nil {
		log.Printf("[GetOrderByID] json.Marshal: %v\n", err)
		delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
		return
	}
	err = s.cache.Set(req.Context(), id.String(), orderRaw)
	if err != nil {
		log.Printf("[GetOrderByID] s.cache.Set: %v\n", err)
		delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
		return
	}

	log.Printf("[GetOrderByID] Got order by ID\n")
	delivery.RenderResponse(w, req, http.StatusOK, delivery.MakeRespOrder(order))
}
