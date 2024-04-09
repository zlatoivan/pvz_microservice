package order

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

const queryUpdateOrder = `
UPDATE orders 
SET client_id = $1, weight = $2, cost = $3, stores_till = $4 
WHERE id = $5
  AND is_deleted = FALSE;`

// UpdateOrder updates Order in repo
func (repo Repo) UpdateOrder(ctx context.Context, updOrder model.Order) error {
	t, err := repo.db.Exec(ctx, queryUpdateOrder, updOrder.ClientID, updOrder.Weight, updOrder.Cost, updOrder.StoresTill, updOrder.ID)
	if err != nil {
		return fmt.Errorf("repo.db.Exec: %w", err)
	}
	if t.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
