# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is an experimental **Agentic SDLC** repository—a microservices-based calculator API used as a testbed for exploring AI-driven software development workflows. The API provides basic arithmetic operations (add, subtract, multiply) with an OpenAPI specification, built for Cloudflare Workers.

## Repository Structure

```
.
├── services/
│   └── calculator/           # Calculator microservice (Cloudflare Workers)
│       ├── src/
│       │   ├── index.ts      # Worker entry point
│       │   ├── routes/       # HTTP handlers
│       │   ├── services/     # Business logic
│       │   └── types/        # TypeScript interfaces
│       ├── test/             # Vitest tests
│       ├── wrangler.toml     # Cloudflare Workers config
│       ├── package.json
│       ├── openapi.yaml      # API contract
│       └── README.md
├── scripts/agents/           # CI/CD automation agents
└── CLAUDE.md
```

## Build and Run Commands

```bash
# Install dependencies
cd services/calculator
npm install

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

## API Endpoints

All arithmetic operations accept POST requests with JSON body `{"a": number, "b": number}`:
- `POST /add` - Addition
- `POST /multiply` - Multiplication
- `POST /subtract` - Subtraction
- `GET /health` - Health check

## Architecture

The calculator service is built with TypeScript and Hono framework:
- **src/index.ts**: Worker entry point and app configuration
- **src/routes/**: HTTP request handling with Hono
- **src/services/**: Core business logic (arithmetic operations)
- **src/types/**: TypeScript interfaces

The service includes:
- TypeScript with strict type checking
- Hono framework for fast, lightweight routing
- Cloudflare Workers for edge deployment
- Input validation (rejects NaN and Infinity)
- Comprehensive test coverage with Vitest

## Custom Slash Commands

This repo includes three specialized Node/TypeScript agents as custom commands:

- `/project:node-architect` - Reviews API design, TypeScript patterns, error handling, and module structure
- `/project:node-tester` - Generates comprehensive tests with Vitest patterns and HTTP handler tests
- `/project:node-engineer-reviewer` - Conducts thorough code review across quality, security, and performance

## Automated CI/CD

GitHub Actions workflow (`.github/workflows/ai-code-review.yml`) automatically runs tests on PRs:

```bash
# CI runs these commands on each PR
npm ci
npm run typecheck
npm test
```

The workflow triggers on changes to `services/calculator/**` files.
