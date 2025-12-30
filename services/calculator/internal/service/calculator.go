package service

import (
	"errors"
	"math"
)

var (
	ErrInvalidInput = errors.New("invalid input: NaN and Infinity not allowed")
)

// CalculatorService defines the interface for arithmetic operations.
type CalculatorService interface {
	Add(a, b float64) (float64, error)
	Subtract(a, b float64) (float64, error)
	Multiply(a, b float64) (float64, error)
}

// Calculator provides arithmetic operations.
type Calculator struct{}

// Ensure Calculator implements CalculatorService at compile time.
var _ CalculatorService = (*Calculator)(nil)

// NewCalculator creates a new Calculator instance.
func NewCalculator() *Calculator {
	return &Calculator{}
}

// Add returns the sum of two numbers.
func (c *Calculator) Add(a, b float64) (float64, error) {
	if err := validateInputs(a, b); err != nil {
		return 0, err
	}
	return a + b, nil
}

// Subtract returns the difference of two numbers (a - b).
func (c *Calculator) Subtract(a, b float64) (float64, error) {
	if err := validateInputs(a, b); err != nil {
		return 0, err
	}
	return a - b, nil
}

// Multiply returns the product of two numbers.
func (c *Calculator) Multiply(a, b float64) (float64, error) {
	if err := validateInputs(a, b); err != nil {
		return 0, err
	}
	return a * b, nil
}

// validateInputs checks that the inputs are valid numbers.
func validateInputs(a, b float64) error {
	if math.IsNaN(a) || math.IsInf(a, 0) || math.IsNaN(b) || math.IsInf(b, 0) {
		return ErrInvalidInput
	}
	return nil
}
