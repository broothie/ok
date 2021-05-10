
build: generate
	go build cmd/ok/ok.go

install: generate
	go install cmd/ok/ok.go

param/type_string.go:
	go generate ./...

generate: param/type_string.go
