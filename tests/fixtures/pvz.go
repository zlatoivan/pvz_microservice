package fixtures

import (
	"github.com/google/uuid"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

type PVZBuilder struct {
	instance *model.PVZ
}

func PVZ() *PVZBuilder {
	return &PVZBuilder{instance: &model.PVZ{}}
}

func (b *PVZBuilder) ID(id uuid.UUID) *PVZBuilder {
	b.instance.ID = id
	return b
}

func (b *PVZBuilder) Name(name string) *PVZBuilder {
	b.instance.Name = name
	return b
}

func (b *PVZBuilder) Address(address string) *PVZBuilder {
	b.instance.Address = address
	return b
}

func (b *PVZBuilder) Contacts(contacts string) *PVZBuilder {
	b.instance.Contacts = contacts
	return b
}

func (b *PVZBuilder) P() *model.PVZ {
	return b.instance
}

func (b *PVZBuilder) V() model.PVZ {
	return *b.instance
}

func (b *PVZBuilder) Valid() *PVZBuilder {
	return b.ID(ID).
		Name(Name).
		Address(Address).
		Contacts(Contacts)
}
