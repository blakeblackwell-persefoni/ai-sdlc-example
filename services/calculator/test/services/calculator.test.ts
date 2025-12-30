import { describe, it, expect } from "vitest";
import {
  add,
  subtract,
  multiply,
  validateInputs,
  InvalidInputError,
} from "../../src/services/calculator";

describe("Calculator Service", () => {
  describe("add", () => {
    it.each([
      { a: 10, b: 5, expected: 15, name: "positive numbers" },
      { a: -10, b: -5, expected: -15, name: "negative numbers" },
      { a: 10, b: -5, expected: 5, name: "mixed signs" },
      { a: 0, b: 0, expected: 0, name: "zeros" },
      { a: 1.5, b: 2.5, expected: 4, name: "decimals" },
      { a: 0, b: 5, expected: 5, name: "zero first operand" },
      { a: 5, b: 0, expected: 5, name: "zero second operand" },
      { a: -10, b: 10, expected: 0, name: "opposite signs equal zero" },
      { a: 1e10, b: 1e10, expected: 2e10, name: "large numbers" },
      { a: 1e-10, b: 1e-10, expected: 2e-10, name: "small numbers" },
    ])("$name: add($a, $b) = $expected", ({ a, b, expected }) => {
      expect(add(a, b)).toBe(expected);
    });

    it("throws InvalidInputError for NaN first operand", () => {
      expect(() => add(NaN, 5)).toThrow(InvalidInputError);
    });

    it("throws InvalidInputError for NaN second operand", () => {
      expect(() => add(5, NaN)).toThrow(InvalidInputError);
    });

    it("throws InvalidInputError for positive Infinity", () => {
      expect(() => add(Infinity, 5)).toThrow(InvalidInputError);
    });

    it("throws InvalidInputError for negative Infinity", () => {
      expect(() => add(5, -Infinity)).toThrow(InvalidInputError);
    });
  });

  describe("subtract", () => {
    it.each([
      { a: 10, b: 5, expected: 5, name: "positive numbers" },
      { a: -10, b: -5, expected: -5, name: "negative numbers" },
      { a: 10, b: -5, expected: 15, name: "mixed signs" },
      { a: 0, b: 0, expected: 0, name: "zeros" },
      { a: 2.5, b: 1.5, expected: 1, name: "decimals" },
      { a: 0, b: 5, expected: -5, name: "zero first operand" },
      { a: 5, b: 0, expected: 5, name: "zero second operand" },
      { a: 5, b: 5, expected: 0, name: "equal operands" },
    ])("$name: subtract($a, $b) = $expected", ({ a, b, expected }) => {
      expect(subtract(a, b)).toBe(expected);
    });

    it("throws InvalidInputError for NaN", () => {
      expect(() => subtract(NaN, 5)).toThrow(InvalidInputError);
    });

    it("throws InvalidInputError for Infinity", () => {
      expect(() => subtract(Infinity, 5)).toThrow(InvalidInputError);
    });
  });

  describe("multiply", () => {
    it.each([
      { a: 10, b: 5, expected: 50, name: "positive numbers" },
      { a: -10, b: -5, expected: 50, name: "negative numbers" },
      { a: 10, b: -5, expected: -50, name: "mixed signs" },
      { a: 0, b: 5, expected: 0, name: "zero first operand" },
      { a: 5, b: 0, expected: 0, name: "zero second operand" },
      { a: 1.5, b: 2, expected: 3, name: "decimals" },
      { a: -3, b: 4, expected: -12, name: "negative times positive" },
    ])("$name: multiply($a, $b) = $expected", ({ a, b, expected }) => {
      expect(multiply(a, b)).toBe(expected);
    });

    it("throws InvalidInputError for NaN", () => {
      expect(() => multiply(NaN, 5)).toThrow(InvalidInputError);
    });

    it("throws InvalidInputError for Infinity", () => {
      expect(() => multiply(5, Infinity)).toThrow(InvalidInputError);
    });
  });

  describe("validateInputs", () => {
    it("does not throw for valid inputs", () => {
      expect(() => validateInputs(10, 5)).not.toThrow();
      expect(() => validateInputs(-10, -5)).not.toThrow();
      expect(() => validateInputs(0, 0)).not.toThrow();
      expect(() => validateInputs(1.5, 2.5)).not.toThrow();
    });

    it("throws for NaN first operand", () => {
      expect(() => validateInputs(NaN, 5)).toThrow(InvalidInputError);
    });

    it("throws for NaN second operand", () => {
      expect(() => validateInputs(5, NaN)).toThrow(InvalidInputError);
    });

    it("throws for both NaN", () => {
      expect(() => validateInputs(NaN, NaN)).toThrow(InvalidInputError);
    });

    it("throws for positive Infinity first operand", () => {
      expect(() => validateInputs(Infinity, 5)).toThrow(InvalidInputError);
    });

    it("throws for negative Infinity first operand", () => {
      expect(() => validateInputs(-Infinity, 5)).toThrow(InvalidInputError);
    });

    it("throws for positive Infinity second operand", () => {
      expect(() => validateInputs(5, Infinity)).toThrow(InvalidInputError);
    });

    it("throws for negative Infinity second operand", () => {
      expect(() => validateInputs(5, -Infinity)).toThrow(InvalidInputError);
    });

    it("validates error message is correct", () => {
      expect(() => validateInputs(NaN, 5)).toThrow(
        "invalid input: NaN and Infinity not allowed"
      );
    });
  });
});
