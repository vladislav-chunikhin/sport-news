ifneq ("$(wildcard .env)","")
    include .env
endif

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
deploy: ## Deploy the service with its dependencies. You can deploy custom containers using the 'services' variable. Example: make deploy services='feed-fetcher rabbitmq'
	@echo 'Deploying containers: $(services)'
	docker-compose -f ./deployment/docker-compose.yaml --env-file .env -p ${APP_NAME} up -d $(services)

.PHONY: delete
delete: ## Delete the deployed service
	@docker-compose -f ./deployment/docker-compose.yaml --env-file .env -p ${APP_NAME} rm -sf

.PHONY: test-unit
test-unit: ## Run unit tests
	go test -v -race -short  ./... -coverprofile cover.out
	@go tool cover -func cover.out | grep total | grep -Eo '[0-9]+\.[0-9]+' | xargs -I'{}' echo total '{}'%

.PHONY: coverage
coverage: ## Show coverage
	@go tool cover -html=cover.out

# Delete image
.PHONY: delete-image
delete-image: ## Delete image
	@echo 'Deleting the image'
	docker rmi ${PROJECT_NAME}/${APP_NAME}:${APP_VERSION}

