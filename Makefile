GO ?= go
COVER_FILE ?= coverage.out

.DEFAULT_GOAL := help

.PHONY: help
help:
	@echo "Usage: make <target>"
	@echo
	@echo "Targets:"
	@echo "  build         go build ./..."
	@echo "  vet           go vet ./..."
	@echo "  test          go test -race -count=1 ./..."
	@echo "  test-cover    go test with coverage profile"
	@echo "  cover-html    open coverage HTML"
	@echo "  cgo0-build    CGO_ENABLED=0 go build ./..."
	@echo "  tidy          go mod tidy"
	@echo "  check         build + vet + test-cover + cgo0-build"
	@echo "  clean         remove coverage and build artifacts"

.PHONY: build
build:
	$(GO) build ./...

.PHONY: vet
vet:
	$(GO) vet ./...

.PHONY: test
test:
	$(GO) test -race -count=1 ./...

.PHONY: test-cover
test-cover:
	$(GO) test -race -count=1 -coverprofile=$(COVER_FILE) ./...
	@$(GO) tool cover -func=$(COVER_FILE) | tail -1

.PHONY: cover-html
cover-html: test-cover
	$(GO) tool cover -html=$(COVER_FILE)

.PHONY: cgo0-build
cgo0-build:
	CGO_ENABLED=0 $(GO) build ./...

.PHONY: tidy
tidy:
	$(GO) mod tidy

.PHONY: check
check: build vet test-cover cgo0-build

.PHONY: clean
clean:
	rm -f $(COVER_FILE) coverage.html
