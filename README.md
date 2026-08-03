# Card Validation API

A Go-based API for validating payment cards.

## Requirements
- [Go 1.26+](https://go.dev/)
- [Docker](https://www.docker.com/) (Optional)
- `make` (Optional)

## Docker Setup

### Build the Image
```bash
docker build -t card-validation-api .
```

### Run the Container
```bash
docker run -p 5000:5000 card-validation-api
```
*Note: The application defaults to port 5000. Change it by passing `--env-file .env`.*

## Make Setup (Local Development)

### Run the Application
```bash
make run
```

### Build the Binary
Compiles the application into the `bin/` directory.
```bash
make build
```

### Clean Build Artifacts
```bash
make clean
```

## Testing

Execute the test suite with race detection and coverage reporting:
```bash
make test
```
*Alternatively, without Make:*
```bash
go test -v -race -cover ./...
```
