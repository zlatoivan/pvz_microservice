package repo

import (
	"context"
	"fmt"
	"sync"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
	"gitlab.ozon.dev/zlatoivan4/homework/pkg/db/postgres"
)

type PVZRepo struct {
	db *postgres.Database
	mu sync.Mutex
}

func NewRepo(database *postgres.Database) (*PVZRepo, error) {
	return &PVZRepo{db: database}, nil
}

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

func (repo *PVZRepo) GetListOfPVZ(ctx context.Context) ([]model.PVZ, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	query := `SELECT id, name, address, contacts FROM pvz;`
	rows, err := repo.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("repo.db.Query: %w", err)
	}
	defer rows.Close()

	pvzs := make([]model.PVZ, 0)
	for rows.Next() {
		var pvz model.PVZ
		err = rows.Scan(&pvz.ID, &pvz.Name, &pvz.Address, &pvz.Contacts)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}
		pvzs = append(pvzs, pvz)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("rows.Err: %w", err)
	}

	return pvzs, nil
}

func (repo *PVZRepo) GetPVZByID(ctx context.Context, id int) (model.PVZ, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	query := `SELECT id, name, address, contacts FROM pvz WHERE id = $1;`
	var pvz model.PVZ
	err := repo.db.QueryRow(ctx, query, id).Scan(&pvz.ID, &pvz.Name, &pvz.Address, &pvz.Contacts)
	if err != nil {
		return model.PVZ{}, fmt.Errorf("repo.db.QueryRow: %w", err)
	}

	return pvz, nil
}

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
