# repo-ready build tooling (FOUND-FR-01).
#
# Every target wraps the Go toolchain directly so a clean checkout builds with
# nothing but Go installed. The module is dependency-free at this stage.

GO      ?= go
BINARY  ?= repo-ready
CMD     ?= ./cmd/repo-ready
BIN_DIR ?= bin

# Build metadata is injected into internal/version.Version (FOUND-FR-03).
# VERSION may be overridden on the command line, e.g. `make build VERSION=v1.2.3`.
VERSION ?= dev
LDFLAGS := -X github.com/mcfuzzysquirrel/repo-ready/internal/version.Version=$(VERSION)

.PHONY: build test vet fmt

# build compiles a single static binary (CGO_ENABLED=0, CONST-03) with the
# version injected via -ldflags.
build:
	CGO_ENABLED=0 $(GO) build -trimpath -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY) $(CMD)

# test runs the full unit test suite.
test:
	$(GO) test ./...

# vet runs the Go static analyzer over every package.
vet:
	$(GO) vet ./...

# fmt rewrites Go source files to canonical gofmt form.
fmt:
	gofmt -w .
