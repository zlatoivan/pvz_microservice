package pvz

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

func (s Service) CreatePVZ(ctx context.Context, pvz model.PVZ) (uuid.UUID, error) {
	id, err := s.repo.CreatePVZ(ctx, pvz)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("s.CreatePVZ: %w", err)
	}
	return id, nil
}
