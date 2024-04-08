package model

import "github.com/google/uuid"

// PVZ = Пункт Выдачи Заказов
type PVZ struct {
	ID       uuid.UUID `json:"id"`       // the id of PVZ
	Name     string    `json:"name"`     // the name of PVZ
	Address  string    `json:"address"`  // the address of PVZ
	Contacts string    `json:"contacts"` // the contacts of PVZ
}
