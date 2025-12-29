# AI Coding Agent Instructions

## Project Overview
This is an experimental **Agentic SDLC** repository exploring AI-driven software development. The codebase implements a microservices-based calculator API in Go, serving as a testbed for AI-first development workflows.

## Architecture Patterns

### Clean Architecture Structure
Follow the established clean architecture pattern:
```
services/{service-name}/
├── cmd/server/main.go          # Application entry point with graceful shutdown
├── internal/
│   ├── handler/                # HTTP request handling and routing
│   ├── models/                 # Request/response types with JSON tags
│   └── service/                # Business logic with input validation
├── go.mod                      # Minimal dependencies (Go 1.21+)
├── openapi.yaml                # API contract definition
└── README.md                   # Service documentation
```

**Key conventions:**
- Use `internal/` packages for private APIs
- Dependency injection: `service.NewCalculator()` → `handler.NewHandler(calc)`
- Generic handler pattern: `handleOperation` function reduces duplication across endpoints

### HTTP Server Setup
```go
server := &http.Server{
    Addr:         ":8080",
    Handler:      mux,
    ReadTimeout:  5 * time.Second,
    WriteTimeout: 10 * time.Second,
    IdleTimeout:  120 * time.Second,
}
```
- Graceful shutdown with 30-second timeout
- Request body size limit: 1MB
- Method validation and proper HTTP status codes

## Testing Patterns

### Table-Driven Unit Tests
```go
tests := []struct {
    name    string
    a, b    float64
    want    float64
    wantErr bool
}{
    {"positive numbers", 10, 5, 15, false},
    {"NaN input", math.NaN(), 5, 0, true},
}
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) { /* test logic */ })
}
```
- Comprehensive edge case coverage (NaN, Infinity validation)
- Use `math.IsNaN()` and `math.IsInf()` for input validation

### HTTP Handler Tests
```go
req := httptest.NewRequest(http.MethodPost, "/add", bytes.NewReader(body))
rec := httptest.NewRecorder()
mux.ServeHTTP(rec, req)
// Assert status codes and JSON responses
```
- Test both success and error responses
- Validate JSON marshaling/unmarshaling

## API Design Principles

### Request/Response Models
```go
type OperationRequest struct {
    A float64 `json:"a"`
    B float64 `json:"b"`
}

type OperationResponse struct {
    Result float64 `json:"result"`
}

type ErrorResponse struct {
    Error string `json:"error"`
}
```
- Consistent JSON field naming
- Separate error response type

### Endpoint Patterns
- POST for operations: `/add`, `/subtract`, `/multiply`
- GET for health: `/health`
- Input validation rejects NaN and Infinity values
- Return 400 for invalid input, 405 for wrong methods

## Development Workflow

### Build Commands
```bash
# Run service
cd services/calculator && go run cmd/server/main.go

# Test with coverage
go test ./... -cover

# Build binary
go build -o calculator cmd/server/main.go
```

### AI Agent Integration
- Python agents in `scripts/agents/` provide automated code review
- Requires `ANTHROPIC_API_KEY` for CI/CD automation
- Agents review: architecture, testing, code quality
- GitHub Actions workflow posts review comments on PRs

## Code Quality Standards

### Error Handling
- Custom error types: `ErrInvalidInput`
- Proper error wrapping and propagation
- HTTP handlers return appropriate status codes

### Security & Validation
- Input validation in service layer (not just handlers)
- Request body size limits
- Reject invalid numeric inputs (NaN, Infinity)

### Performance
- Minimal allocations in hot paths
- Reasonable timeouts and limits
- Clean separation of concerns for maintainability

## Agentic Development Guidelines

This project embraces AI-driven development:
- AI agents should propose architectural changes
- Focus on test-driven development patterns
- Maintain OpenAPI contract accuracy
- Follow established naming and structure conventions
- Prioritize code readability and maintainability

Reference files: [services/calculator/internal/service/calculator.go](services/calculator/internal/service/calculator.go), [services/calculator/internal/handler/handler.go](services/calculator/internal/handler/handler.go), [services/calculator/openapi.yaml](services/calculator/openapi.yaml)</content>
<parameter name="filePath">/workspaces/ai-sdlc-example/.github/copilot-instructions.md