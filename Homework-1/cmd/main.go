package main

import (
	"errors"
	"log"
	"os"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/storage/order"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/storage/pvz"
)

func validateArgs() error {
	args := os.Args[1:]
	if len(args) < 1 {
		return errors.New("not enough arguments")
	}
	return nil
}

func main() {
	err := validateArgs()
	if err != nil {
		log.Fatalf("validateArgs: %v", err)
	}

	orderStore, err := order.New()
	if err != nil {
		log.Fatalf("orderStore.New: %v", err)
	}

	pvzStore, err := pvz.New()
	if err != nil {
		log.Fatalf("pvzStore.New: %v", err)
	}

	app, err := app.New(orderStore, pvzStore)
	if err != nil {
		log.Fatalf("service.New: %v", err)
	}

	err = app.Work()
	if err != nil {
		log.Fatalf("app.Work: %v", err)
	}
}
