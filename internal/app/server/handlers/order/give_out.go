package order

import (
	"errors"
	"log"
	"net/http"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handlers/delivery"
)

// GiveOutOrders gives out a list of orders
func (s Handler) GiveOutOrders(w http.ResponseWriter, req *http.Request) {
	clientID, ids, err := delivery.GetDataForGiveOutFromReq(req)
	if err != nil {
		log.Printf("[GiveOutOrders] GetDataForGiveOutFromReq: %v", err)
		delivery.RenderResponse(w, req, http.StatusBadRequest, delivery.MakeRespErrInvalidData(err))
		return
	}

	err = s.service.GiveOutOrders(req.Context(), clientID, ids)
	if err != nil {
		log.Printf("[GiveOutOrders] s.Service.GiveOutOrders: %v\n", err)
		if errors.Is(err, ErrNotFound) {
			delivery.RenderResponse(w, req, http.StatusNotFound, delivery.MakeRespErrNotFoundByID(err))
			return
		}
		delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
		return
	}

	log.Println("[GiveOutOrders] Orders are given out")
	delivery.RenderResponse(w, req, http.StatusOK, delivery.MakeRespComment("Orders are given out"))
}
