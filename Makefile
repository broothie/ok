
build: generate
	go build cmd/ok/ok.go

install: generate
	go install cmd/ok/ok.go

task/type_string.go: task/type.go
	go generate ./...

generate: task/type_string.go
