package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

type Storage struct {
	storage *os.File
	mutex   sync.Mutex
	pvzs    map[string]model.PVZ
}

const storagePath = "db/db.txt"
const PVZStoragePath = "db/pvz_db.json"

func New() (*Storage, error) {
	// For hw-1
	file, err := os.OpenFile(storagePath, os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	store := &Storage{
		storage: file,
		pvzs:    make(map[string]model.PVZ),
	}

	// For hw-2
	file, err = os.OpenFile(PVZStoragePath, os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	bytes, err := os.ReadFile(PVZStoragePath)
	if err != nil {
		return nil, err
	}
	if len(bytes) == 0 {
		return store, nil
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&store.pvzs)
	return store, err
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

	return orders, err
}

func rewriteStorageFile(all []model.Order) error {
	bytes, err := json.MarshalIndent(all, "", "\t")
	if err != nil {
		return err
	}

	err = os.WriteFile(storagePath, bytes, 0644)

	return err
}

func (s *Storage) Create(order model.Order) error {
	all, err := s.listAll()
	if err != nil {
		return err
	}

	// Проверка, что заказ уже принят
	for _, v := range all {
		if v.ID == order.ID {
			if v.IsDeleted {
				return fmt.Errorf("this order has already been deleted")
			}
			return fmt.Errorf("this order has already been accepted")
		}
	}

	// Проверка, что срок хранения заказа не истек
	if order.StoresTill.Before(time.Now()) {
		return fmt.Errorf("the shelf life of this order has expired")
	}

	all = append(all, order)
	err = rewriteStorageFile(all)

	return err
}

func (s *Storage) Delete(id int) error {
	all, err := s.listAll()
	if err != nil {
		return err
	}

	ok := false
	for i := range all {
		if all[i].ID == id {
			if all[i].IsDeleted {
				return fmt.Errorf("this order has already been deleted")
			}
			if !all[i].GiveOutTime.IsZero() {
				return fmt.Errorf("this order has been given out to the client")
			}
			if time.Now().Before(all[i].StoresTill) {
				return fmt.Errorf("the storage period of this order has not expired")
			}
			all[i].IsDeleted = true
			ok = true
		}
	}

	if !ok {
		return fmt.Errorf("this order was not found in the storage")
	}

	err = rewriteStorageFile(all)

	return err
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
			if v.ID == id {
				ok = true
				// Проверка того, что срок заказа не истек
				if v.StoresTill.Before(time.Now()) {
					return fmt.Errorf("the storage period of order " + strconv.Itoa(v.ID) + " has expired")
				}
				// Проверка того, что все заказы принадлежат одному клиенту
				if clientId == -1 {
					clientId = v.ClientID
				} else if clientId != v.ClientID {
					return fmt.Errorf("not all orders belong to the same client")
				}
			}
		}
		// Проверка того, что все заказы найдены в хранилище
		if !ok {
			return fmt.Errorf("orders were not found in the storage")
		}
	}

	for i, a := range all {
		for _, id := range ids {
			if a.ID == id {
				all[i].GiveOutTime = time.Now()
			}
		}
	}

	err = rewriteStorageFile(all)

	return err
}

func (s *Storage) List(id int, lastN int, inPVZ bool) ([]int, error) {
	all, err := s.listAll()
	if err != nil {
		return nil, err
	}

	list := make([]int, 0)
	for i := len(all) - 1; i >= 0; i-- {
		// Выбрать только последние n заказов клиента, если есть такое уточнение
		if lastN != -1 && len(list) == lastN {
			break
		}
		v := all[i]
		if v.ClientID == id {
			// Выбрать только те заказы, которые находятся в нашем ПВЗ, если есть такое уточнение
			if inPVZ && (!v.GiveOutTime.IsZero() || v.IsDeleted) {
				continue
			}
			list = append(list, v.ID)
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
		if v.ID == id && v.ClientID == clientId {
			if v.IsReturned {
				return fmt.Errorf("this order has already been returned")
			}
			// Проверка, что заказ был выдан с нашего ПВЗ
			if v.GiveOutTime.IsZero() {
				return fmt.Errorf("this order has not been given out to the client")
			}
			// Проверка, что заказ возвращен в течение двух дней с момента выдачи
			today := time.Now()
			daysBetween := today.Sub(v.GiveOutTime).Hours() / 24
			if daysBetween > 2 {
				return fmt.Errorf("the storage period of this order is less than two days")
			}
			all[i].IsReturned = true
			all[i].GiveOutTime = time.Time{}
			ok = true
		}
	}

	if !ok {
		return fmt.Errorf("this order was not found in the storage")
	}

	err = rewriteStorageFile(all)

	return err
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
		if all[i].IsReturned && !all[i].IsDeleted {
			list = append(list, all[i].ID)
		}
	}

	return list, nil
}

func (s *Storage) CreatePVZ(pvz model.PVZ) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for k := range s.pvzs {
		if k == pvz.Title {
			return fmt.Errorf("PVZ with this title is already in the storage")
		}
	}
	s.pvzs[pvz.Title] = pvz

	file, err := os.Open(PVZStoragePath)
	if err != nil {
		return fmt.Errorf("s.store.CreatePVZ: %w", err)
	}
	defer file.Close()
	bytes, err := json.MarshalIndent(s.pvzs, "", "\t")
	if err != nil {
		return err
	}
	err = os.WriteFile(PVZStoragePath, bytes, 0644)
	return err
}

func (s *Storage) GetPVZ(title string) (model.PVZ, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.pvzs[title], nil
}
