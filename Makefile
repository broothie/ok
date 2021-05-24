
# Check all is okay
check: test build clean

build: generate
	go build -o _ok ok.go

test: generate
	go test -cover ./...

# Installs ok locally
install: generate
	go install ok.go

task/type_string.go: task/type.go
	go generate ./...

generate: task/type_string.go

clean:
	go clean
	rm -rf _ok dist
	go mod tidy

other things:
	echo other things
