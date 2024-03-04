package main

import (
	"fmt"
	"os"
	"route_256/homework/Homework-1/internal/service"
	"route_256/homework/Homework-1/internal/storage"
)

func main() {
	// Проверка аргументов
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("Недостаточно аргументов")
		return
	}

	// Создание бд
	stor, err := storage.New()
	if err != nil {
		fmt.Println("Не удалось подключиться к хранилищу")
		return
	}

	// Создание сервиса
	serv := service.New(&stor)

	// Работа сервиса
	err = serv.Work()
	if err != nil {
		fmt.Println(err)
		return
	}
}
