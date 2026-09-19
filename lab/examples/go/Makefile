BINARY      := orderctl
IMAGE       := hexslice-order-example
GO_VERSION  := 1.27.1
GOLANGCI    := v2.13.2

DEV_IMAGE    := $(IMAGE):dev
EXPORT_IMAGE := $(IMAGE):export
BUILD_ARGS   := --build-arg GO_VERSION=$(GO_VERSION) --build-arg GOLANGCI_VERSION=$(GOLANGCI)

# Everything runs in Docker — no local Go toolchain and no bind mounts. Source
# is baked into the image via COPY; build artifacts are streamed out over stdout
# by the `export` stage's entrypoint (`cat`).

.DEFAULT_GOAL := help

# Digest-pin the architecture gate here (overrides the fallback baked into
# a-check.mk). Bump this per release; override locally with
#   make a-check A_CHECK_IMAGE=a-check:dev
# Pinned to a-check v0.15.0 (adds the HexSlice rules lateral-slice & port-locality).
A_CHECK_IMAGE ?= ghcr.io/pt9912/a-check@sha256:6425c93a9a4359ef28c4da231a2d1db6f421fdaa8f96877ac89d201827c42d09

include a-check.mk

.PHONY: help dev-image image build run test vet fmt lint tidy check clean

# `grep -h` is load-bearing: $(MAKEFILE_LIST) holds more than one file once
# a-check.mk is included, and without -h every hit would carry a filename prefix
# that lands in the target column instead of the target name.
help: ## Show this help
	@grep -hE '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | \
		awk 'BEGIN{FS=":.*?## "}{printf "  \033[36m%-13s\033[0m %s\n", $$1, $$2}'

dev-image: ## Build the dev image (toolchain + source + linter)
	docker build --target dev $(BUILD_ARGS) -t $(DEV_IMAGE) .

image: ## Build the minimal runtime image
	docker build --target run $(BUILD_ARGS) -t $(IMAGE) .

build: ## Compile the binary and extract it to ./bin (no bind mount)
	@mkdir -p bin
	docker build --target export $(BUILD_ARGS) -t $(EXPORT_IMAGE) .
	docker run --rm $(EXPORT_IMAGE) > bin/$(BINARY).tar
	tar -xf bin/$(BINARY).tar -C bin
	@rm -f bin/$(BINARY).tar
	@echo "built bin/$(BINARY)"

run: image ## Build and run the demo scenario in the runtime image
	docker run --rm $(IMAGE) demo

test: dev-image ## Run all tests in Docker
	docker run --rm $(DEV_IMAGE) go test ./...

vet: dev-image ## Run go vet in Docker
	docker run --rm $(DEV_IMAGE) go vet ./...

fmt: dev-image ## Check formatting with gofmt in Docker
	docker run --rm $(DEV_IMAGE) sh -c 'gofmt -l . && test -z "$$(gofmt -l .)"'

lint: dev-image ## Run golangci-lint (pinned $(GOLANGCI)) in Docker
	docker run --rm $(DEV_IMAGE) golangci-lint run

tidy: dev-image ## Verify go.mod is tidy in Docker
	docker run --rm $(DEV_IMAGE) sh -c 'cp go.mod /tmp/go.mod.bak && go mod tidy && diff -u /tmp/go.mod.bak go.mod'

check: fmt vet lint test ## Run fmt, vet, lint and tests

clean: ## Remove build artifacts
	rm -rf bin
