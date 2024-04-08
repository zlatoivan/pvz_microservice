package order

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

func (s Service) ListReturnedOrders(ctx context.Context) ([]model.Order, error) {
	list, err := s.repo.ListReturnedOrders(ctx)
	if err != nil {
		return nil, fmt.Errorf("s.repo.ListReturnedOrders: %w", err)
	}
	return list, nil
}
