package order

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
	repo2 "gitlab.ozon.dev/zlatoivan4/homework/internal/repo/order"
)

func (s Service) CreateOrder(ctx context.Context, order model.Order) (uuid.UUID, error) {
	_, err := s.repo.GetOrderByID(ctx, order.ID)
	if !errors.Is(err, repo2.ErrNotFound) {
		if err == nil {
			return uuid.UUID{}, repo2.ErrAlreadyExists
		} else {
			return uuid.UUID{}, fmt.Errorf("s.repo.GetOrderByID: %w", err)
		}
	}

	if order.StoresTill.Before(time.Now()) {
		return uuid.UUID{}, fmt.Errorf("the stores period of the order has expired")
	}

	newOrder, err := ApplyPackaging(order)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("ApplyPackaging: %w", err)
	}

	id, err := s.repo.CreateOrder(ctx, newOrder)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("s.repo.CreateOrder: %w", err)
	}

	return id, nil
}
