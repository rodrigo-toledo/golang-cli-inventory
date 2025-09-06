package openapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/legacy"
)

// Validator handles OpenAPI specification validation
type Validator struct {
	doc     *openapi3.T
	router  routers.Router
	enabled bool
}

// NewValidator creates a new OpenAPI validator
func NewValidator(specPath string) (*Validator, error) {
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	// Get absolute path
	absPath, err := filepath.Abs(specPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Load OpenAPI specification
	doc, err := loader.LoadFromFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load OpenAPI spec: %w", err)
	}

	// Validate the document
	if err := doc.Validate(loader.Context); err != nil {
		return nil, fmt.Errorf("OpenAPI spec validation failed: %w", err)
	}

	// Create router
	router, err := legacy.NewRouter(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to create router: %w", err)
	}

	return &Validator{
		doc:     doc,
		router:  router,
		enabled: true,
	}, nil
}

// Middleware returns an HTTP middleware for OpenAPI validation
func (v *Validator) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !v.enabled {
				next.ServeHTTP(w, r)
				return
			}

			// Skip validation for non-API routes
			if !strings.HasPrefix(r.URL.Path, "/api/") {
				next.ServeHTTP(w, r)
				return
			}

			// Find route
			route, pathParams, err := v.router.FindRoute(r)
			if err != nil {
				// Route not found in OpenAPI spec
				next.ServeHTTP(w, r)
				return
			}

			// Read request body
			var body []byte
			if r.Body != nil {
				body, err = io.ReadAll(r.Body)
				if err != nil {
					http.Error(w, "Failed to read request body", http.StatusBadRequest)
					return
				}
				// Restore the body for subsequent handlers
				r.Body = io.NopCloser(strings.NewReader(string(body)))
			}

			// Create request validation input
			input := &openapi3filter.RequestValidationInput{
				Request:    r,
				PathParams: pathParams,
				Route:      route,
			}

			// Validate request
			if err := openapi3filter.ValidateRequest(r.Context(), input); err != nil {
				v.sendErrorResponse(w, http.StatusBadRequest, "Request validation failed", err)
				return
			}

			// Create response recorder to capture response
			recorder := &responseRecorder{ResponseWriter: w}

			// Call next handler
			next.ServeHTTP(recorder, r)

			// Validate response if status code is documented
			if recorder.statusCode > 0 {
				if err := v.validateResponse(r, route, recorder.statusCode, recorder.body); err != nil {
					// Log response validation error but don't fail the request
					// In production, you might want to log this to monitoring
					fmt.Printf("Response validation error: %v\n", err)
				}
			}
		})
	}
}

// validateResponse validates the response against the OpenAPI specification
func (v *Validator) validateResponse(r *http.Request, route *routers.Route, statusCode int, body []byte) error {
	// Find the operation
	operation := route.Operation
	if operation == nil {
		return fmt.Errorf("operation not found")
	}

	// Find the response schema
	response := operation.Responses.Status(statusCode)
	if response == nil {
		// Status code not documented in OpenAPI spec
		return nil
	}

	// If no content or schema, skip validation
	if response.Ref == "" && response.Value.Content == nil {
		return nil
	}

	// Get JSON schema
	content := response.Value.Content.Get("application/json")
	if content == nil || content.Schema == nil {
		// No JSON content expected
		return nil
	}

	// Parse response body
	var data interface{}
	if len(body) > 0 {
		if err := json.Unmarshal(body, &data); err != nil {
			return fmt.Errorf("failed to parse response JSON: %w", err)
		}
	} else {
		data = nil
	}

	// Validate against schema
	if err := content.Schema.Value.VisitJSON(data); err != nil {
		return fmt.Errorf("response schema validation failed: %w", err)
	}

	return nil
}

// sendErrorResponse sends a standardized error response
func (v *Validator) sendErrorResponse(w http.ResponseWriter, statusCode int, message string, err error) {
	errorResponse := map[string]interface{}{
		"error": message,
	}

	if err != nil {
		errorResponse["details"] = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse)
}

// Enable enables OpenAPI validation
func (v *Validator) Enable() {
	v.enabled = true
}

// Disable disables OpenAPI validation
func (v *Validator) Disable() {
	v.enabled = false
}

