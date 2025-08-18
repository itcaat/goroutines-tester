# Variables
DOCKER_IMAGE = goroutines-tester
DOCKER_TAG = latest
DOCKER_USERNAME = itcaat
VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT ?= $(shell git rev-parse --short HEAD)
DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

.PHONY: help build run stop clean test docker-build docker-run docker-push monitoring

help: ## Show help
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build Go application
	go build -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)" -o goroutines-tester .

test: ## Run tests
	go test ./...

clean: ## Clean build artifacts
	rm -f goroutines-tester *.out *.test

docker-build: ## Build Docker image
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg DATE=$(DATE) \
		-t $(DOCKER_IMAGE):$(DOCKER_TAG) \
		-t $(DOCKER_IMAGE):$(VERSION) .

docker-run: ## Run application in Docker
	docker-compose up -d goroutines-tester

docker-stop: ## Stop Docker containers
	docker-compose down

docker-logs: ## Show Docker container logs
	docker-compose logs -f goroutines-tester

docker-push: docker-build ## Push image to Docker Hub
	docker tag $(DOCKER_IMAGE):$(DOCKER_TAG) $(DOCKER_USERNAME)/$(DOCKER_IMAGE):$(DOCKER_TAG)
	docker tag $(DOCKER_IMAGE):$(VERSION) $(DOCKER_USERNAME)/$(DOCKER_IMAGE):$(VERSION)
	docker push $(DOCKER_USERNAME)/$(DOCKER_IMAGE):$(DOCKER_TAG)
	docker push $(DOCKER_USERNAME)/$(DOCKER_IMAGE):$(VERSION)

monitoring: ## Run full monitoring stack (app + Prometheus + Grafana)
	docker-compose --profile monitoring up -d

monitoring-stop: ## Stop monitoring
	docker-compose --profile monitoring down

run: build ## Run application locally
	./goroutines-tester -t 10 --metrics

run-debug: build ## Run application locally with debug
	./goroutines-tester -t 5 -d --metrics

# Development commands
dev-setup: ## Setup development environment
	go mod tidy
	go mod download

lint: ## Run linter
	golangci-lint run

format: ## Format code
	go fmt ./...

# Docker utilities
docker-shell: ## Enter container shell
	docker-compose exec goroutines-tester sh

docker-clean: ## Clean Docker resources
	docker-compose down -v --remove-orphans
	docker system prune -f

# Release commands
release: ## Create and push git tag (usage: make release VERSION=v1.0.11)
ifndef VERSION
	$(error VERSION is required. Usage: make release VERSION=v1.0.11)
endif
	@echo "Creating release $(VERSION)..."
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push origin $(VERSION)
	@echo "Release $(VERSION) created and pushed successfully!"

release-patch: ## Create patch release (auto-increment patch version)
	$(eval CURRENT_TAG := $(shell git tag -l | grep "^v[0-9]" | sort -V | tail -1))
	$(eval CURRENT_VERSION := $(if $(CURRENT_TAG),$(CURRENT_TAG),v0.0.0))
	$(eval NEW_VERSION := $(shell echo $(CURRENT_VERSION) | awk -F. '{printf "v%d.%d.%d", substr($$1,2), $$2, $$3+1}'))
	@echo "Current version: $(CURRENT_VERSION)"
	@echo "New version: $(NEW_VERSION)"
	$(MAKE) release VERSION=$(NEW_VERSION)

release-minor: ## Create minor release (auto-increment minor version)
	$(eval CURRENT_TAG := $(shell git tag -l | grep "^v[0-9]" | sort -V | tail -1))
	$(eval CURRENT_VERSION := $(if $(CURRENT_TAG),$(CURRENT_TAG),v0.0.0))
	$(eval NEW_VERSION := $(shell echo $(CURRENT_VERSION) | awk -F. '{printf "v%d.%d.%d", substr($$1,2), $$2+1, 0}'))
	@echo "Current version: $(CURRENT_VERSION)"
	@echo "New version: $(NEW_VERSION)"
	$(MAKE) release VERSION=$(NEW_VERSION)

release-major: ## Create major release (auto-increment major version)
	$(eval CURRENT_TAG := $(shell git tag -l | grep "^v[0-9]" | sort -V | tail -1))
	$(eval CURRENT_VERSION := $(if $(CURRENT_TAG),$(CURRENT_TAG),v0.0.0))
	$(eval NEW_VERSION := $(shell echo $(CURRENT_VERSION) | awk -F. '{printf "v%d.%d.%d", substr($$1,2)+1, 0, 0}'))
	@echo "Current version: $(CURRENT_VERSION)"
	@echo "New version: $(NEW_VERSION)"
	$(MAKE) release VERSION=$(NEW_VERSION)
