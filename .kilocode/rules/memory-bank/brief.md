# Project Brief: CLI Inventory Management System

## Purpose

This project is a command-line interface (CLI) application for managing inventory. It is intended to be a small, maintainable, and well-tested application that serves as both a practical tool and a learning resource for Go developers.

## Core Features

- **Product Management**: Add, list, and find products by SKU.
- **Stock Management**: Add stock to locations and move stock between locations atomically.
- **Reporting**: Generate reports, such as a low-stock report.
- **HTTP API**: Exposes a RESTful API for all core functionalities.

## Key Technical Decisions

- **Language**: Go (version 1.25), using the experimental JSON v2 package.
- **Architecture**: Clean, layered architecture (CLI -> Service -> Repository -> Database).
- **Database**: PostgreSQL, managed with Docker.
- **Database Access**: `sqlc` for generating type-safe Go code from raw SQL queries.
- **CLI Framework**: `cobra` for building the command-line interface.
- **Testing**: `testify` for assertions and `dockertest` for integration tests against a real database.

## Project Goals

- Provide reliable CLI commands for all inventory operations.
- Maintain a strong separation of concerns between layers.
- Achieve high test coverage with both unit and integration tests.
- Ensure a simple and reproducible developer setup using Docker and Makefiles.