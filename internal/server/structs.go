package server

import (
	"github.com/google/uuid"
)

type ResponseID struct {
	ID uuid.UUID
}

type responseComment struct {
	Comment string
}

type requestOrder struct {
	ClientID   uuid.UUID `db:"client_id"`   // the id of client
	StoresTill string    `db:"stores_till"` // the storage period of order
}
