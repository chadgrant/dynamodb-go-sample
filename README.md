# dynamodb-go-sample

[Endpoint Documentation](docs/docs.md) or when running [http://localhost:5000](http://localhost:5000)

## Running
```bash
make docker-run
curl http://localhost:5000/product/hats
```

## Building Docker images
```bash
make docker-build
```

## Running infrastructure locally to work/test against
```bash
make docker-infra
make test
```

## Running Integration Tests in Docker Compose
```bash
make docker-test
```

## Clean / delete docker images
```bash
make docker-clean
```
