BIN ?= glide
GO ?= go
CMD_DIR ?= ./cmd/glide
LDFLAGS ?= $(if $(VERSION),-X github.com/muniere/glide/internal/cli.Version=$(VERSION))

.PHONY: all help build run fmt tidy test check clean

all: build

help:
	@printf '%s\n' \
		'make build            Build ./$(BIN)' \
		'make run ARGS="..."   Run $(CMD_DIR) with arguments' \
		'make fmt              Format Go code' \
		'make tidy             Tidy Go modules' \
		'make test             Run tests' \
		'make check            fmt + test' \
		'make clean            Remove build output'

build:
	$(GO) build -ldflags "$(LDFLAGS)" -o $(BIN) $(CMD_DIR)

run:
	$(GO) run $(CMD_DIR) $(ARGS)

fmt:
	$(GO) fmt ./...

tidy:
	$(GO) mod tidy

test:
	$(GO) test ./...

check: fmt test

clean:
	rm -f $(BIN)
