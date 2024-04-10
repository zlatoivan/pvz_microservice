package order

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

const querySelectOrder = `
SELECT id, client_id, weight, cost, stores_till, give_out_time, is_returned, packaging_type 
FROM orders 
WHERE is_deleted = FALSE;`

// ListOrders gets list of Order from repo
func (repo Repo) ListOrders(ctx context.Context) ([]model.Order, error) {
	var orders []model.Order
	err := repo.db.Select(ctx, &orders, querySelectOrder)
	if err != nil {
		return nil, fmt.Errorf("repo.db.Select: %w", err)
	}

	return orders, nil
}
