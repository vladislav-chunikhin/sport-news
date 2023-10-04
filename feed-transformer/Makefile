ifneq ("$(wildcard .env)","")
    include .env
    export
endif

APP_NAME ?= feed-transformer
PROJECT_NAME ?= sport-news
APP_VERSION ?= local
APP_PORT ?= 8088
APP_MAINTENANCE_PORT ?= 8090

BUILD_DIR ?= "./cmd/app"
OUTPUT_BINARY ?= "build/output/main"
DOCKERFILE_NAME ?= build/Dockerfile

MONGO_USERNAME ?= vlad
MONGO_PASSWORD ?= sport
MONGO_PORT ?= 27017

RABBITMQ_USERNAME ?= vlad
RABBITMQ_PASSWORD ?= sport
RABBITMQ_PORT ?= 5672
RABBITMQ_UI_PORT ?= 15672

REDIS_PORT ?= 6379

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command to run in "$(APP_NAME)":"
	@echo
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@echo

.PHONY: build-image
build-image: ## Build a Docker image of the application
	@echo 'Building the Docker image'
	docker build --no-cache --rm --progress=plain -t $(PROJECT_NAME)/$(APP_NAME):$(APP_VERSION) \
	-f $(DOCKERFILE_NAME) --build-arg OUTPUT_BINARY=$(OUTPUT_BINARY) --build-arg APP_VERSION=$(APP_VERSION) \
	--build-arg APP_NAME=$(APP_NAME) --build-arg BUILD_DIR=$(BUILD_DIR) --build-arg APP_PORT=$(APP_PORT) .

services=
.PHONY: deploy
deploy: ## Deploy the service with its dependencies. You can deploy custom containers using the 'services' variable. Example: make deploy services='rabbitmq mongodb redis'
	@echo 'Deploying containers: $(services)'
	docker-compose -f ./deployment/docker-compose.yaml -p ${APP_NAME} up -d $(services)

.PHONY: delete
delete: ## Delete the deployed service
	@docker-compose -f ./deployment/docker-compose.yaml -p ${APP_NAME} rm -sf

.PHONY: test-unit
test-unit: ## Run unit tests
	go test -v -race -short  ./... -coverprofile cover.out
	@go tool cover -func cover.out | grep total | grep -Eo '[0-9]+\.[0-9]+' | xargs -I'{}' echo total '{}'%
