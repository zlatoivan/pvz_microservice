package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/config"
	serverErrors "gitlab.ozon.dev/zlatoivan4/homework/internal/errors"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

type repo interface {
	CreatePVZ(ctx context.Context, pvz model.PVZ) (uuid.UUID, error)
	GetListOfPVZ(ctx context.Context) ([]model.PVZ, error)
	GetPVZByID(ctx context.Context, id uuid.UUID) (model.PVZ, error)
	UpdatePVZ(ctx context.Context, updPVZ model.PVZ) error
	DeletePVZ(ctx context.Context, id uuid.UUID) error
}

type Server struct {
	repo repo
}

func New(repo repo) Server {
	return Server{repo: repo}
}

func redirectToHTTPS(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "https://localhost:9001"+req.RequestURI, http.StatusMovedPermanently)
}

// Run starts the server
func (s Server) Run(ctx context.Context, cfg config.Config) error {
	router := s.createRouter(cfg)
	httpsPort := cfg.Server.HttpsPort
	httpPort := cfg.Server.HttpPort
	httpsServer := &http.Server{Addr: "localhost:" + httpsPort, Handler: router}
	httpServer := &http.Server{Addr: "localhost:" + httpPort, Handler: http.HandlerFunc(redirectToHTTPS)}

	go func() {
		log.Printf("[httpsServer] starting on %s\n", httpsPort)
		err := httpsServer.ListenAndServeTLS("internal/server/certs/server.crt", "internal/server/certs/server.key")
		if err != nil {
			log.Printf("[httpsServer] ListenAndServeTLS: %v\n", err)
		}
	}()

	go func() {
		log.Printf("[httpServer] starting on %s\n", httpPort)
		err := httpServer.ListenAndServe()
		if err != nil {
			log.Printf("[httpServer] ListenAndServe: %v\n", err)
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	log.Println("[servers] shutting down")
	err := httpsServer.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("httpsServer.Shutdown: %w", err)
	}
	err = httpServer.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("httpServer.Shutdown: %w", err)
	}
	<-ctx.Done()
	log.Println("[servers] shut down successfully")

	return nil
}

// createRouter creates http router
func (s Server) createRouter(cfg config.Config) *chi.Mux {
	pvzCreds := map[string]string{cfg.Server.Login: cfg.Server.Password}

	r := chi.NewRouter()

	r.NotFound(s.notFound)

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.RedirectSlashes)
	//r.Use(middleware.Logger)

	r.Get("/", s.mainPage)

	r.Route("/api/v1/pvzs", func(r chi.Router) {
		r.Use(middleware.BasicAuth("pvzs", pvzCreds))
		r.With(mwLogger).Post("/", s.createPVZ) // Create
		r.Get("/", s.ListPVZs)                  // List
		r.Route("/{pvzID}", func(r chi.Router) {
			r.Get("/", s.getPVZByID)                  // GetById
			r.With(mwLogger).Put("/", s.updatePVZ)    // Update
			r.With(mwLogger).Delete("/", s.deletePVZ) // Delete
		})
	})

	return r
}

// notFound informs that the page is not found
func (s Server) notFound(w http.ResponseWriter, _ *http.Request) {
	log.Println("Page not found")
	w.WriteHeader(http.StatusNotFound)
}

// mainPage shows the main page
func (s Server) mainPage(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("This is the main page!\n"))
}

func getIDFromURL(req *http.Request) (uuid.UUID, error) {
	idStr := chi.URLParam(req, "pvzID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("uuid.Parse: %w", err)
	}
	return id, nil
}

func getPVZFromReq(req *http.Request) (model.PVZ, error) {
	var pvz model.PVZ
	data, err := io.ReadAll(req.Body)
	if err != nil {
		return model.PVZ{}, fmt.Errorf("io.ReadAll: %w", err)
	}
	defer func() {
		err = req.Body.Close()
		if err != nil {
			log.Printf("[error] req.Body.Close: %v", err)
		}
	}()
	err = json.Unmarshal(data, &pvz)
	if err != nil {
		return model.PVZ{}, fmt.Errorf("json.NewDecoder().Decode: %w", err)
	}
	req.Body = io.NopCloser(bytes.NewBuffer(data))
	return pvz, nil
}

