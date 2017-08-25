## Описание API

Сервис позволяет делать запросы к базе данных Postgres удобным способом, возвращая ответы с помощью JSON.

`GET /databases/`

#### Возвращает список баз данных

```
[
  "heroes",
  "villains"
]
```

`GET /databases/heroes/`

#### Возвращает список таблиц в базе данных

```
[
  "victories",
  "birthdays",
  "couples"
]
```

`GET /databases/heroes/victories`

#### Возвращает всю таблицу

```
[
  {
    "fight_identifier" : 0,
    "first_fighter_identifier" : 12,
    "second_fighter_identifier" : 20,
    "first_fighter_score" : 100,
    "second_fighter_score" : 500
  },
  {
    ...
  },
  ...
]
```

`GET /databases/heroes/victories/0`

#### Возвращает элемент таблицы по идентификатору

```
{
  "fight_identifier" : 0,
  "first_fighter_identifier" : 12,
  "second_fighter_identifier" : 20,
  "first_fighter_score" : 100,
  "second_fighter_score" : 500
}
```

#### Запрос элементов таблицы с условием

Поддерживаемые условия: >, <, >=, <=, ==

`GET /database/heroes/victories/?first_fighter_score>200&second_fighter_score<=400`

```
[
  {
    "fight_identifier" : 0,
    "first_fighter_identifier" : 12,
    "second_fighter_identifier" : 20,
    "first_fighter_score" : 100,
    "second_fighter_score" : 500
  }
]
```

#### Агрегирующий запрос

`GET /databases/heroes/?`

#### Добавление элемента в таблицу

`PUT /databases/heroes/victories/?first_fighter_identifier=12...`

```
{
  "success" : true
}
```

#### Удаление элемента из таблицы

`DELETE /databases/heroes/victories/0`

```
{
  "success" : true
}
```

#### Создание таблицы

`PUT /databases/heroes/failures/?first_fighter_identifier=int&second_fighter_score=int...`

```
{
  "success" : true
}
```

#### Создание базы данных

`PUT /databases/daydreamers/`

```
{
  "success" : true
}
```

#### Удаление таблицы
`DELETE /databases/heroes/failures/`

```
{
  "success" : true
}
```

#### Удаление базы данных
`DELETE /databases/daydreamers/`

```
{
  "success" : true
}
```

## Следующие шаги

* Добавить описание агрегирующего запроса
* Добавить поддержку авторизации пользователя с выдачей токена
