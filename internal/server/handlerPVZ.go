package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/storage/repo"
)

// createPVZ creates PVZ
func (s Server) createPVZ(w http.ResponseWriter, req *http.Request) {
	newPVZ, err := GetPVZWithoutIDFromReq(req)
	if err != nil {
		log.Printf("[createPVZ] GetPVZWithoutIDFromReq: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		writeComment(w, "Invalid data: "+err.Error())
		return
	}

	id, err := s.pvzService.CreatePVZ(req.Context(), newPVZ)
	if err != nil {
		log.Printf("[createPVZ] s.pvzService.CreatePVZ: %v\n", err)
		if errors.Is(err, repo.ErrorAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			writeComment(w, "ID already exists")
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("PVZ created! id = %s", id)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(responseID{ID: id})
	if err != nil {
		log.Printf("[createPVZ] json.NewEncoder().Encode: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// listPVZs gets list of PVZ
func (s Server) listPVZs(w http.ResponseWriter, req *http.Request) {
	list, err := s.pvzService.ListPVZs(req.Context())
	if err != nil {
		log.Printf("[listPVZs] s.pvzService.ListPVZs: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("Got PVZ list! Length = %d.\n", len(list))

	w.Header().Set("Content-Type", "application/json")
	if len(list) == 0 {
		w.WriteHeader(http.StatusOK)
		writeComment(w, "No PVZ in database")
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
	id, err := GetPVZIDFromURL(req)
	if err != nil {
		log.Printf("[getPVZByID] GetPVZIDFromURL: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		writeComment(w, "Invalid data: "+err.Error())
		return
	}

	pvz, err := s.pvzService.GetPVZByID(req.Context(), id)
	if err != nil {
		log.Printf("[getPVZByID] s.pvzService.GetPVZByID: %v\n", err)
		if errors.Is(err, repo.ErrorNotFound) {
			w.WriteHeader(http.StatusNotFound)
			writeComment(w, "PVZ not found by this ID")
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("PVZ by ID:\n" + PrepToPrintPVZ(pvz))

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
	updPVZ, err := GetPVZFromReq(req)
	if err != nil {
		log.Printf("[updatePVZ] GetPVZFromReq: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		writeComment(w, "Invalid data: "+err.Error())
		return
	}

	err = s.pvzService.UpdatePVZ(req.Context(), updPVZ)
	if err != nil {
		log.Printf("[updatePVZ] s.pvzService.UpdatePVZ: %v\n", err)
		if errors.Is(err, repo.ErrorNotFound) {
			w.WriteHeader(http.StatusNotFound)
			writeComment(w, "PVZ not found by this ID")
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("PVZ updated!")

	w.WriteHeader(http.StatusOK)
}

// deletePVZ deletes PVZ
func (s Server) deletePVZ(w http.ResponseWriter, req *http.Request) {
	id, err := GetPVZIDFromURL(req)
	if err != nil {
		log.Printf("[deletePVZ] GetPVZIDFromURL: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		writeComment(w, "Invalid data: "+err.Error())
		return
	}

	err = s.pvzService.DeletePVZ(req.Context(), id)
	if err != nil {
		log.Printf("[deletePVZ] s.pvzService.DeletePVZ: %v\n", err)
		if errors.Is(err, repo.ErrorNotFound) {
			w.WriteHeader(http.StatusNotFound)
			writeComment(w, "PVZ not found by this ID")
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("PVZ deleted!")

	w.WriteHeader(http.StatusOK)
}
