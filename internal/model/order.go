package model

import (
	"time"

	"github.com/google/uuid"
)

// Order of the client
type Order struct {
	ID          uuid.UUID `json:"id"`            // the id of order
	ClientID    uuid.UUID `json:"client_id"`     // the id of client
	Weight      int       `json:"weight"`        // order weight
	Cost        int       `json:"cost"`          // order cost
	StoresTill  time.Time `json:"stores_till"`   // the storage period of order
	GiveOutTime time.Time `json:"give_out_time"` // date and time of order give out
	IsReturned  bool      `json:"is_returned"`   // marker, the order has been returned or not
	IsDeleted   bool      `json:"is_deleted"`    // marker, the order has been deleted or not
}
