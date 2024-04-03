package migration

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/storage/postgres"
)

func New(ctx context.Context, database postgres.Database) error {
	db := stdlib.OpenDBFromPool(database.GetPool(ctx))
	defer func() {
		err := db.Close()
		if err != nil {
			log.Printf("db.Close: %v\n", err)
		}
	}()

	err := goose.RunContext(ctx, "up", db, "migrations")
	if err != nil {
		return fmt.Errorf("goose.RunContext: %w", err)
	}
	return nil
}
