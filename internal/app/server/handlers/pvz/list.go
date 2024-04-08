package pvz

import (
	"log"
	"net/http"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handlers/delivery"
)

// ListPVZs gets list of PVZ
func (s Handler) ListPVZs(w http.ResponseWriter, req *http.Request) {
	list, err := s.service.ListPVZs(req.Context())
	if err != nil {
		log.Printf("[ListPVZs] s.Service.ListPVZs: %v\n", err)
		delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
		return
	}

	log.Printf("[ListPVZs] Got list of PVZs. Length = %d.\n", len(list))
	delivery.RenderResponse(w, req, http.StatusOK, delivery.MakeRespPVZList(list))
}
