package model

import "github.com/google/uuid"

type PVZ struct {
	ID       uuid.UUID `db:"id"`       // the id of PVZ
	Name     string    `db:"name"`     // the name of PVZ
	Address  string    `db:"address"`  // the address of PVZ
	Contacts string    `db:"contacts"` // the contacts of PVZ
}
