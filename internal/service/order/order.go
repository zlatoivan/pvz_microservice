package order

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

type repo interface {
	CreateOrder(ctx context.Context, order model.Order) (uuid.UUID, error)
	ListOrders(ctx context.Context) ([]model.Order, error)
	GetOrderByID(ctx context.Context, id uuid.UUID) (model.Order, error)
	UpdateOrder(ctx context.Context, updOrder model.Order) error
	DeleteOrder(ctx context.Context, id uuid.UUID) error
	ListClientOrders(ctx context.Context, id uuid.UUID) ([]model.Order, error)
	GiveOutOrder(ctx context.Context, clientID uuid.UUID, id uuid.UUID) error
	ReturnOrder(ctx context.Context, clientID uuid.UUID, id uuid.UUID) error
	ListReturnedOrders(ctx context.Context) ([]model.Order, error)
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
		return uuid.UUID{}, fmt.Errorf("s.repo.CreateOrder: %w", err)
	}
	return id, nil
}

func (s Service) ListOrders(ctx context.Context) ([]model.Order, error) {
	list, err := s.repo.ListOrders(ctx)
	if err != nil {
		return nil, fmt.Errorf("s.repo.ListOrders: %w", err)
	}
	return list, nil
}

func (s Service) GetOrderByID(ctx context.Context, id uuid.UUID) (model.Order, error) {
	order, err := s.repo.GetOrderByID(ctx, id)
	if err != nil {
		return model.Order{}, fmt.Errorf("s.repo.GetOrderByID: %w", err)
	}
	return order, nil
}

func (s Service) UpdateOrder(ctx context.Context, updOrder model.Order) error {
	err := s.repo.UpdateOrder(ctx, updOrder)
	if err != nil {
		return fmt.Errorf("s.repo.UpdateOrder: %w", err)
	}
	return nil
}

func (s Service) DeleteOrder(ctx context.Context, id uuid.UUID) error {
	err := s.repo.DeleteOrder(ctx, id)
	if err != nil {
		return fmt.Errorf("s.repo.DeleteOrder: %w", err)
	}
	return nil
}

func (s Service) ListClientOrders(ctx context.Context, id uuid.UUID) ([]model.Order, error) {
	list, err := s.repo.ListClientOrders(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("s.repo.ListClientOrders: %w", err)
	}
	return list, nil
}

func (s Service) GiveOutOrders(ctx context.Context, clientID uuid.UUID, ids []uuid.UUID) error {
	for _, id := range ids {
		// Проверка того, что все заказ найден в хранилище
		order, err := s.repo.GetOrderByID(ctx, id)
		if err != nil {
			return fmt.Errorf("s.repo.GetOrderByID: %w. ID = %s", err, id)
		}

		// Проверка того, что заказ еще не выдан
		if !order.GiveOutTime.IsZero() {
			return fmt.Errorf("this order is already given out. ID = %s", id)
		}

		// Проверка того, что срок хранения заказа не истек
		if order.StoresTill.Before(time.Now()) {
			return fmt.Errorf("the stores period of order has expired. ID = %s", order.ID)
		}
	}

	for _, id := range ids {
		err := s.repo.GiveOutOrder(ctx, clientID, id)
		if err != nil {
			return fmt.Errorf("s.repo.GiveOutOrder: %w. ID = %s", err, id)
		}
	}
	return nil
}

func (s Service) ReturnOrder(ctx context.Context, clientID uuid.UUID, id uuid.UUID) error {
	order, err := s.repo.GetOrderByID(ctx, id)
	if err != nil {
		return fmt.Errorf("s.repo.GetOrderByID: %w. ID = %s", err, id)
	}

	// Проверка того, что заказ был выдан с нашего ПВЗ
	if order.GiveOutTime.IsZero() {
		return fmt.Errorf("this order has not been given out to the client")
	}

	// Проверка, что заказ возвращен в течение двух дней с момента выдачи
	today := time.Now()
	daysBetween := today.Sub(order.GiveOutTime).Hours() / 24
	if daysBetween > 2 {
		return fmt.Errorf("the orders period of this order is less than two days")
	}

	err = s.repo.ReturnOrder(ctx, clientID, id)
	if err != nil {
		return fmt.Errorf("s.repo.ReturnOrder: %w", err)
	}
	return nil
}

func (s Service) ListReturnedOrders(ctx context.Context) ([]model.Order, error) {
	list, err := s.repo.ListReturnedOrders(ctx)
	if err != nil {
		return nil, fmt.Errorf("s.repo.ListReturnedOrders: %w", err)
	}
	return list, nil
}
