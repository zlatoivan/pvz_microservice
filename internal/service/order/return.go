package order

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (s Service) ReturnOrder(ctx context.Context, clientID uuid.UUID, id uuid.UUID) error {
	order, err := s.repo.GetOrderByID(ctx, id)
	if err != nil {
		return fmt.Errorf("s.repo.GetOrderByID: %w. ID = %s", err, id)
	}

	// Проверка того, что заказ был выдан с нашего ПВЗ
	if order.GiveOutTime.IsZero() {
		return fmt.Errorf("this order has not been given out to the client")
	}

	// Проверка, что заказ возвращен в течение двух дней с момента выдачи
	today := time.Now()
	daysBetween := today.Sub(order.GiveOutTime).Hours() / 24
	if daysBetween > 2 {
		return fmt.Errorf("the orders period of this order is less than two days")
	}

	err = s.repo.ReturnOrder(ctx, clientID, id)
	if err != nil {
		return fmt.Errorf("s.repo.ReturnOrder: %w", err)
	}
	return nil
}
