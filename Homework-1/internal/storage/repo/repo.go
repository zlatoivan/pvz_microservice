package repo

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"sync"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

type postgres interface {
	GetPool(_ context.Context) *pgxpool.Pool
	QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
	Query(ctx context.Context, query string) (pgx.Rows, error)
	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type PVZRepo struct {
	db postgres
	mu sync.Mutex
}

func NewRepo(database postgres) (*PVZRepo, error) {
	return &PVZRepo{db: database}, nil
}

// CreatePVZ creates PVZ in postgres
func (repo *PVZRepo) CreatePVZ(ctx context.Context, pvz model.PVZ) (int64, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	var id int64
	query := `INSERT INTO pvz (name, address, contacts) VALUES ($1, $2, $3) RETURNING id;`
	err := repo.db.QueryRow(ctx, query, pvz.Name, pvz.Address, pvz.Contacts).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("repo.db.QueryRow().Scan: %w", err)
	}

	return id, nil
}

// GetListOfPVZ gets list of PVZ from postgres
func (repo *PVZRepo) GetListOfPVZ(ctx context.Context) ([]model.PVZ, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	var pvzs []model.PVZ
	query := `SELECT id, name, address, contacts FROM pvz;`
	err := repo.db.Select(ctx, &pvzs, query)
	if err != nil {
		return nil, fmt.Errorf("repo.db.Select: %w", err)
	}

	return pvzs, nil
}

// GetPVZByID gets PVZ by ID from postgres
func (repo *PVZRepo) GetPVZByID(ctx context.Context, id int) (model.PVZ, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	query := `SELECT id, name, address, contacts FROM pvz WHERE id = $1;`
	var pvz model.PVZ
	err := repo.db.Get(ctx, &pvz, query, id)
	if err != nil {
		return model.PVZ{}, fmt.Errorf("repo.db.Get: %w", err)
	}

	return pvz, nil
}

// UpdatePVZ updates PVZ in postgres
func (repo *PVZRepo) UpdatePVZ(ctx context.Context, id int, updPVZ model.PVZ) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	query := `UPDATE pvz SET name = $2, address = $3, contacts = $4 WHERE id = $1;`
	t, err := repo.db.Exec(ctx, query, id, updPVZ.Name, updPVZ.Address, updPVZ.Contacts)
	if err != nil {
		return fmt.Errorf("repo.db.Exec: %w", err)
	}
	if t.RowsAffected() == 0 {
		return fmt.Errorf("no PVZ with this id")
	}

	return nil
}

// DeletePVZ deletes PVZ from postgres
func (repo *PVZRepo) DeletePVZ(ctx context.Context, id int) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	query := `DELETE FROM pvz WHERE id = $1;`
	t, err := repo.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("repo.db.Exec: %w", err)
	}
	if t.RowsAffected() == 0 {
		return fmt.Errorf("no PVZ with this id")
	}

	return nil
}
