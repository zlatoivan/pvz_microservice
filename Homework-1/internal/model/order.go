package model

import "time"

type Order struct {
	Id          int
	ClientId    int
	ShelfLife   time.Time
	IsDeleted   bool
	IsGaveOut   bool
	GiveOutTime time.Time
	IsReturned  bool
}
