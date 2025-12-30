# AI Coding Agent Instructions

## Project Overview
This is an experimental **Agentic SDLC** repository exploring AI-driven software development. The codebase implements a microservices-based calculator API in TypeScript for Cloudflare Workers, serving as a testbed for AI-first development workflows.

## Architecture Patterns

### Clean Architecture Structure
Follow the established clean architecture pattern:
```
services/{service-name}/
├── src/
│   ├── index.ts              # Worker entry point
│   ├── routes/               # HTTP request handling with Hono
│   ├── services/             # Business logic with input validation
│   └── types/                # TypeScript interfaces
├── test/
│   ├── routes/               # HTTP integration tests
│   └── services/             # Unit tests
├── wrangler.toml             # Cloudflare Workers config
├── package.json              # Dependencies and scripts
├── tsconfig.json             # TypeScript configuration
├── vitest.config.ts          # Test runner config
├── openapi.yaml              # API contract definition
└── README.md                 # Service documentation
```

**Key conventions:**
- Use separate directories for routes, services, and types
- Dependency injection: services imported into route handlers
- Type guards for runtime validation of request bodies

### HTTP Server Setup (Hono + Cloudflare Workers)
```typescript
import { Hono } from "hono";
import { calculator } from "./routes/calculator";

const app = new Hono();
app.route("/", calculator);

app.notFound((c) => c.json({ error: "Not found" }, 404));
app.onError((err, c) => c.json({ error: "Internal server error" }, 500));

export default app;
```
- Hono framework for lightweight, fast routing
- Edge deployment on Cloudflare Workers
- Proper error handling with consistent JSON responses

## Testing Patterns

### Table-Driven Unit Tests (Vitest)
```typescript
describe("add", () => {
  it.each([
    { a: 10, b: 5, expected: 15, name: "positive numbers" },
    { a: -10, b: -5, expected: -15, name: "negative numbers" },
    { a: 0, b: 0, expected: 0, name: "zeros" },
  ])("$name: add($a, $b) = $expected", ({ a, b, expected }) => {
    expect(add(a, b)).toBe(expected);
  });

  it("throws InvalidInputError for NaN", () => {
    expect(() => add(NaN, 5)).toThrow(InvalidInputError);
  });
});
```
- Comprehensive edge case coverage (NaN, Infinity validation)
- Use `Number.isFinite()` for input validation

### HTTP Handler Tests
```typescript
async function makeRequest(path: string, options?: RequestInit) {
  const request = new Request(`http://localhost${path}`, options);
  return app.fetch(request);
}

describe("POST /add", () => {
  it("returns correct sum for valid inputs", async () => {
    const response = await makeRequest("/add", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ a: 10, b: 5 }),
    });

    expect(response.status).toBe(200);
    const json = await response.json();
    expect(json).toEqual({ result: 15 });
  });
});
```
- Test both success and error responses
- Validate HTTP status codes and JSON responses

## API Design Principles

### Request/Response Models
```typescript
export interface OperationRequest {
  a: number;
  b: number;
}

export interface OperationResponse {
  result: number;
}

export interface ErrorResponse {
  error: string;
}

export function isOperationRequest(obj: unknown): obj is OperationRequest {
  return (
    typeof obj === "object" &&
    obj !== null &&
    "a" in obj &&
    "b" in obj &&
    typeof (obj as OperationRequest).a === "number" &&
    typeof (obj as OperationRequest).b === "number"
  );
}
```
- Consistent JSON field naming
- Type guards for runtime validation
- Separate error response type

### Endpoint Patterns
- POST for operations: `/add`, `/subtract`, `/multiply`
- GET for health: `/health`
- Input validation rejects NaN and Infinity values
- Return 400 for invalid input, 405 for wrong methods

## Development Workflow

### Build Commands
```bash
# Install dependencies
cd services/calculator && npm install

# Start local development server (port 8787)
npm run dev

# Run tests
npm test

# Run tests with coverage
npm run test:coverage

# Type check
npm run typecheck

# Deploy to Cloudflare Workers
npm run deploy
```

### AI Agent Integration
- Custom Claude commands provide specialized reviews:
  - `/project:node-architect` - Reviews API design, TypeScript patterns, error handling, and module structure
  - `/project:node-tester` - Generates comprehensive tests with Vitest patterns and HTTP handler tests
  - `/project:node-engineer-reviewer` - Conducts thorough code review across quality, security, and performance
- GitHub Actions workflow runs tests on PRs

## Code Quality Standards

### Error Handling
```typescript
export class InvalidInputError extends Error {
  constructor(message: string = "invalid input: NaN and Infinity not allowed") {
    super(message);
    this.name = "InvalidInputError";
  }
}
```
- Custom error classes with proper inheritance
- HTTP handlers return appropriate status codes
- Consistent error response format

### Security & Validation
- Input validation in service layer (not just handlers)
- Type guards for runtime type checking
- Reject invalid numeric inputs (NaN, Infinity)

### TypeScript Best Practices
- Strict mode enabled
- No implicit any
- Proper type imports (`import type`)
- Interface-first design

## Agentic Development Guidelines

This project embraces AI-driven development:
- AI agents should propose architectural changes
- Focus on test-driven development patterns
- Maintain OpenAPI contract accuracy
- Follow established naming and structure conventions
- Prioritize code readability and maintainability

Reference files: [services/calculator/src/services/calculator.ts](services/calculator/src/services/calculator.ts), [services/calculator/src/routes/calculator.ts](services/calculator/src/routes/calculator.ts), [services/calculator/openapi.yaml](services/calculator/openapi.yaml)
