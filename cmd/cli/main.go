package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app"
	orderService_ "gitlab.ozon.dev/zlatoivan4/homework/internal/service/cli/order"
	pvzService_ "gitlab.ozon.dev/zlatoivan4/homework/internal/service/cli/pvz"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/storage/cli/order"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/storage/cli/pvz"
)

func main() {
	ctx := context.Background()

	err := bootstrap(ctx)
	if err != nil {
		log.Fatalf("[main] bootstrap: %v", err)
	}
}

func validateArgs() error {
	args := os.Args[1:]
	if len(args) < 1 {
		return errors.New("not enough arguments")
	}
	return nil
}

func bootstrap(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	err := validateArgs()
	if err != nil {
		return fmt.Errorf("validateArgs: %v", err)
	}

	orderStore, err := order.New()
	if err != nil {
		return fmt.Errorf("orderStore.New: %v", err)
	}
	defer func() {
		err = orderStore.Close()
		if err != nil {
			log.Printf("pvzStore.Close: %v", err)
		}
	}()

	pvzStore, err := pvz.New()
	if err != nil {
		return fmt.Errorf("pvzStore.New: %v", err)
	}
	defer func() {
		err = pvzStore.Close()
		if err != nil {
			log.Printf("pvzStore.Close: %v", err)
		}
	}()

	orderService, err := orderService_.New(orderStore)
	if err != nil {
		return fmt.Errorf("orderService_.New: %v", err)
	}

	pvzService, err := pvzService_.New(pvzStore)
	if err != nil {
		return fmt.Errorf("pvzService_.New: %v", err)
	}

	app, err := app.New(orderService, pvzService)
	if err != nil {
		return fmt.Errorf("app.New: %v", err)
	}

	err = app.Work(ctx, cancel)
	if err != nil {
		return fmt.Errorf("app.Work: %v", err)
	}

	return nil
}
