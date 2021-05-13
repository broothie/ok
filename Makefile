
# Check all is okay
check: test build clean

build: generate
	go build -o _ok ok.go

test: generate
	go test -cover ./...

install: generate
	go install ok.go

task/type_string.go: task/type.go
	go generate ./...

generate: task/type_string.go

clean:
	go clean
	rm -rf _ok dist
