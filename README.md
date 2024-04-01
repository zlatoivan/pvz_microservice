# Утилита для работы с ПВЗ

## UML диаграмма

Для описания архитектуры решения по проверке типа упаковки был выбран стандарт "Sequence diagram"
Документация: https://mermaid.js.org/syntax/sequenceDiagram.html

Финальная версия диаграммы находится в папке /uml:
[GOHW-4, UML.mmd](uml%2FGOHW-4%2C%20UML.mmd)

И имеет следующий вид:
![GOHW-4, UML dark.svg](uml%2FGOHW-4%2C%20UML%20dark.svg)
![GOHW-4, UML dark.svg](uml%2FGOHW-4%2C%20UML%20dark.svg)
![GOHW-4, UML.png](uml%2FGOHW-4%2C%20UML.png)
![GOHW-4, UML.svg](uml%2FGOHW-4%2C%20UML.svg)


## Запуск программы

    make compose-up
    make gen-ssl-cert  (если нужны свежие сертификаты)
    make migration-up
    make run


## Запросы к серверу

Первый запрос - для HTTP.
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
curl http://localhost:9000/api/v1/pvzs/9967bb48-bd6f-4ad0-924d-8c9094c4d8c2 -i -k --location-trusted -u ivan:pvz_best_pass
````
````
curl https://localhost:9001/api/v1/pvzs/9967bb48-bd6f-4ad0-924d-8c9094c4d8c2 -i -k -u ivan:pvz_best_pass
````

##### Update PVZ
````
curl -X PUT http://localhost:9000/api/v1/pvzs/86595598-f70d-4ffa-bc2b-29e11de41df8 -i -k --location-trusted -u ivan:pvz_best_pass -H 'Content-Type: application/json' -d '{
    "name": "Ozon Company",
    "address": "Moscow, Arbat, 27",
    "contacts": "+7 999 888 11 11"
}'
````
````
curl -X PUT https://localhost:9001/api/v1/pvzs/9967bb48-bd6f-4ad0-924d-8c9094c4d8c2 -i -k -u ivan:pvz_best_pass -H 'Content-Type: application/json' -d '{
    "name": "Ozon Company",
    "address": "Moscow, Arbat, 27",
    "contacts": "+7 999 888 11 11"
}'
````

##### Delete PVZ
````
curl -X DELETE http://localhost:9000/api/v1/pvzs/3bdc65d0-3e6a-406f-9ed1-b52962b5faf8 -i -k --location-trusted -u ivan:pvz_best_pass
````
````
curl -X DELETE https://localhost:9001/api/v1/pvzs/9967bb48-bd6f-4ad0-924d-8c9094c4d8c2 -i -k -u ivan:pvz_best_pass
````

##### Create Order
````
curl --post301 http://localhost:9000/api/v1/orders -i -k --location-trusted -u ivan:order_best_pass -H 'Content-Type: application/json' -d '{
    "client_id": "9967bb48-bd6f-4ad0-924d-8c9094c4d8c2",
    "stores_till": "22.04.2024 13:14",
    "weight": 29,
    "cost": 1100,
    "packaging_type": "box"
}'
````
````    
curl -X POST http://localhost:9000/api/v1/orders -i -k -u ivan:order_best_pass -H 'Content-Type: application/json' -d '{
    "client_id": "9967bb48-bd6f-4ad0-924d-8c9094c4d8c2",
    "stores_till": "22.04.2024 13:14",
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
curl http://localhost:9000/api/v1/orders/165dbe30-d0c4-4727-9504-827db76d214e -i -k --location-trusted -u ivan:order_best_pass
````
````
curl https://localhost:9001/api/v1/orders/165dbe30-d0c4-4727-9504-827db76d214e -i -k -u ivan:order_best_pass
````

##### Update order
````
curl -X PUT http://localhost:9000/api/v1/orders/165dbe30-d0c4-4727-9504-827db76d214e -i -k --location-trusted -u ivan:order_best_pass -H 'Content-Type: application/json' -d '{
    "client_id": "9967bb48-bd6f-4ad0-924d-8c9094c4d8c2",
    "stores_till": "22.04.2024 13:14",
    "weight": 15,
    "cost": 500,
    "packaging_type": "tape"
}'
````
````
curl -X PUT https://localhost:9001/api/v1/orders/165dbe30-d0c4-4727-9504-827db76d214e -i -k -u ivan:order_best_pass -H 'Content-Type: application/json' -d '{
    "client_id": "9967bb48-bd6f-4ad0-924d-8c9094c4d8c2",
    "stores_till": "22.04.2024 13:14",
    "weight": 15,
    "cost": 500,
    "packaging_type": "tape"
}'
````

##### Delete order
````
curl -X DELETE http://localhost:9000/api/v1/orders/5ae15592-7ef2-41b2-a1f1-959d6e935c -i -k --location-trusted -u ivan:order_best_pass
````
````
curl -X DELETE https://localhost:9001/api/v1/orders/5ae15592-7ef2-41b2-a1f1-959d6e935c -i -k -u ivan:order_best_pass
````

##### List of client orders
````
curl http://localhost:9000/api/v1/orders/client/9967bb48-bd6f-4ad0-924d-8c9094c4d8c2 -i -k --location-trusted -u ivan:order_best_pass
````
````
curl https://localhost:9001/api/v1/orders/client/9967bb48-bd6f-4ad0-924d-8c9094c4d8c2 -i -k -u ivan:order_best_pass
````

##### Give out orders
````
curl -X PUT http://localhost:9000/api/v1/orders/client/9967bb48-bd6f-4ad0-924d-8c9094c4d8c2 -i -k --location-trusted -u ivan:order_best_pass -H 'Content-Type: application/json' -d '{
    "ids": [
        "427bf09a-59ff-4e2d-b55f-19582037456d",        
        "06db932a-f3b1-49bc-9928-dd5838b38d76"
    ]
}'
````
````
curl -X PUT https://localhost:9001/api/v1/orders/client/9967bb48-bd6f-4ad0-924d-8c9094c4d8c2 -i -k -u ivan:order_best_pass -H 'Content-Type: application/json' -d '{
    "ids": [
        "427bf09a-59ff-4e2d-b55f-19582037456d",        
        "06db932a-f3b1-49bc-9928-dd5838b38d76"
    ]
}'
````

##### Return order
````
curl -X PUT http://localhost:9000/api/v1/orders/client/9967bb48-bd6f-4ad0-924d-8c9094c4d8c2/return -i -k --location-trusted -u ivan:order_best_pass -H 'Content-Type: application/json' -d '{
    "id": "427bf09a-59ff-4e2d-b55f-19582037456d"
}'
````
````
curl -X PUT https://localhost:9001/api/v1/orders/client/9967bb48-bd6f-4ad0-924d-8c9094c4d8c2/return -i -k -u ivan:order_best_pass -H 'Content-Type: application/json' -d '{
    "id": "427bf09a-59ff-4e2d-b55f-19582037456d"
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