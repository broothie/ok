
check: test build
	rm -rf ok

build: generate
	go build cmd/ok/ok.go

test: generate
	go test -cover ./...

install: generate
	go install cmd/ok/ok.go

task/type_string.go: task/type.go
	go generate ./...

generate: task/type_string.go

clean:
	go clean
	rm -rf ok
