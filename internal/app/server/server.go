package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/metrics"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/pkg/pb"

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

	// HTTP
	httpServer := &http.Server{Addr: "localhost:" + cfg.HttpPort, Handler: router} // http.HandlerFunc(redirectToHTTPS)

	// HTTPS
	httpsServer := &http.Server{Addr: "localhost:" + cfg.HttpsPort, Handler: router}

	// HTTP gateway to gRPC
	const gatewayPort = "9002"
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := pb.RegisterApiV1HandlerFromEndpoint(ctx, mux, "localhost:"+cfg.GrpcPort, opts)
	if err != nil {
		log.Printf("[httpServerForGrpc] pb.RegisterGatewayHandlerFromEndpoint: %v", err)
	}
	httpGatewayToGRPCServer := &http.Server{Addr: "localhost:" + gatewayPort, Handler: mux}

	// HTTP metrics
	const metricsPort = "9003"
	grpcMetrics := grpc_prometheus.NewServerMetrics()
	reg := prometheus.NewRegistry()
	reg.MustRegister(
		grpcMetrics,
		metrics.GivenOutOrdersCounterMetric,
		metrics.ClientGivenOutOrdersCounterMetric,
		metrics.ReturnedOrdersCounterMetric,
		metrics.DeletedPVZsCounterMetric,
	)

	httpMetricsServer := &http.Server{
		Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}),
		Addr:    "localhost:" + metricsPort,
	}

	// gRPC
	lis, err := net.Listen("tcp", ":"+cfg.GrpcPort)
	if err != nil {
		return fmt.Errorf("[grpcServer] net.Listen: %v\n", err)
	}
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpcMetrics.StreamServerInterceptor()),
		grpc.UnaryInterceptor(grpcMetrics.UnaryServerInterceptor()),
	)
	pb.RegisterApiV1Server(grpcServer, s.controllerGRPC)

	// Run
	wg := sync.WaitGroup{}

	log.Printf("[httpServer] starting on %s\n", cfg.HttpPort)
	wg.Add(1)
	go func() {
		httpServerRun(httpServer)
		wg.Done()
	}()

	log.Printf("[httpsServer] starting on %s\n", cfg.HttpsPort)
	wg.Add(1)
	go func() {
		httpsServerRun(httpsServer)
		wg.Done()
	}()

	log.Printf("[httpForGrpcServer] starting on %s\n", gatewayPort)
	wg.Add(1)
	go func() {
		httpGatewayToGrpcServerRun(httpGatewayToGRPCServer)
		wg.Done()
	}()

	log.Printf("[httpMetricsServer] starting on %s\n", metricsPort)
	wg.Add(1)
	go func() {
		httpMetricsServerRun(httpMetricsServer)
		wg.Done()
	}()

	log.Printf("[grpcServer] starting on %s\n", cfg.GrpcPort)
	wg.Add(1)
	go func() {
		grpcServerRun(grpcServer, lis)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		gracefulShutdown(ctx, httpsServer, httpServer, grpcServer, httpGatewayToGRPCServer, httpMetricsServer)
		wg.Done()
	}()

	wg.Wait()

	return nil
}

func httpServerRun(httpServer *http.Server) {
	err := httpServer.ListenAndServe()
	if err != nil {
		log.Printf("[httpServer] ListenAndServe: %v\n", err)
	}
}

func httpsServerRun(httpsServer *http.Server) {
	err := httpsServer.ListenAndServeTLS("internal/app/server/certificate/server.crt", "internal/app/server/certificate/server.key")
	if err != nil {
		log.Printf("[httpsServer] ListenAndServeTLS: %v\n", err)
	}
}

func httpGatewayToGrpcServerRun(httpGatewayToGRPCServer *http.Server) {
	err := httpGatewayToGRPCServer.ListenAndServe()
	if err != nil {
		log.Printf("[httpGatewayToGRPCServer] http.ListenAndServe: %v", err)
	}
}

func httpMetricsServerRun(httpMetricsServer *http.Server) {
	err := httpMetricsServer.ListenAndServe()
	if err != nil {
		log.Printf("[httpMetricsServer] http.ListenAndServe: %v", err)
	}
}

func grpcServerRun(grpcServer *grpc.Server, lis net.Listener) {
	err := grpcServer.Serve(lis)
	if err != nil {
		log.Printf("[grpcServer] grpcServer.Serve: %v\n", err)
	}
}

func gracefulShutdown(ctx context.Context, httpsServer *http.Server, httpServer *http.Server, grpcServer *grpc.Server, httpForGRPCServer *http.Server, httpMetricsServer *http.Server) {
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

	err = httpForGRPCServer.Shutdown(ctxTo)
	if err != nil {
		log.Printf("httpForGRPCServer.Shutdown: %v\n", err)
	}

	err = httpMetricsServer.Shutdown(ctxTo)
	if err != nil {
		log.Printf("httpMetricsServer.Shutdown: %v\n", err)
	}

	grpcServer.GracefulStop()

	log.Println("[gracefulShutdown] shut down successfully")
}
