
build: generate
	go build cmd/now/now.go

install: generate
	go install cmd/now/now.go

param/type_string.go:
	go generate ./...

generate: param/type_string.go
