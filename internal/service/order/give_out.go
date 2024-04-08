package order

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (s Service) GiveOutOrders(ctx context.Context, clientID uuid.UUID, ids []uuid.UUID) error {
	for _, id := range ids {
		// Проверка того, что заказ найден в хранилище
		order, err := s.repo.GetOrderByID(ctx, id)
		if err != nil {
			return fmt.Errorf("s.repo.GetOrderByID: %w. ID = %s", err, id)
		}

		// Проверка того, что заказ еще не выдан
		if !order.GiveOutTime.IsZero() {
			return fmt.Errorf("this order is already given out. ID = %s", id)
		}

		// Проверка того, что срок хранения заказа не истек
		if order.StoresTill.Before(time.Now()) {
			return fmt.Errorf("the stores period of the order has expired. ID = %s", order.ID)
		}
	}

	for _, id := range ids {
		err := s.repo.GiveOutOrder(ctx, clientID, id)
		if err != nil {
			return fmt.Errorf("s.repo.GiveOutOrder: %w. ID = %s", err, id)
		}
	}
	return nil
}
