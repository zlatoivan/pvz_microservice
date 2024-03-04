package service

import "route_256/homework/Homework-1/internal/model"

type storage interface {
	Create(order model.Order) error
	Delete(id int) error
	Giveout(ids []int) error
	List(id int, lastn int, inpvz bool) ([]int, error)
	Return(id int, clientId int) error
	ListOfReturned(pagenum int, itemsonpage int) ([]int, error)
}

type Service struct {
	stor storage
}

func New(s storage) Service {
	return Service{stor: s}
}

func (s *Service) Create(order model.Order) error {
	return s.stor.Create(order)
}

func (s *Service) Delete(id int) error {
	return s.stor.Delete(id)
}

func (s *Service) Giveout(ids []int) error {
	return s.stor.Giveout(ids)
}

func (s *Service) List(id int, lastn int, inpvz bool) ([]int, error) {
	return s.stor.List(id, lastn, inpvz)
}

func (s *Service) Return(id int, clientId int) error {
	return s.stor.Return(id, clientId)
}

func (s *Service) ListOfRuturned(pagenum int, itemsonpage int) ([]int, error) {
	return s.stor.ListOfReturned(pagenum, itemsonpage)
}
