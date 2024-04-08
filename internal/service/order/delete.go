package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s Service) DeleteOrder(ctx context.Context, id uuid.UUID) error {
	err := s.repo.DeleteOrder(ctx, id)
	if err != nil {
		return fmt.Errorf("s.repo.DeleteOrder: %w", err)
	}
	return nil
}
