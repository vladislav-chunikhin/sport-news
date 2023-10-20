# Include .env file if it exists
ifneq ("$(wildcard .env)","")
    include .env
endif

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
deploy-local: build-images ## Deploy all services locally (default target)
	@echo 'Deploying containers'
	$(DOCKER_COMPOSE) up -d

# Delete target
.PHONY: delete-local
delete-local: ## Delete the deployed services locally
	$(DOCKER_COMPOSE) rm -sf

# Delete images
.PHONY: delete-images
delete-images: ## Delete images
	@echo 'Deleting images'
	@for service in $(SERVICES); do \
	    docker rmi ${PROJECT_NAME}/$$service:${APP_VERSION}; \
	done

# Set the default target to 'deploy-local'
.DEFAULT_GOAL := deploy-local
