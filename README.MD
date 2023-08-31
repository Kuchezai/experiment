# Сервис динамического сегментирования пользователей
## Description
Сервис, хранящий пользователя и сегменты, в которых он состоит. Сервис позволяет создавать, удалять и добавлять сегменты, а также присваивать их пользователям.

#### Стек
- Golang, Gin
- PostgreSQL
- Docker
- Swagger
- JWT

#### Выполненные задания:
1. Сохранение истории попадания/выбывания пользователя из сегмента с возможностью получения отчета по пользователю за определенный период
2. Возможность задавать TTL
3. Возможность добавлять пользователя в сегмент не в ручную, а автоматически. В сегмент попадает заданный процент пользователей

#### Особенности:
1. Покрытие кода тестами (% покрытия >50)
2. В репозитории настроены GitHub Actions для автоматических прогонов тестов при push'e в main ветку
3. Документация Swagger, доступная при запуске по адресу  http://localhost:8888/#/
4. Добавление аутентификации c использованием JWT





## Service starting
1. Переименуйте `.env.example ` в `.env`, поменяйте переменные окружения при необходимости.
2. Запустите контейнеры командой ``
make compose-up
``

## Requests

###  Регистрация пользователя

`Request:`

``` 
curl --location 'http://localhost:8080/api/v1/registration' \
--header 'Content-Type: application/json' \
--data '{
    "name": "username",
    "pass": "superpass"
}'
```
`Response:`

```json 
{
    "id": 1
}
``` 

###  Аутентификация пользователя

`Request:`

``` 
curl --location 'http://localhost:8080/api/v1/login' \
--header 'Content-Type: application/json' \
--data '{
    "name": "username",
    "pass": "superpass"
}'
```

`Response:`

```json 
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoidXNlcm5hbWUxMjIyIn0.OXn8NT4w1xKPydDlnJUIWt6I9RBqtFbCbz17-otdibQ"
}
```

### Создание сегмента

`Request:`

``` 
curl --location 'http://localhost:8080/api/v1/segments' \
--header 'Content-Type: application/json' \
--data  '{
	"slug": "AVITO_PERFORMANCE_VAS"
}'
```

`Response:`

```
201 Created
```

###  Создание сегмента с автоматическим присвоением

`Request:`

``` 
curl --location 'http://localhost:8080/api/v1/segments/auto-assign' \
--header 'Content-Type: application/json' \
--data '{
    "slug": "AVITO_DISCOUNT_30",
    "percent": 70
}'
```

`Response:`

```json 
{
	"ids":  [9,10,7,1,5]
}
```

### Удаление сегмента

`Request:`

``` 
curl --location --request DELETE 'http://localhost:8080/api/v1/segments/AVITO_DISCOUNT_30' \
--data ''
```

`Response:`

```
200 OK
```

### Получение сегментов пользователя

`Request:`

``` 
curl --location 'http://localhost:8080/api/v1/users/1/segments' \
--data ''
```

`Response:`

```json
[
    {
        "slug": "AVITO_DISCOUNT_50",
        "expired_date": "9999-12-31T23:59:59.999999999Z"
    },
    {
        "slug": "AVITO_PERFORMANCE_VAS",
        "expired_date": "2024-08-31T17:53:24.527156Z"
    }
]
```

### Редактирование сегментов пользователя

`Request:`

``` 
curl --location --request PATCH 'http://localhost:8080/api/v1/users/1/segments' \
--header 'Content-Type: application/json' \
--data '{
    "add_segments": [
        {
            "slug": "AVITO_PERFORMANCE_VAS",
            "ttl": 10
        }
    ],
    "remove_segments": [
        "AVITO_DISCOUNT_50"
    ]
}'
```

`Response:`

```
200 OK
```

### Создать CSV файл с историей добавления/выбывания сегментов

`Request:`

``` 
curl --location 'http://localhost:8080/api/v1/users/segments/history' \
--header 'Content-Type: application/json' \
--data '{
    "year" : 2023,
    "month" : 8
}'
```

`Response:`

```json
{
    "link": "./history/user_segments_history-2023-8.csv"
}
```

### Скачать CSV файл с историей добавления/выбывания сегментов
```diff
! Это защищенная ручка, чтобы скачать файл нужно указать токен, получаемый при аутентификации
```

`Request:`



``` 
curl --location 'http://localhost:8080/history/user_segments_history-2023-8.csv' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoidGVzdCJ9.Poy32jrZkt3gvJL6uzfE4buxhbrNjTwvTZGPoQntvgk'
```

`Response*:`

```
    operation_id,user_id,segment_slug,is_added,date
1,9,AVITO_DISCOUNT_30,true,2023-08-31T17:47:24Z
2,10,AVITO_DISCOUNT_30,true,2023-08-31T17:47:24Z
3,7,AVITO_DISCOUNT_30,true,2023-08-31T17:47:24Z
4,1,AVITO_DISCOUNT_30,true,2023-08-31T17:47:24Z
5,5,AVITO_DISCOUNT_30,true,2023-08-31T17:47:24Z
6,9,AVITO_DISCOUNT_30,false,2023-08-31T17:49:16Z
7,10,AVITO_DISCOUNT_30,false,2023-08-31T17:49:16Z
```
`*При запросе из браузера будет скачен .csv файл`
