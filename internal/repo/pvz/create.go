package pvz

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

const queryInsertPVZ = `INSERT INTO pvzs (name, address, contacts) VALUES ($1, $2, $3) RETURNING id;`

// CreatePVZ creates PVZ in repo
func (repo Repo) CreatePVZ(ctx context.Context, pvz model.PVZ) (uuid.UUID, error) {
	options := pgx.TxOptions{
		IsoLevel:   pgx.Serializable,
		AccessMode: pgx.ReadWrite,
	}
	tx, err := repo.db.BeginTx(ctx, options)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("repo.db.BeginTx: %w", err)
	}

	var id uuid.UUID
	err = tx.QueryRow(ctx, queryInsertPVZ, pvz.Name, pvz.Address, pvz.Contacts).Scan(&id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("repo.db.QueryRow().Scan: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("tx.Commit: %w", err)
	}

	repo.cache.Set(id, pvz, 5*time.Minute)

	return id, nil
}
