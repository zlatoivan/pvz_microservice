package middleware

import (
	"log"
	"net/http"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handlers/delivery"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			log.Printf("[MW]: GET request.\n")
		case http.MethodPost:
			pvz, err := delivery.GetPVZFromReq(req)
			if err != nil {
				log.Printf("[Logger] getPVZFromReq: %v", err)
				delivery.RenderResponse(w, req, http.StatusBadRequest, delivery.MakeRespErrInvalidData(err))
				return
			}
			log.Printf("[MW]: POST request:\n" + delivery.PrepToPrintPVZ(pvz))
		case http.MethodPut:
			pvz, err := delivery.GetPVZFromReq(req)
			if err != nil {
				log.Printf("[Logger] getDataFromReq: %v", err)
				delivery.RenderResponse(w, req, http.StatusBadRequest, delivery.MakeRespErrInvalidData(err))
				return
			}
			log.Printf("[MW]: PUT request:\n" + delivery.PrepToPrintPVZ(pvz))
		case http.MethodDelete:
			id, err := delivery.GetPVZFromReq(req)
			if err != nil {
				log.Printf("[Logger] getIDFromURL: %v", err)
				delivery.RenderResponse(w, req, http.StatusBadRequest, delivery.MakeRespErrInvalidData(err))
				return
			}
			log.Printf("[MW]: DELETE request:\nid = %s\n", id)
		}
		next.ServeHTTP(w, req)
	})
}
