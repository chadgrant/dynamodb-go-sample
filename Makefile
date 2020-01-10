APPLICATION?=dynamodb-go-sample
FRIENDLY?=DynamoDB and Go Service
DESCRIPTION?=Sample service using Go and DynamoDB
VENDOR?=Chad Grant
BINARY_NAME?=$(shell basename $(PWD))

REPO_URL?=https://github.com/chadgrant/docker-tools/dynamodb-go-sample
DOCKER_REGISTRY?=docker.io
DOCKER_TAG?=chadgrant/dynamodb-go-sample

BUILD_NUMBER?=1.0.0
BUILD_GROUP?=sample-group
BUILD_BRANCH?=$(shell git rev-parse --abbrev-ref HEAD)
BUILD_HASH?=$(shell git rev-parse HEAD)
BUILD_DATE?=$(shell date -u +%s)
BUILD_USER?=$(USER)

ifdef BUILD_HASH
	BUILD_USER?=$(shell git --no-pager show -s --format='%ae' $(BUILD_HASH))
endif

PKG=github.com/chadgrant/go-http-infra/infra
LDFLAGS="-w -s \
		-X '$(PKG).Application=$(APPLICATION)' \
		-X '$(PKG).Friendly=$(FRIENDLY)' \
		-X '${PKG}.BuildNumber=$(BUILD_NUMBER)' \
		-X '$(PKG).BuiltBy=$(BUILD_USER)' \
		-X '$(PKG).BuiltWhen=$(BUILD_DATE)' \
		-X '$(PKG).GitSha1=$(BUILD_HASH)' \
		-X '$(PKG).GitBranch=$(BUILD_BRANCH)' \
		-X '$(PKG).GroupID=$(BUILD_GROUP)' \
		-X '$(PKG).CompilerVersion=$(shell go version)'"

.PHONY: build
.DEFAULT_GOAL := help
.EXPORT_ALL_VARIABLES:

clean:
	go clean -i
	rm -f $(OUT_DIR)$(BINARY_NAME)
	
build:
	go build -o $(OUT_DIR)$(BINARY_NAME) -ldflags $(LDFLAGS)

test:
	CGO_ENABLED=1 go test -v -race ./...

test-integration:
	CGO_ENABLED=1 TEST_INTEGRATION=1 go test -race -v ./...

tidy:
ifeq (,$(shell type goimports 2>/dev/null))
	go get golang.org/x/tools/cmd/goimports
endif
	go fmt ./...
	goimports -w $(shell go list -f {{.Dir}} ./... | grep -v /vendor/)

lint:
ifeq (,$(shell type golangci-lint 2>/dev/null))
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
		| sh -s -- -b $(shell go env GOPATH)/bin v1.22.2
endif
	golangci-lint run --timeout=300s --skip-dirs-use-default --exclude="should have comment or be unexported"  ./...

docker-build:
	docker-compose build

docker-build-api:
	docker-compose build api

docker-push: docker-build
	docker-compose push api

docker-infra:
	docker-compose up --no-start
	docker-compose start data

docker-infra-api:
	docker-compose up --no-start
	docker-compose start data
	docker-compose start api

docker-run:
	docker-compose up --no-start
	docker-compose start data
	docker-compose up -d

docker-test:
	docker-compose up --no-start
	docker-compose start data
	sleep 5 #wait for infra to come up
	docker-compose run tests

docker-stop:
	-docker container stop `docker container ls -q --filter name=sample_api*`

docker-rm: docker-stop
	-docker container rm `docker container ls -aq --filter name=sample_api*`

docker-clean: docker-rm
	-docker rmi `docker images --format '{{.Repository}}:{{.Tag}}' | grep "${DOCKER_TAG}"` -f
	-docker system prune -f --volumes
