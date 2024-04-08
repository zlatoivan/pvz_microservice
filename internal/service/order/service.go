//go:generate minimock -i repo -o mock/repo_mock.go -p mock -g

package order

import (
	"context"

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
