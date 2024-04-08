package order

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

func (s Service) UpdateOrder(ctx context.Context, updOrder model.Order) error {
	err := s.repo.UpdateOrder(ctx, updOrder)
	if err != nil {
		return fmt.Errorf("s.repo.UpdateOrder: %w", err)
	}
	return nil
}
