package pvz

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

const querySelectPVZ = `SELECT id, name, address, contacts FROM pvzs;`

// ListPVZs gets list of PVZ from repo
func (repo Repo) ListPVZs(ctx context.Context) ([]model.PVZ, error) {
	var pvzs []model.PVZ
	err := repo.db.Select(ctx, &pvzs, querySelectPVZ)
	if err != nil {
		return nil, fmt.Errorf("repo.db.Select: %w", err)
	}

	return pvzs, nil
}
