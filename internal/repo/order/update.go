package order

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

const queryUpdateOrder = `
UPDATE orders 
SET client_id = $1, weight = $2, cost = $3, stores_till = $4 
WHERE id = $5
  AND is_deleted = FALSE;`

// UpdateOrder updates Order in repo
func (repo Repo) UpdateOrder(ctx context.Context, updOrder model.Order) error {
	options := pgx.TxOptions{
		IsoLevel:   pgx.Serializable,
		AccessMode: pgx.ReadOnly,
	}
	tx, err := repo.db.BeginTx(ctx, options)
	if err != nil {
		return fmt.Errorf("repo.db.BeginTx: %w", err)
	}

	t, err := tx.Exec(ctx, queryUpdateOrder, updOrder.ClientID, updOrder.Weight, updOrder.Cost, updOrder.StoresTill, updOrder.ID)
	if err != nil {
		return fmt.Errorf("repo.db.Exec: %w", err)
	}
	if t.RowsAffected() == 0 {
		return ErrNotFound
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("tx.Commit: %w", err)
	}

	repo.cache.Set(updOrder.ID, updOrder, 5*time.Minute)

	return nil
}
