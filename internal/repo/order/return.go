package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

const queryUpdateReturn = `
UPDATE orders 
SET is_returned = TRUE 
WHERE client_id = $1 
  AND id = $2 
  AND is_returned = FALSE 
  AND is_deleted = FALSE;`

// ReturnOrder gives out a client order
func (repo Repo) ReturnOrder(ctx context.Context, clientID uuid.UUID, id uuid.UUID) error {
	t, err := repo.db.Exec(ctx, queryUpdateReturn, clientID, id)
	if err != nil {
		return fmt.Errorf("repo.db.Exec: %w", err)
	}
	if t.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
