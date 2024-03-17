package main

import (
	"context"
	"errors"
	"fmt"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/app"
	orderService_ "gitlab.ozon.dev/zlatoivan4/homework/internal/service/order"
	pvzService_ "gitlab.ozon.dev/zlatoivan4/homework/internal/service/pvz"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/storage/order"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/storage/pvz"
	"log"
	"os"
)

func main() {
	//ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bootstrap(ctx, cancel)
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

func bootstrap(ctx context.Context, cancel context.CancelFunc) error {
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
