package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

const queryUpdateSoftDelete = `
UPDATE orders
SET is_deleted = TRUE
WHERE id = $1;`

// DeleteOrder deletes Order from repo
func (repo Repo) DeleteOrder(ctx context.Context, id uuid.UUID) error {
	t, err := repo.db.Exec(ctx, queryUpdateSoftDelete, id)
	if err != nil {
		return fmt.Errorf("repo.db.Exec: %w", err)
	}
	if t.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
