# Cryptotrade Ecommerce API

## Overview
Cryptotrade is a layered ecommerce REST API built with [Gin](https://gin-gonic.com/). It demonstrates how to structure a Go service around clearly separated domain, repository, service, and transport concerns while remaining lightweight enough to run with in-memory storage. The application exposes endpoints for managing products, registering customers, and placing orders with simple inventory checks.

The project intentionally ships with an in-memory data store so it can be explored without external dependencies. You can swap in persistent repositories later without touching the HTTP layer or business logic.

## Architecture
The codebase follows a hexagonal-inspired organization where each layer has a single responsibility:

| Layer | Location | Responsibilities |
| --- | --- | --- |
| Domain | [`internal/domain`](internal/domain) | Defines core entities (`Product`, `User`, `Order`) and validation rules that protect invariants before data is persisted. |
| Repository | [`internal/repository`](internal/repository) | Declares storage interfaces and provides an in-memory implementation guarded by mutexes for safe concurrent access. |
| Service | [`internal/service`](internal/service) | Contains business use cases such as enforcing uniqueness, applying validation, managing stock levels, and translating errors into domain-specific failures. |
| HTTP Handlers | [`internal/handler`](internal/handler) | Maps services onto Gin routes, handles input binding, and normalizes error responses for clients. |
| Router | [`internal/router`](internal/router/router.go) | Centralizes Gin engine creation, middleware, API grouping, and the `/health` endpoint. |
| Configuration & Bootstrap | [`internal/config`](internal/config) & [`main.go`](main.go) | Loads environment-based configuration, wires dependencies, and starts the HTTP server with graceful shutdown. |

Because each layer depends only on the layer directly beneath it, you can swap implementations (for example, replacing the memory repositories with a database-backed package) without rewriting business or transport logic.

## API Surface
All business endpoints are namespaced under `/api/v1`:

| Method | Path | Description |
| --- | --- | --- |
| `GET` | `/health` | Liveness probe returning application status. |
| `GET` | `/api/v1/products` | List all products. |
| `POST` | `/api/v1/products` | Create a product (requires `name`, `price`, optional `description`, `stock`). |
| `GET` | `/api/v1/products/:id` | Fetch a product by ID. |
| `PUT` | `/api/v1/products/:id` | Update product details. |
| `DELETE` | `/api/v1/products/:id` | Remove a product. |
| `GET` | `/api/v1/users` | List registered users. |
| `POST` | `/api/v1/users` | Create a user (valid email required). |
| `GET` | `/api/v1/users/:id` | Fetch a user by ID. |
| `GET` | `/api/v1/orders` | List orders. |
| `POST` | `/api/v1/orders` | Create an order for an existing user with product line items. |
| `GET` | `/api/v1/orders/:id` | Fetch an order by ID. |

Orders automatically validate the requesting user, confirm product availability, reserve stock, and calculate totals before persisting the purchase. All persistence happens in-memory, so restarting the service clears state.

## Running Locally
### Prerequisites
* Go 1.23+ (Go toolchain 1.24 is configured in [`go.mod`](go.mod))

### Using the Makefile
This repository includes a Makefile to streamline common tasks:

```bash
make help        # Show available commands
make run         # Start the Gin server on http://localhost:8080
make test        # Execute the full test suite
make fmt         # Run gofmt via go fmt on all packages
make vet         # Static analysis using go vet
make tidy        # Clean go.mod/go.sum
```

Environment variables can be overridden per invocation, for example:

```bash
PORT=9090 APP_ENV=production make run
```

### Manual commands
If you prefer not to use Make, you can run the service and tests directly:

```bash
go run .
go test ./...
```

## Configuration
Runtime configuration is read from environment variables with sensible defaults:

| Variable | Default | Purpose |
| --- | --- | --- |
| `APP_ENV` | `development` | Controls Gin mode (release mode when set to `production`). |
| `PORT` | `8080` | Port the HTTP server listens on (prefixed with `:` internally). |

## Sample Workflow
1. Start the server (`make run`).
2. Create a user:
   ```bash
   curl -X POST http://localhost:8080/api/v1/users \
     -H 'Content-Type: application/json' \
     -d '{"name":"Ada Lovelace","email":"ada@example.com"}'
   ```
3. Create a product:
   ```bash
   curl -X POST http://localhost:8080/api/v1/products \
     -H 'Content-Type: application/json' \
     -d '{"name":"Laptop","description":"Developer laptop","price":1999.99,"stock":5}'
   ```
4. Place an order using the IDs returned above:
   ```bash
   curl -X POST http://localhost:8080/api/v1/orders \
     -H 'Content-Type: application/json' \
     -d '{"user_id":"<user-id>","items":[{"product_id":"<product-id>","quantity":1}]}'
   ```

Because the repositories are in-memory, repeating the process from a clean start ensures consistent results without lingering state.

## Development Notes
* Error responses follow a consistent JSON contract and map validation failures, conflicts, and missing resources to appropriate HTTP status codes.
* Graceful shutdown waits up to 10 seconds for in-flight requests before terminating the server.
* The service layer composes repositories rather than accessing them directly from handlers, simplifying future upgrades to persistent storage or background processing.

Happy building!
