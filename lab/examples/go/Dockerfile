# syntax=docker/dockerfile:1

ARG GO_VERSION=1.27.1
ARG GOLANGCI_VERSION=v2.13.2

# ---- base: toolchain + module cache ---------------------------------------
FROM golang:${GO_VERSION}-alpine AS base
WORKDIR /src
COPY go.mod ./
RUN go mod download

# ---- lint tool: pinned golangci-lint binary (no compile in dev) -----------
FROM golangci/golangci-lint:${GOLANGCI_VERSION} AS golangci

# ---- dev: linter + full source; runs test/vet/fmt/lint (no bind mounts) ---
# The source is baked in via COPY, so `docker run <dev-image> go test ./...`
# works without mounting anything.
FROM base AS dev
COPY --from=golangci /usr/bin/golangci-lint /usr/local/bin/golangci-lint
COPY . .

# ---- build: compile a fully static binary ---------------------------------
FROM dev AS build
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" \
	-o /out/orderctl ./cmd/orderctl

# ---- export: hand the binary out of the build via `cat` (no bind mount) ----
# `docker run --rm <export-image> > orderctl.tar` streams the artifact to the
# host over stdout — nothing is mounted into the container.
FROM build AS export
RUN tar -C /out -cf /orderctl.tar orderctl
ENTRYPOINT ["cat", "/orderctl.tar"]

# ---- run: minimal runtime image -------------------------------------------
FROM scratch AS run
COPY --from=build /out/orderctl /orderctl
ENTRYPOINT ["/orderctl"]
CMD ["demo"]
