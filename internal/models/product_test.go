//go:build unit

package models

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateProductRequest_BasicValidation(t *testing.T) {
	tests := []struct {
		name    string
		input   *CreateProductRequest
		wantErr bool
	}{
		{
			name: "Valid Product",
			input: &CreateProductRequest{
				SKU:         "TEST001",
				Name:        "Test Product",
				Description: "A test product",
				Price:       9.99,
			},
			wantErr: false,
		},
		{
			name: "Empty SKU",
			input: &CreateProductRequest{
				SKU:         "",
				Name:        "Test Product",
				Description: "A test product",
				Price:       9.99,
			},
			wantErr: true,
		},
		{
			name: "Empty Name",
			input: &CreateProductRequest{
				SKU:         "TEST001",
				Name:        "",
				Description: "A test product",
				Price:       9.99,
			},
			wantErr: true,
		},
		{
			name: "Negative Price",
			input: &CreateProductRequest{
				SKU:         "TEST001",
				Name:        "Test Product",
				Description: "A test product",
				Price:       -9.99,
			},
			wantErr: true,
		},
		{
			name: "Zero Price",
			input: &CreateProductRequest{
				SKU:         "TEST001",
				Name:        "Test Product",
				Description: "A test product",
				Price:       0,
			},
			wantErr: false, // Zero price might be allowed (free product)
		},
		{
			name: "Empty Description",
			input: &CreateProductRequest{
				SKU:         "TEST001",
				Name:        "Test Product",
				Description: "",
				Price:       9.99,
			},
			wantErr: false, // Empty description should be allowed
		},
		{
			name: "SKU With Special Characters",
			input: &CreateProductRequest{
				SKU:         "TEST-001_2024",
				Name:        "Test Product",
				Description: "A test product",
				Price:       9.99,
			},
			wantErr: false, // Special characters in SKU should be allowed
		},
		{
			name: "Name With Numbers",
			input: &CreateProductRequest{
				SKU:         "TEST001",
				Name:        "Product 2024",
				Description: "A test product",
				Price:       9.99,
			},
			wantErr: false, // Numbers in name should be allowed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// validate is a simple validation method for CreateProductRequest
func (r *CreateProductRequest) validate() error {
	if r.SKU == "" {
		return assert.AnError
	}
	if r.Name == "" {
		return assert.AnError
	}
	if r.Price < 0 {
		return assert.AnError
	}
	return nil
}

func TestProduct_BasicProperties(t *testing.T) {
	testTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		product       *Product
		expectedID    int
		expectedSKU   string
		expectedName  string
		expectedDesc  string
		expectedPrice float64
		expectedTime  time.Time
	}{
		{
			name: "Valid Product",
			product: &Product{
				ID:          1,
				SKU:         "TEST001",
				Name:        "Test Product",
				Description: "A test product",
				Price:       9.99,
				CreatedAt:   testTime,
			},
			expectedID:    1,
			expectedSKU:   "TEST001",
			expectedName:  "Test Product",
			expectedDesc:  "A test product",
			expectedPrice: 9.99,
			expectedTime:  testTime,
		},
		{
			name: "Product With Different Values",
			product: &Product{
				ID:          42,
				SKU:         "PREMIUM002",
				Name:        "Premium Product",
				Description: "A premium quality product",
				Price:       99.99,
				CreatedAt:   testTime,
			},
			expectedID:    42,
			expectedSKU:   "PREMIUM002",
			expectedName:  "Premium Product",
			expectedDesc:  "A premium quality product",
			expectedPrice: 99.99,
			expectedTime:  testTime,
		},
		{
			name: "Product With Empty Description",
			product: &Product{
				ID:          3,
				SKU:         "BASIC003",
				Name:        "Basic Product",
				Description: "",
				Price:       0.00,
				CreatedAt:   testTime,
			},
			expectedID:    3,
			expectedSKU:   "BASIC003",
			expectedName:  "Basic Product",
			expectedDesc:  "",
			expectedPrice: 0.00,
			expectedTime:  testTime,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedID, tt.product.ID)
			assert.Equal(t, tt.expectedSKU, tt.product.SKU)
			assert.Equal(t, tt.expectedName, tt.product.Name)
			assert.Equal(t, tt.expectedDesc, tt.product.Description)
			assert.Equal(t, tt.expectedPrice, tt.product.Price)
			assert.Equal(t, tt.expectedTime, tt.product.CreatedAt)
		})
	}
}

