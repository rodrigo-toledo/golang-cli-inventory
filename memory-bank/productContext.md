# Product Context

What this project is
- A CLI-first inventory management tool to create and track products, locations, stock levels, and stock movements with an audit log.
- Intended to be lightweight, maintainable, and a learning resource demonstrating clean architecture and test strategies.

Primary users
- Developers (primary): explore, extend, and test the codebase.
- Power users / ops (secondary): CLI users who need a local inventory tool or who want to script inventory tasks.

Problems it solves
- Provide a simple, auditable way to track inventory across multiple locations.
- Offer a reproducible developer environment to test DB interactions and business rules.
- Reduce boilerplate around SQL queries via sqlc-generated type-safe code.

Core user flows
1. Add a product (sku, name, description, price).
2. Add stock to a specific location for an existing product.
3. Move stock atomically between locations with audit entries.
4. Find product details by SKU and list stock by location.
5. Generate reports (e.g., low-stock) for operational awareness.

User experience goals
- Clear, concise CLI output for commands (successful operations and descriptive errors).
- Fast developer feedback via unit and integration tests.
- Minimal friction for local setup: Docker + docker-compose, `make` targets, and `sqlc generate`.
