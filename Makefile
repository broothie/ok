# Makefile

# Bump then tag
release: bump tag

# Bump VERSION
bump:
	bump VERSION

# Tag then push
tag:
	git tag "$$(cat VERSION)"
	git push --tags

# Build bin
build:
	go build ./cmd/ok
