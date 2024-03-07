package model

import "time"

type Order struct {
	ID          int       // the id of order
	ClientID    int       // the id of client
	StoresTill  time.Time // the storage period of order
	IsDeleted   bool      // marker, the order has been deleted or not
	IsGivenOut  bool      // marker, the order has been given out or not
	GiveOutTime time.Time // date and time of order give out
	IsReturned  bool      // marker, the order has been returned or not
}
