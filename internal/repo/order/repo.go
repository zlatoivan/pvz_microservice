//go:generate minimock -i postgres -o mock/postgres_mock.go -p mock -g

package order

import (
	"context"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

type postgres interface {
	GetPool(_ context.Context) *pgxpool.Pool
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row
	Query(ctx context.Context, query string) (pgx.Rows, error)
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	Get(ctx context.Context, querier pgxscan.Querier, dest any, query string, args ...any) error
	Select(ctx context.Context, dest any, query string, args ...any) error
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

type inMemoryCache interface {
	Set(key uuid.UUID, value model.Order, ttl time.Duration)
	Get(key uuid.UUID) (model.Order, bool)
	Delete(key uuid.UUID)
}

type Repo struct {
	db    postgres
	cache inMemoryCache
}

func New(database postgres, cache inMemoryCache) Repo {
	return Repo{db: database, cache: cache}
}
