package postgres

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/delivery"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/config"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/repo/postgres"
)

func SetUp(ctx context.Context) (postgres.Database, error) {
	//cfg := config.Config{
	//	Pg: config.Pg{
	//		Host:     "localhost",
	//		Port:     "5431",
	//		DBname:   "test",
	//		User:     "postgres",
	//		Password: "postgres",
	//	},
	//}
	cfg, err := config.New()
	if err != nil {
		return postgres.Database{}, fmt.Errorf("config.New: %w", err)
	}

	db, err := postgres.New(ctx, cfg)
	if err != nil {
		return postgres.Database{}, fmt.Errorf("postgres.SeptUp: %w", err)
	}

	return db, nil
}

func TearDown(ctx context.Context, db postgres.Database) {
	db.GetPool(ctx).Close()
}

const queryInsertOrder = `
INSERT INTO orders (client_id, weight, cost, stores_till, give_out_time, packaging_type)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;`

// CreateOrder создает тестовый заказ, чтоб взять его ID
func CreateOrder(t *testing.T, ctx context.Context, db postgres.Database, order delivery.RequestOrder) uuid.UUID {
	var id uuid.UUID
	var timeNull time.Time
	err := db.QueryRow(ctx, queryInsertOrder, order.ClientID, order.Weight, order.Cost, order.StoresTill, timeNull, order.PackagingType).Scan(&id)
	require.NoError(t, err)
	return id
}

const queryDeleteOrder = `
DELETE FROM orders
WHERE id = $1;`

// DeleteOrder удаляет тестовый заказ
func DeleteOrder(t *testing.T, ctx context.Context, db postgres.Database, id uuid.UUID) {
	_, err := db.Exec(ctx, queryDeleteOrder, id)
	require.NoError(t, err)
}

const querySelectOrderByID = `
SELECT id, client_id, weight, cost, stores_till, give_out_time, is_returned, packaging_type, is_deleted
FROM orders
WHERE id = $1`

// GetByIDOrder получает заказ из базы данных
func GetByIDOrder(t *testing.T, ctx context.Context, db postgres.Database, id uuid.UUID) model.Order {
	var order model.Order
	err := db.Get(ctx, &order, querySelectOrderByID, id)
	require.NoError(t, err)
	return order
}
