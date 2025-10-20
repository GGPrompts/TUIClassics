.PHONY: all minesweeper solitaire classics install clean test run-minesweeper run-solitaire run-classics

all: minesweeper solitaire classics

minesweeper:
	@echo "Building minesweeper..."
	@go build -o bin/minesweeper ./cmd/minesweeper

solitaire:
	@echo "Building solitaire..."
	@go build -o bin/solitaire ./cmd/solitaire

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

run-classics:
	@go run ./cmd/classics

dev-minesweeper:
	@go run ./cmd/minesweeper

build-all: clean all
	@echo "All games built successfully!"
