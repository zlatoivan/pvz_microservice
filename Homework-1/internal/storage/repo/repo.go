package repo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

type postgres interface {
	GetPool(_ context.Context) *pgxpool.Pool
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row
	Query(ctx context.Context, query string) (pgx.Rows, error)
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	Get(ctx context.Context, dest any, query string, args ...any) error
	Select(ctx context.Context, dest any, query string, args ...any) error
}

type Repo struct {
	db postgres
}

func New(database postgres) Repo {
	return Repo{db: database}
}

const queryInsertPVZ = `INSERT INTO pvz (name, address, contacts) VALUES ($1, $2, $3) RETURNING id;`

// CreatePVZ creates PVZ in postgres
func (repo Repo) CreatePVZ(ctx context.Context, pvz model.PVZ) (uuid.UUID, error) {
	var id uuid.UUID
	err := repo.db.QueryRow(ctx, queryInsertPVZ, pvz.Name, pvz.Address, pvz.Contacts).Scan(&id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("repo.db.QueryRow().Scan: %w", err)
	}

	return id, nil
}

const querySelectPVZ = `SELECT id, name, address, contacts FROM pvz;`

// GetListOfPVZ gets list of PVZ from postgres
func (repo Repo) GetListOfPVZ(ctx context.Context) ([]model.PVZ, error) {
	var pvzs []model.PVZ
	err := repo.db.Select(ctx, &pvzs, querySelectPVZ)
	if err != nil {
		return nil, fmt.Errorf("repo.db.Select: %w", err)
	}

	return pvzs, nil
}

const querySelectPBZByID = `SELECT id, name, address, contacts FROM pvz WHERE id = $1;`

// GetPVZByID gets PVZ by ID from postgres
func (repo Repo) GetPVZByID(ctx context.Context, id uuid.UUID) (model.PVZ, error) {
	var pvz model.PVZ
	err := repo.db.Get(ctx, &pvz, querySelectPBZByID, id)
	if err != nil {
		return model.PVZ{}, ErrorNotFound
	}

	return pvz, nil
}

const queryUpdatePVZ = `UPDATE pvz SET name = $2, address = $3, contacts = $4 WHERE id = $1;`

// UpdatePVZ updates PVZ in postgres
func (repo Repo) UpdatePVZ(ctx context.Context, updPVZ model.PVZ) error {
	t, err := repo.db.Exec(ctx, queryUpdatePVZ, updPVZ.ID, updPVZ.Name, updPVZ.Address, updPVZ.Contacts)
	if err != nil {
		return fmt.Errorf("repo.db.Exec: %w", err)
	}
	if t.RowsAffected() == 0 {
		return ErrorNotFound
	}

	return nil
}

const queryDeletePVZ = `DELETE FROM pvz WHERE id = $1;`

// DeletePVZ deletes PVZ from postgres
func (repo Repo) DeletePVZ(ctx context.Context, id uuid.UUID) error {
	t, err := repo.db.Exec(ctx, queryDeletePVZ, id)
	if err != nil {
		return fmt.Errorf("repo.db.Exec: %w", err)
	}
	if t.RowsAffected() == 0 {
		return ErrorNotFound
	}

	return nil
}
