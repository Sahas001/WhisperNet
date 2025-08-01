# Makefile - WhisperNet Project

# Directories
BACKEND_DIR := backend
FRONTEND_DIR := frontend

# Backend build output
BACKEND_BIN := $(BACKEND_DIR)/whisperd

.PHONY: all build build-backend build-frontend run-backend run-frontend run clean

all: build

# Build everything
build: build-backend build-frontend
	@echo "All parts built successfully."

# Build backend Go daemon
build-backend:
	@echo "🔧 Building backend..."
	cd $(BACKEND_DIR) && go build -o whisperd ./cmd/whisperd
	@echo "Backend built."

# Build frontend Next.js app (static export)
build-frontend:
	@echo "Building frontend..."
	cd $(FRONTEND_DIR) && npm install
	cd $(FRONTEND_DIR) && npm run build
	@echo "Frontend built."

# Run backend server
run-backend:
	@echo "Starting backend..."
	cd $(BACKEND_DIR) && go run ./cmd/whisperd/main.go

# Run frontend dev server
run-frontend:
	@echo "Starting frontend dev server..."
	cd $(FRONTEND_DIR) && npm run dev

# Run both backend and frontend concurrently (requires 'concurrently' npm package globally or installed)
run:
	@echo "Starting both backend and frontend concurrently..."
	concurrently "make run-backend" "make run-frontend"

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	rm -f $(BACKEND_BIN)
	rm -rf $(FRONTEND_DIR)/.next
	rm -rf $(FRONTEND_DIR)/node_modules
	@echo "Clean complete."

