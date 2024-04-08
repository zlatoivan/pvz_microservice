package pvz

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

func (s Service) UpdatePVZ(ctx context.Context, updPVZ model.PVZ) error {
	err := s.repo.UpdatePVZ(ctx, updPVZ)
	if err != nil {
		return fmt.Errorf("[service] s.UpdatePVZ: %w", err)
	}
	return nil
}
