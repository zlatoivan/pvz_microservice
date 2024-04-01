package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/config"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

type pvzService interface {
	CreatePVZ(ctx context.Context, pvz model.PVZ) (uuid.UUID, error)
	ListPVZs(ctx context.Context) ([]model.PVZ, error)
	GetPVZByID(ctx context.Context, id uuid.UUID) (model.PVZ, error)
	UpdatePVZ(ctx context.Context, updPVZ model.PVZ) error
	DeletePVZ(ctx context.Context, id uuid.UUID) error
}

type orderService interface {
	CreateOrder(ctx context.Context, packagingType string, order model.Order) (uuid.UUID, error)
	ListOrders(ctx context.Context) ([]model.Order, error)
	GetOrderByID(ctx context.Context, id uuid.UUID) (model.Order, error)
	UpdateOrder(ctx context.Context, updPVZ model.Order) error
	DeleteOrder(ctx context.Context, id uuid.UUID) error
	ListClientOrders(ctx context.Context, id uuid.UUID) ([]model.Order, error)
	GiveOutOrders(ctx context.Context, id uuid.UUID, ids []uuid.UUID) error
	ReturnOrder(ctx context.Context, clientID uuid.UUID, id uuid.UUID) error
	ListReturnedOrders(ctx context.Context) ([]model.Order, error)
}

type Server struct {
	PvzService   pvzService
	OrderService orderService
}

func New(pvzService pvzService, orderService orderService) Server {
	server := Server{
		PvzService:   pvzService,
		OrderService: orderService,
	}
	return server
}

//func redirectToHTTPS(w http.ResponseWriter, req *http.Request) {
//	http.Redirect(w, req, "https://localhost:9001"+req.RequestURI, http.StatusMovedPermanently)
//}

// Run starts the server
func (s Server) Run(ctx context.Context, cfg config.Config) error {
	router := s.createRouter(cfg)
	httpsPort := cfg.Server.HttpsPort
	httpPort := cfg.Server.HttpPort
	httpsServer := &http.Server{Addr: "localhost:" + httpsPort, Handler: router}
	httpServer := &http.Server{Addr: "localhost:" + httpPort, Handler: router} // http.HandlerFunc(redirectToHTTPS)

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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
	log.Println("[servers] shut down successfully")

	return nil
}
