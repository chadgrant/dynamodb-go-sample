ARG builder_img="golang:1.13.5-alpine"
ARG runtime_img="alpine:3.11.2"

FROM $builder_img AS builder
ENV CGO_ENABLED=0
RUN apk add --update --no-cache make

ARG application
ARG friendly
ARG build_hash
ARG build_branch
ARG build_user
ARG build_number
ARG build_group
ENV APPLICATION=$application FRIENDLY=$friendly BUILD_HASH=$build_hash BUILD_BRANCH=$build_branch BUILD_USER=$build_user BUILD_NUMER=$build_number BUILD_GROUP=$build_group

WORKDIR /go/src/github.com/chadgrant/$application/

COPY go.mod go.sum ./
RUN go mod download

COPY Makefile main.go ./
COPY store ./store/

RUN BUILDOUT=/go/bin/goapp make build

FROM $runtime_img
RUN apk add --no-cache ca-certificates libc6-compat 
WORKDIR /app
COPY docs /app/docs/
COPY --from=builder /go/bin/goapp /app/
ENTRYPOINT ./goapp