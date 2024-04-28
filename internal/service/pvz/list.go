package pvz

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

func (s Service) ListPVZs(ctx context.Context) ([]model.PVZ, error) {
	list, err := s.repo.ListPVZs(ctx)
	if err != nil {
		return nil, fmt.Errorf("s.ListPVZs: %w", err)
	}
	return list, nil
}
