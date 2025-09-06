// Package main provides the command-line interface for the inventory management system.
// It initializes the application, sets up dependencies, and delegates command handling
// to the CLI package which uses Cobra for robust command parsing.
package main

import (
	"cli-inventory/internal/cli"
)

// main initializes the application and starts the CLI.
// The database connection and services are initialized lazily when needed.
func main() {
	// Execute the CLI application
	// Database and services will be initialized when commands are executed
	cli.Execute()
}
