package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

func (s Service) GetOrderByID(ctx context.Context, id uuid.UUID) (model.Order, error) {
	order, err := s.repo.GetOrderByID(ctx, id)
	if err != nil {
		return model.Order{}, fmt.Errorf("s.repo.GetOrderByID: %w", err)
	}
	return order, nil
}
