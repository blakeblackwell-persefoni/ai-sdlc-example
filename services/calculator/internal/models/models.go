package models

// OperationRequest represents a request for an arithmetic operation.
type OperationRequest struct {
	A float64 `json:"a"`
	B float64 `json:"b"`
}

// OperationResponse represents the result of an arithmetic operation.
type OperationResponse struct {
	Result float64 `json:"result"`
}

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Error string `json:"error"`
}

// HealthResponse represents a health check response.
type HealthResponse struct {
	Status string `json:"status"`
}
