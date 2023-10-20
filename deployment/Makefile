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

# Validate the project structure
.PHONY: validate-structure
validate-structure: ## Validate the existence of directories and files
	@echo 'Validating directories and files'
	@for service in $(SERVICES); do \
		if [ ! -d "../$$service" ]; then \
			echo "Directory ../$$service does not exist"; \
			exit 1; \
		fi; \
		if [ ! -f "../$$service/Makefile" ]; then \
        	echo "File ../$$service/Makefile does not exist"; \
        	exit 1; \
        fi; \
	done
	@echo 'All directories and files exist'

# Build images
.PHONY: build-images
build-images: validate-structure ## Build all images
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
