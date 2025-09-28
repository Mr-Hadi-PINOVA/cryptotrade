APP_NAME := cryptotrade
GO ?= go
APP_ENV ?= development
PORT ?= 8080

.PHONY: run test fmt vet tidy help

run:
	APP_ENV=$(APP_ENV) PORT=$(PORT) $(GO) run .

test:
	$(GO) test ./...

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

tidy:
	$(GO) mod tidy

help:
	@echo "Available targets:"
	@echo "  run   - Start the API server (override APP_ENV and PORT as needed)"
	@echo "  test  - Run go test across all packages"
	@echo "  fmt   - Format Go source files with gofmt"
	@echo "  vet   - Static analysis with go vet"
	@echo "  tidy  - Update go.mod and go.sum dependencies"
