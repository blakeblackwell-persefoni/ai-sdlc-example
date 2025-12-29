# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is an experimental **Agentic SDLC** repository—a microservices-based calculator API used as a testbed for exploring AI-driven software development workflows. The API provides basic arithmetic operations (add, subtract, multiply) with an OpenAPI specification.

## Repository Structure

```
.
├── services/
│   └── calculator/           # Calculator microservice
│       ├── cmd/server/       # Application entry point
│       ├── internal/
│       │   ├── handler/      # HTTP handlers
│       │   ├── models/       # Request/response types
│       │   └── service/      # Business logic
│       ├── go.mod
│       ├── openapi.yaml      # API contract
│       └── README.md
├── scripts/agents/           # CI/CD automation agents
└── CLAUDE.md
```

## Build and Run Commands

```bash
# Run the calculator service
cd services/calculator
go run cmd/server/main.go

# Build the binary
go build -o calculator cmd/server/main.go

# Run tests
go test ./... -v

# Run tests with coverage
go test ./... -cover
```

The server starts on port 8080 with graceful shutdown support.

## API Endpoints

All arithmetic operations accept POST requests with JSON body `{"a": number, "b": number}`:
- `POST /add` - Addition
- `POST /multiply` - Multiplication
- `POST /subtract` - Subtraction
- `GET /health` - Health check

## Architecture

The calculator service follows clean architecture principles:
- **cmd/server**: Application bootstrap and server configuration
- **internal/handler**: HTTP request handling and routing
- **internal/service**: Core business logic (arithmetic operations)
- **internal/models**: Shared data structures

The service includes:
- Graceful shutdown with configurable timeout
- Request body size limits (1MB)
- Input validation (rejects NaN and Infinity)
- Comprehensive test coverage

## Custom Slash Commands

This repo includes three specialized Go agents as custom commands:

- `/project:go-architect` - Reviews API design, concurrency patterns, error handling, and package structure
- `/project:go-tester` - Generates comprehensive tests with table-driven patterns and HTTP handler tests
- `/project:go-engineer-reviewer` - Conducts thorough code review across quality, security, and performance

## Automated CI/CD Agents

Python scripts in `scripts/agents/` run as GitHub Actions on PRs:

```bash
# Install dependencies
pip install -r scripts/agents/requirements.txt

# Run locally (requires ANTHROPIC_API_KEY)
python scripts/agents/go_architect.py
python scripts/agents/go_tester.py
python scripts/agents/go_engineer_reviewer.py --diff
```

GitHub Actions workflow (`.github/workflows/ai-code-review.yml`) automatically posts review comments on PRs. Requires `ANTHROPIC_API_KEY` secret in repository settings.
