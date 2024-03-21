# Домашнее задание 3
Работу продолжаем в репозитории домашнего задания 2.

1.  __Сохранение данных в бд__
    Подключить работу с базой данных, добавить конфиг подключения, инициализацию коннекта.

Работу с мапкой/файлом оставляем неизменными, доменные структуры можно переиспользовать

Написать CRUD операции для работы с бд
Должны быть реализованы методы записи и чтения данных простой системы  хранения ПВЗ

###### _Подсказка:_
_Вам могут пригодится следующие методы_
- _GetByID_
- _List_
- _Update_
- _Create_
- _Delete_

_Так же можете реализовать те методы, которые вы делали с файлом._

2. __Разработать HTTP сервер.__

- Необходимо реализовать HTTP сервер, который будет работать с методами базы данных, реализованными в 1 пункте.

- Методы должны позволять манипулировать данными(create,read,update,delete) для системы хранения пвз

- Методы должны быть выполнены в restful стиле. Необходимо корректно обрабатывать все коды ошибок

- Входящие параметры должны быть представлены либо в формате json либо в query параметрах(зависит от метода)

- Сервис должен использовать порт 9000

3. __В ридми приложить curl запросы, на каждую ручку. Запросы должны быть валидными и возвращать нужный код ответа__
4. __Необходимо реализовать middleware, который будет логгировать поля POST,PUT,DELETE запросов__


###### _Подсказка:_
_Посмотрите на результат выполнения дз 2. Сервис должен делать похожий flow, но используя бд как хранилище и http как интерфейс взаимодействия с пользователем_
## Дополнительно:
1. Поддержать https. Можно использовать самоподписный сертификат от Let's Encrypt
2. Реализовать middleware с basic auth. Юзер/пароль можно задать как в конфиге, так и хранить в базе(создать круд юзеров)

за дополнительные задания - 10 баллов

## Ограничения дз:
- Нельзя использовать orm или sql билдеры
- Для реализации http сервера можно использовать как net/http так и gin/fasthttp и прочее
- Коды ошибок должны соответствовать поведению сервиса. Хендлеры, которые отдают только 500 в случае ошибки - не принимаются
- В хендлерах должна быть базовая валидация данных, соответствующая бизнес-логике

# Дедлайн:

23 марта, 23:59 (сдача) / 26 марта, 23:59 (проверка)

---

# Утилита для работы с ПВЗ

## Подготовка

##### Сгенерировать приватный ключ (.key):

    openssl genrsa -out server.key 2048

##### Сгенерировать публичный ключ (.crt), основанный на приватном ключе:

    openssl req -new -x509 -sha256 -key server.key -out server.crt -days 365 -nodes

##### Поместить оба файла в папку /configs

    mv -f server.key server.crt configs/

##### Запустить демон docker

    systemctl --user start docker-desktop

##### Поднять базу данных docker compose

    docker compose up --build

##### Поднять миграцию

    make migration-up

---

## Запуск программы

    go run cmd/server/main.go

---

## Входные данные ДЗ-3

##### Main Page

    curl http://localhost:9000 -u ivan:the_best_pass

##### Create

    curl -X POST http://localhost:9000/api/v1/pvzs -u ivan:the_best_pass -H 'Content-Type: application/json' -d '{
        "Name": "Ozon Tech",
        "Address": "Moscow, Presnenskaya nab. 10, block С",
        "Contacts": "+7 958 400-00-05, add 76077"
    }'

##### List

    curl http://localhost:9000/api/v1/pvzs -u ivan:the_best_pass

##### GetById

    curl http://localhost:9000/api/v1/pvzs/2 -u ivan:the_best_pass

##### Update

    curl -X PUT http://localhost:9000/api/v1/pvzs/2 -u ivan:the_best_pass -H 'Content-Type: application/json' -d '{
        "Name": "Ozon Company",
        "Address": "Moscow, Arbat, 27",
        "Contacts": "+7 999 888 11 11"
    }'

##### Delete

    curl -X DELETE http://localhost:9000/api/v1/pvzs/2 -u ivan:the_best_pass

##### HTTPS
    
    Запросы аналогичные, но надо изменить в запросе:
    http -> https
    9000 -> 9001
    Добавить флаг -k

---

## Входные данные ДЗ-2

##### Запустить интерактивный режим

        go run cmd/cli/main.go interactive_mode

Далее следовать его командам.

---

## Входные данные ДЗ-1

    Принять заказ (создать):
        go run cmd/cli/main.go -id=1212 -clientid=9886 -storestill=15.09.2024 create
    Вернуть заказ курьеру (удалить):
        go run cmd/cli/main.go -id=1212 delete
    Выдать заказ клиенту:
        go run cmd/cli/main.go -ids=[1212,1214] giveout
    Получить список заказов:
        go run cmd/cli/main.go -clientid=9886 -lastn=2 -inpvz=true list
    Возврат заказа клиентом:
        go run cmd/cli/main.go -id=1212 -clientid=9886 return
    Получить список возвращенных товаров:
        go run cmd/cli/main.go -pagenum=1 -itemsonpage=2 listofreturned

---

## Информация о функционале

go run cmd/server/main.go

    Затем curl запросы выше

go run cmd/cli/main.go help

    Это утилита для управления ПВЗ.

    Применение:
        go run cmd/main.go [flags] [command]
    
    command:            Описание:                                flags:
        create            Принять заказ (создать).                 -id=1212 -clientid=9886 -storestill=15.09.2024
        delete            Вернуть заказ курьеру (удалить).         -id=1212
        giveout           Выдать заказ клиенту.                    -ids=[1212,1214]
        list              Получить список заказов клиента.         -clientid=9886 -lastn=2 -inpvz=true  (последние два опциональные)
        return            Возврат заказа клиентом.                 -id=1212 -clientid=9886
        listofreturned    Получить список возвращенных заказов.    -pagenum=1 -itemsonpage=2

        interactive_mode  Interactive mode to add and get PVZ      No flags. Enter command and follow the instructions
            command:
                add       Create PVZ
                get       Get the information about PVZ