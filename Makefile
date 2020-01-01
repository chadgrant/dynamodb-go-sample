APPLICATION?=application_name
FRIENDLY?=Friendly Name
BUILD_NUMBER?=1.0.0
BUILD_GROUP?=sample-group
BUILD_BRANCH?=$(shell git rev-parse --abbrev-ref HEAD)
BUILD_HASH?=$(shell git rev-parse HEAD)
BUILD_DATE?=$(shell date -u +%s)

ifdef BUILD_HASH
	BUILD_USER?=$(shell git --no-pager show -s --format='%ae' $(BUILD_HASH))
endif

ifdef BUILDOUT
	OUTPUT=-o ${BUILDOUT}
endif

PKG=github.com/chadgrant/go-http-infra
LDFLAGS += -X '${PKG}.Application=${APPLICATION}'
LDFLAGS += -X '${PKG}.Friendly=${FRIENDLY}'
LDFLAGS += -X '${PKG}.BuildNumber=${BUILD_NUMBER}'
LDFLAGS += -X '${PKG}.BuiltBy=${BUILD_USER}'
LDFLAGS += -X '${PKG}.BuiltWhen=${BUILD_DATE}'
LDFLAGS += -X '${PKG}.GitSha1=${BUILD_HASH}'
LDFLAGS += -X '${PKG}.GitBranch=${BUILD_BRANCH}'
LDFLAGS += -X '${PKG}.GroupID=${BUILD_GROUP}'
LDFLAGS += -X '${PKG}.CompilerVersion=$(shell go version)'

.PHONY: build
.DEFAULT_GOAL := help
.EXPORT_ALL_VARIABLES:

clean:
	-rm dynamodb-go-sample

get:
	go get -u ./...

build:
	go build ${OUTPUT} -ldflags "-s ${LDFLAGS}"

test: get
	go test ./... -v

docker-build:
	docker-compose build

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
	-docker rmi `docker images --format '{{.Repository}}:{{.Tag}}' | grep "chadgrant/sample"` -f
	-docker rmi `docker images -qf dangling=true`
	#-docker volume rm `docker volume ls -qf dangling=true`
