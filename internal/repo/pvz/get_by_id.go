package pvz

import (
	"context"

	"github.com/google/uuid"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

const querySelectPVZByID = `SELECT id, name, address, contacts FROM pvzs WHERE id = $1;`

// GetPVZByID gets PVZ by ID from repo
func (repo Repo) GetPVZByID(ctx context.Context, id uuid.UUID) (model.PVZ, error) {
	var pvz model.PVZ
	err := repo.db.Get(ctx, &pvz, querySelectPVZByID, id)
	if err != nil {
		return model.PVZ{}, ErrNotFound
	}

	return pvz, nil
}