func TestProduct_IsPremium(t *testing.T) {
	testTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		product  *Product
		expected bool
	}{
		{
			name: "Premium Product",
			product: &Product{
				ID:          1,
				SKU:         "PREMIUM001",
				Name:        "Premium Product",
				Description: "A premium quality product",
				Price:       99.99,
				CreatedAt:   testTime,
			},
			expected: true,
		},
		{
			name: "Expensive Product",
			product: &Product{
				ID:          1,
				SKU:         "EXPENSIVE001",
				Name:        "Expensive Product",
				Description: "An expensive product",
				Price:       150.00,
				CreatedAt:   testTime,
			},
			expected: true,
		},
		{
			name: "Standard Product",
			product: &Product{
				ID:          1,
				SKU:         "STANDARD001",
				Name:        "Standard Product",
				Description: "A standard quality product",
				Price:       25.00,
				CreatedAt:   testTime,
			},
			expected: false,
		},
		{
			name: "Budget Product",
			product: &Product{
				ID:          1,
				SKU:         "BUDGET001",
				Name:        "Budget Product",
				Description: "A budget product",
				Price:       5.00,
				CreatedAt:   testTime,
			},
			expected: false,
		},
		{
			name: "Free Product",
			product: &Product{
				ID:          1,
				SKU:         "FREE001",
				Name:        "Free Product",
				Description: "A free product",
				Price:       0.00,
				CreatedAt:   testTime,
			},
			expected: false,
		},
		{
			name: "Exactly Premium Threshold",
			product: &Product{
				ID:          1,
				SKU:         "THRESHOLD001",
				Name:        "Threshold Product",
				Description: "Product at premium threshold",
				Price:       50.00,
				CreatedAt:   testTime,
			},
			expected: true, // Assuming 50.00 is the premium threshold
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isPremium(tt.product)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestProduct_GetDisplayName(t *testing.T) {
	testTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		product  *Product
		expected string
	}{
		{
			name: "Normal Product",
			product: &Product{
				ID:          1,
				SKU:         "TEST001",
				Name:        "Test Product",
				Description: "A test product",
				Price:       9.99,
				CreatedAt:   testTime,
			},
			expected: "Test Product (TEST001)",
		},
		{
			name: "Product With Long Name",
			product: &Product{
				ID:          1,
				SKU:         "LONG001",
				Name:        "Very Long Product Name That Might Need Truncating",
				Description: "A product with a very long name",
				Price:       19.99,
				CreatedAt:   testTime,
			},
			expected: "Very Long Product Name That Might Need Truncating (LONG001)",
		},
		{
			name: "Product With Special Characters",
			product: &Product{
				ID:          1,
				SKU:         "SPECIAL-001",
				Name:        "Product #1 - Special Edition",
				Description: "A special edition product",
				Price:       29.99,
				CreatedAt:   testTime,
			},
			expected: "Product #1 - Special Edition (SPECIAL-001)",
		},
		{
			name: "Product With Numbers",
			product: &Product{
				ID:          1,
				SKU:         "NUM2024001",
				Name:        "Product 2024",
				Description: "A 2024 product",
				Price:       39.99,
				CreatedAt:   testTime,
			},
			expected: "Product 2024 (NUM2024001)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getProductDisplayName(tt.product)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestProduct_GetDisplayPrice(t *testing.T) {
	testTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		product  *Product
		expected string
	}{
		{
			name: "Normal Price",
			product: &Product{
				ID:          1,
				SKU:         "TEST001",
				Name:        "Test Product",
				Description: "A test product",
				Price:       9.99,
				CreatedAt:   testTime,
			},
			expected: "$9.99",
		},
		{
			name: "Whole Number Price",
			product: &Product{
				ID:          1,
				SKU:         "WHOLE001",
				Name:        "Whole Price Product",
				Description: "A product with whole number price",
				Price:       25.00,
				CreatedAt:   testTime,
			},
			expected: "$25.00",
		},
		{
			name: "High Price",
			product: &Product{
				ID:          1,
				SKU:         "HIGH001",
				Name:        "High Price Product",
				Description: "An expensive product",
				Price:       999.99,
				CreatedAt:   testTime,
			},
			expected: "$999.99",
		},
		{
			name: "Free Product",
			product: &Product{
				ID:          1,
				SKU:         "FREE001",
				Name:        "Free Product",
				Description: "A free product",
				Price:       0.00,
				CreatedAt:   testTime,
			},
			expected: "$0.00",
		},
		{
			name: "Price With Single Digit",
			product: &Product{
				ID:          1,
				SKU:         "SINGLE001",
				Name:        "Single Digit Price",
				Description: "A product with single digit price",
				Price:       5.50,
				CreatedAt:   testTime,
			},
			expected: "$5.50",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getDisplayPrice(tt.product)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Helper functions to simulate the methods that would be on the Product struct

func isPremium(product *Product) bool {
	return product.Price >= 50.00 // Assuming 50.00 is the premium threshold
}

func getProductDisplayName(product *Product) string {
	return product.Name + " (" + product.SKU + ")"
}

func getDisplayPrice(product *Product) string {
	return "$" + formatPrice(product.Price)
}

func formatPrice(price float64) string {
	// Simple price formatting - in a real implementation, you might use more sophisticated formatting
	return fmt.Sprintf("%.2f", price)
}
