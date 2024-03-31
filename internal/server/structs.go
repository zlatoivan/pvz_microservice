package server

import (
	"time"

	"github.com/google/uuid"
)

type requestOrder struct {
	ClientID   uuid.UUID `json:"ClientID"`   // the id of client
	StoresTill string    `json:"StoresTill"` // the storage period of order
}

type requestClientOrders struct {
	IDs []uuid.UUID `json:"ids"`
}

type requestOrderID struct {
	ID uuid.UUID `json:"id"`
}

type responseID struct {
	ID uuid.UUID
}

type responseComment struct {
	Comment string
}

type responseOrder struct {
	ID          uuid.UUID `json:"id"`            // the id of order
	ClientID    uuid.UUID `json:"client_id"`     // the id of client
	StoresTill  time.Time `json:"stores_till"`   // the storage period of order
	GiveOutTime time.Time `json:"give_out_time"` // date and time of order give out
	IsReturned  bool      `json:"is_returned"`   // marker, the order has been returned or not
}
