SERVICE?=dynamodb-go-sample
SERVICE_FRIENDLY?=DynamoDB and Go Service
SERVICE_DESCRIPTION?=Sample service using Go and DynamoDB
SERVICE_URL?=http://localhost
VENDOR?=chadgrant
GROUP?=sample-group
BINARY_NAME?=$(shell basename $(PWD))

BUILD_REPO?=https://github.com/chadgrant/dynamodb-go-sample
BUILD_NUMBER?=$(subst v,,$(shell git describe --tags --dirty --match=v* 2> /dev/null || echo 1.0.0))
BUILD_BRANCH?=$(shell git rev-parse --abbrev-ref HEAD)
BUILD_HASH?=$(shell git rev-parse HEAD)
BUILD_DATE?=$(shell TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ')

DOCKER_REGISTRY?=docker.io
DOCKER_IMG?=$(VENDOR)/$(SERVICE)
DOCKER_TAG?=$(BUILD_NUMBER)
DOCKER_TEST_IMG?=$(DOCKER_IMG)-test
DOCKER_UI_IMG?=$(DOCKER_IMG)-ui
DOCKER_DB_IMG?=$(DOCKER_IMG)-db

ifdef BUILD_HASH
 ifndef BUILD_USER
	BUILD_USER?=$(shell git --no-pager show -s --format='%ae' $(BUILD_HASH) 2> /dev/null || echo $(USER))
 endif
else
 BUILD_USER?=$(USER)
endif

PKG=github.com/chadgrant/go-http-infra/infra
LDFLAGS="-w -s \
		-X '$(PKG).Service=$(SERVICE)' \
		-X '$(PKG).Friendly=$(SERVICE_FRIENDLY)' \
		-X '$(PKG).Description=$(SERVICE_DESCRIPTION)' \
		-X '${PKG}.BuildNumber=$(BUILD_NUMBER)' \
		-X '$(PKG).BuiltBy=$(BUILD_USER)' \
		-X '$(PKG).BuiltWhen=$(BUILD_DATE)' \
		-X '$(PKG).GitSha1=$(BUILD_HASH)' \
		-X '$(PKG).GitBranch=$(BUILD_BRANCH)' \
		-X '$(PKG).GroupID=$(GROUP)' \
		-X '$(PKG).CompilerVersion=$(shell go version)'"

.DEFAULT_GOAL := help
.EXPORT_ALL_VARIABLES:

.PHONY: help
help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: build
build: ## Builds go binary
	@go build -o $(OUT_DIR)$(BINARY_NAME) -ldflags $(LDFLAGS)

.PHONY: test
test: ## Run Tests
	@CGO_ENABLED=1 go test -v -race ./...

.PHONY: test-integration
test-integration: # Run integration tests
	@CGO_ENABLED=1 TEST_INTEGRATION=1 go test -race -v ./...
	
clean: ## Cleans directory of temp files
	-@GO111MODULE=off go clean -i
	-@rm -f $(OUT_DIR)$(BINARY_NAME)
	-@rm -f coverage.html profile.out cpu.prof coverage.txt

build-vars: ## Echo's build variables
	@echo "VENDOR=$(VENDOR)"
	@echo "GROUP=$(GROUP)"
	@echo "SERVICE=$(SERVICE)"
	@echo "SERVICE_FRIENDLY=$(SERVICE_FRIENDLY)"
	@echo "SERVICE_DESCRIPTION=$(SERVICE_DESCRIPTION)"
	@echo "SERVICE_URL=$(SERVICE_URL)"
	@echo "BINARY_NAME=$(BINARY_NAME)"
	@echo "DOCKER_REGISTRY=$(DOCKER_REGISTRY)"
	@echo "DOCKER_IMG=$(DOCKER_IMG)"
	@echo "DOCKER_TAG=$(DOCKER_TAG)"
	@echo "BUILD_USER=$(BUILD_USER)"
	@echo "BUILD_REPO=$(BUILD_REPO)"	
	@echo "BUILD_NUMBER=$(BUILD_NUMBER)"
	@echo "BUILD_BRANCH=$(BUILD_BRANCH)"
	@echo "BUILD_HASH=$(BUILD_HASH)"
	@echo "BUILD_DATE=$(BUILD_DATE)"

tidy: ## Run goimports and go fmt on all *.go files
ifeq (,$(shell type goimports 2>/dev/null))
	go get golang.org/x/tools/cmd/goimports
endif
	@go fmt ./...
	@goimports -w $(shell go list -f {{.Dir}} ./... | grep -v /vendor/)

lint: ## Execute linter
ifeq (,$(shell type golangci-lint 2>/dev/null))
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
		| sh -s -- -b $(shell go env GOPATH)/bin v1.22.2
endif
	@golangci-lint run --timeout=300s --skip-dirs-use-default --exclude="should have comment or be unexported"  ./...

reportcard: ## Display go report card status
ifeq (,$(shell type gometalinter 2>/dev/null))
	## has not transitioned to golangci-lint yet
	cd $(GOPATH); curl -L https://git.io/vp6lP | sh
endif
ifeq (,$(shell type goreportcard-cli 2>/dev/null))
	go get github.com/gojp/goreportcard/cmd/goreportcard-cli
endif
	@goreportcard-cli

cover: ## Run go code coverage tool
	@go test -covermode=atomic -coverpkg=./... -coverprofile=profile.out ./...
	@go tool cover -func=profile.out
	@go tool cover -html=profile.out -o coverage.html

compose-build: ## Build containers with docker-compose
	@docker-compose build

compose-build-api: ## Build only the API container with docker-compose
	@docker-compose build api

compose-push: docker-build ## Push the docker images with docker-compose
	@docker-compose push api

compose-infra: ## Start the docker infrastructure databases, etc with docker-compose
	@docker-compose up --no-start
	@docker-compose start data

compose-infra-api: ## Start the docker infrastructure databases and the API etc with docker-compose
	@docker-compose up --no-start
	@docker-compose start data
	@docker-compose start api

compose-up: ## Start all docker containers in docker-compose.yml
	@docker-compose up --no-start
	@docker-compose start data
	@docker-compose up -d

compose-test:  ## Run tests in docker containers with docker-compose
	@docker-compose up --no-start
	@docker-compose start data
	@sleep 5 #wait for infra to come up
	@docker-compose run tests

docker-stop: ## Stop all containers
	-@docker container stop `docker container ls -q --filter name=sample_api*`

docker-rm: docker-stop ## Stop and remove all containers
	-@docker container rm `docker container ls -aq --filter name=sample_api*`

docker-clean: docker-rm ## Stop, remove all containers and remove images
	-@docker rmi `docker images --format '{{.Repository}}:{{.Tag}}' | grep "${DOCKER_IMG}"` -f
	-@docker system prune -f --volumes

docker-all: docker docker-tests docker-db docker-ui ## Builds all the docker images without docker compose (faster)

docker: DOCKER_FILE=Dockerfile
docker: IMG=$(DOCKER_IMG)
docker: docker-internal-normal ## Builds the api docker image without docker compose (faster)

docker-tests: DOCKER_FILE=Dockerfile.tests
docker-tests: IMG=$(DOCKER_TEST_IMG)
docker-tests: docker-internal-tests ## Builds the test docker image without docker compose (faster)

docker-ui: CTX=ui/
docker-ui: DOCKER_FILE=ui/Dockerfile
docker-ui: IMG=$(DOCKER_UI_IMG)
docker-ui: docker-internal-ui ## Builds the ui docker image without docker compose (faster)

docker-db: CTX=db/
docker-db: DOCKER_FILE=$(CTX)/Dockerfile
docker-db: IMG=$(DOCKER_DB_IMG)
docker-db: docker-internal-db ## Builds the data docker image without docker compose (faster)

docker-internal-%:
	@BUILDKIT=1 docker build -f $(DOCKER_FILE) \
		--build-arg "VENDOR=$(VENDOR)" \
		--build-arg "GROUP=$(GROUP)" \
		--build-arg "SERVICE=$(SERVICE)" \
		--build-arg "SERVICE_FRIENDLY=$(SERVICE_FRIENDLY)" \
		--build-arg "SERVICE_URL=$(SERVICE_URL)" \
		--build-arg "SERVICE_DESCRIPTION=$(SERVICE_DESCRIPTION)" \
		--build-arg "BUILD_NUMBER=$(BUILD_NUMBER)" \		
		--build-arg "BUILD_USER=$(BUILD_USER)" \
		--build-arg "BUILD_DATE=$(BUILD_DATE)" \
		--build-arg "BUILD_BRANCH=$(BUILD_BRANCH)" \
		--build-arg "BUILD_HASH=${BUILD_HASH}" \
		--build-arg "BUILD_REPO=$(BUILD_REPO)" \
		--tag $(DOCKER_REGISTRY)/$(IMG):$(DOCKER_TAG) $(CTX).

docker-push-all: | docker-all docker-push docker-push-ui

docker-push: IMG=$(DOCKER_IMG)
docker-push: docker-push-internal-normal ## Builds/Pushes the api docker image without docker compose (faster)

docker-push-ui: IMG=$(DOCKER_IMG)
docker-push-ui: docker-push-internal-ui ## Builds/Pushes the ui docker image without docker compose (faster)

docker-push-internal-%:
	@BUILDKIT=1 docker push $(DOCKER_REGISTRY)/$(IMG):$(DOCKER_TAG)