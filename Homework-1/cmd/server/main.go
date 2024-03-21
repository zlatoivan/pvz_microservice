package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"gitlab.ozon.dev/zlatoivan4/homework/configs"
	server_ "gitlab.ozon.dev/zlatoivan4/homework/internal/server"
	repo_ "gitlab.ozon.dev/zlatoivan4/homework/internal/storage/repo"
	"gitlab.ozon.dev/zlatoivan4/homework/pkg/db/postgres"
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

	go func() {}()
	cfg, err := configs.NewConfig()
	if err != nil {
		return fmt.Errorf("configs.NewConfig: %w", err)
	}

	database, err := postgres.NewDB(ctx, cfg)
	if err != nil {
		return fmt.Errorf("repo.NewDB: %w", err)
	}
	defer database.GetPool(ctx).Close()

	repo, err := repo_.NewRepo(database)
	if err != nil {
		return fmt.Errorf("repo.NewRepo: %w", err)
	}

	server, err := server_.NewServer(repo)
	if err != nil {
		return fmt.Errorf("server.NewServer: %w", err)
	}

	err = server.Run(ctx, cfg)
	if err != nil {
		return fmt.Errorf("server.Run: %w", err)
	}

	return nil
}
