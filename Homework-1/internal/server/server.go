package server

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
	"time"

	"gitlab.ozon.dev/zlatoivan4/homework/configs"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

type repo interface {
	CreatePVZ(ctx context.Context, pvz model.PVZ) (int, error)
	GetListOfPVZ(ctx context.Context) ([]model.PVZ, error)
	GetPVZByID(ctx context.Context, id int) (model.PVZ, error)
	UpdatePVZ(ctx context.Context, id int, updPVZ model.PVZ) error
	DeletePVZ(ctx context.Context, id int) error
}

type Server struct {
	repo repo
}

func NewServer(repo repo) (*Server, error) {
	return &Server{repo: repo}, nil
}

// Run starts the server
func (s *Server) Run(ctx context.Context, cfg configs.Config) error {
	router := s.createRouter(cfg)

	go func() {
		httpsPort := cfg.HttpsPort
		fmt.Printf("Server with HTTPS is running at port %s ...\n\n", httpsPort)
		err := http.ListenAndServeTLS(":"+httpsPort, "configs/server.crt", "configs/server.key", router)
		if err != nil {
			log.Printf("http.ListenAndServeTLS: %v", err)
		}
	}()

	go func() {
		httpPort := cfg.HttpPort
		fmt.Printf("Server with HTTP is running at port %s ...\n\n", httpPort)
		err := http.ListenAndServe(":"+httpPort, router)
		if err != nil {
			log.Printf("http.ListenAndServe: %v", err)
		}
	}()

	<-ctx.Done()
	log.Printf("\nThe program termination signal has been received.\nShutting down the tool...\n\n")
	time.Sleep(1 * time.Second)

	return nil
}

// createRouter creates http router
func (s *Server) createRouter(cfg configs.Config) *chi.Mux {
	r := chi.NewRouter()

	r.Use(mwBasicAuth(map[string]string{cfg.Auth.Login: cfg.Auth.Password}))

	r.Get("/", s.mainPage)

	r.Route("/api/v1/pvzs", func(r chi.Router) {
		r.With(mwGetData).Post("/", s.createPVZ) // Create
		r.Get("/", s.getListOfPVZ)               // List
		r.Route("/{pvzID}", func(r chi.Router) {
			r.Use(mwGetData)
			r.Get("/", s.getPVZByID)   // GetById
			r.Put("/", s.updatePVZ)    // Update
			r.Delete("/", s.deletePVZ) // Delete
		})
	})

	return r
}

func mwBasicAuth(creds map[string]string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, pass, ok := r.BasicAuth()
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				log.Printf("Wrong format of creds for auth\n\n")
				return
			}

			credPass, credUserOk := creds[user]
			if !credUserOk || subtle.ConstantTimeCompare([]byte(pass), []byte(credPass)) != 1 {
				w.WriteHeader(http.StatusUnauthorized)
				log.Printf("Wrong creds for auth\n\n")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func getIDFromURL(req *http.Request) (int, error) {
	idStr := chi.URLParam(req, "pvzID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("strconv.Atoi: %w", err)
	}
	return id, nil
}

func getPVZFromReq(req *http.Request) (model.PVZ, error) {
	var pvz model.PVZ
	err := json.NewDecoder(req.Body).Decode(&pvz)
	if err != nil {
		return model.PVZ{}, fmt.Errorf("json.NewDecoder().Decode: %w", err)
	}
	return pvz, nil
}

func getDataFromReq(req *http.Request) (int, model.PVZ, error) {
	id, err := getIDFromURL(req)
	if err != nil {
		return 0, model.PVZ{}, fmt.Errorf("getIDFromURL: %w", err)
	}
	pvz, err := getPVZFromReq(req)
	if err != nil {
		return 0, model.PVZ{}, fmt.Errorf("getPVZFromReq: %w", err)
	}
	pvz.ID = id
	return id, pvz, nil
}

func prepToPrint(pvz model.PVZ) string {
	return fmt.Sprintf("   Id: %d\n   Name: %s\n   Address: %s\n   Contacts: %s\n\n", pvz.ID, pvz.Name, pvz.Address, pvz.Contacts)
}

func mwGetData(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			id, err := getIDFromURL(req)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			ctx := context.WithValue(req.Context(), "id", id)
			next.ServeHTTP(w, req.WithContext(ctx))
		case http.MethodPost:
			pvz, err := getPVZFromReq(req)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			log.Printf("[MW]: POST request:\n" + prepToPrint(pvz))
			ctx := context.WithValue(req.Context(), "pvz", pvz)
			next.ServeHTTP(w, req.WithContext(ctx))
		case http.MethodPut:
			id, pvz, err := getDataFromReq(req)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			log.Printf("[MW]: PUT request:\n" + prepToPrint(pvz))
			ctx := context.WithValue(req.Context(), "pvz", pvz)
			ctx = context.WithValue(ctx, "id", id)
			next.ServeHTTP(w, req.WithContext(ctx))
		case http.MethodDelete:
			id, err := getIDFromURL(req)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			log.Printf("[MW]: DELETE request:\nid = %d\n\n", id)
			ctx := context.WithValue(req.Context(), "id", id)
			next.ServeHTTP(w, req.WithContext(ctx))
		}
	})
}

// mainPage shows main page
func (s *Server) mainPage(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("This is the main page!\n"))
}

// createPVZ creates PVZ
func (s *Server) createPVZ(w http.ResponseWriter, req *http.Request) {
	newPVZ := req.Context().Value("pvz").(model.PVZ)

	id, err := s.repo.CreatePVZ(req.Context(), newPVZ)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("[createPVZ][s.repo.CreatePVZ]: %v\n\n", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Printf("PVZ created!\nid = %d\n\n", id)
}

// getListOfPVZ gets list of PVZ
func (s *Server) getListOfPVZ(w http.ResponseWriter, req *http.Request) {
	list, err := s.repo.GetListOfPVZ(req.Context())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("[getListOfPVZ][s.repo.GetListOfPVZ]: %v\n\n", err)
		return
	}

	fmt.Println("PVZ list:")
	for i, p := range list {
		fmt.Printf("%d) Id: %d\n   Name: %s\n   Address: %s\n   Contacts: %s\n", i+1, p.ID, p.Name, p.Address, p.Contacts)
	}
	fmt.Println()

	w.WriteHeader(http.StatusOK)
}

// getPVZByID gets PVZ by ID
func (s *Server) getPVZByID(w http.ResponseWriter, req *http.Request) {
	id := req.Context().Value("id").(int)

	pvz, err := s.repo.GetPVZByID(req.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("[getPVZByID][s.repo.GetPVZByID]: %v\n\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Printf("PVZ by ID:\n" + prepToPrint(pvz))
}

// updatePVZ updates PVZ
func (s *Server) updatePVZ(w http.ResponseWriter, req *http.Request) {
	id := req.Context().Value("id").(int)
	updPVZ := req.Context().Value("pvz").(model.PVZ)

	err := s.repo.UpdatePVZ(req.Context(), id, updPVZ)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("[updatePVZ][s.repo.UpdatePVZ]: %v\n\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Printf("PVZ updated!\n\n")
}

// deletePVZ deletes PVZ
func (s *Server) deletePVZ(w http.ResponseWriter, req *http.Request) {
	id := req.Context().Value("id").(int)

	err := s.repo.DeletePVZ(req.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("[deletePVZ][s.repo.DeletePVZ]: %v\n\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Printf("PVZ deleted!\n\n")
}
