entrypoint=main.go

# Check all is okay
check: test build clean

# Build locally into _ok
build: generate
	go build -o _ok $(entrypoint)

# Run tests with coverage
test: generate
	go test -cover ./...

# Install ok locally
install: generate
	go install $(entrypoint)

# Generate target
task/type_string.go: task/type.go
	go generate ./...

# go generate
generate: task/type_string.go

# Cleans local dir
clean:
	go clean
	rm -rf _ok dist
	go mod tidy

# Bumps VERSION file
bump:
	bump VERSION
