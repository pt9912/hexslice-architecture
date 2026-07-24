# HexSlice Example — Go

A small, dependency-free Go project that demonstrates the
[HexSlice Architecture](../../../docs/architecture/hexslice-architecture.md):
vertical use-case slices inside a hexagonal application core, with adapters
outside and ports owned by the application.

The example models a single business area — **order** — with two vertical
slices: **create order** and **cancel order**. An in-memory repository, a
sequential id generator, and a stdout notifier act as outbound adapters, and a
CLI acts as the inbound adapter.

## Requirements

* **Docker** only — no local Go toolchain is needed.

Every target runs inside Docker. The Go toolchain (pinned to **1.26.5**) and
`golangci-lint` (pinned to **v2.12.2**) live in the images; the source is baked
in via `COPY` — nothing is bind-mounted. Build artifacts are streamed out of the
container over stdout by the `export` stage's entrypoint (see `make build`).

## Quick start

```bash
make run        # build the runtime image and run the create -> cancel demo
make check      # fmt + vet + lint + test, all in Docker
make build      # compile and extract a static binary to ./bin/orderctl
make image      # build the minimal (scratch) runtime image

# architecture gate
make a-check        # fails on any hexagonal-layering violation
make a-check-graph  # prints the declared layers as a Mermaid flowchart
```

Run `make help` to list all targets. Individual checks are also available:
`make test`, `make vet`, `make fmt`, `make lint`, `make tidy`.

## Architecture gate (a-check)

The example is validated by [a-check](https://github.com/pt9912/a-check), a
containerized gate that verifies the hexagonal layering. The layer/edge mapping
lives in [`.a-check.yml`](.a-check.yml); the make targets come from the generated
[`a-check.mk`](a-check.mk), and the image digest is pinned in the `Makefile`.

`make a-check` returns exit code `0` (no violations) for this example. Unlike the
Go toolchain targets — which bake the source in via `COPY` and use no mounts —
a-check reads the tree through a **read-only** bind mount (`:ro`) and runs
**network-isolated** (`--network none`); this is a-check's own contract and it
never writes to the repository.

## CLI usage

Build the runtime image once (`make image`), then invoke subcommands against it,
or run the extracted binary from `make build` directly:

```bash
# self-contained lifecycle demo (default when no subcommand is given)
docker run --rm hexslice-order-example demo

# single use-case invocations
docker run --rm hexslice-order-example create -customer cust-1 \
  -line BOOK-1:2:1999:EUR -line PEN-7:5:150:EUR
docker run --rm hexslice-order-example cancel -id ORD-000001 -reason "changed mind"

# or, after `make build`, run the native binary
./bin/orderctl demo
```

> The repository is in-memory, so state does not survive across separate process
> invocations. Use the `demo` subcommand to see a full create → cancel lifecycle
> in one run.

## How the folders map to HexSlice

```text
cmd/orderctl/                     composition root (wires adapters to slices)
internal/
  hexagon/                        the application core
    domain/order/                 business core: Order aggregate, value objects,
                                  rules (New, Cancel, Total) — depends on nothing
    application/order/            the order business area
      ports/                      business-area shared port (OrderRepository)
      createorder/                vertical slice: command, handler, validator,
        ports/                    result + use-case-local port (IDGenerator)
      cancelorder/                vertical slice: command, handler, validator,
        ports/                    result + use-case-local port (Notifier)
  adapters/
    inbound/cli/order/            inbound adapter: CLI drives the use cases
    outbound/memory/order/        outbound adapter: implements OrderRepository
    outbound/id/                  outbound adapter: implements IDGenerator
    outbound/notify/              outbound adapter: implements Notifier
```

The three port scopes from the architecture are all present:

| Scope | Port | Location |
| --- | --- | --- |
| Use-case-local | `IDGenerator` | `application/order/createorder/ports` |
| Use-case-local | `Notifier` | `application/order/cancelorder/ports` |
| Business-area shared | `OrderRepository` | `application/order/ports` |

## Dependency direction

Dependencies always point inward. Adapters depend on the application core and on
ports; the core never imports an adapter. `internal/` enforces the module
boundary, and the composition root in `cmd/orderctl` is the only place where the
core and infrastructure are wired together.

```text
CLI (inbound)  ->  use-case slice  ->  domain
                          |
                          v
                        ports  <-  adapters (outbound)
```
