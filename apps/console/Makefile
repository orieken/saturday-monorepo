# Simple Makefile for final/backend

SHELL := /bin/bash

# Config
BINARY      ?= bin/server
MAIN_PKG    ?= ./cmd/server
PKG         ?= ./...
COVER_FILE  ?= coverage.out
COVER_HTML  ?= coverage.html

.PHONY: help deps build run test cover cover-html fmt vet tidy clean

help:
	@echo "Available targets:"
	@echo "  deps        - Ensure module deps are up to date (go mod tidy)"
	@echo "  build       - Build the backend server binary ($(BINARY))"
	@echo "  run         - Run the backend server (go run ./cmd/server)"
	@echo "  test        - Run unit tests with coverage (writes $(COVER_FILE))"
	@echo "  cover       - Print coverage summary from $(COVER_FILE)"
	@echo "  cover-html  - Generate HTML coverage report ($(COVER_HTML))"
	@echo "  fmt         - go fmt ./..."
	@echo "  vet         - go vet ./..."
	@echo "  tidy        - go mod tidy"
	@echo "  clean       - Remove build artifacts"

deps: tidy

tidy:
	go mod tidy

build:
	@mkdir -p $(dir $(BINARY))
	go build -o $(BINARY) $(MAIN_PKG)

run:
	go run $(MAIN_PKG)

test:
	go test $(PKG) -cover -coverprofile=$(COVER_FILE)

cover:
	@test -f $(COVER_FILE) || (echo "$(COVER_FILE) not found. Run 'make test' first." && exit 1)
	go tool cover -func=$(COVER_FILE)

cover-html:
	@test -f $(COVER_FILE) || (echo "$(COVER_FILE) not found. Run 'make test' first." && exit 1)
	go tool cover -html=$(COVER_FILE) -o $(COVER_HTML)
	@echo "HTML report: $(COVER_HTML)"

fmt:
	go fmt $(PKG)

vet:
	go vet $(PKG)

clean:
	rm -rf bin $(COVER_FILE) $(COVER_HTML)
