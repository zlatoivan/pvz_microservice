package main

import (
	"errors"
	"log"
	"os"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/service"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/storage"
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

	store, err := storage.New()
	if err != nil {
		log.Fatalf("storage.New: %v", err)
	}

	serv, err := service.New(store)
	if err != nil {
		log.Fatalf("service.New: %v", err)
	}

	err = serv.Work()
	if err != nil {
		log.Fatalf("serv.Work: %v", err)
	}
}
