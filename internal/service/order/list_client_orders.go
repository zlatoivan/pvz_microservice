package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

func (s Service) ListClientOrders(ctx context.Context, id uuid.UUID) ([]model.Order, error) {
	list, err := s.repo.ListClientOrders(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("s.repo.ListClientOrders: %w", err)
	}
	return list, nil
}
