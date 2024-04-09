package pvz

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

const queryUpdatePVZ = `UPDATE pvzs SET name = $2, address = $3, contacts = $4 WHERE id = $1;`

// UpdatePVZ updates PVZ in repo
func (repo Repo) UpdatePVZ(ctx context.Context, updPVZ model.PVZ) error {
	t, err := repo.db.Exec(ctx, queryUpdatePVZ, updPVZ.ID, updPVZ.Name, updPVZ.Address, updPVZ.Contacts)
	if err != nil {
		return fmt.Errorf("repo.db.Exec: %w", err)
	}
	if t.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
