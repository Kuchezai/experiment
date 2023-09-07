# Dynamic User Segmentation service
## Differences from the main branch
- Покрытие тестами увеличено до 70%
## Description
Сервис, хранящий пользователя и сегменты, в которых он состоит. Сервис позволяет создавать, удалять и добавлять сегменты, а также присваивать их пользователям.

#### Стек
- Golang, Gin
- PostgreSQL
- Docker
- Swagger
- JWT

#### Выполненные задания:
1. Сохранение истории попадания/выбывания пользователя из сегмента с возможностью получения отчета по пользователям за определенный период
2. Возможность задавать TTL в днях. При истечении срока жизни сегмент перестает учитываться
3. Возможность добавлять пользователя в сегмент не в ручную, а автоматически. В сегмент попадает заданный процент пользователей. В случае малого количества пользователей сегмент может быть не присвоен (например при проценте < 50, и числу пользователей равному 1, сегмент останется неприсвоенным, однако он присвоится единственному пользователю, если процент будет >= 50. Округление идет в соответствие с математическими правилами)

#### Особенности:
1. Покрытие кода тестами (% покрытия >50)
2. В репозитории настроены GitHub Actions для автоматических прогонов тестов при push'e в main ветку
3. Документация Swagger, доступная при запуске по адресу  http://localhost:8888/#/
4. Добавление аутентификации c использованием JWT


## Service starting
1. Переименуйте `.env.example ` в `.env`, поменяйте переменные окружения при необходимости
2. Запустите контейнеры командой ``
make compose_up
``

## Requests
* [Регистрация пользователя](#registration)
* [Аутентификация пользователя](#login)
* [Создание сегмента](#create-segment)
* [Создание сегмента с автоматическим присвоением](#create-segment-auto)
* [Удаление сегмента](#delete-segment)
* [Получение сегментов пользователя](#get-segments)
* [Редактирование сегментов пользователя](#edit-segments)
* [Создание CSV файл с историей добавления/выбывания сегментов](#create-csv)
* [Получение CSV файл с историей добавления/выбывания сегментов](#download-csv)
### <a name="registration"></a>Регистрация пользователя

Request:

``` 
curl --location 'http://localhost:8080/api/v1/registration' \
--header 'Content-Type: application/json' \
--data '{
    "name": "username",
    "pass": "superpass"
}'
```
Response:

```json 
{
    "id": 1
}
``` 

### <a name="login"></a>Аутентификация пользователя

Request:

``` 
curl --location 'http://localhost:8080/api/v1/login' \
--header 'Content-Type: application/json' \
--data '{
    "name": "username",
    "pass": "superpass"
}'
```

Response:

```json 
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoidXNlcm5hbWUxMjIyIn0.OXn8NT4w1xKPydDlnJUIWt6I9RBqtFbCbz17-otdibQ"
}
```

### <a name="create-segment"></a>Создание сегмента

Request:

``` 
curl --location 'http://localhost:8080/api/v1/segments' \
--header 'Content-Type: application/json' \
--data  '{
	"slug": "AVITO_PERFORMANCE_VAS"
}'
```

Response:

```
201 Created
```

###  <a name="create-segment-auto"></a>Создание сегмента с автоматическим присвоением

Request:

``` 
curl --location 'http://localhost:8080/api/v1/segments/auto-assign' \
--header 'Content-Type: application/json' \
--data '{
    "slug": "AVITO_DISCOUNT_30",
    "percent": 70
}'
```

Response:

```json 
{
    "ids":  [9,10,7,1,5]
}
```

### <a name="delete-segment"></a>Удаление сегмента

Request:

``` 
curl --location --request DELETE 'http://localhost:8080/api/v1/segments/AVITO_DISCOUNT_30' \
--data ''
```

Response:

```
200 OK
```

### <a name="get-segments"></a>Получение сегментов пользователя

Request:

``` 
curl --location 'http://localhost:8080/api/v1/users/1/segments'
```

Response:

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

### <a name="edit-segments"></a>Редактирование сегментов пользователя

Request:

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

Response:

```
200 OK
```

### <a name="create-csv"></a>Создать CSV файл с историей добавления/выбывания сегментов

Request:

``` 
curl --location 'http://localhost:8080/api/v1/users/segments/history' \
--header 'Content-Type: application/json' \
--data '{
    "year" : 2023,
    "month" : 8
}'
```

Response:

```json
{
    "link": "./history/user_segments_history-2023-8.csv"
}
```

### <a name="download-csv"></a>Скачать CSV файл с историей добавления/выбывания сегментов

⚠️  Это защищенная ручка, чтобы скачать файл нужно указать токен, получаемый при аутентификации


Request:



``` 
curl --location 'http://localhost:8080/history/user_segments_history-2023-8.csv' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoidGVzdCJ9.Poy32jrZkt3gvJL6uzfE4buxhbrNjTwvTZGPoQntvgk'
```

Response*:

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
*При запросе из браузера будет скачен .csv файл


## Questions

##### Какой первичный ключ использовать для сегментов: суррогатный (id) или естественный (slug)?
Т.к. все операции (поиск сегментов, вставка и удаление) по ТЗ должны выполняться с помощью slug, было решено использовать именно это поле как первичный ключ, потому что он точно останется уникальным и необходимости в суррогатном ключе нет

##### Позволять ли добавлять/удалять пользователю сегменты, если один из них оказался невалидным?
Т.к. в ТЗ указанно, что сегменты пользователя очень важны, было принято решение, что в случае невалидности одного из сегментов, все остальные изменения применяться не будут, о чем клиент будет уведомлен (с помощью кода ответа). В ином случае клиент мог бы не догадываться о том, что какой-либо из сегментов не был добавлен/удален.

##### Как выбирать случайных пользователей для авто-присвоения сегмента?
Изначально был выбран вариант с ORDER BY RAND(), однако в ходе изучения темы выяснилось, что этот подход может быть неоптимальным по сравнению с аналогами. Из-за сложности и неочевидности реализации других способов было решено оставить первоначальный, однако в случае увеличение нагрузки его можно будет легко поменять


## P.S.
Все коммиты в этой ветке были сделаны до дедлайна (31.08), однако у меня в планах увеличить % покрытия тестами. Поэтому, вероятно, до момента проверки я внесу изменения в проект, и если они будут, то уже в другой ветке. Если разрешено вносит изменения после дедлайна, то вы можете посмотреть самую последнюю версию в соседней ветке
