APP_NAME := card-validation-api
MAIN_PKG := ./cmd/$(APP_NAME)
BUILD_DIR := bin

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY: build run test clean

.DEFAULT_GOAL: run

build:
	@echo "==> Building $(APP_NAME)..."
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PKG)

run:
	@echo "==> Running $(APP_NAME)..."
	@go run $(MAIN_PKG)

test:
	@echo "==> Running tests..."
	@go test -v -race -cover ./...

clean:
	@echo "==> Cleaning up..."
	@rm -rf $(BUILD_DIR)
