package postgres

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"gitlab.ozon.dev/zlatoivan4/homework/configs"
)

type Database struct {
	pool *pgxpool.Pool
}

func generateDsn(cfg configs.Config) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Pg.Host, cfg.Pg.Port, cfg.Pg.User, cfg.Pg.Password, cfg.Pg.DBname)
}

func NewDB(ctx context.Context, cfg configs.Config) (*Database, error) {
	pool, err := pgxpool.Connect(ctx, generateDsn(cfg))
	if err != nil {
		return nil, err
	}
	return &Database{pool: pool}, nil
}

func (db *Database) GetPool(_ context.Context) *pgxpool.Pool {
	return db.pool
}

func (db *Database) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	return db.pool.Exec(ctx, query, args...)
}

func (db *Database) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return db.pool.QueryRow(ctx, query, args...)
}

func (db *Database) Query(ctx context.Context, query string) (pgx.Rows, error) {
	return db.pool.Query(ctx, query)
}

func (db *Database) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return pgxscan.Get(ctx, db.pool, dest, query, args...)
}

func (db *Database) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return pgxscan.Select(ctx, db.pool, dest, query, args...)
}
