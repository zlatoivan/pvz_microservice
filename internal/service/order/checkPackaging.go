package order

import (
	"fmt"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

type tape struct{}

func (_ tape) check(order model.Order) (model.Order, error) {
	order.Cost = order.Cost + 1
	return order, nil
}

type pack struct{}

func (_ pack) check(order model.Order) (model.Order, error) {
	if order.Weight > 10 {
		return model.Order{}, fmt.Errorf("the weight of the order is too large for such packaging")
	}
	order.Cost = order.Cost + 5
	return order, nil
}

type box struct{}

func (_ box) check(order model.Order) (model.Order, error) {
	if order.Weight > 30 {
		return model.Order{}, fmt.Errorf("the weight of the order is too large for such packaging")
	}
	order.Cost = order.Cost + 20
	return order, nil
}

type checker interface {
	check(order model.Order) (model.Order, error)
}

func checkPackaging(packagingType string, order model.Order) (model.Order, error) {
	var newChecker checker
	switch packagingType {
	case "box":
		newChecker = box{}
	case "pack":
		newChecker = pack{}
	case "tape":
		newChecker = tape{}
	}

	newOrder, err := newChecker.check(order)
	if err != nil {
		return model.Order{}, fmt.Errorf("newChecker.check: %w", err)
	}

	return newOrder, nil
}
