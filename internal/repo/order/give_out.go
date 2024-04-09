package order

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const queryUpdateGiveOut = `
UPDATE orders 
SET give_out_time = $3 
WHERE client_id = $1 
  AND id = $2 
  AND is_deleted = FALSE;`

// GiveOutOrder gives out a client order
func (repo Repo) GiveOutOrder(ctx context.Context, clientID uuid.UUID, id uuid.UUID) error {
	t, err := repo.db.Exec(ctx, queryUpdateGiveOut, clientID, id, time.Now())
	if err != nil {
		return fmt.Errorf("repo.db.Exec: %w", err)
	}
	if t.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
