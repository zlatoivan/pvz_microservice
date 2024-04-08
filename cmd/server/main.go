package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/config"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/repo"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/repo/postgres"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/service/order"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/service/pvz"
)

func main() {
	ctx := context.Background()

	err := bootstrap(ctx)
	if err != nil {
		log.Fatalf("[main] bootstrap: %v", err)
	}
}

func bootstrap(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("config.New: %w", err)
	}

	database, err := postgres.New(ctx, cfg)
	if err != nil {
		return fmt.Errorf("postgres.New: %w", err)
	}
	defer database.GetPool(ctx).Close()

	mainRepo := repo.New(database)

	pvzService := pvz.New(mainRepo)
	orderService := order.New(mainRepo)

	mainServer := server.New(pvzService, orderService)

	err = mainServer.Run(ctx, cfg)
	if err != nil {
		return fmt.Errorf("mainServer.Run: %w", err)
	}

	return nil
}
