package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/server/delivery"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/storage/repo"
)

// createPVZ creates PVZ
func (s Server) createPVZ(w http.ResponseWriter, req *http.Request) {
	newPVZ, err := delivery.GetPVZWithoutIDFromReq(req)
	if err != nil {
		log.Printf("[createPVZ] GetPVZWithoutIDFromReq: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		WriteComment(w, "Invalid data: "+err.Error())
		return
	}

	id, err := s.PvzService.CreatePVZ(req.Context(), newPVZ)
	if err != nil {
		log.Printf("[createPVZ] s.PvzService.CreatePVZ: %v\n", err)
		if errors.Is(err, repo.ErrorAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			WriteComment(w, "ID already exists")
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("PVZ created! id = %s", id)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(ResponseID{ID: id})
	if err != nil {
		log.Printf("[createPVZ] json.NewEncoder().Encode: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// listPVZs gets list of PVZ
func (s Server) listPVZs(w http.ResponseWriter, req *http.Request) {
	list, err := s.PvzService.ListPVZs(req.Context())
	if err != nil {
		log.Printf("[listPVZs] s.PvzService.ListPVZs: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("Got PVZ list! Length = %d.\n", len(list))

	w.Header().Set("Content-Type", "application/json")
	if len(list) == 0 {
		w.WriteHeader(http.StatusOK)
		WriteComment(w, "No PVZ in database")
		return
	}

	err = json.NewEncoder(w).Encode(list)
	if err != nil {
		log.Printf("[listPVZs] json.NewEncoder().Encode: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// getPVZByID gets PVZ by ID
func (s Server) getPVZByID(w http.ResponseWriter, req *http.Request) {
	id, err := delivery.GetPVZIDFromURL(req)
	if err != nil {
		log.Printf("[getPVZByID] GetPVZIDFromURL: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		WriteComment(w, "Invalid data: "+err.Error())
		return
	}

	pvz, err := s.PvzService.GetPVZByID(req.Context(), id)
	if err != nil {
		log.Printf("[getPVZByID] s.PvzService.GetPVZByID: %v\n", err)
		if errors.Is(err, repo.ErrorNotFound) {
			w.WriteHeader(http.StatusNotFound)
			WriteComment(w, "PVZ not found by this ID")
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("PVZ by ID:\n" + delivery.PrepToPrintPVZ(pvz))

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(pvz)
	if err != nil {
		log.Printf("[getPVZByID] json.NewEncoder().Encode: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// updatePVZ updates PVZ
func (s Server) updatePVZ(w http.ResponseWriter, req *http.Request) {
	updPVZ, err := delivery.GetPVZFromReq(req)
	if err != nil {
		log.Printf("[updatePVZ] GetPVZFromReq: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		WriteComment(w, "Invalid data: "+err.Error())
		return
	}

	err = s.PvzService.UpdatePVZ(req.Context(), updPVZ)
	if err != nil {
		log.Printf("[updatePVZ] s.PvzService.UpdatePVZ: %v\n", err)
		if errors.Is(err, repo.ErrorNotFound) {
			w.WriteHeader(http.StatusNotFound)
			WriteComment(w, "PVZ not found by this ID")
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("PVZ updated!")

	w.WriteHeader(http.StatusOK)
}

// deletePVZ deletes PVZ
func (s Server) deletePVZ(w http.ResponseWriter, req *http.Request) {
	id, err := delivery.GetPVZIDFromURL(req)
	if err != nil {
		log.Printf("[deletePVZ] GetPVZIDFromURL: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		WriteComment(w, "Invalid data: "+err.Error())
		return
	}

	err = s.PvzService.DeletePVZ(req.Context(), id)
	if err != nil {
		log.Printf("[deletePVZ] s.PvzService.DeletePVZ: %v\n", err)
		if errors.Is(err, repo.ErrorNotFound) {
			w.WriteHeader(http.StatusNotFound)
			WriteComment(w, "PVZ not found by this ID")
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("PVZ deleted!")

	w.WriteHeader(http.StatusOK)
}
