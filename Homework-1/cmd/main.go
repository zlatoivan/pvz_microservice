package main

import (
	"flag"
	"fmt"
	"os"
	"route_256/homework/Homework-1/internal/model"
	"route_256/homework/Homework-1/internal/service"
	"route_256/homework/Homework-1/internal/storage"
	"strconv"
	"strings"
	"time"
)

func main() {
	// Проверка аргументов
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("Недостаточно аргументов")
		return
	}
	command := args[len(args)-1]

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

	// Создание бд
	stor, err := storage.New()
	if err != nil {
		fmt.Println("не удалось подключиться к хранилищу")
		return
	}

	// Создание сервиса
	serv := service.New(&stor)
	_ = serv

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
			fmt.Println("Неправильный формат входных данных")
			return
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
		err = serv.Create(newOrder)
		if err != nil {
			fmt.Println("Ошибка создания заказа:", err)
			return
		}

		fmt.Println("Заказ принят")

	case "delete":
		if *id == -1 {
			fmt.Println("Неправильный формат входных данных")
			return
		}

		err = serv.Delete(*id)
		if err != nil {
			fmt.Println("Ошибка удаления заказа:", err)
			return
		}

		fmt.Println("Заказ удален")

	case "giveout":
		idsToSplit := (*idsStr)[1 : len(*idsStr)-1]
		idsToInt := strings.Split(idsToSplit, ",")
		ids := make([]int, len(idsToInt))
		for i := range idsToInt {
			idInt, err := strconv.Atoi(idsToInt[i])
			if err != nil {
				fmt.Println("Неверные id заказов", idsToInt[i])
			}
			ids[i] = idInt
		}

		err = serv.Giveout(ids)
		if err != nil {
			fmt.Println("Ошибка выдачи заказов клиенту:", err)
			return
		}

		fmt.Println("Заказы выданы клиенту")

	case "list":
		if *clientId == -1 {
			fmt.Println("Неправильный формат входных данных")
			return
		}

		list, err := serv.List(*clientId, *lastn, *inPvz)
		if err != nil {
			fmt.Println("Ошибка получения списка заказов клиента:", err)
			return
		}

		if len(list) == 0 {
			fmt.Println("У данного пользователя нет заказов с такими параметрами")
			return
		}

		fmt.Println("Заказы клиента:", list)

	case "return":
		if *id == -1 || *clientId == -1 {
			fmt.Println("Неправильный формат входных данных")
		}

		err = serv.Return(*id, *clientId)
		if err != nil {
			fmt.Println("Ошибка возврата заказа:", err)
			return
		}

		fmt.Println("Возврат заказа принят")

	case "listofreturned":
		if *pagenum == -1 || *itemsonpage == -1 {
			fmt.Println("Неправильный формат входных данных")
			return
		}

		list, err := serv.ListOfRuturned(*pagenum, *itemsonpage)
		if err != nil {
			fmt.Println("Ошибка получения списка возвращенных заказов:", err)
			return
		}

		if len(list) == 0 {
			fmt.Println("Возвращенных заказов нет")
			return
		}

		fmt.Printf("Возвращенные заказы (стр.%d; %d заказа на стр.):\n%v\n", *pagenum, *itemsonpage, list)

	default:
		fmt.Println("Неизвестная команда")
		return
	}
}
