SHELL := /bin/bash
.DEFAULT_GOAL := help

APP_NAME := echo
MAIN_PACKAGE := ./cmd/echo
BIN_DIR := bin
BIN_TARGET := $(BIN_DIR)/$(APP_NAME)
SOURCE_PACKAGES := ./cmd/... ./internal/...
TEST_PACKAGES := ./tests/...
ALL_PACKAGES := $(SOURCE_PACKAGES) $(TEST_PACKAGES)
GO ?= go
ARGS ?=

.PHONY: help build run test test-all test-verbose fmt vet tidy clean check

help:
	@printf "Available targets:\n"
	@printf "  make build        Build the binary into %s/\n" "$(BIN_DIR)"
	@printf "  make run          Run the app with go run\n"
	@printf "  make test         Run tests in %s\n" "$(TEST_PACKAGES)"
	@printf "  make test-all     Run tests in all packages\n"
	@printf "  make test-verbose Run tests in %s with verbose output\n" "$(TEST_PACKAGES)"
	@printf "  make fmt          Format Go source files\n"
	@printf "  make vet          Run go vet on all packages\n"
	@printf "  make tidy         Tidy module dependencies\n"
	@printf "  make clean        Remove build artifacts\n"
	@printf "  make check        Run fmt, vet, and test\n"

build:
	@mkdir -p "$(BIN_DIR)"
	@printf "Building %s...\n" "$(APP_NAME)"
	"$(GO)" build -o "$(BIN_TARGET)" "$(MAIN_PACKAGE)"

run:
	@printf "Running %s...\n" "$(APP_NAME)"
	"$(GO)" run "$(MAIN_PACKAGE)" $(ARGS)

test:
	@printf "Running tests in %s\n" "$(TEST_PACKAGES)"
	"$(GO)" test $(TEST_PACKAGES)

test-all:
	@printf "Running tests in all packages...\n"
	"$(GO)" test ./...

test-verbose:
	@printf "Running tests in %s with verbose output\n" "$(TEST_PACKAGES)"
	"$(GO)" test -v $(TEST_PACKAGES)

fmt:
	@printf "Formatting Go files...\n"
	"$(GO)" fmt $(ALL_PACKAGES)

vet:
	@printf "Running go vet...\n"
	"$(GO)" vet $(ALL_PACKAGES)

tidy:
	@printf "Tidying module dependencies...\n"
	"$(GO)" mod tidy

clean:
	@printf "Removing build artifacts...\n"
	rm -rf "$(BIN_DIR)"

check: fmt vet test
