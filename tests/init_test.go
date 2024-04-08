package tests

import (
	"context"
	"log"
	"os"
	"testing"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/config"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/repo/postgres"
)

var (
	DB postgres.Database
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	//cfg, err := config.New()
	//if err != nil {
	//	log.Fatal("config.New: ", err)
	//}
	//cfg.Pg.DBname = "test"
	cfg := config.Config{
		Pg: config.Pg{
			Host:     "localhost",
			Port:     "5431",
			DBname:   "test",
			User:     "postgres",
			Password: "postgres",
		},
	}

	db, err := postgres.New(ctx, cfg)
	if err != nil {
		log.Fatal("postgres.New: ", err)
	}
	DB = db

	exitCode := m.Run()

	DB.GetPool(ctx).Close()

	os.Exit(exitCode)
}
