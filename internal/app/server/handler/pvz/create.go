package pvz

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/delivery"
)

// CreatePVZ creates PVZ
func (s Handler) CreatePVZ(w http.ResponseWriter, req *http.Request) {
	newPVZ, err := delivery.GetPVZFromReq(req)
	if err != nil {
		log.Printf("[CreatePVZ] GetPVZFromReq: %v", err)
		delivery.RenderResponse(w, req, http.StatusBadRequest, delivery.MakeRespErrInvalidData(err))
		return
	}

	id, err := s.service.CreatePVZ(req.Context(), newPVZ)
	if err != nil {
		log.Printf("[CreatePVZ] s.Service.CreatePVZ: %v\n", err)
		if errors.Is(err, ErrAlreadyExists) {
			delivery.RenderResponse(w, req, http.StatusConflict, delivery.MakeRespErrAlreadyExists(err))
			return
		}
		delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
		return
	}

	pvzRaw, err := json.Marshal(newPVZ)
	if err != nil {
		log.Printf("[CreatePVZ] json.Marshal: %v\n", err)
		delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
		return
	}
	err = s.cache.Set(req.Context(), id.String(), pvzRaw)
	if err != nil {
		log.Printf("[CreatePVZ] s.cache.Set: %v\n", err)
		delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
		return
	}

	log.Printf("[CreatePVZ] PVZ created. id = %s\n", id)
	delivery.RenderResponse(w, req, http.StatusCreated, delivery.MakeRespId(id))
}
