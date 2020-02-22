# dynamodb-go-sample

[Endpoint Documentation](docs/docs.md) or when running [http://localhost:5000](http://localhost:5000)

## Building Docker images

Using docker-compose:

```bash
make compose-build
```

Using docker (buildkit is faster)

```bash
make docker-build
```

## Running
```bash
make compose-build
make compose-up
curl http://localhost:5000/products/hats
```

## UI
The UI is available at http://localhost:4000

## Running infrastructure locally to work/test against
```bash
make compose-infra
make test
```

## Running Integration Tests in Docker Compose
```bash
make compose-build
make compose-test
```

## Running Integration Tests locally
```bash
make compose-infra

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

## Dynamodb UI
```bash
make compose-up
```

available at http://localhost:8001
