package testutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"cli-inventory/internal/openapi"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// OpenAPITestHelper provides utilities for OpenAPI-based testing
type OpenAPITestHelper struct {
	validator *openapi.Validator
	t         *testing.T
}

// NewOpenAPITestHelper creates a new OpenAPI test helper
func NewOpenAPITestHelper(t *testing.T, specPath string) *OpenAPITestHelper {
	validator, err := openapi.NewValidator(specPath)
	require.NoError(t, err, "Failed to create OpenAPI validator")

	return &OpenAPITestHelper{
		validator: validator,
		t:         t,
	}
}

// ValidateRequest validates that a request conforms to OpenAPI specification
func (h *OpenAPITestHelper) ValidateRequest(method, path string, body interface{}) {
	h.t.Helper()

	// Get operation from OpenAPI spec
	operation, err := h.validator.GetOperation(method, path)
	require.NoError(h.t, err, "Failed to get operation from OpenAPI spec")

	// If there's a request body, validate its schema
	if body != nil {
		requestBody := operation.RequestBody
		if requestBody != nil && requestBody.Value != nil {
			content := requestBody.Value.Content.Get("application/json")
			if content != nil && content.Schema != nil {
				// Marshal body to JSON
				jsonBody, err := json.Marshal(body)
				require.NoError(h.t, err, "Failed to marshal request body")

				// Validate against schema
				err = h.validator.ValidateRequestPayload(nil, content.Schema, jsonBody)
				assert.NoError(h.t, err, "Request body does not conform to OpenAPI schema")
			}
		}
	}
}

// ValidateResponse validates that a response conforms to OpenAPI specification
func (h *OpenAPITestHelper) ValidateResponse(method, path, statusCode string, body []byte) {
	h.t.Helper()

	// Get operation from OpenAPI spec
	operation, err := h.validator.GetOperation(method, path)
	require.NoError(h.t, err, "Failed to get operation from OpenAPI spec")

	// Find response schema
	response := operation.Responses.Value(statusCode)
	if response == nil {
		// Response not documented in OpenAPI spec - skip validation
		h.t.Logf("Response with status code '%s' not documented in OpenAPI spec for %s %s", statusCode, method, path)
		return
	}

	// Get JSON schema
	content := response.Value.Content.Get("application/json")
	if content != nil && content.Schema != nil {
		// Validate against schema
		err = h.validator.ValidateResponsePayload(content.Schema, body)
		assert.NoError(h.t, err, "Response body does not conform to OpenAPI schema")
	}
}

// ValidateHTTPResponse validates an HTTP response against OpenAPI specification
func (h *OpenAPITestHelper) ValidateHTTPResponse(method, path string, resp *httptest.ResponseRecorder) {
	h.t.Helper()

	// Convert numeric status code to string as expected by OpenAPI spec
	statusCode := fmt.Sprintf("%d", resp.Code)

	h.ValidateResponse(method, path, statusCode, resp.Body.Bytes())
}

// CreateTestRequest creates a test request and validates it against OpenAPI spec
func (h *OpenAPITestHelper) CreateTestRequest(method, path string, body interface{}) *http.Request {
	h.t.Helper()

	h.ValidateRequest(method, path, body)

	var reqBody *bytes.Buffer
	if body != nil {
		jsonData, err := json.Marshal(body)
		require.NoError(h.t, err, "Failed to marshal request body")
		reqBody = bytes.NewBuffer(jsonData)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, err := http.NewRequest(method, path, reqBody)
	require.NoError(h.t, err, "Failed to create request")

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req
}

// AssertOpenAPICompliance asserts that a test response is OpenAPI compliant
func (h *OpenAPITestHelper) AssertOpenAPICompliance(method, path string, resp *httptest.ResponseRecorder) {
	h.t.Helper()

	// Validate response against OpenAPI spec
	h.ValidateHTTPResponse(method, path, resp)

	// Additional compliance checks
	h.t.Run("Content-Type", func(t *testing.T) {
		contentType := resp.Header().Get("Content-Type")
		if resp.Body.Len() > 0 {
			assert.Contains(t, contentType, "application/json", "Response with body should have JSON content type")
		}
	})

	h.t.Run("Status Code", func(t *testing.T) {
		// Check if status code is valid HTTP status code
		assert.True(t, resp.Code >= 100 && resp.Code < 600, "Status code should be valid HTTP status code")
	})
}

// GetSchema returns a schema from the OpenAPI specification
func (h *OpenAPITestHelper) GetSchema(schemaName string) *openapi3.SchemaRef {
	h.t.Helper()

	schema, err := h.validator.GetSchema(schemaName)
	require.NoError(h.t, err, "Failed to get schema from OpenAPI spec")

	return schema
}

// ValidateJSONAgainstSchema validates JSON data against a specific schema
func (h *OpenAPITestHelper) ValidateJSONAgainstSchema(schemaName string, data []byte) {
	h.t.Helper()

	schema := h.GetSchema(schemaName)
	err := h.validator.ValidateResponsePayload(schema, data)
	assert.NoError(h.t, err, "Data does not conform to schema '%s'", schemaName)
}

// RequireOpenAPIField asserts that a JSON response contains required fields from OpenAPI schema
func (h *OpenAPITestHelper) RequireOpenAPIField(method, path, statusCode, fieldName string, resp *httptest.ResponseRecorder) {
	h.t.Helper()

	// Get operation and response schema
	operation, err := h.validator.GetOperation(method, path)
	require.NoError(h.t, err, "Failed to get operation from OpenAPI spec")

	response := operation.Responses.Value(statusCode)
	require.NotNil(h.t, response, "Response not documented in OpenAPI spec")

	content := response.Value.Content.Get("application/json")
	require.NotNil(h.t, content, "JSON content not documented in OpenAPI spec")
	require.NotNil(h.t, content.Schema, "Schema not documented in OpenAPI spec")

	// Parse response body
	var responseData map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &responseData)
	require.NoError(h.t, err, "Failed to parse response JSON")

	// Check if field exists in response
	_, exists := responseData[fieldName]
	assert.True(h.t, exists, "Required field '%s' missing from response", fieldName)
}

// RequireOpenAPIFields asserts that a JSON response contains all required fields from OpenAPI schema
func (h *OpenAPITestHelper) RequireOpenAPIFields(method, path, statusCode string, resp *httptest.ResponseRecorder) {
	h.t.Helper()

	// Get operation and response schema
	operation, err := h.validator.GetOperation(method, path)
	require.NoError(h.t, err, "Failed to get operation from OpenAPI spec")

	response := operation.Responses.Value(statusCode)
	require.NotNil(h.t, response, "Response not documented in OpenAPI spec")

	content := response.Value.Content.Get("application/json")
	require.NotNil(h.t, content, "JSON content not documented in OpenAPI spec")
	require.NotNil(h.t, content.Schema, "Schema not documented in OpenAPI spec")

	schema := content.Schema.Value
	require.NotNil(h.t, schema, "Schema value is nil")

	// Parse response body
	var responseData map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &responseData)
	require.NoError(h.t, err, "Failed to parse response JSON")

	// Check all required fields exist
	for _, requiredField := range schema.Required {
		_, exists := responseData[requiredField]
		assert.True(h.t, exists, "Required field '%s' missing from response", requiredField)
	}
}