func getDataFromReq(req *http.Request) (model.PVZ, error) {
	id, err := getIDFromURL(req)
	if err != nil {
		return model.PVZ{}, fmt.Errorf("getIDFromURL: %w", err)
	}
	pvz, err := getPVZFromReq(req)
	if err != nil {
		return model.PVZ{}, fmt.Errorf("getPVZFromReq: %w", err)
	}
	pvz.ID = id
	return pvz, nil
}

func prepToPrint(pvz model.PVZ) string {
	if pvz.ID == uuid.Nil {
		return fmt.Sprintf("   Name: %s\n   Address: %s\n   Contacts: %s\n", pvz.Name, pvz.Address, pvz.Contacts)
	}
	return fmt.Sprintf("   Id: %s\n   Name: %s\n   Address: %s\n   Contacts: %s\n", pvz.ID, pvz.Name, pvz.Address, pvz.Contacts)
}

func mwLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			pvz, err := getPVZFromReq(req)
			if err != nil {
				log.Printf("[mwLogger] getPVZFromReq: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			log.Printf("[MW]: POST request:\n" + prepToPrint(pvz))
			next.ServeHTTP(w, req)
		case http.MethodPut:
			pvz, err := getDataFromReq(req)
			if err != nil {
				log.Printf("[mwLogger] getDataFromReq: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			log.Printf("[MW]: PUT request:\n" + prepToPrint(pvz))
			next.ServeHTTP(w, req)
		case http.MethodDelete:
			id, err := getIDFromURL(req)
			if err != nil {
				log.Printf("[mwLogger] getIDFromURL: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			log.Printf("[MW]: DELETE request:\nid = %d\n", id)
			next.ServeHTTP(w, req)
		}
	})
}

// createPVZ creates PVZ
func (s Server) createPVZ(w http.ResponseWriter, req *http.Request) {
	newPVZ, err := getPVZFromReq(req)
	if err != nil {
		log.Printf("[createPVZ] getPVZFromReq: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := s.repo.CreatePVZ(req.Context(), newPVZ)
	if err != nil {
		log.Printf("[createPVZ] s.repo.CreatePVZ: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("PVZ created! id = %s", id)

	w.WriteHeader(http.StatusCreated)
}

// ListPVZs gets list of PVZ
func (s Server) ListPVZs(w http.ResponseWriter, req *http.Request) {
	list, err := s.repo.GetListOfPVZ(req.Context())
	if err != nil {
		log.Printf("[ListPVZs] s.repo.GetListOfPVZ: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("Got PVZ list!")

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(list)
	if err != nil {
		log.Printf("[ListPVZs] json.NewEncoder().Encode: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// getPVZByID gets PVZ by ID
func (s Server) getPVZByID(w http.ResponseWriter, req *http.Request) {
	id, err := getIDFromURL(req)
	if err != nil {
		log.Printf("[getPVZByID] getIDFromURL: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pvz, err := s.repo.GetPVZByID(req.Context(), id)
	if err != nil {
		log.Printf("[getPVZByID] s.repo.GetPVZByID: %v\n", err)
		if errors.Is(err, serverErrors.ErrorNotFound) {
			w.WriteHeader(http.StatusConflict)
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("PVZ by ID:\n" + prepToPrint(pvz))

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(pvz)
	if err != nil {
		log.Printf("[getPVZByID] json.NewEncoder().Encode: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// updatePVZ updates PVZ
func (s Server) updatePVZ(w http.ResponseWriter, req *http.Request) {
	updPVZ, err := getDataFromReq(req)
	if err != nil {
		log.Printf("[updatePVZ] getDataFromReq: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.repo.UpdatePVZ(req.Context(), updPVZ)
	if err != nil {
		log.Printf("[updatePVZ] s.repo.UpdatePVZ: %v\n", err)
		if errors.Is(err, serverErrors.ErrorNotFound) {
			w.WriteHeader(http.StatusConflict)
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("PVZ updated!")

	w.WriteHeader(http.StatusOK)
}

// deletePVZ deletes PVZ
func (s Server) deletePVZ(w http.ResponseWriter, req *http.Request) {
	id, err := getIDFromURL(req)
	if err != nil {
		log.Printf("[deletePVZ] getIDFromURL: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.repo.DeletePVZ(req.Context(), id)
	if err != nil {
		log.Printf("[deletePVZ] s.repo.DeletePVZ: %v\n", err)
		if errors.Is(err, serverErrors.ErrorNotFound) {
			w.WriteHeader(http.StatusConflict)
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("PVZ deleted!")

	w.WriteHeader(http.StatusOK)
}
