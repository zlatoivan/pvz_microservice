package migrate

import (
	"context"
	"fmt"
	"log"
	"os"

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
	provider, err := goose.NewProvider(
		goose.DialectPostgres,
		db,
		os.DirFS("migrations"),
	)
	if err != nil {
		return fmt.Errorf("goose.NewProvider: %w", err)
	}
	results, err := provider.Up(ctx)
	if err != nil {
		return fmt.Errorf("provider.Up: %w", err)
	}
	log.Println("migration result:", results)
	return nil
}
