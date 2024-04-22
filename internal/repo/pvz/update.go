package pvz

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

const queryUpdatePVZ = `UPDATE pvzs SET name = $2, address = $3, contacts = $4 WHERE id = $1;`

// UpdatePVZ updates PVZ in repo
func (repo Repo) UpdatePVZ(ctx context.Context, updPVZ model.PVZ) error {
	options := pgx.TxOptions{
		IsoLevel:   pgx.Serializable,
		AccessMode: pgx.ReadWrite,
	}
	tx, err := repo.db.BeginTx(ctx, options)
	if err != nil {
		return fmt.Errorf("repo.db.BeginTx: %w", err)
	}

	t, err := tx.Exec(ctx, queryUpdatePVZ, updPVZ.ID, updPVZ.Name, updPVZ.Address, updPVZ.Contacts)
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

	repo.cache.Set(updPVZ.ID, updPVZ, 5*time.Minute)

	return nil
}
