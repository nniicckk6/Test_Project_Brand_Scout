# Brand Scout Quotes

`
Этот проект - тестовое задание от Бренд Скаут ИИ
`

<img width=90% alt="Картинка тестового задания" src="https://github.com/nniicckk6/Test_Project_Brand_Scout/blob/master/image_of_test.png?raw=true">



REST API для хранения и управления цитатами. Поддерживает хранение в памяти или в базе данных SQLite.

## Быстрый старт

1. **Клонируйте репозиторий и перейдите в папку проекта:**

```sh
cd ~/Test_Project_Brand_Scout
```

2. **Установите зависимости:**

```sh
go mod tidy
```

3. **Настройте переменные окружения в файле `.env`:**

Пример содержимого:

```
PORT=8080           # На каком порту запускать сервер
DB_MODE=sqlite      # 'sqlite' — хранить в базе, любое другое значение или пусто — хранить в памяти
DB_PATH=quotes.db   # Путь к файлу SQLite (например, quotes.db или :memory: для тестов)
```

- `PORT` — порт, на котором будет работать HTTP-сервер (по умолчанию 8080).
- `DB_MODE` — режим хранения: `sqlite` для хранения в базе, любое другое значение — хранение в памяти.
- `DB_PATH` — путь к файлу базы данных SQLite. Если не задан, используется quotes.db. Для тестов можно указать `:memory:`.

4. **Запустите сервер:**

```sh
go run main.go
```

Сервер будет доступен по адресу: http://localhost:8080 (или на порту, который вы указали).

## API эндпоинты

- `POST   /quotes` — добавить цитату (JSON: `{ "author": "Автор", "quote": "Текст" }`)
- `GET    /quotes` — получить список всех цитат
- `GET    /quotes/random` — получить случайную цитату
- `GET    /quotes?author=Имя` — получить все цитаты по автору
- `DELETE /quotes/{id}` — удалить цитату по id

## Запуск тестов

Выполните в корне проекта:

```sh
go test ./... -v
```

Будут запущены все unit-тесты для всех режимов (в том числе SQLite с in-memory базой, не затрагивая ваш quotes.db).

## Примечания

- Для работы с SQLite требуется установленная библиотека `github.com/mattn/go-sqlite3` (устанавливается автоматически через go mod tidy).
- Для тестов и CI рекомендуется использовать `DB_PATH=:memory:`.

## Пример .env

```
PORT=8080
DB_MODE=sqlite
DB_PATH=quotes.db
```

