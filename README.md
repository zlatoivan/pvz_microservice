# Бэкенд для работы с ПВЗ

## Архитектура

Демонстрация архитектуры сервера показана на примере этапов запроса создания нового заказа: <br/>
![GOHW-4, UML.svg](uml%2FGOHW-4%2C%20UML.svg)


## Запуск программы

Локально приложение, остальное в докере (пока только так работает):

    make compose-up
    make migration-up-test
    make run-test
    make test-integration (в другой консоли)

Локально приложение, остальное в докере (устаревшее):

    В config/config.yaml postgres: host: сделать localhost.
    Оставить в docker-compose.yaml только pg_db (для тестов pg_db_test)
    Затем:

    docker compose up -d pg_db   (для тестов pg_tb_test)
	docker compose build
    make gen-ssl-cert            (если нужны свежие сертификаты)
    make migration-up            (для тестов migration-test-up)
    make run                     (для тестов run-test)

## Остановка программы

    doocker compose down


## Запросы к серверу

Первый запрос - для HTTP. <br/>
Второй запрос - для HTTPS.

##### Main Page
````
curl http://localhost:9000 -i -k -L
````
````
curl https://localhost:9001 -i -k
````

##### Create PVZ
````
curl --post301 http://localhost:9000/api/v1/pvzs -i -k --location-trusted -u ivan:pvz_best_pass -H 'Content-Type: application/json' -d '{
    "name": "Ozon Tech",
    "address": "Moscow, Presnenskaya nab. 10, block С",
    "contacts": "+7 958 400-00-05, add 76077"
}'
````
````    
curl -X POST https://localhost:9001/api/v1/pvzs -i -k -u ivan:pvz_best_pass -H 'Content-Type: application/json' -d '{
    "name": "Ozon Tech",
    "address": "Moscow, Presnenskaya nab. 10, block С",
    "contacts": "+7 958 400-00-05, add 76077"
}'
````

##### List of PVZs
````
curl http://localhost:9000/api/v1/pvzs -i -k --location-trusted -u ivan:pvz_best_pass
````
````
curl https://localhost:9001/api/v1/pvzs -i -k -u ivan:pvz_best_pass
````

##### Get PVZ by ID
````
curl http://localhost:9000/api/v1/pvzs/id -i -k --location-trusted -u ivan:pvz_best_pass -H 'Content-Type: application/json' -d '{
    "id": "9967bb48-bd6f-4ad0-924d-8c9094c4d8c2"
}'
````
````
curl https://localhost:9001/api/v1/pvzs/id -i -k -u ivan:pvz_best_pass -H 'Content-Type: application/json' -d '{
    "id": "9967bb48-bd6f-4ad0-924d-8c9094c4d8c2"
}'
````

##### Update PVZ
````
curl -X PUT http://localhost:9000/api/v1/pvzs/id -i -k --location-trusted -u ivan:pvz_best_pass -H 'Content-Type: application/json' -d '{
    "id": "9967bb48-bd6f-4ad0-924d-8c9094c4d8c2",
    "name": "Ozon Company",
    "address": "Moscow, Arbat, 27",
    "contacts": "+7 999 888 11 11"
}'
````
````
curl -X PUT https://localhost:9001/api/v1/pvzs/id -i -k -u ivan:pvz_best_pass -H 'Content-Type: application/json' -d '{
    "id": "9967bb48-bd6f-4ad0-924d-8c9094c4d8c2",
    "name": "Ozon Company",
    "address": "Moscow, Arbat, 27",
    "contacts": "+7 999 888 11 11"
}'
````

##### Delete PVZ
````
curl -X DELETE http://localhost:9000/api/v1/pvzs/id -i -k --location-trusted -u ivan:pvz_best_pass -H 'Content-Type: application/json' -d '{
    "id": "9967bb48-bd6f-4ad0-924d-8c9094c4d8c2"
}'
````
````
curl -X DELETE https://localhost:9001/api/v1/pvzs/id -i -k -u ivan:pvz_best_pass -H 'Content-Type: application/json' -d '{
    "id": "9967bb48-bd6f-4ad0-924d-8c9094c4d8c2"
}'
````

##### Create Order
````
curl --post301 http://localhost:9000/api/v1/orders -i -k --location-trusted -u ivan:order_best_pass -H 'Content-Type: application/json' -d '{
    "client_id": "9967bb48-bd6f-4ad0-924d-8c9094c4d8c2",
    "stores_till": "2024-04-22T13:14:01Z",
    "weight": 29,
    "cost": 1100,
    "packaging_type": "box"
}'
````
````    
curl -X POST http://localhost:9000/api/v1/orders -i -k -u ivan:order_best_pass -H 'Content-Type: application/json' -d '{
    "client_id": "9967bb48-bd6f-4ad0-924d-8c9094c4d8c2",
    "stores_till": "2024-04-22T13:14:01Z",
    "weight": 29,
    "cost": 1100,
    "packaging_type": "box"
}'
````

##### List of orders
````
curl http://localhost:9000/api/v1/orders -i -k --location-trusted -u ivan:order_best_pass
````
````
curl https://localhost:9001/api/v1/pvzs -i -k -u ivan:order_best_pass
````

