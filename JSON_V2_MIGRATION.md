# JSON v2 Migration Guide

This document provides guidance on how to migrate the existing JSON usage in the project to the new experimental JSON v2 package introduced in Go 1.25.

## Overview

Go 1.25 introduces a new experimental JSON implementation with significant performance improvements and enhanced features. The new package is available at `encoding/json/v2` and can be used alongside the existing `encoding/json` package.

## Migration Steps

### 1. Enable the JSON v2 Experiment

To use the new JSON v2 package, you need to enable the `jsonv2` experiment:

```bash
GOEXPERIMENT=jsonv2 go build .
```

Or set it as an environment variable:

```bash
export GOEXPERIMENT=jsonv2
go build .
```

### 2. Update Import Statements

Replace the existing `encoding/json` import with `encoding/json/v2`:

```go
// Before
import (
    "encoding/json"
)

// After
import (
    "encoding/json/v2"
)
```

### 3. Update Code Usage

The API for the new JSON v2 package is largely compatible with the existing package. Most code should work without changes, but there are some enhancements and new features available.

#### Basic Usage (Decoder)

```go
// Before
var req models.CreateProductRequest
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    // handle error
}

// After (same code works with v2)
var req models.CreateProductRequest
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    // handle error
}
```

#### Basic Usage (Encoder)

```go
// Before
if err := json.NewEncoder(w).Encode(product); err != nil {
    // handle error
}

// After (same code works with v2)
if err := json.NewEncoder(w).Encode(product); err != nil {
    // handle error
}
```

### 4. Using New Features (Optional)

The JSON v2 package includes several new features that you can optionally use:

#### Streaming Parser

The new package provides more efficient streaming capabilities:

```go
import "encoding/json/v2/jsontext"

dec := jsontext.NewDecoder(r.Body)
for dec.PeekKind() != 0 { // 0 indicates EOF
    // Process JSON tokens incrementally
    tok, err := dec.ReadToken()
    if err != nil {
        // handle error
    }
    // Process token
}
```

#### Improved Error Handling

The new package provides more detailed error messages:

```go
err := json.Unmarshal(data, &v)
if err != nil {
    var serr *json.SemanticError
    if errors.As(err, &serr) {
        // Handle semantic errors specifically
        log.Printf("Semantic error at %s: %v", serr.JSONPointer(), serr.Error())
    }
}
```

## Files Updated

The following files in the project have been updated to use the new JSON v2 package:

1. `internal/handlers/error_handler.go`
2. `internal/handlers/location_handler.go`
3. `internal/handlers/location_handler_test.go`
4. `internal/handlers/product_handler.go`
5. `internal/handlers/product_handler_test.go`
6. `internal/handlers/stock_handler.go`
7. `internal/handlers/stock_handler_test.go`
8. `internal/openapi/validator.go`
9. `internal/testutils/openapi.go`

## Testing

After migration, run the tests to ensure everything works correctly:

```bash
GOEXPERIMENT=jsonv2 go test ./...
```

## Performance Benefits

The JSON v2 package provides significant performance improvements:

- Up to 2x faster encoding
- Up to 1.5x faster decoding
- Reduced memory allocations
- Better streaming performance

## Compatibility

The JSON v2 package maintains backward compatibility with the existing `encoding/json` package. All existing code should work without modification when using the new package.

However, some edge cases might behave differently:

1. Slight differences in error messages
2. More strict validation in some cases
3. Different handling of certain struct tags

## Rollback

If you need to rollback to the previous version:

1. Remove `GOEXPERIMENT=jsonv2` from your environment
2. Change import statements back to `encoding/json`
3. Rebuild the project

## References

- [Go 1.25 Release Notes - JSON v2](https://go.dev/doc/go1.25#json_v2)
- [JSON v2 Package Documentation](https://pkg.go.dev/encoding/json/v2)
- [JSON Text Package Documentation](https://pkg.go.dev/encoding/json/v2/jsontext)