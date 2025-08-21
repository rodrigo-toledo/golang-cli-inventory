# Product Context

## Why this project exists

This project serves two primary purposes:
1.  It is a practical, lightweight, and auditable command-line tool for managing inventory across multiple locations.
2.  It acts as a high-quality learning resource and reference implementation for Go developers, demonstrating best practices in clean architecture, testing, and database management.

## Problems it solves

-   **Inventory Tracking**: Provides a simple, scriptable interface to track products, stock levels, and movements, which is often a manual or spreadsheet-driven process in smaller operations.
-   **Developer Education**: Offers a clear, real-world example of a well-structured Go application, reducing the learning curve for concepts like layered architecture, `sqlc` for database access, and integration testing with Docker.
-   **Reproducibility**: Eliminates "it works on my machine" issues by providing a fully containerized development environment, ensuring consistency for all contributors.

## How it should work

The system functions as both a CLI tool and an HTTP server.

-   **CLI**: The primary interface for interaction. Users execute commands like `add-product`, `move-stock`, and `generate-report`. The CLI is designed to be clear, concise, and easily scriptable.
-   **HTTP API**: Exposes all core inventory functionalities through a RESTful API. This allows for integration with other systems or building a web-based frontend in the future.

All operations, especially stock movements, must be atomic and auditable to ensure data integrity.

## User experience goals

-   **For CLI Users**: The experience should be fast and intuitive. Commands must provide clear success messages and descriptive, actionable error messages.
-   **For Developers**: The setup process should be frictionless, with a "clone and run" simplicity. The codebase should be easy to navigate, with a logical structure that makes it obvious where to find and add functionality. The test suite should be fast and reliable, providing quick feedback during development.