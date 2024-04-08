//go:generate minimock -i service -o mock/service_mock.go -p mock -g

package pvz

import (
	"context"

	"github.com/google/uuid"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

type Service interface {
	CreatePVZ(ctx context.Context, pvz model.PVZ) (uuid.UUID, error)
	ListPVZs(ctx context.Context) ([]model.PVZ, error)
	GetPVZByID(ctx context.Context, id uuid.UUID) (model.PVZ, error)
	UpdatePVZ(ctx context.Context, updPVZ model.PVZ) error
	DeletePVZ(ctx context.Context, id uuid.UUID) error
}

type Handler struct {
	service Service
}

func New(service Service) Handler {
	return Handler{service: service}
}
