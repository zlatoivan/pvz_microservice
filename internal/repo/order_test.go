package repo

import (
	"context"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/repo/mock"
	"gitlab.ozon.dev/zlatoivan4/homework/tests/fixtures"
)

func Test_CreateOrder(t *testing.T) {
	t.Parallel()
}

func Test_ListOrders(t *testing.T) {
	t.Parallel()
}

func Test_GetOrderByID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := fixtures.ID

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		//arrange
		mc := minimock.NewController(t)
		const querySelectOrderByID = `SELECT id, client_id, weight, cost, stores_till, give_out_time, is_returned FROM "order" WHERE id = $1 AND is_deleted = FALSE;`
		dbMock := mock.NewPostgresMock(mc).GetMock.Expect(ctx, &model.Order{}, querySelectOrderByID, id).Return(nil)
		repo := New(dbMock)

		// act
		order, err := repo.GetOrderByID(ctx, id)

		// assert
		require.NoError(t, err)
		assert.Equal(t, order.ID, id)
	})

	t.Run("order not found", func(t *testing.T) {
		t.Parallel()

		//arrange
		mc := minimock.NewController(t)
		const querySelectOrderByID = `SELECT id, client_id, weight, cost, stores_till, give_out_time, is_returned FROM "order" WHERE id = $1 AND is_deleted = FALSE;`
		pgMock := mock.NewPostgresMock(mc).GetMock.Expect(ctx, &model.Order{}, querySelectOrderByID, id).Return(ErrNotFound)
		repo := New(pgMock)

		// act
		order, err := repo.GetOrderByID(ctx, id)

		// assert
		require.Equal(t, "repo.db.Get: not found", err.Error())
		assert.Equal(t, order.ID, id)
	})
}
