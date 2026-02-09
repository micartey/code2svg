default:
    @just --list

# Run the server
run: build
    go run ./cmd/code2svg

# Run tests
test:
    go test -v ./...

# Build the server
build:
    go build -o code2svg ./cmd/code2svg
