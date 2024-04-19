package fixtures

import (
	"time"

	"github.com/google/uuid"
)

var (
	ID, _  = uuid.Parse("64c6b9a4-b872-4d04-a1d2-6d072c7a4e2d")
	ID2, _ = uuid.Parse("74c6b9a4-b872-4d04-a1d2-6d072c7a4e2e")

	Name     = "Ozon Tech"
	Address  = "Moscow, Presnenskaya nab. 10, block С"
	Contacts = "+7 958 400-00-05, add 76077"

	ClientID, _   = uuid.Parse("88cda6c0-36fc-4be4-b976-e11a8a7a8f7e")
	Weight        = 29
	Cost          = 1100
	StoresTill, _ = time.Parse(time.RFC3339, "2024-04-22T12:12:00Z")
	PackagingType = "box"
	IsReturned    = false
	IsDeleted     = false
)
