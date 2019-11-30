ifndef BUILD_GROUP
	BUILD_GROUP="sample-group"
endif

ifndef BUILD_NUMBER
	BUILD_NUMBER="1.0.0"
endif

ifndef BRANCH
	BRANCH=$(git symbolic-ref -q HEAD)
	BRANCH=${BRANCH##refs/heads/}
	BRANCH=${BRANCH:-HEAD}
endif

ifndef HASH
	HASH=$(git rev-parse HEAD)
endif

ifndef BUILD_USER
	BUILD_USER=$(git --no-pager show -s --format='<mailto:%ae|%an>' $HASH)
endif

.PHONY: help

help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

clean:
	rm dynamodb-go-sample

get:
	go get -u ./...

build: get
	go build

test: get
	go test ./... -v

docker-build:
	docker-compose build

docker-push: docker-build
	docker-compose push api

docker-run:
	docker-compose up -d

docker-test: docker-run
	#optional sleep 15 #wait for infra to come up
	docker-compose run tests

docker-clean:
	docker container stop `docker container ls -q --filter name=sample_api*`
	docker container rm `docker container ls -aq --filter name=sample_api*`
	docker rmi `docker images --format '{{.Repository}}:{{.Tag}}' | grep "chadgrant/sample"` -f
	docker rmi `docker images -qf dangling=true`
	#docker volume rm `docker volume ls -qf dangling=true`