package order

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

const queryInsertOrder = `
INSERT INTO orders (client_id, weight, cost, stores_till, give_out_time, packaging_type)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;`

// CreateOrder creates Order in repo
func (repo Repo) CreateOrder(ctx context.Context, order model.Order) (uuid.UUID, error) {
	var id uuid.UUID
	var timeNull time.Time
	err := repo.db.QueryRow(ctx, queryInsertOrder, order.ClientID, order.Weight, order.Cost, order.StoresTill, timeNull, order.PackagingType).Scan(&id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("repo.db.QueryRow().Scan: %w", err)
	}

	return id, nil
}
