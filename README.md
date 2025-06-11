# Task Manager

![Go](https://img.shields.io/badge/Go-1.23.4-blue) ![PostgreSQL](https://img.shields.io/badge/Swagger-16-green)

HTTP API-сервис для создания, мониторинга и удаления задач с хранением состояния **в памяти**.

---

## Содержание

1. [Обзор](#обзор)  
2. [Требования](#требования)  
3. [Установка](#установка)  
4. [Запуск сервиса](#запуск-сервиса)  
5. [Makefile (команды)](#makefile-команды)  
6. [API-эндпоинты](#api-эндпоинты)  
7. [Swagger (OpenAPI)](#swagger-openapi)  
8. [Тестирование](#тестирование)  
9. [Структура проекта](#структура-проекта)  
10. [Автор](#автор)


---

## Обзор

- **Парадигма**: DDD + Hexagonal (Clean) Architecture   
- **Хранение задач**: in-memory (map + mutex)  
- **Асинхронность**: goroutine + context.WithCancel  
- **HTTP-сервер**: net/http + ручной mux  

## Требования

- Go 1.23+  
- make (рекомендуется)  
- swag (для генерации Swagger-документации)   

---

## Установка

1. Клонировать репозиторий:
   ```bash
   git clone <URL-репозитория>
   cd task-manager
   ```
2. Привести зависимости:
   ```bash
   go mod tidy
   ```
3. (опционально) Установить `swag`:
   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

---

## Запуск сервиса

```bash
make run
```

или без Makefile:

```bash
go run main.go
```

Сервис стартует на порту `:8080`.

---

## Makefile (команды)

```makefile
all: 
    deps fmt swagger test
deps:
	go mod tidy
fmt:
	go fmt ./...
swagger:
	swag init -g internal/adapter/inbound/http/handler.go -o docs
test:
	go test ./... -race -coverprofile=coverage.out
coverage:
	go tool cover -func=coverage.out
run:
	go run main.go
build:
	go build -o bin/task-manager main.go
```

---

## API-эндпоинты

| Метод    | Путь          | Описание                                         | Код  | Ответ                |
| -------- | ------------- | ------------------------------------------------ | ---- | -------------------- |
| POST     | /tasks        | Создать новую задачу                             | 201  | CreateTaskResponse   |
| GET      | /tasks        | Список всех задач (ID, статус, createdAt, duration) | 200 | []TaskSummaryDTO     |
| GET      | /tasks/{id}   | Подробный статус задачи                          | 200  | TaskDTO              |
| DELETE   | /tasks/{id}   | Отменить или удалить задачу                      | 204  | —                    |

---

## Swagger (OpenAPI)

Сгенерировать:
```bash
make swagger
```
Запустить UI:
```go
router.Handle("/swagger/", httpSwagger.WrapHandler)
```
Открыть: `http://localhost:8080/swagger/index.html`

---

## Тестирование

```bash
make test
make coverage
```

Покрытие ≥ ~60%.

---

## Структура проекта

```
.
├── cmd/               # реальная точка входа
├── docs/              # Swagger-документация
├── internal/
│   ├── adapter/       # HTTP, memstore, idgen
│   ├── application/   # ports, services
│   ├── common/logger/
│   └── domain/        # Task, Status, Repo interface
├── main.go            # точка входа, для обхода путей
├── Makefile
├── go.mod
└── README.md
```

---

## Автор

- Жаксылыков Азамат

- <a href="https://t.me/hmlssdeus" target="_blank"><img src="https://img.shields.io/badge/telegram-@hmlssdeus-blue?logo=Telegram" alt="Status" /></a>
