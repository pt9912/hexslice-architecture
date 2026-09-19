# HexSlice Example — Go

A small, dependency-free Go project that demonstrates the
[HexSlice Architecture](../../../docs/architecture/hexslice-architecture.md):
vertical use-case slices inside a hexagonal application core, with adapters
outside and ports owned by the application.

The example models a single business area — **order** — with two vertical
slices: **create order** and **cancel order**. An in-memory repository, a
sequential id generator, and a stdout notifier act as driven adapters, and a
CLI acts as the driving adapter.

## Requirements

* **Docker** only — no local Go toolchain is needed.

Every target runs inside Docker. The Go toolchain (pinned to **1.27.1**) and
`golangci-lint` (pinned to **v2.13.2**) live in the images; the source is baked
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

The example is validated by [a-check](https://github.com/pt9912/a-check)
(**v0.19.0**), a containerized gate that verifies the hexagonal layering. The
layer/edge mapping lives in [`.a-check.yml`](.a-check.yml); the make targets come
from the generated [`a-check.mk`](a-check.mk), and the image digest is pinned in
the `Makefile`.

Because the config uses clean per-slice / per-port directory globs, a-check also
enforces the two **vertical-slice** rules on top of the classic hexagonal ones:

- `lateral-slice` — one use-case slice may not import another slice of the same
  application layer (cross-slice contracts must go through a shared port).
- `port-locality` — a use-case-local port may only be used inside its own slice;
  shared contracts must be promoted to the business-area (`order/ports`) level.

`make a-check` returns exit code `0` (no violations) for this example. Both
slice rules are demonstrably **active** here: injecting a cross-slice import and
an import of a foreign slice-local port into a copy of the tree makes a-check
report `lateral-slice` and `port-locality` and exit `1`.

Two boundaries of the gate are worth stating, because both are silent when they
bite:

- The port globs end at the `ports` segment, **not** at the direction folder
  below it. `port-locality` derives a port's scope by stripping the last segment
  of the longest matching glob prefix, so a glob on `.../ports/outbound/**`
  would scope the port to `.../ports` instead of `.../createorder` and the rule
  would go inert — measured: a genuine cross-slice import then reports nothing.
  A port layer therefore cannot carry a `direction` while its slice holds both
  inbound and outbound ports.
- `port-direction-mismatch` pairs `driving` with `inbound` and `driven` with
  `outbound`, and needs both sides declared. Only the adapter layers declare a
  direction here; the rule is opt-in and stays inert until the port side can
  carry one, exactly as the dimension specifies.

Unlike the Go toolchain targets — which bake the source in via `COPY` and use no
mounts — a-check reads the tree through a **read-only** bind mount (`:ro`) and
runs **network-isolated** (`--network none`); this is a-check's own contract and
it never writes to the repository.

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
      ports/outbound/             business-area shared port (OrderRepository)
      createorder/                vertical slice: handler, validator + its ports
        ports/inbound/            use-case interface + request/result
                                  (CreateOrder)
        ports/outbound/           use-case-local port (IDGenerator)
      cancelorder/                vertical slice: handler, validator + its ports
        ports/inbound/            use-case interface + request/result
                                  (CancelOrder)
        ports/outbound/           use-case-local port (Notifier)
  adapters/
    driving/cli/order/            driving adapter: the CLI calls the inbound ports
    driven/memory/order/          driven adapter: implements OrderRepository
    driven/id/                    driven adapter: implements IDGenerator
    driven/notify/                driven adapter: implements Notifier
```

Both port directions are present, across all three port scopes:

| Direction | Scope | Port | Location |
| --- | --- | --- | --- |
| Inbound | Use-case-local | `CreateOrder` | `application/order/createorder/ports/inbound` |
| Inbound | Use-case-local | `CancelOrder` | `application/order/cancelorder/ports/inbound` |
| Outbound | Use-case-local | `IDGenerator` | `application/order/createorder/ports/outbound` |
| Outbound | Use-case-local | `Notifier` | `application/order/cancelorder/ports/outbound` |
| Outbound | Business-area shared | `OrderRepository` | `application/order/ports/outbound` |

The inbound port owns its contract: the interface *and* the request/result types
belong to the port package, which is why a driving adapter can import the port
alone and never the slice. The slice implements the port and imports
its own port for those types — inward, like every other dependency. The
composition root is the one place that knows both sides, so it is also where
"the handler satisfies the port" is checked, by the compiler.

## Dependency direction

Dependencies always point inward. Adapters depend on the application core and on
ports; the core never imports an adapter. `internal/` enforces the module
boundary, and the composition root in `cmd/orderctl` is the only place where the
core and infrastructure are wired together.

```text
CLI (driving)  ->  inbound ports  ->  use-case slice  ->  domain
                                            |
                                            v
                          outbound ports  <-  adapters (driven)
```
