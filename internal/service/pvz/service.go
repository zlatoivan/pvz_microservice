//go:generate minimock -i repo -o mock/repo_mock.go -p mock -g

package pvz

import (
	"context"

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
