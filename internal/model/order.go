package model

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID          uuid.UUID `db:"id"`            // the id of order
	ClientID    uuid.UUID `db:"client_id"`     // the id of client
	StoresTill  time.Time `db:"stores_till"`   // the storage period of order
	GiveOutTime time.Time `db:"give_out_time"` // date and time of order give out
	IsReturned  bool      `db:"is_returned"`   // marker, the order has been returned or not
	IsDeleted   bool      `db:"is_deleted"`    // marker, the order has been deleted or not
}
