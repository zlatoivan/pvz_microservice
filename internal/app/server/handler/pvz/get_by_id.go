package pvz

import (
	"errors"
	"log"
	"net/http"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/delivery"
)

// GetPVZByID gets PVZ by ID
func (s Handler) GetPVZByID(w http.ResponseWriter, req *http.Request) {
	id, err := delivery.GetIDFromReq(req)
	if err != nil {
		log.Printf("[GetPVZByID] GetIDFromReq: %v", err)
		delivery.RenderResponse(w, req, http.StatusBadRequest, delivery.MakeRespErrInvalidData(err))
		return
	}

	pvz, err := s.service.GetPVZByID(req.Context(), id)
	if err != nil {
		log.Printf("[GetPVZByID] s.Service.GetPVZByID: %v\n", err)
		if errors.Is(err, ErrNotFound) {
			delivery.RenderResponse(w, req, http.StatusNotFound, delivery.MakeRespErrNotFoundByID(err))
			return
		}
		delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
		return
	}

	log.Printf("[GetPVZByID] Got PVZ by ID\n")
	delivery.RenderResponse(w, req, http.StatusOK, delivery.MakeRespPVZ(pvz))
}
