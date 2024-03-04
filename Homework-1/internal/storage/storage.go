package storage

import (
	"encoding/json"
	"errors"
	"os"
	"route_256/homework/Homework-1/internal/model"
	"strconv"
	"time"
)

type Storage struct {
	storage *os.File
}

const storagePath = "internal/storage/storage.txt"

func New() (Storage, error) {
	file, err := os.OpenFile(storagePath, os.O_CREATE, 0777)
	if err != nil {
		return Storage{}, err
	}
	return Storage{storage: file}, nil
}

func (s *Storage) listAll() ([]model.Order, error) {
	bytes, err := os.ReadFile(storagePath)
	if err != nil {
		return nil, err
	}

	var orders []model.Order
	if len(bytes) == 0 {
		return orders, nil
	}
	err = json.Unmarshal(bytes, &orders)
	if err != nil {
		return orders, err
	}

	return orders, nil
}

func rewriteStorageFile(all []model.Order) error {
	bytes, err := json.MarshalIndent(all, "", "\t")
	if err != nil {
		return err
	}

	err = os.WriteFile(storagePath, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) Create(order model.Order) error {
	all, err := s.listAll()
	if err != nil {
		return err
	}

	// Проверка, что заказ уже принят
	for _, v := range all {
		if v.Id == order.Id {
			if v.IsDeleted == true {
				return errors.New("данный заказ уже удален")
			}
			return errors.New("данный заказ уже принят")
		}
	}

	// Проверка, что срок хранения заказа не истек
	if order.ShelfLife.Before(time.Now()) {
		return errors.New("срок хранения данного заказа истек")
	}

	all = append(all, order)
	err = rewriteStorageFile(all)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) Delete(id int) error {
	all, err := s.listAll()
	if err != nil {
		return err
	}

	ok := false
	for i := range all {
		if all[i].Id == id {
			if all[i].IsDeleted == true {
				return errors.New("данный заказ уже удален")
			}
			if all[i].IsGaveOut == true {
				return errors.New("данный заказ выдан клиенту")
			}
			if time.Now().Before(all[i].ShelfLife) {
				return errors.New("у данного заказа не истек срок хранения")
			}
			all[i].IsDeleted = true
			ok = true
		}
	}

	if !ok {
		return errors.New("данный заказ не найден в хранилище")
	}

	err = rewriteStorageFile(all)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GiveOut(ids []int) error {
	all, err := s.listAll()
	if err != nil {
		return err
	}

	clientId := -1
	for _, id := range ids {
		ok := false
		for _, v := range all {
			if v.Id == id {
				ok = true
				// Проверка того, что срок заказа не истек
				if v.ShelfLife.Before(time.Now()) {
					return errors.New("срок заказа " + strconv.Itoa(v.Id) + " истек")
				}
				// Проверка того, что все заказы принадлежат одному клиенту
				if clientId == -1 {
					clientId = v.ClientId
				} else if clientId != v.ClientId {
					return errors.New("не все заказы принадлежат одному клиенту")
				}
			}
		}
		// Проверка того, что все заказы найдены в хранилище
		if !ok {
			return errors.New("заказы не найдены в хранилище")
		}
	}

	for i, a := range all {
		for _, id := range ids {
			if a.Id == id {
				all[i].IsGaveOut = true
				all[i].GiveOutTime = time.Now()
			}
		}
	}

	err = rewriteStorageFile(all)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) List(id int, lastn int, inpvz bool) ([]int, error) {
	all, err := s.listAll()
	if err != nil {
		return nil, err
	}

	list := make([]int, 0)
	for i := len(all) - 1; i >= 0; i-- {
		// Выбрать только последние n заказов клиента, если есть такое уточнение
		if lastn != -1 && len(list) == lastn {
			break
		}
		v := all[i]
		if v.ClientId == id {
			// Выбрать только те заказы, которые находятся в нашем ПВЗ, если есть такое уточнение
			if inpvz && (v.IsGaveOut || v.IsDeleted) {
				continue
			}
			list = append(list, v.Id)
		}
	}

	return list, nil
}

func (s *Storage) Return(id int, clientId int) error {
	all, err := s.listAll()
	if err != nil {
		return err
	}

	ok := false
	for i, v := range all {
		if v.Id == id && v.ClientId == clientId {
			if v.IsReturned {
				return errors.New("данный заказ уже возвращен")
			}
			// Проверка, что заказ был выдан с нашего ПВЗ
			if v.IsGaveOut == false {
				return errors.New("данный заказ не был выдан клиенту")
			}
			// Проверка, что заказ возвращен в течение двух дней с момента выдачи
			today := time.Now()
			daysBetween := today.Sub(v.GiveOutTime).Hours() / 24
			if daysBetween > 2 {
				return errors.New("срок хранения данного заказа менее двух дней")
			}
			all[i].IsReturned = true
			all[i].IsGaveOut = false
			ok = true
		}
	}

	if !ok {
		return errors.New("данный заказ не найден в хранилище")
	}

	err = rewriteStorageFile(all)
	if err == nil {
		return err
	}

	return nil
}

func (s *Storage) ListOfReturned(pagenum int, itemsonpage int) ([]int, error) {
	all, err := s.listAll()
	if err != nil {
		return nil, err
	}

	list := make([]int, 0)
	start := (pagenum - 1) * itemsonpage
	for i := start; i < start+itemsonpage; i++ {
		if i == len(all) {
			break
		}
		if all[i].IsDeleted {
			continue
		}
		list = append(list, all[i].Id)
	}

	return list, nil
}
