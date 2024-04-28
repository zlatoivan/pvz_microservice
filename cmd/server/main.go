package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler_grpc"
	orderinmemorycache "gitlab.ozon.dev/zlatoivan4/homework/internal/cache/in_memory/order"
	pvzinmemorycache "gitlab.ozon.dev/zlatoivan4/homework/internal/cache/in_memory/pvz"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/cache/redis"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/config"
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

	inMemPVZCache := pvzinmemorycache.New(5*time.Minute, 10*time.Minute)
	inMemOrderCache := orderinmemorycache.New(5*time.Minute, 10*time.Minute)
	redisPVZCache := redis.New(cfg.Redis, 5*time.Minute)
	defer func() {
		err = redisPVZCache.Close()
		if err != nil {
			log.Printf("redisPVZCache.Close: %v", err)
		}
	}()
	redisOrderCache := redis.New(cfg.Redis, 5*time.Minute)
	defer func() {
		err = redisOrderCache.Close()
		if err != nil {
			log.Printf("redisPVZCache.Close: %v", err)
		}
	}()

	pvzRepo := pvzrepo.New(database, inMemPVZCache)
	orderRepo := orderrepo.New(database, inMemOrderCache)

	pvzService := pvz.New(pvzRepo)
	orderService := order.New(orderRepo)

	controllerGRPC := handler_grpc.New(pvzService, orderService)

	mainServer := server.New(pvzService, orderService, controllerGRPC)

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

	err = mainServer.Run(ctx, cfg.Server, producer, redisPVZCache, redisOrderCache)
	if err != nil {
		return fmt.Errorf("mainServer.Run: %w", err)
	}

	<-ctx.Done()

	return nil
}
