export class InvalidInputError extends Error {
  constructor(
    message: string = "invalid input: NaN and Infinity not allowed"
  ) {
    super(message);
    this.name = "InvalidInputError";
  }
}

export function validateInputs(a: number, b: number): void {
  if (!Number.isFinite(a) || !Number.isFinite(b)) {
    throw new InvalidInputError();
  }
}

export function add(a: number, b: number): number {
  validateInputs(a, b);
  return a + b;
}

export function subtract(a: number, b: number): number {
  validateInputs(a, b);
  return a - b;
}

export function multiply(a: number, b: number): number {
  validateInputs(a, b);
  return a * b;
}
