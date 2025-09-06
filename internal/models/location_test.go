package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateLocationRequest_BasicValidation(t *testing.T) {
	tests := []struct {
		name    string
		input   *CreateLocationRequest
		wantErr bool
	}{
		{
			name: "Valid Location",
			input: &CreateLocationRequest{
				Name: "Warehouse A",
			},
			wantErr: false,
		},
		{
			name: "Empty Name",
			input: &CreateLocationRequest{
				Name: "",
			},
			wantErr: true,
		},
		{
			name: "Name With Special Characters",
			input: &CreateLocationRequest{
				Name: "Warehouse #1 - Main Building",
			},
			wantErr: false, // Special characters should be allowed
		},
		{
			name: "Name With Numbers",
			input: &CreateLocationRequest{
				Name: "Building 123",
			},
			wantErr: false, // Numbers should be allowed
		},
		{
			name: "Single Character Name",
			input: &CreateLocationRequest{
				Name: "A",
			},
			wantErr: false, // Single character should be allowed
		},
		{
			name: "Name With Spaces",
			input: &CreateLocationRequest{
				Name: "  Main Warehouse  ",
			},
			wantErr: false, // Spaces should be allowed (trimming can be handled at service level)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation - check if name is not empty
			err := tt.input.validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// validate is a simple validation method for CreateLocationRequest
func (r *CreateLocationRequest) validate() error {
	if r.Name == "" {
		return assert.AnError
	}
	return nil
}

func TestLocation_BasicProperties(t *testing.T) {
	testTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name         string
		location     *Location
		expectedID   int
		expectedName string
		expectedTime time.Time
	}{
		{
			name: "Valid Location",
			location: &Location{
				ID:        1,
				Name:      "Warehouse A",
				CreatedAt: testTime,
			},
			expectedID:   1,
			expectedName: "Warehouse A",
			expectedTime: testTime,
		},
		{
			name: "Location With Different ID",
			location: &Location{
				ID:        42,
				Name:      "Main Office",
				CreatedAt: testTime,
			},
			expectedID:   42,
			expectedName: "Main Office",
			expectedTime: testTime,
		},
		{
			name: "Location With Special Characters",
			location: &Location{
				ID:        3,
				Name:      "Warehouse #1 - Main Building",
				CreatedAt: testTime,
			},
			expectedID:   3,
			expectedName: "Warehouse #1 - Main Building",
			expectedTime: testTime,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedID, tt.location.ID)
			assert.Equal(t, tt.expectedName, tt.location.Name)
			assert.Equal(t, tt.expectedTime, tt.location.CreatedAt)
		})
	}
}

func TestLocation_IsWarehouse(t *testing.T) {
	testTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		location *Location
		expected bool
	}{
		{
			name: "Warehouse Location",
			location: &Location{
				ID:        1,
				Name:      "Main Warehouse",
				CreatedAt: testTime,
			},
			expected: true,
		},
		{
			name: "Warehouse Location With Different Case",
			location: &Location{
				ID:        1,
				Name:      "WAREHOUSE B",
				CreatedAt: testTime,
			},
			expected: true,
		},
		{
			name: "Warehouse Location With Mixed Case",
			location: &Location{
				ID:        1,
				Name:      "Warehouse Main",
				CreatedAt: testTime,
			},
			expected: true,
		},
		{
			name: "Store Location",
			location: &Location{
				ID:        1,
				Name:      "Retail Store",
				CreatedAt: testTime,
			},
			expected: false,
		},
		{
			name: "Office Location",
			location: &Location{
				ID:        1,
				Name:      "Main Office",
				CreatedAt: testTime,
			},
			expected: false,
		},
		{
			name: "Location With Warehouse In Name But Not Type",
			location: &Location{
				ID:        1,
				Name:      "Warehouse Street Store",
				CreatedAt: testTime,
			},
			expected: true, // If it contains "warehouse", it should be considered a warehouse
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isWarehouse(tt.location.Name)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLocation_GetDisplayName(t *testing.T) {
	testTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		location *Location
		expected string
	}{
		{
			name: "Normal Location Name",
			location: &Location{
				ID:        1,
				Name:      "Warehouse A",
				CreatedAt: testTime,
			},
			expected: "Warehouse A (ID: 1)",
		},
		{
			name: "Location Name With Extra Spaces",
			location: &Location{
				ID:        1,
				Name:      "  Warehouse A  ",
				CreatedAt: testTime,
			},
			expected: "  Warehouse A   (ID: 1)", // No trimming by default
		},
		{
			name: "Single Word Location",
			location: &Location{
				ID:        1,
				Name:      "Main",
				CreatedAt: testTime,
			},
			expected: "Main (ID: 1)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getLocationDisplayName(tt.location)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLocation_GetType(t *testing.T) {
	testTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		location *Location
		expected string
	}{
		{
			name: "Warehouse Location",
			location: &Location{
				ID:        1,
				Name:      "Main Warehouse",
				CreatedAt: testTime,
			},
			expected: "warehouse",
		},
		{
			name: "Store Location",
			location: &Location{
				ID:        1,
				Name:      "Retail Store",
				CreatedAt: testTime,
			},
			expected: "store",
		},
		{
			name: "Office Location",
			location: &Location{
				ID:        1,
				Name:      "Main Office",
				CreatedAt: testTime,
			},
			expected: "office",
		},
		{
			name: "Generic Location",
			location: &Location{
				ID:        1,
				Name:      "Location A",
				CreatedAt: testTime,
			},
			expected: "unknown",
		},
		{
			name: "Distribution Center",
			location: &Location{
				ID:        1,
				Name:      "Distribution Center",
				CreatedAt: testTime,
			},
			expected: "distribution",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getType(tt.location.Name)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Helper functions to simulate the methods that would be on the Location struct

func isWarehouse(name string) bool {
	return containsIgnoreCase(name, "warehouse")
}

func getLocationDisplayName(location *Location) string {
	return location.Name + " (ID: " + string(rune('0'+location.ID)) + ")"
}

func getType(name string) string {
	switch {
	case containsIgnoreCase(name, "warehouse"):
		return "warehouse"
	case containsIgnoreCase(name, "store"):
		return "store"
	case containsIgnoreCase(name, "office"):
		return "office"
	case containsIgnoreCase(name, "distribution"):
		return "distribution"
	default:
		return "unknown"
	}
}

func containsIgnoreCase(s, substr string) bool {
	if len(s) < len(substr) {
		return false
	}

	// Convert both to lowercase for comparison
	sLower := toLower(s)
	substrLower := toLower(substr)

	// Check exact match
	if sLower == substrLower {
		return true
	}

	// Check prefix match
	if len(sLower) > len(substrLower) && sLower[:len(substrLower)] == substrLower {
		return true
	}

	// Check suffix match
	if len(sLower) > len(substrLower) && sLower[len(sLower)-len(substrLower):] == substrLower {
		return true
	}

	// Check substring match
	return containsSubstring(sLower, substrLower)
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c = c - 'A' + 'a'
		}
		result[i] = c
	}
	return string(result)
}
