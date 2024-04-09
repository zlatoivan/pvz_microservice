//go:build integration
// +build integration

package order

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handlers/delivery"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/repo/postgres"
)

var (
	client http.Client
)

const queryInsertOrder = `
INSERT INTO orders (client_id, weight, cost, stores_till, give_out_time, packaging_type)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;`

// Создать тестовый заказ, чтоб взять его ID
func dbCreateOrder(t *testing.T, ctx context.Context, db postgres.Database, order delivery.RequestOrder) uuid.UUID {
	var id uuid.UUID
	var timeNull time.Time
	err := db.QueryRow(ctx, queryInsertOrder, order.ClientID, order.Weight, order.Cost, order.StoresTill, timeNull, order.PackagingType).Scan(&id)
	require.NoError(t, err)
	return id
}

const queryDeleteOrder = `
DELETE FROM orders
WHERE id = $1;`

func dbDeleteOrder(t *testing.T, ctx context.Context, db postgres.Database, id uuid.UUID) {
	_, err := db.Exec(ctx, queryDeleteOrder, id)
	require.NoError(t, err)
}

const querySelectOrderByID = `
SELECT id, client_id, weight, cost, stores_till, give_out_time, is_returned, packaging_type, is_deleted
FROM orders
WHERE id = $1`

func dbGetByIDOrder(t *testing.T, ctx context.Context, db postgres.Database, id uuid.UUID) model.Order {
	var order model.Order
	err := db.Get(ctx, &order, querySelectOrderByID, id)
	require.NoError(t, err)
	return order
}

func getOrderIDFromRespOrder(t *testing.T, res *http.Response) uuid.UUID {
	defer func() {
		err := res.Body.Close()
		require.NoError(t, err)
	}()
	var respID delivery.ResponseID
	err := json.NewDecoder(res.Body).Decode(&respID)
	require.NoError(t, err)
	return respID.ID
}

func getResp(t *testing.T, res *http.Response, respType string) (int, interface{}) {
	defer func() {
		err := res.Body.Close()
		require.NoError(t, err)
	}()

	var respJSON interface{}

	switch res.StatusCode {
	case http.StatusCreated:
		respID := delivery.ResponseID{}
		err := json.NewDecoder(res.Body).Decode(&respID)
		require.NoError(t, err)
		respJSON = respID

	case http.StatusOK:
		switch respType {
		case "ID":
			respID := delivery.ResponseID{}
			err := json.NewDecoder(res.Body).Decode(&respID)
			require.NoError(t, err)
			respJSON = respID
		case "Comment":
			var respComment delivery.ResponseComment
			err := json.NewDecoder(res.Body).Decode(&respComment)
			require.NoError(t, err)
			respJSON = respComment
		case "Order":
			var respOrder delivery.ResponseOrder
			err := json.NewDecoder(res.Body).Decode(&respOrder)
			require.NoError(t, err)
			respJSON = respOrder
		case "ListOrders":
			var respOrders []delivery.ResponseOrder
			err := json.NewDecoder(res.Body).Decode(&respOrders)
			require.NoError(t, err)
			respJSON = respOrders
		}

	default:
		var respErr delivery.ResponseError
		err := json.NewDecoder(res.Body).Decode(&respErr)
		require.NoError(t, err)
		respJSON = respErr
	}

	return res.StatusCode, respJSON
}
