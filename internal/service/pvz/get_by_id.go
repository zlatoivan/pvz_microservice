package pvz

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

func (s Service) GetPVZByID(ctx context.Context, id uuid.UUID) (model.PVZ, error) {
	pvz, err := s.repo.GetPVZByID(ctx, id)
	if err != nil {
		return model.PVZ{}, fmt.Errorf("[service] s.GetPVZByID: %w", err)
	}
	return pvz, nil
}
