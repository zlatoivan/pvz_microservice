package repo

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

const queryInsertOrder = `INSERT INTO "order" (client_id, stores_till, is_deleted, give_out_time, is_returned) VALUES ($1, $2, $3, $4, $5) RETURNING id;`
const queryCheckInsertOrder = `SELECT COUNT(*) FROM "order" WHERE id = $1;`

// CreateOrder creates Order in repo
func (repo Repo) CreateOrder(ctx context.Context, order model.Order) (uuid.UUID, error) {
	var id uuid.UUID
	err := repo.db.QueryRow(ctx, queryInsertOrder, order.ClientID, order.StoresTill, order.IsDeleted, order.GiveOutTime, order.IsReturned).Scan(&id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("repo.db.QueryRow().Scan: %w", err)
	}

	t, err := repo.db.Exec(ctx, queryCheckInsertOrder, id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("repo.db.Exec: %w", err)
	}
	if t.RowsAffected() == 0 {
		return uuid.UUID{}, ErrorAlreadyExists
	}

	return id, nil
}

const querySelectOrder = `SELECT id, client_id, stores_till, is_deleted, give_out_time, is_returned FROM "order";`

// ListOrders gets list of Order from repo
func (repo Repo) ListOrders(ctx context.Context) ([]model.Order, error) {
	var orders []model.Order
	err := repo.db.Select(ctx, &orders, querySelectOrder)
	if err != nil {
		return nil, fmt.Errorf("repo.db.Select: %w", err)
	}

	return orders, nil
}

const querySelectPBZByID = `SELECT id, client_id, stores_till, is_deleted, give_out_time, is_returned FROM "order" WHERE id = $1;`

// GetOrderByID gets Order by ID from repo
func (repo Repo) GetOrderByID(ctx context.Context, id uuid.UUID) (model.Order, error) {
	var order model.Order
	err := repo.db.Get(ctx, &order, querySelectPBZByID, id)
	if err != nil {
		return model.Order{}, ErrorNotFound
	}

	return order, nil
}

const queryUpdateOrder = `UPDATE "order" SET client_id = $2, stores_till = $3, is_deleted = $4, give_out_time = $5, is_returned = $6 WHERE id = $1;`

// UpdateOrder updates Order in repo
func (repo Repo) UpdateOrder(ctx context.Context, updOrder model.Order) error {
	t, err := repo.db.Exec(ctx, queryUpdateOrder, updOrder.ID, updOrder.ClientID, updOrder.StoresTill, updOrder.IsDeleted, updOrder.GiveOutTime, updOrder.IsReturned)
	if err != nil {
		return fmt.Errorf("repo.db.Exec: %w", err)
	}
	if t.RowsAffected() == 0 {
		return ErrorNotFound
	}

	return nil
}

const queryDeleteOrder = `DELETE FROM "order" WHERE id = $1;`

// DeleteOrder deletes Order from repo
func (repo Repo) DeleteOrder(ctx context.Context, id uuid.UUID) error {
	t, err := repo.db.Exec(ctx, queryDeleteOrder, id)
	if err != nil {
		return fmt.Errorf("repo.db.Exec: %w", err)
	}
	if t.RowsAffected() == 0 {
		return ErrorNotFound
	}

	return nil
}
