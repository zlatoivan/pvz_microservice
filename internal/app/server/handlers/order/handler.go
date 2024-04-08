//go:generate minimock -i service -o mock/service_mock.go -p mock -g

package order

import (
	"context"

	"github.com/google/uuid"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

type Service interface {
	CreateOrder(ctx context.Context, order model.Order) (uuid.UUID, error)
	ListOrders(ctx context.Context) ([]model.Order, error)
	GetOrderByID(ctx context.Context, id uuid.UUID) (model.Order, error)
	UpdateOrder(ctx context.Context, updPVZ model.Order) error
	DeleteOrder(ctx context.Context, id uuid.UUID) error
	ListClientOrders(ctx context.Context, id uuid.UUID) ([]model.Order, error)
	GiveOutOrders(ctx context.Context, id uuid.UUID, ids []uuid.UUID) error
	ReturnOrder(ctx context.Context, clientID uuid.UUID, id uuid.UUID) error
	ListReturnedOrders(ctx context.Context) ([]model.Order, error)
}

type Handler struct {
	service Service
}

func New(service Service) Handler {
	return Handler{service: service}
}
