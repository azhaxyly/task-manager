APP_NAME := task-manager
COVERAGE_FILE := coverage.out

.PHONY: all deps fmt swagger test coverage run build

# по дефолту
all: deps fmt swagger test

# зависимости
deps:
	@echo ">> installing dependencies..."
	go mod tidy

# гофампт
fmt:
	@echo ">> go fmt ./..."
	go fmt ./...

# генерация сваггер документации
swagger:
	@echo ">> swag init"
	swag init -g internal/adapter/inbound/http/handler.go -o docs

# запуск всех тестов
test:
	@echo ">> go test ./... -race -coverprofile=$(COVERAGE_FILE)"
	go test ./... -race -coverprofile=$(COVERAGE_FILE)

# процент и отчет покрытия
coverage:
	@echo ">> go tool cover -func=$(COVERAGE_FILE)"
	go tool cover -func=$(COVERAGE_FILE)

# запуск приложения
run:
	@echo ">> go run main.go"
	go run main.go

# собрать бинарник
build:
	@echo ">> go build -o bin/$(APP_NAME) main.go"
	go build -o bin/$(APP_NAME) main.go
