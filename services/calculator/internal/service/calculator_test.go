package service

import (
	"math"
	"testing"
)

func TestCalculator_Add(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name    string
		a       float64
		b       float64
		want    float64
		wantErr bool
	}{
		{
			name: "positive numbers",
			a:    10,
			b:    5,
			want: 15,
		},
		{
			name: "negative numbers",
			a:    -10,
			b:    -5,
			want: -15,
		},
		{
			name: "mixed signs",
			a:    10,
			b:    -5,
			want: 5,
		},
		{
			name: "zeros",
			a:    0,
			b:    0,
			want: 0,
		},
		{
			name: "decimals",
			a:    1.5,
			b:    2.5,
			want: 4,
		},
		{
			name:    "NaN input a",
			a:       math.NaN(),
			b:       5,
			wantErr: true,
		},
		{
			name:    "NaN input b",
			a:       5,
			b:       math.NaN(),
			wantErr: true,
		},
		{
			name:    "Infinity input",
			a:       math.Inf(1),
			b:       5,
			wantErr: true,
		},
		{
			name:    "negative Infinity input",
			a:       5,
			b:       math.Inf(-1),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calc.Add(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculator_Subtract(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name    string
		a       float64
		b       float64
		want    float64
		wantErr bool
	}{
		{
			name: "positive numbers",
			a:    10,
			b:    5,
			want: 5,
		},
		{
			name: "negative result",
			a:    5,
			b:    10,
			want: -5,
		},
		{
			name: "negative numbers",
			a:    -10,
			b:    -5,
			want: -5,
		},
		{
			name: "zeros",
			a:    0,
			b:    0,
			want: 0,
		},
		{
			name: "decimals",
			a:    5.5,
			b:    2.5,
			want: 3,
		},
		{
			name:    "NaN input",
			a:       math.NaN(),
			b:       5,
			wantErr: true,
		},
		{
			name:    "Infinity input",
			a:       math.Inf(1),
			b:       5,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calc.Subtract(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Subtract() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Subtract() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculator_Multiply(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name    string
		a       float64
		b       float64
		want    float64
		wantErr bool
	}{
		{
			name: "positive numbers",
			a:    10,
			b:    5,
			want: 50,
		},
		{
			name: "negative result",
			a:    -10,
			b:    5,
			want: -50,
		},
		{
			name: "both negative",
			a:    -10,
			b:    -5,
			want: 50,
		},
		{
			name: "multiply by zero",
			a:    10,
			b:    0,
			want: 0,
		},
		{
			name: "decimals",
			a:    2.5,
			b:    4,
			want: 10,
		},
		{
			name:    "NaN input",
			a:       math.NaN(),
			b:       5,
			wantErr: true,
		},
		{
			name:    "Infinity input",
			a:       5,
			b:       math.Inf(1),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calc.Multiply(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Multiply() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Multiply() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateInputs(t *testing.T) {
	tests := []struct {
		name    string
		a       float64
		b       float64
		wantErr bool
	}{
		{
			name:    "valid inputs",
			a:       10,
			b:       5,
			wantErr: false,
		},
		{
			name:    "NaN a",
			a:       math.NaN(),
			b:       5,
			wantErr: true,
		},
		{
			name:    "NaN b",
			a:       5,
			b:       math.NaN(),
			wantErr: true,
		},
		{
			name:    "positive infinity a",
			a:       math.Inf(1),
			b:       5,
			wantErr: true,
		},
		{
			name:    "negative infinity b",
			a:       5,
			b:       math.Inf(-1),
			wantErr: true,
		},
		{
			name:    "both NaN",
			a:       math.NaN(),
			b:       math.NaN(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateInputs(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateInputs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
