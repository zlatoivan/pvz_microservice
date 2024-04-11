package order

import (
	"errors"
	"log"
	"net/http"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/delivery"
)

// ReturnOrder returns order
func (s Handler) ReturnOrder(w http.ResponseWriter, req *http.Request) {
	clientID, id, err := delivery.GetDataForReturnFromReq(req)
	if err != nil {
		log.Printf("[ReturnOrder] GetDataForGiveOut: %v", err)
		delivery.RenderResponse(w, req, http.StatusBadRequest, delivery.MakeRespErrInvalidData(err))
		return
	}

	err = s.service.ReturnOrder(req.Context(), clientID, id)
	if err != nil {
		log.Printf("[ReturnOrder] s.Service.ReturnOrder: %v\n", err)
		if errors.Is(err, ErrNotFound) {
			delivery.RenderResponse(w, req, http.StatusNotFound, delivery.MakeRespErrNotFoundByID(err))
			return
		}
		delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
		return
	}

	log.Println("[ReturnOrder] Order is returned")
	delivery.RenderResponse(w, req, http.StatusOK, delivery.MakeRespComment("Order returned"))
}
