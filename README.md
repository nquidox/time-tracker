# Тестовое задание
## Тайм треккер

Создать rest api для подсчета трудозатрат пользователей. Включает в себя:
1) CRUD для пользователей, фильтрация по всем полям и пагинация.
2) Получение части данных о новом пользователе от другого API.
3) CRUD для задач, плюс методы для старта и окончания выполнения задачи.
4) Вывод сводки выполненных задач для пользователя с сортировкой по убыванию.
5) Создать таблицы в базе данных путем миграции.
6) Добавить логгирование, сваггер, параметры запуска вынести в .env

Примечания к реализации в конце описания.

## Запуск
Создать .env файл и заполнить следующие значения:
```sh
# DB params
DB_HOST=localhost
DB_PORT=5432
DB_USER=timetracker
DB_PASSWORD=timetracker
DB_NAME=timetracker
DB_SSLMODE=disable

# HTTP Server params
HTTP_HOST=localhost
HTTP_PORT=9000

# External API params
EXTERNAL_API_URL=http://localhost:9001

#Log levels
APP_LOG_LEVEL=Info
DB_LOG_LEVEL=Silent
```
Для получения данных с external API можете использовать сервис https://github.com/nquidox/mock-api

## Лицензия
```
MIT
```

## Примечания

- В проекте подразумевается, что серия + номер паспорта уникальны и не могут принадлежать разным людям.
- Серия + номер могут быть повторно присвоены другому пользователю только в случае если аналогичная пара ранее принадлежала удаленному из БД пользователю.
- Механизм авторизации отсутствует, в запросах, где требуется указать принадлежность задачи конкретному пользователю, нужно ввести uuid этого пользователя вручную.
- Если параметры пагинации не указаны или указаны некорректно, принимаются значения по умолчанию: page=1, perPage=10.
- Параметры GET запросов не валидируются, в случае некорректных значений будут приняты значения по умолчанию (если есть), либо сервер вернет ответ 404.
- Намеренно допускаются одинаковые имена задач (Title).
- Время, затраченное на выполнение задачи, подситывается только в момент запроса finish, механизм расчета промежуточных значений, например, пауза, не реализован.
- Сортировка длительности трудозатрат происходит по полю duration по убыванию.
- При указании дат периода указываются только дни в формате дд-мм-гггг. Время при этом нулевое, поэтому для того, чтобы вывести данные о задачах по текущий день включительно, нужно указать конец периода на 1 день больше. По умолчанию выводятся задачи за все время.


