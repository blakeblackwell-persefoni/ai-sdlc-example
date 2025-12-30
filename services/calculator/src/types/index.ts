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

export interface HealthResponse {
  status: string;
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
