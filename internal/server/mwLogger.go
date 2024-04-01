package server

import (
	"log"
	"net/http"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/server/delivery"
)

func mwLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			log.Printf("[MW]: GET request.\n")
		case http.MethodPost:
			pvz, err := delivery.GetPVZWithoutIDFromReq(req)
			if err != nil {
				log.Printf("[mwLogger] getPVZFromReq: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				WriteComment(w, "Invalid data: "+err.Error())
				return
			}
			log.Printf("[MW]: POST request:\n" + delivery.PrepToPrintPVZ(pvz))
		case http.MethodPut:
			pvz, err := delivery.GetPVZFromReq(req)
			if err != nil {
				log.Printf("[mwLogger] getDataFromReq: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				WriteComment(w, "Invalid data: "+err.Error())
				return
			}
			log.Printf("[MW]: PUT request:\n" + delivery.PrepToPrintPVZ(pvz))
		case http.MethodDelete:
			id, err := delivery.GetPVZIDFromURL(req)
			if err != nil {
				log.Printf("[mwLogger] getIDFromURL: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				WriteComment(w, "Invalid data: "+err.Error())
				return
			}
			log.Printf("[MW]: DELETE request:\nid = %s\n", id)
		}
		next.ServeHTTP(w, req)
	})
}
