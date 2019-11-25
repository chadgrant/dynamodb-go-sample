FROM golang:1.13 AS builder

ARG application
ARG friendly
ARG build_hash
ARG build_branch
ARG build_user
ARG build_number
ARG build_group

COPY main.go /go/src/github.com/chadgrant/$application/
COPY store /go/src/github.com/chadgrant/$application/store/

WORKDIR /go/src/github.com/chadgrant/$application/

RUN go get ./... && \
    go build -o /go/bin/goapp && \
    echo "application=$application">/go/bin/metadata.txt && \
    echo "friendly=$friendly">>/go/bin/metadata.txt && \
    echo "build_compiler=$(go version)">>/go/bin/metadata.txt && \
    echo "build_date=$(date -u)">>/go/bin/metadata.txt && \
    echo "build_hash=$build_hash">>/go/bin/metadata.txt && \
    echo "build_branch=$build_branch">>/go/bin/metadata.txt && \
    echo "build_group=$build_group">>/go/bin/metadata.txt && \
    echo "build_user=$build_user">>/go/bin/metadata.txt && \   
    echo "build_number=$build_number">>/go/bin/metadata.txt
FROM alpine:3.10.3
RUN apk add --no-cache ca-certificates libc6-compat 
WORKDIR /app
COPY docs /app/docs/
COPY data /app/data/
COPY --from=builder /go/bin/goapp /app/
COPY --from=builder /go/bin/metadata.txt /app/
ENTRYPOINT ./goapp