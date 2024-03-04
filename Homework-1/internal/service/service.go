package service

import (
	"errors"
	"flag"
	"fmt"
	"route_256/homework/Homework-1/internal/model"
	"strconv"
	"strings"
	"time"
)

const helpConst = `Это утилита для управления ПВЗ.

Применение:
        go run cmd/main.go [flags] [command]

command:            Описание:                                flags:
        create            Принять заказ (создать).                 -id=1212 -clientid=9886 -shelflife=15.09.2024
        delete            Вернуть заказ курьеру (удалить).         -id=1212
        giveout           Выдать заказ клиенту.                    -ids=[1212,1214]
        list              Получить список заказов клиента.         -clientid=9886 -lastn=2 -inpvz=true  (последние два опциональные)
        return            Возврат заказа клиентом.                 -id=1212 -clientid=9886
        listofreturned    Получить список возвращенных заказов.    -pagenum=1 -itemsonpage=2
`

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

func (s *Service) Help() error {
	fmt.Print(helpConst)
	return nil
}

func (s *Service) Create(id int, clientId int, shelfLifeStr string) error {
	if id == -1 || clientId == -1 || shelfLifeStr == "-" {
		return errors.New("Неправильный формат входных данных")
	}

	// Привести срок хранения к типу даты
	datePattern := "02.01.2006"
	shelfLife, err := time.Parse(datePattern, shelfLifeStr)
	if err != nil {
		return errors.New("Неправльный формат даты")
	}

	newOrder := model.Order{
		Id:          id,
		ClientId:    clientId,
		ShelfLife:   shelfLife,
		IsDeleted:   false,
		IsGaveOut:   false,
		GiveOutTime: time.Date(1, 0, 0, 0, 0, 0, 0, time.UTC),
		IsReturned:  false,
	}
	err = s.stor.Create(newOrder)
	if err != nil {
		return errors.New("Ошибка создания заказа: " + err.Error())
	}

	fmt.Println("Заказ принят")

	return nil
}

func (s *Service) Delete(id int) error {
	if id == -1 {
		return errors.New("Неправильный формат входных данных")
	}

	err := s.stor.Delete(id)
	if err != nil {
		return errors.New("Ошибка удаления заказа: " + err.Error())
	}

	fmt.Println("Заказ удален")

	return nil
}

func (s *Service) Giveout(idsStr string) error {
	idsToSplit := (idsStr)[1 : len(idsStr)-1]
	idsToInt := strings.Split(idsToSplit, ",")
	ids := make([]int, len(idsToInt))
	for i := range idsToInt {
		idInt, err := strconv.Atoi(idsToInt[i])
		if err != nil {
			return errors.New("Неверные id заказов " + idsToInt[i])
		}
		ids[i] = idInt
	}

	err := s.stor.Giveout(ids)
	if err != nil {
		return errors.New("Ошибка выдачи заказов клиенту: " + err.Error())
	}

	fmt.Println("Заказы выданы клиенту")

	return nil
}

func (s *Service) List(clientId int, lastn int, inPvz bool) error {
	if clientId == -1 {
		return errors.New("Неправильный формат входных данных")
	}

	list, err := s.stor.List(clientId, lastn, inPvz)
	if err != nil {
		return errors.New("Ошибка получения списка заказов клиента: " + err.Error())
	}

	if len(list) == 0 {
		return errors.New("У данного пользователя нет заказов с такими параметрами")
	}

	fmt.Println("Заказы клиента:", list)

	return nil
}

func (s *Service) Return(id int, clientId int) error {
	if id == -1 || clientId == -1 {
		return errors.New("Неправильный формат входных данных")
	}

	err := s.stor.Return(id, clientId)
	if err != nil {
		return errors.New("Ошибка возврата заказа: " + err.Error())
	}

	fmt.Println("Возврат заказа принят")

	return nil
}

func (s *Service) ListOfReturned(pagenum int, itemsonpage int) error {
	if pagenum == -1 || itemsonpage == -1 {
		return errors.New("Неправильный формат входных данных")
	}

	list, err := s.stor.ListOfReturned(pagenum, itemsonpage)
	if err != nil {
		return errors.New("Ошибка получения списка возвращенных заказов: " + err.Error())
	}

	if len(list) == 0 {
		return errors.New("Возвращенных заказов нет")
	}

	fmt.Printf("Возвращенные заказы (стр.%d; %d заказа на стр.):\n%v\n", pagenum, itemsonpage, list)

	return nil
}

func (s *Service) Work() error {
	// Считывание флагов
	id := flag.Int("id", -1, "id of order")
	clientId := flag.Int("clientid", -1, "id of client")
	shelfLifeStr := flag.String("shelflife", "-", "shelf life of order")
	idsStr := flag.String("ids", "-", "ids of orders to give out")
	lastn := flag.Int("lastn", -1, "last n orders of client")
	inPvz := flag.Bool("inpvz", false, "client's orders that are in pvz")
	pagenum := flag.Int("pagenum", -1, "number of pages")
	itemsonpage := flag.Int("itemsonpage", -1, "number of items on page")
	flag.Parse()
	command := flag.Args()[len(flag.Args())-1]

	// Бизнес-логика
	switch command {
	case "help":
		return s.Help()
	case "create":
		return s.Create(*id, *clientId, *shelfLifeStr)
	case "delete":
		return s.Delete(*id)
	case "giveout":
		return s.Giveout(*idsStr)
	case "list":
		return s.List(*clientId, *lastn, *inPvz)
	case "return":
		return s.Return(*id, *clientId)
	case "listofreturned":
		return s.ListOfReturned(*pagenum, *itemsonpage)
	default:
		return errors.New("Неизвестная команда")
	}
}
