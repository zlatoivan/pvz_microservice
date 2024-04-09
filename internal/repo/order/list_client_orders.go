package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

const querySelectClientOrders = `
SELECT id, client_id, weight, cost, stores_till, give_out_time, is_returned 
FROM orders 
WHERE client_id = $1 
  AND is_deleted = FALSE;`

// ListClientOrders gets list of Order from repo
func (repo Repo) ListClientOrders(ctx context.Context, id uuid.UUID) ([]model.Order, error) {
	var orders []model.Order
	err := repo.db.Select(ctx, &orders, querySelectClientOrders, id)
	if err != nil {
		return nil, fmt.Errorf("repo.db.Select: %w", err)
	}

	return orders, nil
}
