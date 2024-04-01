package delivery

import (
	"time"

	"github.com/google/uuid"
)

type RequestOrder struct {
	ClientID      uuid.UUID `json:"client_id"`      // the id of client
	StoresTill    string    `json:"stores_till"`    // the storage period of order
	Weight        int       `json:"weight"`         // order weight
	Cost          int       `json:"cost"`           // order cost
	PackagingType string    `json:"packaging_type"` // packaging of the order
}

type RequestClientOrders struct {
	IDs []uuid.UUID `json:"ids"`
}

type RequestOrderID struct {
	ID uuid.UUID `json:"id"`
}

type ResponseID struct {
	ID uuid.UUID
}

type ResponseComment struct {
	Comment string
}

type ResponseOrder struct {
	ID          uuid.UUID `json:"id"`            // the id of order
	ClientID    uuid.UUID `json:"client_id"`     // the id of client
	Weight      int       `json:"weight"`        // order weight
	Cost        int       `json:"cost"`          // order cost
	StoresTill  time.Time `json:"stores_till"`   // the storage period of order
	GiveOutTime time.Time `json:"give_out_time"` // date and time of order give out
	IsReturned  bool      `json:"is_returned"`   // marker, the order has been returned or not
}
