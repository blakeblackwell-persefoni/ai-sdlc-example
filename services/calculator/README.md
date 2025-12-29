# Calculator Service

A simple REST API service providing basic arithmetic operations.

## Architecture

This service follows a clean architecture pattern:

```
services/calculator/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── handler/
│   │   ├── handler.go       # HTTP handlers
│   │   └── handler_test.go  # Handler tests
│   ├── models/
│   │   └── models.go        # Request/response types
│   └── service/
│       ├── calculator.go    # Business logic
│       └── calculator_test.go
├── go.mod
├── openapi.yaml             # API contract
└── README.md
```

## Build and Run

```bash
# From the services/calculator directory
go run cmd/server/main.go

# Or build a binary
go build -o calculator cmd/server/main.go
./calculator
```

The server starts on port 8080 with graceful shutdown support.

## API Endpoints

All arithmetic operations accept POST requests with JSON body:

```json
{"a": number, "b": number}
```

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/add` | POST | Returns a + b |
| `/subtract` | POST | Returns a - b |
| `/multiply` | POST | Returns a * b |
| `/health` | GET | Health check |

### Example

```bash
curl -X POST http://localhost:8080/add \
  -H "Content-Type: application/json" \
  -d '{"a": 10, "b": 5}'

# Response: {"result": 15}
```

## Testing

```bash
# Run all tests
go test ./...

# Run with verbose output
go test ./... -v

# Run with coverage
go test ./... -cover
```

## Features

- Clean architecture with separation of concerns
- Graceful shutdown with configurable timeout
- Request body size limits (1MB)
- Input validation (rejects NaN and Infinity)
- Comprehensive test coverage with table-driven tests
