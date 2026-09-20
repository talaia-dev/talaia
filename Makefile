BINARY  := talaia
PREFIX  ?= /usr/local
GOBIN   := $(shell go env GOPATH)/bin

.PHONY: build install install-user test fmt vet clean

build:
	go build -o bin/$(BINARY) ./cmd/talaia

install: build
	install -m 0755 bin/$(BINARY) $(PREFIX)/bin/$(BINARY)

install-user:
	go install ./cmd/talaia
	@echo "Installed to $(GOBIN)/$(BINARY) (make sure it is in your PATH)"

test:
	go test ./...

fmt:
	gofmt -w ./cmd ./internal

vet:
	go vet ./...

clean:
	rm -rf bin
