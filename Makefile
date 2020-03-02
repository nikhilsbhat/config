GOFMT_FILES?=$$(find . -not -path "./vendor/*" -type f -name '*.go')

default: check

# bin generates the releaseable binaries for config
build: check
	go build -ldflags="-s -w"

# Validates the dependencies exists.
check: fmt
	go mod vendor

# Lints all the go files if they are not.
fmt:
	gofmt -w $(GOFMT_FILES)

# Runs the application by building a binary of it.
run: build
	./config