// responseRecorder captures the response for validation
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       []byte
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	if r.statusCode == 0 {
		r.statusCode = http.StatusOK
	}
	r.body = append(r.body, b...)
	return r.ResponseWriter.Write(b)
}

// ValidateRequestPayload validates a request payload against a schema
func (v *Validator) ValidateRequestPayload(r *http.Request, schemaRef *openapi3.SchemaRef, payload []byte) error {
	if schemaRef == nil {
		return nil
	}

	var data interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return fmt.Errorf("failed to parse payload: %w", err)
	}

	if err := schemaRef.Value.VisitJSON(data); err != nil {
		return fmt.Errorf("payload validation failed: %w", err)
	}

	return nil
}

// ValidateResponsePayload validates a response payload against a schema
func (v *Validator) ValidateResponsePayload(schemaRef *openapi3.SchemaRef, payload []byte) error {
	if schemaRef == nil {
		return nil
	}

	var data interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if err := schemaRef.Value.VisitJSON(data); err != nil {
		return fmt.Errorf("response validation failed: %w", err)
	}

	return nil
}

// GetSchema returns a schema by reference
func (v *Validator) GetSchema(ref string) (*openapi3.SchemaRef, error) {
	if v.doc.Components == nil || v.doc.Components.Schemas == nil {
		return nil, fmt.Errorf("no schemas found in OpenAPI document")
	}
	schema, ok := v.doc.Components.Schemas[ref]
	if !ok {
		return nil, fmt.Errorf("schema not found: %s", ref)
	}
	return schema, nil
}

// GetOperation returns an operation by method and path
func (v *Validator) GetOperation(method, path string) (*openapi3.Operation, error) {
	// First try direct path lookup (for paths without parameters)
	pathItem := v.doc.Paths.Value(path)
	if pathItem != nil {
		var operation *openapi3.Operation
		switch strings.ToUpper(method) {
		case http.MethodGet:
			operation = pathItem.Get
		case http.MethodPost:
			operation = pathItem.Post
		case http.MethodPut:
			operation = pathItem.Put
		case http.MethodDelete:
			operation = pathItem.Delete
		case http.MethodPatch:
			operation = pathItem.Patch
		default:
			return nil, fmt.Errorf("unsupported method: %s", method)
		}

		if operation == nil {
			return nil, fmt.Errorf("operation not found: %s %s", method, path)
		}

		return operation, nil
	}

	// If direct lookup fails, try path template matching
	return v.findOperationByPathTemplate(method, path)
}

// findOperationByPathTemplate implements simple path template matching
func (v *Validator) findOperationByPathTemplate(method, path string) (*openapi3.Operation, error) {
	// Split the request path into segments
	requestSegments := strings.Split(strings.Trim(path, "/"), "/")

	// Iterate through all paths in the spec
	for specPath, pathItem := range v.doc.Paths.Map() {
		// Split the spec path into segments
		specSegments := strings.Split(strings.Trim(specPath, "/"), "/")

		// If segment counts don't match, skip this path
		if len(requestSegments) != len(specSegments) {
			continue
		}

		// Check if segments match (considering path parameters)
		segmentsMatch := true
		for i := 0; i < len(specSegments); i++ {
			specSegment := specSegments[i]
			requestSegment := requestSegments[i]

			// If spec segment is a path parameter (starts with {), it matches anything
			if strings.HasPrefix(specSegment, "{") && strings.HasSuffix(specSegment, "}") {
				continue
			}

			// Otherwise, segments must match exactly
			if specSegment != requestSegment {
				segmentsMatch = false
				break
			}
		}

		if segmentsMatch {
			// Found matching path template, now get the operation
			var operation *openapi3.Operation
			switch strings.ToUpper(method) {
			case http.MethodGet:
				operation = pathItem.Get
			case http.MethodPost:
				operation = pathItem.Post
			case http.MethodPut:
				operation = pathItem.Put
			case http.MethodDelete:
				operation = pathItem.Delete
			case http.MethodPatch:
				operation = pathItem.Patch
			default:
				return nil, fmt.Errorf("unsupported method: %s", method)
			}

			if operation == nil {
				return nil, fmt.Errorf("operation not found: %s %s", method, specPath)
			}

			return operation, nil
		}
	}

	return nil, fmt.Errorf("path not found: %s", path)
}

// Doc returns the OpenAPI document
func (v *Validator) Doc() *openapi3.T {
	return v.doc
}
