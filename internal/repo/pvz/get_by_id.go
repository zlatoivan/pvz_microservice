package pvz

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

const querySelectPVZByID = `SELECT id, name, address, contacts FROM pvzs WHERE id = $1;`

// GetPVZByID gets PVZ by ID from repo
func (repo Repo) GetPVZByID(ctx context.Context, id uuid.UUID) (model.PVZ, error) {
	pvz, ok := repo.cache.Get(id)
	if ok {
		repo.cache.Set(id, pvz, 5*time.Minute)
		return pvz, nil
	}

	err := repo.db.Get(ctx, &pvz, querySelectPVZByID, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.PVZ{}, ErrNotFound
		}
		return model.PVZ{}, fmt.Errorf("repo.db.Get: %w", err)
	}

	repo.cache.Set(id, pvz, 5*time.Minute)
	return pvz, nil
}
