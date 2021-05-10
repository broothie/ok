
build: generate
	go build cmd/now/now.go

generate:
	go generate ./...
