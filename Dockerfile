FROM golang:1.13 AS builder

ARG application
ARG friendly
ARG build_hash
ARG build_branch
ARG build_user
ARG build_number
ARG build_group
ENV APPLICATION=$application FRIENDLY=$friendly BUILD_HASH=$build_hash BUILD_BRANCH=$build_branch BUILD_USER=$build_user BUILD_NUMER=$build_number BUILD_GROUP=$build_group

WORKDIR /go/src/github.com/chadgrant/$application/

COPY makefile .
COPY main.go . 
COPY store ./store/

RUN go get ./... && \
    BUILDOUT=/go/bin/goapp make build

FROM alpine:3.10.3
RUN apk add --no-cache ca-certificates libc6-compat 
WORKDIR /app
COPY docs /app/docs/
COPY data /app/data/
COPY --from=builder /go/bin/goapp /app/
ENTRYPOINT ./goapp