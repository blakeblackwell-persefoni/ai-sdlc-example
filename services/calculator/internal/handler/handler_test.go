package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blakeblackwell-persefoni/ai-sdlc-example/services/calculator/internal/models"
	"github.com/blakeblackwell-persefoni/ai-sdlc-example/services/calculator/internal/service"
)

func setupHandler() *Handler {
	calc := service.NewCalculator()
	return NewHandler(calc)
}

func TestHandler_Add(t *testing.T) {
	h := setupHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	tests := []struct {
		name           string
		method         string
		body           interface{}
		wantStatus     int
		wantResult     float64
		wantError      bool
		wantErrorMsg   string
	}{
		{
			name:       "valid addition",
			method:     http.MethodPost,
			body:       models.OperationRequest{A: 10, B: 5},
			wantStatus: http.StatusOK,
			wantResult: 15,
		},
		{
			name:       "addition with negatives",
			method:     http.MethodPost,
			body:       models.OperationRequest{A: -10, B: 5},
			wantStatus: http.StatusOK,
			wantResult: -5,
		},
		{
			name:         "method not allowed",
			method:       http.MethodGet,
			body:         nil,
			wantStatus:   http.StatusMethodNotAllowed,
			wantError:    true,
			wantErrorMsg: "Method not allowed",
		},
		{
			name:       "invalid JSON",
			method:     http.MethodPost,
			body:       "invalid",
			wantStatus: http.StatusBadRequest,
			wantError:  true,
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

			req := httptest.NewRequest(tt.method, "/add", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}

			if tt.wantError {
				var errResp models.ErrorResponse
				json.NewDecoder(rec.Body).Decode(&errResp)
				if tt.wantErrorMsg != "" && errResp.Error != tt.wantErrorMsg {
					t.Errorf("error message = %q, want %q", errResp.Error, tt.wantErrorMsg)
				}
			} else {
				var resp models.OperationResponse
				json.NewDecoder(rec.Body).Decode(&resp)
				if resp.Result != tt.wantResult {
					t.Errorf("result = %v, want %v", resp.Result, tt.wantResult)
				}
			}
		})
	}
}

func TestHandler_Subtract(t *testing.T) {
	h := setupHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	tests := []struct {
		name       string
		method     string
		body       models.OperationRequest
		wantStatus int
		wantResult float64
	}{
		{
			name:       "valid subtraction",
			method:     http.MethodPost,
			body:       models.OperationRequest{A: 10, B: 5},
			wantStatus: http.StatusOK,
			wantResult: 5,
		},
		{
			name:       "negative result",
			method:     http.MethodPost,
			body:       models.OperationRequest{A: 5, B: 10},
			wantStatus: http.StatusOK,
			wantResult: -5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(tt.method, "/subtract", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}

			var resp models.OperationResponse
			json.NewDecoder(rec.Body).Decode(&resp)
			if resp.Result != tt.wantResult {
				t.Errorf("result = %v, want %v", resp.Result, tt.wantResult)
			}
		})
	}
}

func TestHandler_Multiply(t *testing.T) {
	h := setupHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	tests := []struct {
		name       string
		method     string
		body       models.OperationRequest
		wantStatus int
		wantResult float64
	}{
		{
			name:       "valid multiplication",
			method:     http.MethodPost,
			body:       models.OperationRequest{A: 10, B: 5},
			wantStatus: http.StatusOK,
			wantResult: 50,
		},
		{
			name:       "multiply by zero",
			method:     http.MethodPost,
			body:       models.OperationRequest{A: 10, B: 0},
			wantStatus: http.StatusOK,
			wantResult: 0,
		},
		{
			name:       "negative numbers",
			method:     http.MethodPost,
			body:       models.OperationRequest{A: -10, B: -5},
			wantStatus: http.StatusOK,
			wantResult: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(tt.method, "/multiply", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}

			var resp models.OperationResponse
			json.NewDecoder(rec.Body).Decode(&resp)
			if resp.Result != tt.wantResult {
				t.Errorf("result = %v, want %v", resp.Result, tt.wantResult)
			}
		})
	}
}

func TestHandler_Health(t *testing.T) {
	h := setupHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	tests := []struct {
		name         string
		method       string
		wantStatus   int
		wantError    bool
		wantErrorMsg string
	}{
		{
			name:       "valid health check",
			method:     http.MethodGet,
			wantStatus: http.StatusOK,
		},
		{
			name:         "method not allowed",
			method:       http.MethodPost,
			wantStatus:   http.StatusMethodNotAllowed,
			wantError:    true,
			wantErrorMsg: "Method not allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/health", nil)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}

			if tt.wantError {
				var errResp models.ErrorResponse
				json.NewDecoder(rec.Body).Decode(&errResp)
				if errResp.Error != tt.wantErrorMsg {
					t.Errorf("error message = %q, want %q", errResp.Error, tt.wantErrorMsg)
				}
			} else {
				var resp models.HealthResponse
				json.NewDecoder(rec.Body).Decode(&resp)
				if resp.Status != "ok" {
					t.Errorf("status = %q, want %q", resp.Status, "ok")
				}
			}
		})
	}
}

func TestHandler_ContentType(t *testing.T) {
	h := setupHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	body, _ := json.Marshal(models.OperationRequest{A: 1, B: 2})
	req := httptest.NewRequest(http.MethodPost, "/add", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Content-Type = %q, want %q", contentType, "application/json")
	}
}

func TestHandler_EmptyBody(t *testing.T) {
	h := setupHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodPost, "/add", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}
