package order

import (
	"fmt"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

type tape struct{}

func (_ tape) apply(order model.Order) (model.Order, error) {
	order.Cost = order.Cost + 1
	return order, nil
}

type pack struct{}

func (_ pack) apply(order model.Order) (model.Order, error) {
	if order.Weight > 10 {
		return model.Order{}, fmt.Errorf("the weight of the order is too large for such packaging")
	}
	order.Cost = order.Cost + 5
	return order, nil
}

type box struct{}

func (_ box) apply(order model.Order) (model.Order, error) {
	if order.Weight > 30 {
		return model.Order{}, fmt.Errorf("the weight of the order is too large for such packaging")
	}
	order.Cost = order.Cost + 20
	return order, nil
}

type applyer interface {
	apply(order model.Order) (model.Order, error)
}

func applyPackaging(order model.Order) (model.Order, error) {
	var newApplyer applyer
	switch order.PackagingType {
	case "box":
		newApplyer = box{}
	case "pack":
		newApplyer = pack{}
	case "tape":
		newApplyer = tape{}
	}

	newOrder, err := newApplyer.apply(order)
	if err != nil {
		return model.Order{}, fmt.Errorf("newApplyer.apply: %w", err)
	}

	return newOrder, nil
}
