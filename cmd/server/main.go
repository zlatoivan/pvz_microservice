package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/config"
	orderinmemorycache "gitlab.ozon.dev/zlatoivan4/homework/internal/in_memory_cache/order"
	pvzinmemorycache "gitlab.ozon.dev/zlatoivan4/homework/internal/in_memory_cache/pvz"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/kafka"
	orderrepo "gitlab.ozon.dev/zlatoivan4/homework/internal/repo/order"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/repo/postgres"
	pvzrepo "gitlab.ozon.dev/zlatoivan4/homework/internal/repo/pvz"
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

	database, err := postgres.New(ctx, cfg.Pg)
	if err != nil {
		return fmt.Errorf("postgres.New: %w", err)
	}
	defer database.GetPool(ctx).Close()

	pvzCache := pvzinmemorycache.New(5*time.Minute, 10*time.Minute)
	orderCache := orderinmemorycache.New(5*time.Minute, 10*time.Minute)

	pvzRepo := pvzrepo.New(database, pvzCache)
	orderRepo := orderrepo.New(database, orderCache)

	pvzService := pvz.New(pvzRepo)
	orderService := order.New(orderRepo)

	mainServer := server.New(pvzService, orderService)

	producer, err := kafka.NewProducer(cfg.Brokers, cfg.Topic)
	if err != nil {
		return fmt.Errorf("kafka.NewProducer: %w", err)
	}
	defer func() {
		err = producer.Close()
		if err != nil {
			log.Printf("producer.Close: %v", err)
		}
	}()
	consumer, err := kafka.NewConsumer(cfg.Brokers)
	if err != nil {
		return fmt.Errorf("kafka.NewProducer: %w", err)
	}
	handler := kafka.GetLogHandler()
	err = consumer.Subscribe(cfg.Topic, handler)
	if err != nil {
		return fmt.Errorf("consumer.Subscribe: %w", err)
	}

	err = mainServer.Run(ctx, cfg.Server, producer)
	if err != nil {
		return fmt.Errorf("mainServer.Run: %w", err)
	}

	return nil
}
