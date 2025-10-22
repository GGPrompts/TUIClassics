.PHONY: all minesweeper solitaire 2048 classics install clean test run-minesweeper run-solitaire run-2048 run-classics

all: minesweeper solitaire 2048 classics

minesweeper:
	@echo "Building minesweeper..."
	@go build -o bin/minesweeper ./cmd/minesweeper

solitaire:
	@echo "Building solitaire..."
	@go build -o bin/solitaire ./cmd/solitaire

2048:
	@echo "Building 2048..."
	@go build -o bin/2048 ./cmd/2048

classics:
	@echo "Building classics launcher..."
	@go build -o bin/classics ./cmd/classics

install:
	@echo "Installing all games..."
	@go install ./cmd/...

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/

test:
	@echo "Running tests..."
	@go test ./...

run-minesweeper:
	@go run ./cmd/minesweeper

run-solitaire:
	@go run ./cmd/solitaire

run-2048:
	@go run ./cmd/2048

run-classics:
	@go run ./cmd/classics

dev-minesweeper:
	@go run ./cmd/minesweeper

build-all: clean all
	@echo "All games built successfully!"
