package service

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
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
        listofreturned    Получить список возвращенных заказов.    -pagenum=1 -itemsonpage=2`

const dateLayoutConst = "02.01.2006"

type storage interface {
	Create(order model.Order) error
	Delete(id int) error
	GiveOut(ids []int) error
	List(id int, lastn int, inpvz bool) ([]int, error)
	Return(id int, clientId int) error
	ListOfReturned(pagenum int, itemsonpage int) ([]int, error)
}

type Service struct {
	store storage
}

func New(s storage) (Service, error) {
	return Service{store: s}, nil
}

func (s *Service) Help() error {
	fmt.Println(helpConst)
	fmt.Println()
	return nil
}

func (s *Service) Create(id int, clientId int, shelfLifeStr string) error {
	if id == -1 || clientId == -1 || shelfLifeStr == "-" {
		return fmt.Errorf("incorrect input data format")
	}

	// Привести срок хранения к типу даты
	shelfLife, err := time.Parse(dateLayoutConst, shelfLifeStr)
	if err != nil {
		return fmt.Errorf("incorrect date format")
	}

	newOrder := model.Order{
		ID:          id,
		ClientID:    clientId,
		StoresTill:  shelfLife,
		IsDeleted:   false,
		GiveOutTime: time.Time{}, // zero value
		IsReturned:  false,
	}
	err = s.store.Create(newOrder)
	if err != nil {
		return fmt.Errorf("s.store.Create: %w", err)
	}

	fmt.Println("The order is accepted")

	return nil
}

func (s *Service) Delete(id int) error {
	if id == -1 {
		return fmt.Errorf("incorrect input data format")
	}

	err := s.store.Delete(id)
	if err != nil {
		return fmt.Errorf("s.store.Delete: %w", err)
	}

	fmt.Println("The order has been deleted")

	return nil
}

func (s *Service) GiveOut(idsStr string) error {
	idsToSplit := idsStr[1 : len(idsStr)-1]
	idsToInt := strings.Split(idsToSplit, ",")
	ids := make([]int, len(idsToInt))
	for i := range idsToInt {
		idInt, err := strconv.Atoi(idsToInt[i])
		if err != nil {
			return fmt.Errorf("invalid order IDs " + idsToInt[i])
		}
		ids[i] = idInt
	}

	err := s.store.GiveOut(ids)
	if err != nil {
		return fmt.Errorf("s.store.GiveOut: %w", err)
	}

	fmt.Println("Orders have been given out to the client")

	return nil
}

func (s *Service) List(clientId int, lastn int, inPvz bool) error {
	if clientId == -1 {
		return fmt.Errorf("incorrect input data format")
	}

	list, err := s.store.List(clientId, lastn, inPvz)
	if err != nil {
		return fmt.Errorf("s.store.List: %w", err)
	}

	if len(list) == 0 {
		return fmt.Errorf("this client does not have orders with such parameters")
	}

	fmt.Println("Customer orders:", list)

	return nil
}

func (s *Service) Return(id int, clientId int) error {
	if id == -1 || clientId == -1 {
		return fmt.Errorf("incorrect input data format")
	}

	err := s.store.Return(id, clientId)
	if err != nil {
		return fmt.Errorf("s.store.Return: %w", err)
	}

	fmt.Println("Order return accepted")

	return nil
}

func (s *Service) ListOfReturned(pagenum int, itemsonpage int) error {
	if pagenum == -1 || itemsonpage == -1 {
		return fmt.Errorf("incorrect input data format")
	}

	list, err := s.store.ListOfReturned(pagenum, itemsonpage)
	if err != nil {
		return fmt.Errorf("s.store.ListOfReturned: %w", err)
	}

	if len(list) == 0 {
		return fmt.Errorf("there are no returned orders")
	}

	fmt.Printf("Returned orders (page %d; %d orders on page):\n%v\n", pagenum, itemsonpage, list)

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
		return s.GiveOut(*idsStr)
	case "list":
		return s.List(*clientId, *lastn, *inPvz)
	case "return":
		return s.Return(*id, *clientId)
	case "listofreturned":
		return s.ListOfReturned(*pagenum, *itemsonpage)
	default:
		return fmt.Errorf("unknown command")
	}
}
