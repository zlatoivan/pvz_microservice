package pvz

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/delivery"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

// GetPVZByID gets PVZ by ID
func (s Handler) GetPVZByID(w http.ResponseWriter, req *http.Request) {
	id, err := delivery.GetIDFromReq(req)
	if err != nil {
		log.Printf("[GetPVZByID] GetIDFromReq: %v\n", err)
		delivery.RenderResponse(w, req, http.StatusBadRequest, delivery.MakeRespErrInvalidData(err))
		return
	}

	var pvz model.PVZ
	pvzRaw, err := s.cache.Get(req.Context(), id.String())
	if err == nil {
		err = json.Unmarshal(pvzRaw, &pvz)
		if err != nil {
			log.Printf("[GetPVZByID] json.Unmarshal: %v\n", err)
			delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
			return
		}
		err = s.cache.Set(req.Context(), id.String(), pvzRaw)
		if err != nil {
			log.Printf("[GetPVZByID] s.cache.Set: %v\n", err)
			delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
			return
		}
		log.Printf("[GetPVZByID] Got PVZ by ID\n")
		delivery.RenderResponse(w, req, http.StatusOK, delivery.MakeRespPVZ(pvz))
		return
	}

	pvz, err = s.service.GetPVZByID(req.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			log.Printf("[GetPVZByID] s.service.GetPVZByID: %v\n", err)
			delivery.RenderResponse(w, req, http.StatusNotFound, delivery.MakeRespErrNotFoundByID(err))
			return
		}
		log.Printf("[GetPVZByID] s.service.GetPVZByID: %v\n", err)
		delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
		return
	}

	pvzRaw, err = json.Marshal(pvz)
	if err != nil {
		log.Printf("[GetPVZByID] json.Marshal: %v\n", err)
		delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
		return
	}
	err = s.cache.Set(req.Context(), id.String(), pvzRaw)
	if err != nil {
		log.Printf("[GetPVZByID] s.cache.Set: %v\n", err)
		delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
		return
	}

	log.Printf("[GetPVZByID] Got PVZ by ID\n")
	delivery.RenderResponse(w, req, http.StatusOK, delivery.MakeRespPVZ(pvz))
}
