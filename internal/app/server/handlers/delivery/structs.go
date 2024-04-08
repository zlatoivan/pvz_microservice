package delivery

import (
	"time"

	"github.com/google/uuid"
)

type RequestPVZ struct {
	ID       uuid.UUID `json:"id"`       // the id of PVZ
	Name     string    `json:"name"`     // the name of PVZ
	Address  string    `json:"address"`  // the address of PVZ
	Contacts string    `json:"contacts"` // the contacts of PVZ
}

type RequestOrder struct {
	ID            uuid.UUID `json:"id"`             // the id of order
	ClientID      uuid.UUID `json:"client_id"`      // the id of client
	StoresTill    string    `json:"stores_till"`    // the storage period of order
	Weight        int       `json:"weight"`         // order weight
	Cost          int       `json:"cost"`           // order cost
	PackagingType string    `json:"packaging_type"` // packaging of the order
}

type RequestClientOrders struct {
	IDs []uuid.UUID `json:"ids"`
}

type RequestID struct {
	ID uuid.UUID `json:"id"`
}

type RequestGiveOut struct {
	ClientID uuid.UUID   `json:"client_id"`
	IDs      []uuid.UUID `json:"ids"`
}

type ResponseID struct {
	ID uuid.UUID
}

type ResponseOrder struct {
	ID            uuid.UUID `json:"id"`             // the id of order
	ClientID      uuid.UUID `json:"client_id"`      // the id of client
	Weight        int       `json:"weight"`         // order weight
	Cost          int       `json:"cost"`           // order cost
	StoresTill    time.Time `json:"stores_till"`    // the storage period of order
	GiveOutTime   time.Time `json:"give_out_time"`  // date and time of order give out
	IsReturned    bool      `json:"is_returned"`    // marker, the order has been returned or not
	IsDeleted     bool      `json:"is_deleted"`     // marker, the order has been deleted or not
	PackagingType string    `json:"packaging_type"` // the packaging type of the order
}

type ResponsePVZ struct {
	ID       uuid.UUID `json:"id"`       // the id of PVZ
	Name     string    `json:"name"`     // the name of PVZ
	Address  string    `json:"address"`  // the address of PVZ
	Contacts string    `json:"contacts"` // the contacts of PVZ
}

type ResponseError struct {
	HTTPStatusCode int    `json:"http_status_code"`
	StatusText     string `json:"status_text"`
	ErrorText      string `json:"error_text"`
}

type ResponseComment struct {
	Comment string `json:"comment,omitempty"`
}
