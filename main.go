package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type OperationRequest struct {
	A float64 `json:"a"`
	B float64 `json:"b"`
}

type OperationResponse struct {
	Result float64 `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {
	http.HandleFunc("/add", handleAdd)
	http.HandleFunc("/multiply", handleMultiply)
	http.HandleFunc("/subtract", handleSubtract)
	http.HandleFunc("/health", handleHealth)

	port := ":8080"
	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	req, err := parseRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	result := req.A + req.B
	writeJSON(w, http.StatusOK, OperationResponse{Result: result})
}

func handleMultiply(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	req, err := parseRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	result := req.A * req.B
	writeJSON(w, http.StatusOK, OperationResponse{Result: result})
}

func handleSubtract(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	req, err := parseRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	result := req.A - req.B
	writeJSON(w, http.StatusOK, OperationResponse{Result: result})
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func parseRequest(r *http.Request) (*OperationRequest, error) {
	var req OperationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("invalid request body: %w", err)
	}
	return &req, nil
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}

