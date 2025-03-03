# Build app for production
build:
	@echo "Building..."

	@go build -o bin/chippiphone cmd/main.go

# Run app for development
run:
	@go run cmd/main.go

all: build