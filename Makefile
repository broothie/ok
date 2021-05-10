
build: generate
	go build cmd/now/now.go

generate:
	go generate ./...

install:
	go install cmd/now/now.go
