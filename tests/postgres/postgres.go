package postgres

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/config"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/repo/postgres"
)

func SetUp(ctx context.Context) (postgres.Database, error) {
	//cfg := config.Config{
	//	Pg: config.Pg{
	//		Host:     "localhost",
	//		Port:     "5431",
	//		DBname:   "test",
	//		User:     "postgres",
	//		Password: "postgres",
	//	},
	//}
	cfg, err := config.New()
	if err != nil {
		return postgres.Database{}, fmt.Errorf("config.New: %w", err)
	}

	db, err := postgres.New(ctx, cfg)
	if err != nil {
		return postgres.Database{}, fmt.Errorf("postgres.SeptUp: %w", err)
	}

	return db, nil
}

func TearDown(ctx context.Context, db postgres.Database) {
	db.GetPool(ctx).Close()
}
