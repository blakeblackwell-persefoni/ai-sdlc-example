package models

import (
	"encoding/json"
	"testing"
)

func TestOperationRequest_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name     string
		request  OperationRequest
		wantJSON string
	}{
		{
			name:     "positive integers",
			request:  OperationRequest{A: 10, B: 5},
			wantJSON: `{"a":10,"b":5}`,
		},
		{
			name:     "negative numbers",
			request:  OperationRequest{A: -10.5, B: -5.3},
			wantJSON: `{"a":-10.5,"b":-5.3}`,
		},
		{
			name:     "zeros",
			request:  OperationRequest{A: 0, B: 0},
			wantJSON: `{"a":0,"b":0}`,
		},
		{
			name:     "decimals",
			request:  OperationRequest{A: 3.14159, B: 2.71828},
			wantJSON: `{"a":3.14159,"b":2.71828}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.request)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}
			if string(got) != tt.wantJSON {
				t.Errorf("Marshal() = %s, want %s", got, tt.wantJSON)
			}
		})
	}
}

func TestOperationRequest_JSONUnmarshaling(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    OperationRequest
		wantErr bool
	}{
		{
			name: "valid JSON",
			json: `{"a":10,"b":5}`,
			want: OperationRequest{A: 10, B: 5},
		},
		{
			name: "negative numbers",
			json: `{"a":-10.5,"b":-5.3}`,
			want: OperationRequest{A: -10.5, B: -5.3},
		},
		{
			name: "scientific notation",
			json: `{"a":1.5e2,"b":3.2e-1}`,
			want: OperationRequest{A: 150, B: 0.32},
		},
		{
			name:    "invalid JSON",
			json:    `{"a":10`,
			wantErr: true,
		},
		{
			name:    "wrong types",
			json:    `{"a":"text","b":5}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got OperationRequest
			err := json.Unmarshal([]byte(tt.json), &got)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Unmarshal() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestOperationResponse_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name     string
		response OperationResponse
		wantJSON string
	}{
		{
			name:     "positive result",
			response: OperationResponse{Result: 15},
			wantJSON: `{"result":15}`,
		},
		{
			name:     "negative result",
			response: OperationResponse{Result: -5},
			wantJSON: `{"result":-5}`,
		},
		{
			name:     "zero result",
			response: OperationResponse{Result: 0},
			wantJSON: `{"result":0}`,
		},
		{
			name:     "decimal result",
			response: OperationResponse{Result: 3.14159},
			wantJSON: `{"result":3.14159}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.response)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}
			if string(got) != tt.wantJSON {
				t.Errorf("Marshal() = %s, want %s", got, tt.wantJSON)
			}
		})
	}
}

func TestErrorResponse_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name     string
		response ErrorResponse
		wantJSON string
	}{
		{
			name:     "simple error",
			response: ErrorResponse{Error: "invalid input"},
			wantJSON: `{"error":"invalid input"}`,
		},
		{
			name:     "empty error",
			response: ErrorResponse{Error: ""},
			wantJSON: `{"error":""}`,
		},
		{
			name:     "error with special characters",
			response: ErrorResponse{Error: "error: NaN and Infinity not allowed"},
			wantJSON: `{"error":"error: NaN and Infinity not allowed"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.response)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}
			if string(got) != tt.wantJSON {
				t.Errorf("Marshal() = %s, want %s", got, tt.wantJSON)
			}
		})
	}
}

func TestHealthResponse_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name     string
		response HealthResponse
		wantJSON string
	}{
		{
			name:     "ok status",
			response: HealthResponse{Status: "ok"},
			wantJSON: `{"status":"ok"}`,
		},
		{
			name:     "degraded status",
			response: HealthResponse{Status: "degraded"},
			wantJSON: `{"status":"degraded"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.response)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}
			if string(got) != tt.wantJSON {
				t.Errorf("Marshal() = %s, want %s", got, tt.wantJSON)
			}
		})
	}
}

func TestHealthResponse_JSONUnmarshaling(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    HealthResponse
		wantErr bool
	}{
		{
			name: "ok status",
			json: `{"status":"ok"}`,
			want: HealthResponse{Status: "ok"},
		},
		{
			name: "degraded status",
			json: `{"status":"degraded"}`,
			want: HealthResponse{Status: "degraded"},
		},
		{
			name:    "invalid JSON",
			json:    `{"status":"ok"`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got HealthResponse
			err := json.Unmarshal([]byte(tt.json), &got)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Unmarshal() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
