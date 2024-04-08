package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handlers/order"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handlers/pvz"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/config"
)

type Server struct {
	pvzService   pvz.Service
	orderService order.Service
}

func New(pvzService pvz.Service, orderService order.Service) Server {
	server := Server{
		pvzService:   pvzService,
		orderService: orderService,
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
		err := httpsServer.ListenAndServeTLS("internal/app/server/certificates/server.crt", "internal/app/server/certificates/server.key")
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
