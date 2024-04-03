package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/config"
)

type Database struct {
	pool *pgxpool.Pool
}

func generateDsn(cfg config.Config) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Pg.Host, cfg.Pg.Port, cfg.Pg.User, cfg.Pg.Password, cfg.Pg.DBname)
}

func New(ctx context.Context, cfg config.Config) (Database, error) {
	pool, err := ConnectToPostgres(ctx, generateDsn(cfg))
	if err != nil {
		return Database{}, fmt.Errorf("postgres.New: %w", err)
	}
	return Database{pool: pool}, nil
}

// ConnectToPostgres устанавливает соединение с БД Postgres.
func ConnectToPostgres(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.ParseConfig: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.ConnectConfig: %w", err)
	}

	if err = pingWithRetry(ctx, pool, 10, 5*time.Second); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	return pool, nil
}

// pinger вызова метода Ping (используется в pingWithRetry).
type pinger interface {
	Ping(ctx context.Context) error
}

func pingWithRetry(ctx context.Context, p pinger, retry int, retryDuration time.Duration) error {
	err := p.Ping(ctx)

	ticker := time.NewTicker(retryDuration)
	defer ticker.Stop()

	for err != nil && retry > 0 {
		fmt.Printf("[postgres]: ping error (retry=%d): %s", retry, err.Error())

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			retry--
		}

		err = p.Ping(ctx)
	}

	return err
}

func (db Database) GetPool(_ context.Context) *pgxpool.Pool {
	return db.pool
}

func (db Database) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	return db.pool.Exec(ctx, query, args...)
}

func (db Database) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	return db.pool.QueryRow(ctx, query, args...)
}

func (db Database) Query(ctx context.Context, query string) (pgx.Rows, error) {
	return db.pool.Query(ctx, query)
}

func (db Database) Get(ctx context.Context, dest any, query string, args ...any) error {
	return pgxscan.Get(ctx, db.pool, dest, query, args...)
}

func (db Database) Select(ctx context.Context, dest any, query string, args ...any) error {
	return pgxscan.Select(ctx, db.pool, dest, query, args...)
}
