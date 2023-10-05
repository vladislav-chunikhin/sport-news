ifneq ("$(wildcard .env)","")
    include .env
    export
endif

PROJECT_NAME ?= sport-news
APP_VERSION ?= local

APP_TRANSFORMER_NAME ?= feed-transformer
APP_TRANSFORMER_PORT ?= 8088
APP_TRANSFORMER_MAINTENANCE_PORT ?= 8090

APP_FETCHER_NAME ?= feed-fetcher
APP_FETCHER_PORT ?= 8080
APP_FETCHER_MAINTENANCE_PORT ?= 8082

APP_FEED_NAME ?= feed-service
APP_FEED_PORT ?= 8084
APP_FEED_MAINTENANCE_PORT ?= 8086

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

.PHONY: deploy
deploy: ## Deploy all services
	@echo 'Deploying containers'
	docker compose -f docker-compose.yaml -p ${PROJECT_NAME} up -d

.PHONY: delete
delete: ## Delete the deployed services
	@docker compose -f docker-compose.yaml -p ${PROJECT_NAME} rm -sf
