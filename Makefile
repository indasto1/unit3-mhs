.PHONY: build

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

clean: ## Remove go object files from package source directories.
	go clean -cache -testcache -r

build: ## Build consumer and producer services.
	docker compose build

test: ## Launch all unit tests.
	go test -v ./...

run: ## Start consumer, producer and kafka services.
	docker compose up -d

install: clean test build run ## Test, build and start all services.

down: ## Stop all services.
	docker compose down

restart: down run ## Restart all services.

fmt: ## Run go formatter on all project's files
	gofmt -s -w .
