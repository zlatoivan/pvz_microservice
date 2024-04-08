package fixtures

import (
	"time"

	"github.com/google/uuid"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

type OrderBuilder struct {
	instance *model.Order
}

func Order() *OrderBuilder {
	return &OrderBuilder{instance: &model.Order{}}
}

func (b *OrderBuilder) ID(id uuid.UUID) *OrderBuilder {
	b.instance.ID = id
	return b
}

func (b *OrderBuilder) ClientID(clientID uuid.UUID) *OrderBuilder {
	b.instance.ClientID = clientID
	return b
}

func (b *OrderBuilder) Weight(weight int) *OrderBuilder {
	b.instance.Weight = weight
	return b
}

func (b *OrderBuilder) Cost(cost int) *OrderBuilder {
	b.instance.Cost = cost
	return b
}

func (b *OrderBuilder) StoresTill(storesTill time.Time) *OrderBuilder {
	b.instance.StoresTill = storesTill
	return b
}

func (b *OrderBuilder) PackagingType(packagingType string) *OrderBuilder {
	b.instance.PackagingType = packagingType
	return b
}

func (b *OrderBuilder) IsReturned(isReturned bool) *OrderBuilder {
	b.instance.IsReturned = isReturned
	return b
}

func (b *OrderBuilder) IsDeleted(isDeleted bool) *OrderBuilder {
	b.instance.IsDeleted = isDeleted
	return b
}

func (b *OrderBuilder) P() *model.Order {
	return b.instance
}

func (b *OrderBuilder) V() model.Order {
	return *b.instance
}

func (b *OrderBuilder) Valid() *OrderBuilder {
	return b.ID(ID).
		ClientID(ClientID).
		Weight(Weight).
		Cost(Cost).
		StoresTill(StoresTill).
		PackagingType(PackagingType).
		IsReturned(IsReturned).
		IsDeleted(IsDeleted)
}
