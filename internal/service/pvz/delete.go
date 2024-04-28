package pvz

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s Service) DeletePVZ(ctx context.Context, id uuid.UUID) error {
	err := s.repo.DeletePVZ(ctx, id)
	if err != nil {
		return fmt.Errorf("s.DeletePVZ: %w", err)
	}
	return nil
}
