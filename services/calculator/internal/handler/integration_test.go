package handler

import (
	"bytes"
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blakeblackwell-persefoni/ai-sdlc-example/services/calculator/internal/models"
	"github.com/blakeblackwell-persefoni/ai-sdlc-example/services/calculator/internal/service"
)

// TestIntegration_CompleteWorkflow tests a complete workflow across multiple endpoints
func TestIntegration_CompleteWorkflow(t *testing.T) {
	h := setupHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	// Step 1: Health check
	t.Run("health check", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("health check failed: status = %d", rec.Code)
		}

		var health models.HealthResponse
		json.NewDecoder(rec.Body).Decode(&health)
		if health.Status != "ok" {
			t.Fatalf("health status = %s, want ok", health.Status)
		}
	})

	// Step 2: Perform addition
	var addResult float64
	t.Run("addition operation", func(t *testing.T) {
		body, _ := json.Marshal(models.OperationRequest{A: 100, B: 50})
		req := httptest.NewRequest(http.MethodPost, "/add", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("addition failed: status = %d", rec.Code)
		}

		var resp models.OperationResponse
		json.NewDecoder(rec.Body).Decode(&resp)
		addResult = resp.Result
		if addResult != 150 {
			t.Fatalf("addition result = %v, want 150", addResult)
		}
	})

	// Step 3: Use result in subtraction
	t.Run("subtraction using previous result", func(t *testing.T) {
		body, _ := json.Marshal(models.OperationRequest{A: addResult, B: 25})
		req := httptest.NewRequest(http.MethodPost, "/subtract", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("subtraction failed: status = %d", rec.Code)
		}

		var resp models.OperationResponse
		json.NewDecoder(rec.Body).Decode(&resp)
		if resp.Result != 125 {
			t.Fatalf("subtraction result = %v, want 125", resp.Result)
		}
	})
}

// TestIntegration_ErrorHandling tests error handling across the stack
func TestIntegration_ErrorHandling(t *testing.T) {
	h := setupHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	tests := []struct {
		name         string
		endpoint     string
		method       string
		body         interface{}
		wantStatus   int
		wantErrorMsg string
	}{
		{
			name:         "invalid method on add",
			endpoint:     "/add",
			method:       http.MethodGet,
			body:         nil,
			wantStatus:   http.StatusMethodNotAllowed,
			wantErrorMsg: "Method not allowed",
		},
		{
			name:       "NaN value propagates error",
			endpoint:   "/multiply",
			method:     http.MethodPost,
			body:       models.OperationRequest{A: math.NaN(), B: 10},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "malformed JSON",
			endpoint:   "/add",
			method:     http.MethodPost,
			body:       `{broken json}`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			if tt.body != nil {
				if s, ok := tt.body.(string); ok {
					body = []byte(s)
				} else {
					body, _ = json.Marshal(tt.body)
				}
			}

			req := httptest.NewRequest(tt.method, tt.endpoint, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}

			var errResp models.ErrorResponse
			json.NewDecoder(rec.Body).Decode(&errResp)
			if tt.wantErrorMsg != "" && errResp.Error != tt.wantErrorMsg {
				t.Errorf("error message = %q, want %q", errResp.Error, tt.wantErrorMsg)
			}
		})
	}
}

// TestIntegration_ConcurrentRequests tests handling of concurrent requests
func TestIntegration_ConcurrentRequests(t *testing.T) {
	h := setupHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	// Create multiple concurrent requests
	const numRequests = 50
	done := make(chan bool, numRequests)

	for i := 0; i < numRequests; i++ {
		go func(n int) {
			body, _ := json.Marshal(models.OperationRequest{A: float64(n), B: float64(n)})
			req := httptest.NewRequest(http.MethodPost, "/add", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Errorf("concurrent request %d failed: status = %d", n, rec.Code)
			}

			var resp models.OperationResponse
			json.NewDecoder(rec.Body).Decode(&resp)
			expected := float64(n + n)
			if resp.Result != expected {
				t.Errorf("concurrent request %d: result = %v, want %v", n, resp.Result, expected)
			}

			done <- true
		}(i)
	}

	// Wait for all requests to complete
	for i := 0; i < numRequests; i++ {
		<-done
	}
}

// TestIntegration_AllOperations tests all arithmetic operations in sequence
func TestIntegration_AllOperations(t *testing.T) {
	h := setupHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	operations := []struct {
		name     string
		endpoint string
		a        float64
		b        float64
		want     float64
	}{
		{"add", "/add", 10, 5, 15},
		{"subtract", "/subtract", 10, 5, 5},
		{"multiply", "/multiply", 10, 5, 50},
		{"add negative", "/add", -10, -5, -15},
		{"subtract to negative", "/subtract", 5, 10, -5},
		{"multiply negatives", "/multiply", -10, -5, 50},
		{"add decimals", "/add", 1.5, 2.3, 3.8},
		{"multiply by zero", "/multiply", 100, 0, 0},
	}

	for _, op := range operations {
		t.Run(op.name, func(t *testing.T) {
			body, _ := json.Marshal(models.OperationRequest{A: op.a, B: op.b})
			req := httptest.NewRequest(http.MethodPost, op.endpoint, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("%s failed: status = %d", op.name, rec.Code)
			}

			var resp models.OperationResponse
			json.NewDecoder(rec.Body).Decode(&resp)
			if resp.Result != op.want {
				t.Errorf("%s: result = %v, want %v", op.name, resp.Result, op.want)
			}

			// Verify Content-Type header
			contentType := rec.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("%s: Content-Type = %q, want %q", op.name, contentType, "application/json")
			}
		})
	}
}

// TestIntegration_ServiceInterface verifies the service interface implementation
func TestIntegration_ServiceInterface(t *testing.T) {
	// This test verifies that we can swap implementations of CalculatorService
	var _ service.CalculatorService = &service.Calculator{}

	// Create handler with the Calculator service
	calc := service.NewCalculator()
	h := NewHandler(calc)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	// Test that the handler works correctly with the service
	body, _ := json.Marshal(models.OperationRequest{A: 7, B: 3})
	req := httptest.NewRequest(http.MethodPost, "/add", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("service integration failed: status = %d", rec.Code)
	}

	var resp models.OperationResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Result != 10 {
		t.Errorf("service integration: result = %v, want 10", resp.Result)
	}
}
