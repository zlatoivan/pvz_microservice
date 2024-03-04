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
		fmt.Print(
			"Это утилита для управления ПВЗ.\n\n" +
				"Применение:\n" +
				"\tgo run cmd/main.go [flags] [command]\n\n" +
				"command:            Описание:                                flags:\n" +
				"\tcreate            Принять заказ (создать).                 -id=1212 -clientid=9886 -shelflife=15.09.2024\n" +
				"\tdelete            Вернуть заказ курьеру (удалить).         -id=1212\n" +
				"\tgiveout           Выдать заказ клиенту.                    -ids=[1212,1214]\n" +
				"\tlist              Получить список заказов клиента.         -clientid=9886 -lastn=2 -inpvz=true  (последние два опциональные)\n" +
				"\treturn            Возврат заказа клиентом.                 -id=1212 -clientid=9886\n" +
				"\tlistofreturned    Получить список возвращенных заказов.    -pagenum=1 -itemsonpage=2\n",
		)

	case "create":
		if *id == -1 || *clientId == -1 || *shelfLifeStr == "-" {
			return errors.New("Неправильный формат входных данных")
		}

		// Привести срок хранения к типу даты
		datePattern := "02.01.2006"
		shelfLife, err := time.Parse(datePattern, *shelfLifeStr)
		if err != nil {
			fmt.Println(err)
		}

		newOrder := model.Order{
			Id:          *id,
			ClientId:    *clientId,
			ShelfLife:   shelfLife,
			IsDeleted:   false,
			IsGaveOut:   false,
			GiveOutTime: time.Date(1, 0, 0, 0, 0, 0, 0, time.UTC),
			IsReturned:  false,
		}
		err = s.Create(newOrder)
		if err != nil {
			return errors.New("Ошибка создания заказа: " + err.Error())
		}

		fmt.Println("Заказ принят")

	case "delete":
		if *id == -1 {
			return errors.New("Неправильный формат входных данных")
		}

		err := s.Delete(*id)
		if err != nil {
			return errors.New("Ошибка удаления заказа: " + err.Error())
		}

		fmt.Println("Заказ удален")

	case "giveout":
		idsToSplit := (*idsStr)[1 : len(*idsStr)-1]
		idsToInt := strings.Split(idsToSplit, ",")
		ids := make([]int, len(idsToInt))
		for i := range idsToInt {
			idInt, err := strconv.Atoi(idsToInt[i])
			if err != nil {
				return errors.New("Неверные id заказов " + idsToInt[i])
			}
			ids[i] = idInt
		}

		err := s.Giveout(ids)
		if err != nil {
			return errors.New("Ошибка выдачи заказов клиенту: " + err.Error())
		}

		fmt.Println("Заказы выданы клиенту")

	case "list":
		if *clientId == -1 {
			return errors.New("Неправильный формат входных данных")
		}

		list, err := s.List(*clientId, *lastn, *inPvz)
		if err != nil {
			return errors.New("Ошибка получения списка заказов клиента: " + err.Error())
		}

		if len(list) == 0 {
			return errors.New("У данного пользователя нет заказов с такими параметрами")
		}

		fmt.Println("Заказы клиента:", list)

	case "return":
		if *id == -1 || *clientId == -1 {
			return errors.New("Неправильный формат входных данных")
		}

		err := s.Return(*id, *clientId)
		if err != nil {
			return errors.New("Ошибка возврата заказа: " + err.Error())
		}

		fmt.Println("Возврат заказа принят")

	case "listofreturned":
		if *pagenum == -1 || *itemsonpage == -1 {
			return errors.New("Неправильный формат входных данных")
		}

		list, err := s.ListOfRuturned(*pagenum, *itemsonpage)
		if err != nil {
			return errors.New("Ошибка получения списка возвращенных заказов: " + err.Error())
		}

		if len(list) == 0 {
			return errors.New("Возвращенных заказов нет")
		}

		fmt.Printf("Возвращенные заказы (стр.%d; %d заказа на стр.):\n%v\n", *pagenum, *itemsonpage, list)

	default:
		return errors.New("Неизвестная команда")
	}

	return nil
}
