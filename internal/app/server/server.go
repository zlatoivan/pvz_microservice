package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"google.golang.org/grpc"

	pb "gitlab.ozon.dev/zlatoivan4/homework/internal/pkg/pb"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/order"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/pvz"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler_grpc"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/middleware"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/config"
)

type Server struct {
	pvzService     pvz.Service
	orderService   order.Service
	controllerGRPC handler_grpc.Controller
}

func New(pvzService pvz.Service, orderService order.Service, controllerGRPC handler_grpc.Controller) Server {
	server := Server{
		pvzService:     pvzService,
		orderService:   orderService,
		controllerGRPC: controllerGRPC,
	}
	return server
}

//func redirectToHTTPS(w http.ResponseWriter, req *http.Data) {
//	http.Redirect(w, req, "https://localhost:9001"+req.RequestURI, http.StatusMovedPermanently)
//}

// Run starts the server
func (s Server) Run(ctx context.Context, cfg config.Server, producer middleware.Producer, redisPVZCache pvz.Redis, redisOrderCache order.Redis) error {
	router := s.createRouter(cfg, producer, redisPVZCache, redisOrderCache)
	httpsServer := &http.Server{Addr: "localhost:" + cfg.HttpsPort, Handler: router}
	httpServer := &http.Server{Addr: "localhost:" + cfg.HttpPort, Handler: router} // http.HandlerFunc(redirectToHTTPS)
	lis, err := net.Listen("tcp", ":"+cfg.GrpcPort)
	if err != nil {
		return fmt.Errorf("[grpcServer] net.Listen: %v\n", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterApiServer(grpcServer, s.controllerGRPC)

	wg := sync.WaitGroup{}

	log.Printf("[httpsServer] starting on %s\n", cfg.HttpsPort)
	wg.Add(1)
	go func() {
		httpsServerStart(httpsServer)
		wg.Done()
	}()

	log.Printf("[httpServer] starting on %s\n", cfg.HttpPort)
	wg.Add(1)
	go func() {
		httpServerStart(httpServer)
		wg.Done()
	}()

	log.Printf("[grpcServer] starting on %s\n", cfg.GrpcPort)
	wg.Add(1)
	go func() {
		grpcServerStart(grpcServer, lis)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		gracefulShutdown(ctx, httpsServer, httpServer, grpcServer)
		wg.Done()
	}()

	wg.Wait()

	return nil
}

func httpsServerStart(httpsServer *http.Server) {
	err := httpsServer.ListenAndServeTLS("internal/app/server/certificate/server.crt", "internal/app/server/certificate/server.key")
	if err != nil {
		log.Printf("[httpsServer] ListenAndServeTLS: %v\n", err)
	}
}

func httpServerStart(httpServer *http.Server) {
	err := httpServer.ListenAndServe()
	if err != nil {
		log.Printf("[httpServer] ListenAndServe: %v\n", err)
	}
}

func grpcServerStart(grpcServer *grpc.Server, lis net.Listener) {
	err := grpcServer.Serve(lis)
	if err != nil {
		log.Printf("[grpcServer] grpcServer.Serve: %v\n", err)
	}
}

func gracefulShutdown(ctx context.Context, httpsServer *http.Server, httpServer *http.Server, grpcServer *grpc.Server) {
	<-ctx.Done()
	ctxTo, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("[gracefulShutdown] shutting down")

	err := httpsServer.Shutdown(ctxTo)
	if err != nil {
		log.Printf("httpsServer.Shutdown: %v\n", err)
	}

	err = httpServer.Shutdown(ctxTo)
	if err != nil {
		log.Printf("httpServer.Shutdown: %v\n", err)
	}

	grpcServer.GracefulStop()

	log.Println("[gracefulShutdown] shut down successfully")
}
