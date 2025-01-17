package pvz

import (
	"errors"
	"log"
	"net/http"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/delivery"
)

// DeletePVZ deletes PVZ
func (s Handler) DeletePVZ(w http.ResponseWriter, req *http.Request) {
	id, err := delivery.GetIDFromReq(req)
	if err != nil {
		log.Printf("[DeletePVZ] GetIDFromReq: %v", err)
		delivery.RenderResponse(w, req, http.StatusBadRequest, delivery.MakeRespErrInvalidData(err))
		return
	}

	err = s.service.DeletePVZ(req.Context(), id)
	if err != nil {
		log.Printf("[DeletePVZ] s.Service.DeletePVZ: %v\n", err)
		if errors.Is(err, ErrNotFound) {
			delivery.RenderResponse(w, req, http.StatusNotFound, delivery.MakeRespErrNotFoundByID(err))
			return
		}
		delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
		return
	}

	err = s.cache.Delete(req.Context(), id.String())
	if err != nil {
		log.Printf("[DeletePVZ] s.cache.Delete: %v\n", err)
		delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
		return
	}

	log.Printf("[DeletePVZ] PVZ deleted\n")
	delivery.RenderResponse(w, req, http.StatusOK, delivery.MakeRespComment("PVZ deleted"))
}
