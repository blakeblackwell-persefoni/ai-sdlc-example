package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/blakeblackwell-persefoni/ai-sdlc-example/services/calculator/internal/models"
	"github.com/blakeblackwell-persefoni/ai-sdlc-example/services/calculator/internal/service"
)

// Handler handles HTTP requests for the calculator service.
type Handler struct {
	calculator *service.Calculator
}

// NewHandler creates a new Handler with the given calculator service.
func NewHandler(calc *service.Calculator) *Handler {
	return &Handler{calculator: calc}
}

// RegisterRoutes registers all routes on the given mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/add", h.handleAdd)
	mux.HandleFunc("/subtract", h.handleSubtract)
	mux.HandleFunc("/multiply", h.handleMultiply)
	mux.HandleFunc("/health", h.handleHealth)
}

func (h *Handler) handleAdd(w http.ResponseWriter, r *http.Request) {
	h.handleOperation(w, r, h.calculator.Add)
}

func (h *Handler) handleSubtract(w http.ResponseWriter, r *http.Request) {
	h.handleOperation(w, r, h.calculator.Subtract)
}

func (h *Handler) handleMultiply(w http.ResponseWriter, r *http.Request) {
	h.handleOperation(w, r, h.calculator.Multiply)
}

func (h *Handler) handleOperation(w http.ResponseWriter, r *http.Request, op func(a, b float64) (float64, error)) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	req, err := parseRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := op(req.A, req.B)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, models.OperationResponse{Result: result})
}

func (h *Handler) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	writeJSON(w, http.StatusOK, models.HealthResponse{Status: "ok"})
}

func parseRequest(r *http.Request) (*models.OperationRequest, error) {
	r.Body = http.MaxBytesReader(nil, r.Body, 1<<20) // 1MB limit

	var req models.OperationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("failed to encode JSON response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, models.ErrorResponse{Error: message})
}
