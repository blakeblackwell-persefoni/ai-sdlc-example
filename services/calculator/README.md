# Calculator Service

A simple REST API service providing basic arithmetic operations, built for Cloudflare Workers.

## Architecture

This service is built with TypeScript and Hono framework for Cloudflare Workers:

```
services/calculator/
├── src/
│   ├── index.ts              # Worker entry point
│   ├── routes/
│   │   └── calculator.ts     # HTTP handlers
│   ├── services/
│   │   └── calculator.ts     # Business logic
│   └── types/
│       └── index.ts          # TypeScript interfaces
├── test/
│   ├── routes/
│   │   └── calculator.test.ts
│   └── services/
│       └── calculator.test.ts
├── wrangler.toml             # Cloudflare Workers config
├── package.json
├── tsconfig.json
├── vitest.config.ts
├── openapi.yaml              # API contract
└── README.md
```

## Development

```bash
# Install dependencies
npm install

# Start local development server (port 8787)
npm run dev

# Run tests
npm test

# Run tests with coverage
npm run test:coverage

# Type check
npm run typecheck
```

## Deployment

```bash
# Deploy to Cloudflare Workers (requires wrangler auth)
npm run deploy
```

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
curl -X POST http://localhost:8787/add \
  -H "Content-Type: application/json" \
  -d '{"a": 10, "b": 5}'

# Response: {"result": 15}
```

## Features

- TypeScript with strict type checking
- Hono framework for fast, lightweight routing
- Cloudflare Workers for edge deployment
- Input validation (rejects NaN and Infinity)
- Comprehensive test coverage with Vitest
