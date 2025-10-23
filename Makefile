MAKEFLAGS += --silent

GOOSE_MIGRATION_DIR ?= ./migrations

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## create-migration: create a new migration
.PHONY: create-migration
create-migration:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make create-migration name=your_migration_name"; \
		exit 1; \
	fi; \
	GOOSE_MIGRATION_DIR=$(GOOSE_MIGRATION_DIR) goose create $(name) sql

## migrate-up: run all pending migrations
.PHONY: migrate-up
migrate-up:
	GOOSE_MIGRATION_DIR=$(GOOSE_MIGRATION_DIR) goose postgres "${DATABASE_URL}" up; \

## migrate-up-test: run all pending migrations on test db
.PHONY: migrate-up-test
migrate-up-test:
	GOOSE_MIGRATION_DIR=$(GOOSE_MIGRATION_DIR) goose postgres "${DATABASE_TEST_URL}" up; \

## migrate-down: migrate 1 down
.PHONY: migrate-down
migrate-down:
	GOOSE_MIGRATION_DIR=$(GOOSE_MIGRATION_DIR) goose postgres "${DATABASE_URL}" down; \

## migrate-down-test: migrate 1 down on test db
.PHONY: migrate-down-test
migrate-down-test:
	GOOSE_MIGRATION_DIR=$(GOOSE_MIGRATION_DIR) goose postgres "${DATABASE_TEST_URL}" down; \

## build: build the app
.PHONY: build
build:
	go build -o=./bin/movie-land ./

## run: run the app
.PHONY: run
run: build
	@./bin/movie-land

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## bench: run all benchmarks
.PHONY: bench
bench:
	go test -bench=. ./...

## tidy: tidy modfiles and format .go files
.PHONY: tidy
tidy:
	go mod tidy -v
	go fmt ./...

