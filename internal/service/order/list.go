package order

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

func (s Service) ListOrders(ctx context.Context) ([]model.Order, error) {
	list, err := s.repo.ListOrders(ctx)
	if err != nil {
		return nil, fmt.Errorf("s.repo.ListOrders: %w", err)
	}
	return list, nil
}
