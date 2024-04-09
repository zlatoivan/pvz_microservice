package order

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

const querySelectReturnedOrders = `
SELECT id, client_id, weight, cost, stores_till, give_out_time, is_returned 
FROM orders 
WHERE is_returned = TRUE 
  AND is_deleted = FALSE;`

// ListReturnedOrders gives out a list of returned orders
func (repo Repo) ListReturnedOrders(ctx context.Context) ([]model.Order, error) {
	var orders []model.Order
	err := repo.db.Select(ctx, &orders, querySelectReturnedOrders)
	if err != nil {
		return nil, fmt.Errorf("repo.db.Select: %w", err)
	}

	return orders, nil
}
