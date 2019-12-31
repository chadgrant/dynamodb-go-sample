# dynamodb-go-sample

[Endpoint Documentation](docs/docs.md) or when running [http://localhost:5000](http://localhost:5000)

## Building Docker images
```bash
make docker-build
```

## Running
```bash
make docker-build
make docker-run
curl http://localhost:5000/product/hats
```

## UI ..
The UI is available at http://localhost:4000

## Running infrastructure locally to work/test against
```bash
make docker-infra
make test
```

## Running Integration Tests in Docker Compose
```bash
make docker-build
make docker-test
```

## Running Integration Tests locally
```bash
make docker-infra

export TEST_INTEGRATION=1
go test -v ./...
```

## Clean / delete docker images
```bash
make docker-clean
```

## Building locally
```bash
make get
make build
```