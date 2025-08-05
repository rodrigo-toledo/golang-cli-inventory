package main

import (
	"fmt"
	"log"
	"os"

	"cli-inventory/internal/openapi"
)

func main() {
	validator, err := openapi.NewValidator("api/openapi.yaml")
	if err != nil {
		log.Fatalf("Failed to initialize OpenAPI validator: %v", err)
	}

	fmt.Println("✅ OpenAPI specification is valid")

	// Print some statistics about the spec
	if validator.Doc() != nil && validator.Doc().Paths != nil {
		fmt.Printf("📊 API Statistics:\n")
		fmt.Printf("   - Total paths: %d\n", len(validator.Doc().Paths.Map()))

		operationCount := 0
		for _, pathItem := range validator.Doc().Paths.Map() {
			if pathItem.Get != nil {
				operationCount++
			}
			if pathItem.Post != nil {
				operationCount++
			}
			if pathItem.Put != nil {
				operationCount++
			}
			if pathItem.Delete != nil {
				operationCount++
			}
			if pathItem.Patch != nil {
				operationCount++
			}
		}
		fmt.Printf("   - Total operations: %d\n", operationCount)

		if validator.Doc().Components != nil && validator.Doc().Components.Schemas != nil {
			fmt.Printf("   - Total schemas: %d\n", len(validator.Doc().Components.Schemas))
		}
	}

	os.Exit(0)
}
