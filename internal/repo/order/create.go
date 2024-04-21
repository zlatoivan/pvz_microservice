package order

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

const queryInsertOrder = `
INSERT INTO orders (client_id, weight, cost, stores_till, give_out_time, packaging_type)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;`

// CreateOrder creates Order in repo
func (repo Repo) CreateOrder(ctx context.Context, order model.Order) (uuid.UUID, error) {
	options := pgx.TxOptions{
		IsoLevel:   pgx.Serializable,
		AccessMode: pgx.ReadOnly,
	}
	tx, err := repo.db.BeginTx(ctx, options)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("repo.db.BeginTx: %w", err)
	}

	var id uuid.UUID
	var timeNull time.Time
	err = repo.db.QueryRow(ctx, queryInsertOrder, order.ClientID, order.Weight, order.Cost, order.StoresTill, timeNull, order.PackagingType).Scan(&id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("repo.db.QueryRow().Scan: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("tx.Commit: %w", err)
	}

	repo.cache.Set(id, order, 5*time.Minute)

	return id, nil
}
