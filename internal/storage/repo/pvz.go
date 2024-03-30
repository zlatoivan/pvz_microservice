package repo

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

const queryInsertPVZ = `INSERT INTO pvz (name, address, contacts) VALUES ($1, $2, $3) RETURNING id;`
const queryCheckInsertPVZ = `SELECT COUNT(*) FROM pvz WHERE id = $1;`

// CreatePVZ creates PVZ in repo
func (repo Repo) CreatePVZ(ctx context.Context, pvz model.PVZ) (uuid.UUID, error) {
	var id uuid.UUID
	err := repo.db.QueryRow(ctx, queryInsertPVZ, pvz.Name, pvz.Address, pvz.Contacts).Scan(&id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("repo.db.QueryRow().Scan: %w", err)
	}

	t, err := repo.db.Exec(ctx, queryCheckInsertPVZ, id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("repo.db.Exec: %w", err)
	}
	if t.RowsAffected() == 0 {
		return uuid.UUID{}, ErrorAlreadyExists
	}

	return id, nil
}

const querySelectPVZ = `SELECT id, name, address, contacts FROM pvz;`

// ListPVZs gets list of PVZ from repo
func (repo Repo) ListPVZs(ctx context.Context) ([]model.PVZ, error) {
	var pvzs []model.PVZ
	err := repo.db.Select(ctx, &pvzs, querySelectPVZ)
	if err != nil {
		return nil, fmt.Errorf("repo.db.Select: %w", err)
	}

	return pvzs, nil
}

const querySelectPVZByID = `SELECT id, name, address, contacts FROM pvz WHERE id = $1;`

// GetPVZByID gets PVZ by ID from repo
func (repo Repo) GetPVZByID(ctx context.Context, id uuid.UUID) (model.PVZ, error) {
	var pvz model.PVZ
	err := repo.db.Get(ctx, &pvz, querySelectPVZByID, id)
	if err != nil {
		return model.PVZ{}, ErrorNotFound
	}

	return pvz, nil
}

const queryUpdatePVZ = `UPDATE pvz SET name = $2, address = $3, contacts = $4 WHERE id = $1;`

// UpdatePVZ updates PVZ in repo
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

// DeletePVZ deletes PVZ from repo
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
