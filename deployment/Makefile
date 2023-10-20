# Include .env file if it exists
ifneq ("$(wildcard .env)","")
    include .env
endif

# Default project variables
PROJECT_NAME ?= sport-news
APP_VERSION ?= local

# Transformer variables
APP_TRANSFORMER_NAME ?= feed-transformer
APP_TRANSFORMER_PORT ?= 8088
APP_TRANSFORMER_MAINTENANCE_PORT ?= 8090

# Fetcher variables
APP_FETCHER_NAME ?= feed-fetcher
APP_FETCHER_PORT ?= 8080
APP_FETCHER_MAINTENANCE_PORT ?= 8082

# Feed service variables
APP_FEED_NAME ?= feed-service
APP_FEED_PORT ?= 8084
APP_FEED_MAINTENANCE_PORT ?= 8086

# MongoDB variables
MONGO_USERNAME ?= vlad
MONGO_PASSWORD ?= sport
MONGO_PORT ?= 27017

# RabbitMQ variables
RABBITMQ_USERNAME ?= vlad
RABBITMQ_PASSWORD ?= sport
RABBITMQ_PORT ?= 5672
RABBITMQ_UI_PORT ?= 15672

# Redis variables
REDIS_PORT ?= 6379

# Docker
DOCKER_COMPOSE := docker compose -f docker-compose.yaml -p ${PROJECT_NAME}

# List of service names
SERVICES := \
    feed-fetcher \
    feed-service \
    feed-transformer

# Help target
.PHONY: help
help: Makefile
	@echo
	@echo "Choose a command to run in $(PROJECT_NAME):"
	@echo
	@awk -F ':.*?## ' '/^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo

# Build images
.PHONY: build-images
build-images: ## Build all images
	@echo 'Building images'
	@for service in $(SERVICES); do \
	    echo "Building $$service"; \
	    cd ../$$service && make -f ../$$service/Makefile build-image; \
	done

# Deploy target
.PHONY: deploy-local
deploy-local: ## Deploy all services locally
	@echo 'Deploying containers'
	$(DOCKER_COMPOSE) up -d

# Delete target
.PHONY: delete-local
delete-local: ## Delete the deployed services locally
	$(DOCKER_COMPOSE) rm -sf

# Set the default target to 'deploy-local'
.DEFAULT_GOAL := deploy-local
