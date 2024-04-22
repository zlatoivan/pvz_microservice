package pvz

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const queryDeletePVZ = `DELETE FROM pvzs WHERE id = $1;`

// DeletePVZ deletes PVZ from repo
func (repo Repo) DeletePVZ(ctx context.Context, id uuid.UUID) error {
	options := pgx.TxOptions{
		IsoLevel:   pgx.Serializable,
		AccessMode: pgx.ReadOnly,
	}
	tx, err := repo.db.BeginTx(ctx, options)
	if err != nil {
		return fmt.Errorf("repo.db.BeginTx: %w", err)
	}

	t, err := tx.Exec(ctx, queryDeletePVZ, id)
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
