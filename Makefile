.PHONY: help lint tidy test build fmt

lint: ## Linter
	@golangci-lint run -v ./... --fix

fmt:
	@golangci-lint fmt

build:
	@go build -v ./...

tidy: ## Download latest go module dependencies
	@go mod tidy

test: ## Run all tests
	@go test ./... -cover -count=1 -timeout 5s -race
