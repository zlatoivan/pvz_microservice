package order

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

const querySelectOrderByID = `
SELECT id, client_id, weight, cost, stores_till, give_out_time, is_returned, packaging_type
FROM orders
WHERE id = $1
AND is_deleted = FALSE;`

// GetOrderByID gets Order by ID from repo
func (repo Repo) GetOrderByID(ctx context.Context, id uuid.UUID) (model.Order, error) {
	order, ok := repo.cache.Get(id)
	if ok {
		repo.cache.Set(id, order, 5*time.Minute)
		return order, nil
	}

	err := repo.db.Get(ctx, &order, querySelectOrderByID, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Order{}, ErrNotFound
		}
		return model.Order{}, fmt.Errorf("repo.db.Get: %w", err)
	}

	repo.cache.Set(id, order, 5*time.Minute)
	return order, nil
}
