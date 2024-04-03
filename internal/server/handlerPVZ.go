package server

import (
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
		delivery.RenderResponse(w, req, http.StatusBadRequest, delivery.MakeRespErrInvalidData(err))
		return
	}

	id, err := s.PvzService.CreatePVZ(req.Context(), newPVZ)
	if err != nil {
		log.Printf("[createPVZ] s.PvzService.CreatePVZ: %v\n", err)
		if errors.Is(err, repo.ErrorAlreadyExists) {
			delivery.RenderResponse(w, req, http.StatusConflict, delivery.MakeRespErrAlreadyExists(err))
		}
		delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
		return
	}

	log.Printf("PVZ created. id = %s\n", id)
	delivery.RenderResponse(w, req, http.StatusOK, delivery.MakeRespId(id))
}

// listPVZs gets list of PVZ
func (s Server) listPVZs(w http.ResponseWriter, req *http.Request) {
	list, err := s.PvzService.ListPVZs(req.Context())
	if err != nil {
		log.Printf("[listPVZs] s.PvzService.ListPVZs: %v\n", err)
		delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
		return
	}

	log.Printf("Got list of PVZs. Length = %d.\n", len(list))
	delivery.RenderResponse(w, req, http.StatusOK, delivery.MakeRespPVZList(list))
}

// getPVZByID gets PVZ by ID
func (s Server) getPVZByID(w http.ResponseWriter, req *http.Request) {
	id, err := delivery.GetPVZIDFromURL(req)
	if err != nil {
		log.Printf("[getPVZByID] GetPVZIDFromURL: %v", err)
		delivery.RenderResponse(w, req, http.StatusBadRequest, delivery.MakeRespErrInvalidData(err))
		return
	}

	pvz, err := s.PvzService.GetPVZByID(req.Context(), id)
	if err != nil {
		log.Printf("[getPVZByID] s.PvzService.GetPVZByID: %v\n", err)
		if errors.Is(err, repo.ErrorNotFound) {
			delivery.RenderResponse(w, req, http.StatusNotFound, delivery.MakeRespErrNotFoundByID(err))
			return
		}
		delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
		return
	}

	log.Printf("Got PVZ by ID\n")
	delivery.RenderResponse(w, req, http.StatusOK, delivery.MakeRespPVZ(pvz))
}

// updatePVZ updates PVZ
func (s Server) updatePVZ(w http.ResponseWriter, req *http.Request) {
	updPVZ, err := delivery.GetPVZFromReq(req)
	if err != nil {
		log.Printf("[updatePVZ] GetPVZFromReq: %v", err)
		delivery.RenderResponse(w, req, http.StatusBadRequest, delivery.MakeRespErrInvalidData(err))
		return
	}

	err = s.PvzService.UpdatePVZ(req.Context(), updPVZ)
	if err != nil {
		log.Printf("[updatePVZ] s.PvzService.UpdatePVZ: %v\n", err)
		if errors.Is(err, repo.ErrorNotFound) {
			delivery.RenderResponse(w, req, http.StatusNotFound, delivery.MakeRespErrNotFoundByID(err))
			return
		}
		delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
		return
	}

	log.Println("PVZ updated")
	delivery.RenderResponse(w, req, http.StatusOK, "PVZ updated")
}

// deletePVZ deletes PVZ
func (s Server) deletePVZ(w http.ResponseWriter, req *http.Request) {
	id, err := delivery.GetPVZIDFromURL(req)
	if err != nil {
		log.Printf("[deletePVZ] GetPVZIDFromURL: %v", err)
		delivery.RenderResponse(w, req, http.StatusBadRequest, delivery.MakeRespErrInvalidData(err))
		return
	}

	err = s.PvzService.DeletePVZ(req.Context(), id)
	if err != nil {
		log.Printf("[deletePVZ] s.PvzService.DeletePVZ: %v\n", err)
		if errors.Is(err, repo.ErrorNotFound) {
			delivery.RenderResponse(w, req, http.StatusNotFound, delivery.MakeRespErrNotFoundByID(err))
			return
		}
		delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
		return
	}

	log.Println("PVZ deleted")
	delivery.RenderResponse(w, req, http.StatusOK, "PVZ deleted")
}
