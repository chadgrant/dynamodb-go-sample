FROM golang:1.13 AS builder

COPY . /go/src/github.com/chadgrant/dynamodb-go-sample/

WORKDIR /go/src/github.com/chadgrant/dynamodb-go-sample/

RUN go get ./... && \
    go build -o /go/bin/app

FROM alpine:3.10.3
WORKDIR /app
COPY --from=builder /go/bin/app /app/
ENTRYPOINT ./app