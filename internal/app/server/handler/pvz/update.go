package pvz

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/delivery"
)

// UpdatePVZ updates PVZ
func (s Handler) UpdatePVZ(w http.ResponseWriter, req *http.Request) {
	updPVZ, err := delivery.GetPVZFromReq(req)
	if err != nil {
		log.Printf("[UpdatePVZ] GetPVZFromReq: %v", err)
		delivery.RenderResponse(w, req, http.StatusBadRequest, delivery.MakeRespErrInvalidData(err))
		return
	}

	err = s.service.UpdatePVZ(req.Context(), updPVZ)
	if err != nil {
		log.Printf("[UpdatePVZ] s.Service.UpdatePVZ: %v\n", err)
		if errors.Is(err, ErrNotFound) {
			delivery.RenderResponse(w, req, http.StatusNotFound, delivery.MakeRespErrNotFoundByID(err))
			return
		}
		delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
		return
	}

	pvzRaw, err := json.Marshal(updPVZ)
	if err != nil {
		log.Printf("[UpdatePVZ] json.Marshal: %v\n", err)
		delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
		return
	}
	err = s.cache.Set(req.Context(), updPVZ.ID.String(), pvzRaw)
	if err != nil {
		log.Printf("[UpdatePVZ] s.cache.Set: %v\n", err)
		delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
		return
	}

	log.Println("[UpdatePVZ] PVZ updated")
	delivery.RenderResponse(w, req, http.StatusOK, delivery.MakeRespComment("PVZ updated"))
}
