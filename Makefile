# Переменные
DOCKER_IMAGE = goroutines-tester
DOCKER_TAG = latest
DOCKER_USERNAME = itcaat
VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT ?= $(shell git rev-parse --short HEAD)
DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

.PHONY: help build run stop clean test docker-build docker-run docker-push monitoring

help: ## Показать справку
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Собрать Go приложение
	go build -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)" -o goroutines-tester .

test: ## Запустить тесты
	go test ./...

clean: ## Очистить артефакты сборки
	rm -f goroutines-tester *.out *.test

docker-build: ## Собрать Docker образ
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg DATE=$(DATE) \
		-t $(DOCKER_IMAGE):$(DOCKER_TAG) \
		-t $(DOCKER_IMAGE):$(VERSION) .

docker-run: ## Запустить приложение в Docker
	docker-compose up -d goroutines-tester

docker-stop: ## Остановить Docker контейнеры
	docker-compose down

docker-logs: ## Показать логи Docker контейнера
	docker-compose logs -f goroutines-tester

docker-push: docker-build ## Отправить образ в Docker Hub
	docker tag $(DOCKER_IMAGE):$(DOCKER_TAG) $(DOCKER_USERNAME)/$(DOCKER_IMAGE):$(DOCKER_TAG)
	docker tag $(DOCKER_IMAGE):$(VERSION) $(DOCKER_USERNAME)/$(DOCKER_IMAGE):$(VERSION)
	docker push $(DOCKER_USERNAME)/$(DOCKER_IMAGE):$(DOCKER_TAG)
	docker push $(DOCKER_USERNAME)/$(DOCKER_IMAGE):$(VERSION)

monitoring: ## Запустить полный стек мониторинга (приложение + Prometheus + Grafana)
	docker-compose --profile monitoring up -d

monitoring-stop: ## Остановить мониторинг
	docker-compose --profile monitoring down

run: build ## Запустить приложение локально
	./goroutines-tester -tasks=10 -metrics

run-debug: build ## Запустить приложение локально с debug
	./goroutines-tester -tasks=5 -debug -metrics

# Команды для разработки
dev-setup: ## Настроить окружение для разработки
	go mod tidy
	go mod download

lint: ## Запустить линтер
	golangci-lint run

format: ## Отформатировать код
	go fmt ./...

# Docker утилиты
docker-shell: ## Войти в shell контейнера
	docker-compose exec goroutines-tester sh

docker-clean: ## Очистить Docker ресурсы
	docker-compose down -v --remove-orphans
	docker system prune -f
