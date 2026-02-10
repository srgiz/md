.DEFAULT_GOAL := help
.PHONY: help build test # fix "is up to date"

#GIT_SHORT_SHA := $(shell git rev-parse --short HEAD)

help: ## Show this help (default)
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

migrate: ## Migrate the DB to the most recent version available
	@goose up

down-migration: ## Roll back the version by 1
	@goose down

dev-db-clean: ## Recreate dev database
	@goose down-to 0
	@goose up

test: ## Run tests
	@go test ./...

coverage-html: ## Coverage html
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

gen-mock: ## Generate mocks
	@mockgen -source=internal/lib/kernel/db/conn.go -destination=internal/lib/kernel/db/conn_mock.go -package=db

#wire:
#	@wire ./internal/io/http
#	@wire ./internal/io/cli

app := 'http'

#build: # Build application (app=http|cli)
#	@wire ./internal/io/$(app)
#	@go build -o ./build/$(app) ./cmd/$(app)
#	@go build -ldflags "-X 'exampleapp/internal/infrastructure/di.Version=$(GIT_SHORT_SHA)'" -o ./build/$(app) ./cmd/$(app)

#wirecontext:
#	@wire ./internal/UserCtx/Present/Cmd/

build: ## Build application (app=http|cli)
	@go build -o ./build/$(app) ./cmd/$(app)
