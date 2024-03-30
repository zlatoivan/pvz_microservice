package pvz

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

type repo interface {
	CreatePVZ(ctx context.Context, pvz model.PVZ) (uuid.UUID, error)
	ListPVZs(ctx context.Context) ([]model.PVZ, error)
	GetPVZByID(ctx context.Context, id uuid.UUID) (model.PVZ, error)
	UpdatePVZ(ctx context.Context, updPVZ model.PVZ) error
	DeletePVZ(ctx context.Context, id uuid.UUID) error
}

type Service struct {
	repo repo
}

func New(repo repo) Service {
	return Service{repo: repo}
}

func (s Service) CreatePVZ(ctx context.Context, pvz model.PVZ) (uuid.UUID, error) {
	id, err := s.repo.CreatePVZ(ctx, pvz)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("[service] s.CreatePVZ: %w", err)
	}
	return id, nil
}

func (s Service) ListPVZs(ctx context.Context) ([]model.PVZ, error) {
	list, err := s.repo.ListPVZs(ctx)
	if err != nil {
		return nil, fmt.Errorf("[service] s.ListPVZs: %w", err)
	}
	return list, nil
}

func (s Service) GetPVZByID(ctx context.Context, id uuid.UUID) (model.PVZ, error) {
	pvz, err := s.repo.GetPVZByID(ctx, id)
	if err != nil {
		return model.PVZ{}, fmt.Errorf("[service] s.GetPVZByID: %w", err)
	}
	return pvz, nil
}

func (s Service) UpdatePVZ(ctx context.Context, updPVZ model.PVZ) error {
	err := s.repo.UpdatePVZ(ctx, updPVZ)
	if err != nil {
		return fmt.Errorf("[service] s.UpdatePVZ: %w", err)
	}
	return nil
}

func (s Service) DeletePVZ(ctx context.Context, id uuid.UUID) error {
	err := s.repo.DeletePVZ(ctx, id)
	if err != nil {
		return fmt.Errorf("[service] s.DeletePVZ: %w", err)
	}
	return nil
}