##### Get order by ID
````
curl -X GET http://localhost:9000/api/v1/orders/id -i -k --location-trusted -u ivan:order_best_pass -H 'Content-Type: application/json' -d '{
    "id": "165dbe30-d0c4-4727-9504-827db76d214e",
}'
````
````
curl -X GET https://localhost:9001/api/v1/orders/id -i -k -u ivan:order_best_pass -H 'Content-Type: application/json' -d '{
    "id": "165dbe30-d0c4-4727-9504-827db76d214e",
}'
````

##### Update order
````
curl -X PUT http://localhost:9000/api/v1/orders/id -i -k --location-trusted -u ivan:order_best_pass -H 'Content-Type: application/json' -d '{
    "id": "165dbe30-d0c4-4727-9504-827db76d214e",
    "client_id": "9967bb48-bd6f-4ad0-924d-8c9094c4d8c2",
    "stores_till": "2024-04-22T13:14:01Z",
    "weight": 15,
    "cost": 500,
    "packaging_type": "tape"
}'
````
````
curl -X PUT https://localhost:9001/api/v1/orders/id -i -k -u ivan:order_best_pass -H 'Content-Type: application/json' -d '{
    "id": "165dbe30-d0c4-4727-9504-827db76d214e",
    "client_id": "9967bb48-bd6f-4ad0-924d-8c9094c4d8c2",
    "stores_till": "2024-04-22T13:14:01Z",
    "weight": 15,
    "cost": 500,
    "packaging_type": "tape"
}'
````

##### Delete order
````
curl -X DELETE http://localhost:9000/api/v1/orders/id -i -k --location-trusted -u ivan:order_best_pass -H 'Content-Type: application/json' -d '{
    "id": "165dbe30-d0c4-4727-9504-827db76d214e",
}'
````
````
curl -X DELETE https://localhost:9001/api/v1/orders/id -i -k -u ivan:order_best_pass -H 'Content-Type: application/json' -d '{
    "id": "165dbe30-d0c4-4727-9504-827db76d214e",
}'
````

##### List of client orders
````
curl -X GET http://localhost:9000/api/v1/orders/client/id -i -k --location-trusted -u ivan:order_best_pass -H 'Content-Type: application/json' -d '{
    "id": "165dbe30-d0c4-4727-9504-827db76d214e",
}'
````
````
curl -X GET https://localhost:9001/api/v1/orders/client/id -i -k -u ivan:order_best_pass -H 'Content-Type: application/json' -d '{
    "id": "165dbe30-d0c4-4727-9504-827db76d214e",
}'
````

##### Give out orders
````
curl -X PUT http://localhost:9000/api/v1/orders/client/id -i -k --location-trusted -u ivan:order_best_pass -H 'Content-Type: application/json' -d '{
    "client_id": "165dbe30-d0c4-4727-9504-827db76d214e",
    "ids": [
        "427bf09a-59ff-4e2d-b55f-19582037456d",        
        "06db932a-f3b1-49bc-9928-dd5838b38d76"
    ]
}'
````
````
curl -X PUT https://localhost:9001/api/v1/orders/client/id -i -k -u ivan:order_best_pass -H 'Content-Type: application/json' -d '{
    "client_id": "165dbe30-d0c4-4727-9504-827db76d214e",
    "ids": [
        "427bf09a-59ff-4e2d-b55f-19582037456d",        
        "06db932a-f3b1-49bc-9928-dd5838b38d76"
    ]
}'
````

##### Return order
````
curl -X PUT http://localhost:9000/api/v1/orders/client/id/return -i -k --location-trusted -u ivan:order_best_pass -H 'Content-Type: application/json' -d '{
    "client_id": "427bf09a-59ff-4e2d-b55f-19582037456d",
    "id": "165dbe30-d0c4-4727-9504-827db76d214e"
}'
````
````
curl -X PUT https://localhost:9001/api/v1/orders/client/id/return -i -k -u ivan:order_best_pass -H 'Content-Type: application/json' -d '{
    "client_id": "427bf09a-59ff-4e2d-b55f-19582037456d",
    "id": "165dbe30-d0c4-4727-9504-827db76d214e"
}'
````

##### List of returned orders
````
curl http://localhost:9000/api/v1/orders/returned -i -k --location-trusted -u ivan:order_best_pass
````
````
curl https://localhost:9001/api/v1/orders/returned -i -k -u ivan:order_best_pass
````


##### Дополнительные флаги
    
    -i (--include) Выводит и заголовки, и тело ответа
    -k (--insecure) Игнорирует ошибки SSL сертификата
    -L (--location) Разрешает преадресацию
    --location-trusted Сохраняет данные для аутентификации при переадресации
    -u (--user) Данные для аутентификации
    -d (--data) Данные в теле запроса
    -v (--verbose) Выводит подробную информацию о заголовках и тело ответа


## Консольный режим

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

##### Входные данные для консольного режима

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
    Интерактивный режим (Включить, далее следовать его командам)
        go run cmd/cli/main.go interactive_mode