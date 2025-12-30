import { Hono } from "hono";
import type { Context } from "hono";
import type { ContentfulStatusCode } from "hono/utils/http-status";
import {
  add,
  subtract,
  multiply,
  InvalidInputError,
} from "../services/calculator";
import type {
  OperationRequest,
  OperationResponse,
  ErrorResponse,
  HealthResponse,
} from "../types";
import { isOperationRequest } from "../types";

const calculator = new Hono();

async function parseOperationRequest(c: Context): Promise<OperationRequest> {
  const body = await c.req.json();
  if (!isOperationRequest(body)) {
    throw new Error("Invalid request body");
  }
  return body;
}

function errorResponse(c: Context, status: ContentfulStatusCode, message: string) {
  const error: ErrorResponse = { error: message };
  return c.json(error, status);
}

calculator.post("/add", async (c) => {
  try {
    const { a, b } = await parseOperationRequest(c);
    const result = add(a, b);
    const response: OperationResponse = { result };
    return c.json(response);
  } catch (error) {
    if (error instanceof InvalidInputError) {
      return errorResponse(c, 400, error.message);
    }
    return errorResponse(c, 400, "Invalid request");
  }
});

calculator.post("/subtract", async (c) => {
  try {
    const { a, b } = await parseOperationRequest(c);
    const result = subtract(a, b);
    const response: OperationResponse = { result };
    return c.json(response);
  } catch (error) {
    if (error instanceof InvalidInputError) {
      return errorResponse(c, 400, error.message);
    }
    return errorResponse(c, 400, "Invalid request");
  }
});

calculator.post("/multiply", async (c) => {
  try {
    const { a, b } = await parseOperationRequest(c);
    const result = multiply(a, b);
    const response: OperationResponse = { result };
    return c.json(response);
  } catch (error) {
    if (error instanceof InvalidInputError) {
      return errorResponse(c, 400, error.message);
    }
    return errorResponse(c, 400, "Invalid request");
  }
});

calculator.get("/health", (c) => {
  const response: HealthResponse = { status: "ok" };
  return c.json(response);
});

// Handle wrong HTTP methods
calculator.all("/add", (c) => errorResponse(c, 405, "Method not allowed"));
calculator.all("/subtract", (c) => errorResponse(c, 405, "Method not allowed"));
calculator.all("/multiply", (c) => errorResponse(c, 405, "Method not allowed"));
calculator.all("/health", (c) => errorResponse(c, 405, "Method not allowed"));

export { calculator };
