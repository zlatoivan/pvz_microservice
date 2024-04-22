package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const queryUpdateSoftDelete = `
UPDATE orders
SET is_deleted = TRUE
WHERE id = $1;`

// DeleteOrder deletes Order from repo
func (repo Repo) DeleteOrder(ctx context.Context, id uuid.UUID) error {
	options := pgx.TxOptions{
		IsoLevel:   pgx.Serializable,
		AccessMode: pgx.ReadOnly,
	}
	tx, err := repo.db.BeginTx(ctx, options)
	if err != nil {
		return fmt.Errorf("repo.db.BeginTx: %w", err)
	}

	t, err := tx.Exec(ctx, queryUpdateSoftDelete, id)
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

	repo.cache.Delete(id)

	return nil
}
