package pvz

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

const queryDeletePVZ = `DELETE FROM pvzs WHERE id = $1;`

// DeletePVZ deletes PVZ from repo
func (repo Repo) DeletePVZ(ctx context.Context, id uuid.UUID) error {
	t, err := repo.db.Exec(ctx, queryDeletePVZ, id)
	if err != nil {
		return fmt.Errorf("repo.db.Exec: %w", err)
	}
	if t.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
