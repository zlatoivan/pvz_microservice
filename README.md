# Утилита для работы с ПВЗ

## Запуск программы

    make compose-up
    make gen-ssl-cert  (если нужны свежие сертификаты)
    make migration-up
    make run


## Входные данные ДЗ-3

Первый запрос - для HTTP, второй - для HTTPS.

##### Main Page
````
curl http://localhost:9000 -i -k -L
````
````
curl https://localhost:9001 -i -k
````

##### Create
````
curl --post301 http://localhost:9000/api/v1/pvzs -i -k --location-trusted -u ivan:pvz_best_pass -H 'Content-Type: application/json' -d '{
    "Name": "Ozon Tech",
    "Address": "Moscow, Presnenskaya nab. 10, block С",
    "Contacts": "+7 958 400-00-05, add 76077"
}'
````
````    
curl -X POST https://localhost:9001/api/v1/pvzs -i -k -u ivan:pvz_best_pass -H 'Content-Type: application/json' -d '{
    "Name": "Ozon Tech",
    "Address": "Moscow, Presnenskaya nab. 10, block С",
    "Contacts": "+7 958 400-00-05, add 76077"
}'
````

##### List
````
curl http://localhost:9000/api/v1/pvzs -i -k --location-trusted -u ivan:pvz_best_pass
````
````
curl https://localhost:9001/api/v1/pvzs -i -k -u ivan:pvz_best_pass
````

##### GetById (Вставить UUID)
````
curl http://localhost:9000/api/v1/pvzs/9967bb48-bd6f-4ad0-924d-8c9094c4d8c2 -i -k --location-trusted -u ivan:pvz_best_pass
````
````
curl https://localhost:9001/api/v1/pvzs/9967bb48-bd6f-4ad0-924d-8c9094c4d8c2 -i -k -u ivan:pvz_best_pass
````

##### Update (Вставить UUID)
````
curl -X PUT http://localhost:9000/api/v1/pvzs/86595598-f70d-4ffa-bc2b-29e11de41df8 -i -k --location-trusted -u ivan:pvz_best_pass -H 'Content-Type: application/json' -d '{
    "Name": "Ozon Company",
    "Address": "Moscow, Arbat, 27",
    "Contacts": "+7 999 888 11 11"
}'
````
````
curl -X PUT https://localhost:9001/api/v1/pvzs/9967bb48-bd6f-4ad0-924d-8c9094c4d8c2 -i -k -u ivan:pvz_best_pass -H 'Content-Type: application/json' -d '{
    "Name": "Ozon Company",
    "Address": "Moscow, Arbat, 27",
    "Contacts": "+7 999 888 11 11"
}'
````

##### Delete (Вставить UUID)
````
curl -X DELETE http://localhost:9000/api/v1/pvzs/3bdc65d0-3e6a-406f-9ed1-b52962b5faf8 -i -k --location-trusted -u ivan:pvz_best_pass
````
````
curl -X DELETE https://localhost:9001/api/v1/pvzs/9967bb48-bd6f-4ad0-924d-8c9094c4d8c2 -i -k -u ivan:pvz_best_pass
````

##### Дополнительные флаги
    
    -i (--include) Выводит и заголовки, и тело ответа
    -k (--insecure) Игнорирует ошибки SSL сертификата
    -L (--location) Разрешает преадресацию
    --location-trusted Сохраняет данные для аутентификации при переадресации
    -u (--user) Данные для аутентификации
    -d (--data) Данные в теле запроса
    -v (--verbose) Выводит подробную информацию о заголовках и тело ответа


## Входные данные ДЗ-2

##### Запустить интерактивный режим

    go run cmd/cli/main.go interactive_mode

Далее следовать его командам.


## Входные данные ДЗ-1

    Принять заказ (создать):
        go run cmd/cli/main.go -clientid=9886 -storestill=15.09.2024 create
    Вернуть заказ курьеру (удалить):
        go run cmd/cli/main.go -id=9967bb48-bd6f-4ad0-924d-8c9094c4d8c2 delete
    Выдать заказ клиенту:
        go run cmd/cli/main.go -ids=[9967bb48-bd6f-4ad0-924d-8c9094c4d8d3,9967bb48-bd6f-4ad0-924d-8c9094c4d8c2] giveout
    Получить список заказов:
        go run cmd/cli/main.go -clientid=9967bb48-bd6f-4ad0-924d-8c9094c4d8c2 -lastn=2 -inpvz=true list
    Возврат заказа клиентом:
        go run cmd/cli/main.go -id=9967bb48-bd6f-4ad0-924d-8c9094c4d8c2 -clientid=9967bb48-bd6f-4ad0-924d-8c9094c4d8r1 return
    Получить список возвращенных товаров:
        go run cmd/cli/main.go -pagenum=1 -itemsonpage=2 listofreturned

go run cmd/cli/main.go help

    Это утилита для управления ПВЗ.

    Применение:
        go run cmd/main.go [flags] [command]
    
    command:            Описание:                                flags:
        create            Принять заказ (создать).                 -clientid=9967bb48-bd6f-4ad0-924d-8c9094c4d8c2 -storestill=15.09.2024
        delete            Вернуть заказ курьеру (удалить).         -id=9967bb48-bd6f-4ad0-924d-8c9094c4d8c2
        giveout           Выдать заказ клиенту.                    -ids=[9967bb48-bd6f-4ad0-924d-8c9094c4d8c2,9967bb48-bd6f-4ad0-924d-8c9094c4d8e4]
        list              Получить список заказов клиента.         -clientid=9967bb48-bd6f-4ad0-924d-8c9094c4d8c2 -lastn=2 -inpvz=true  (последние два опциональные)
        return            Возврат заказа клиентом.                 -id=9967bb48-bd6f-4ad0-924d-8c9094c4d8c2 -clientid=9967bb48-bd6f-4ad0-924d-8c9094c4d8r3
        listofreturned    Получить список возвращенных заказов.    -pagenum=1 -itemsonpage=2

        interactive_mode  Interactive mode to add and get PVZ      No flags. Enter command and follow the instructions
            command:
                add       Create PVZ
                get       Get the information about PVZ
