import { describe, it, expect } from "vitest";
import app from "../../src/index";

async function makeRequest(path: string, options?: RequestInit) {
  const request = new Request(`http://localhost${path}`, options);
  return app.fetch(request);
}

describe("Calculator Routes", () => {
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

    it("handles negative numbers", async () => {
      const response = await makeRequest("/add", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ a: -10, b: 5 }),
      });

      expect(response.status).toBe(200);
      const json = await response.json();
      expect(json).toEqual({ result: -5 });
    });

    it("handles decimal numbers", async () => {
      const response = await makeRequest("/add", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ a: 1.5, b: 2.5 }),
      });

      expect(response.status).toBe(200);
      const json = await response.json();
      expect(json).toEqual({ result: 4 });
    });

    it("returns 405 for GET method", async () => {
      const response = await makeRequest("/add", { method: "GET" });

      expect(response.status).toBe(405);
      const json = await response.json();
      expect(json).toEqual({ error: "Method not allowed" });
    });

    it("returns 405 for PUT method", async () => {
      const response = await makeRequest("/add", { method: "PUT" });

      expect(response.status).toBe(405);
    });

    it("returns 405 for DELETE method", async () => {
      const response = await makeRequest("/add", { method: "DELETE" });

      expect(response.status).toBe(405);
    });

    it("returns 400 for invalid JSON", async () => {
      const response = await makeRequest("/add", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: "invalid json",
      });

      expect(response.status).toBe(400);
    });

    it("returns 400 for malformed request body (string instead of number)", async () => {
      const response = await makeRequest("/add", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ a: "text", b: 2 }),
      });

      expect(response.status).toBe(400);
    });

    it("returns 400 for missing field", async () => {
      const response = await makeRequest("/add", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ a: 5 }),
      });

      expect(response.status).toBe(400);
    });

    it("returns 400 for empty body", async () => {
      const response = await makeRequest("/add", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: "",
      });

      expect(response.status).toBe(400);
    });
  });

  describe("POST /subtract", () => {
    it("returns correct difference for valid inputs", async () => {
      const response = await makeRequest("/subtract", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ a: 10, b: 5 }),
      });

      expect(response.status).toBe(200);
      const json = await response.json();
      expect(json).toEqual({ result: 5 });
    });

    it("handles negative result", async () => {
      const response = await makeRequest("/subtract", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ a: 5, b: 10 }),
      });

      expect(response.status).toBe(200);
      const json = await response.json();
      expect(json).toEqual({ result: -5 });
    });

    it("returns 405 for GET method", async () => {
      const response = await makeRequest("/subtract", { method: "GET" });

      expect(response.status).toBe(405);
      const json = await response.json();
      expect(json).toEqual({ error: "Method not allowed" });
    });
  });

  describe("POST /multiply", () => {
    it("returns correct product for valid inputs", async () => {
      const response = await makeRequest("/multiply", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ a: 10, b: 5 }),
      });

      expect(response.status).toBe(200);
      const json = await response.json();
      expect(json).toEqual({ result: 50 });
    });

    it("handles zero", async () => {
      const response = await makeRequest("/multiply", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ a: 10, b: 0 }),
      });

      expect(response.status).toBe(200);
      const json = await response.json();
      expect(json).toEqual({ result: 0 });
    });

    it("handles negative numbers", async () => {
      const response = await makeRequest("/multiply", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ a: -3, b: 4 }),
      });

      expect(response.status).toBe(200);
      const json = await response.json();
      expect(json).toEqual({ result: -12 });
    });

    it("returns 405 for GET method", async () => {
      const response = await makeRequest("/multiply", { method: "GET" });

      expect(response.status).toBe(405);
      const json = await response.json();
      expect(json).toEqual({ error: "Method not allowed" });
    });
  });

  describe("GET /health", () => {
    it("returns ok status", async () => {
      const response = await makeRequest("/health", { method: "GET" });

      expect(response.status).toBe(200);
      const json = await response.json();
      expect(json).toEqual({ status: "ok" });
    });

    it("returns 405 for POST method", async () => {
      const response = await makeRequest("/health", { method: "POST" });

      expect(response.status).toBe(405);
      const json = await response.json();
      expect(json).toEqual({ error: "Method not allowed" });
    });

    it("returns 405 for PUT method", async () => {
      const response = await makeRequest("/health", { method: "PUT" });

      expect(response.status).toBe(405);
    });

    it("returns 405 for DELETE method", async () => {
      const response = await makeRequest("/health", { method: "DELETE" });

      expect(response.status).toBe(405);
    });
  });

  describe("404 Not Found", () => {
    it("returns 404 for unknown endpoints", async () => {
      const response = await makeRequest("/unknown", { method: "GET" });

      expect(response.status).toBe(404);
      const json = await response.json();
      expect(json).toEqual({ error: "Not found" });
    });
  });

  describe("Response headers", () => {
    it("returns application/json content type", async () => {
      const response = await makeRequest("/add", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ a: 1, b: 2 }),
      });

      expect(response.headers.get("content-type")).toContain("application/json");
    });
  });
});
