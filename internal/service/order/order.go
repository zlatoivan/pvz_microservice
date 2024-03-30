package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

type repo interface {
	CreateOrder(ctx context.Context, order model.Order) (uuid.UUID, error)
	ListOrders(ctx context.Context) ([]model.Order, error)
	GetOrderByID(ctx context.Context, id uuid.UUID) (model.Order, error)
	UpdateOrder(ctx context.Context, updOrder model.Order) error
	DeleteOrder(ctx context.Context, id uuid.UUID) error
}

type Service struct {
	repo repo
}

func New(repo repo) Service {
	return Service{repo: repo}
}

func (s Service) CreateOrder(ctx context.Context, order model.Order) (uuid.UUID, error) {
	id, err := s.repo.CreateOrder(ctx, order)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("[service] s.CreateOrder: %w", err)
	}
	return id, nil
}

func (s Service) ListOrders(ctx context.Context) ([]model.Order, error) {
	list, err := s.repo.ListOrders(ctx)
	if err != nil {
		return nil, fmt.Errorf("[service] s.ListOrders: %w", err)
	}
	return list, nil
}

func (s Service) GetOrderByID(ctx context.Context, id uuid.UUID) (model.Order, error) {
	order, err := s.repo.GetOrderByID(ctx, id)
	if err != nil {
		return model.Order{}, fmt.Errorf("[service] s.GetOrderByID: %w", err)
	}
	return order, nil
}

func (s Service) UpdateOrder(ctx context.Context, updOrder model.Order) error {
	err := s.repo.UpdateOrder(ctx, updOrder)
	if err != nil {
		return fmt.Errorf("[service] s.UpdateOrder: %w", err)
	}
	return nil
}

func (s Service) DeleteOrder(ctx context.Context, id uuid.UUID) error {
	err := s.repo.DeleteOrder(ctx, id)
	if err != nil {
		return fmt.Errorf("[service] s.DeleteOrder: %w", err)
	}
	return nil
}
