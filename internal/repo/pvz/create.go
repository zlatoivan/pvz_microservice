package pvz

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

const queryInsertPVZ = `INSERT INTO pvzs (name, address, contacts) VALUES ($1, $2, $3) RETURNING id;`
const queryCheckInsertPVZ = `SELECT COUNT(*) FROM pvzs WHERE id = $1;`

// CreatePVZ creates PVZ in repo
func (repo Repo) CreatePVZ(ctx context.Context, pvz model.PVZ) (uuid.UUID, error) {
	var id uuid.UUID
	err := repo.db.QueryRow(ctx, queryInsertPVZ, pvz.Name, pvz.Address, pvz.Contacts).Scan(&id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("repo.db.QueryRow().Scan: %w", err)
	}

	t, err := repo.db.Exec(ctx, queryCheckInsertPVZ, id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("repo.db.Exec: %w", err)
	}
	if t.RowsAffected() == 0 {
		return uuid.UUID{}, ErrAlreadyExists
	}

	return id, nil
}
