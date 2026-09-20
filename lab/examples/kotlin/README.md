# HexSlice Example — Kotlin

A small Kotlin/JVM project that demonstrates the
[HexSlice Architecture](../../../docs/architecture/hexslice-architecture.md):
vertical use-case slices inside a hexagonal application core, with adapters
outside and ports owned by the application.

It is the twin of the [Go example](../go/): the same business area, the same
slices, the same ports and the same CLI surface — gated by the same
[a-check](https://github.com/pt9912/a-check) rules. What differs is the language
*and* the build: this example is a Gradle multi-module project, one module per
layer and one per adapter, so the coarse dependency direction is enforced by the
compiler. The section [What differs from the Go example](#what-differs-from-the-go-example)
lists it all.

The example models a single business area — **order** — with two vertical
slices: **create order** and **cancel order**. An in-memory repository, a
sequential id generator and a stdout notifier act as driven adapters, and a CLI
acts as the driving adapter.

## Requirements

* **Docker** only — no local JVM, no local Gradle.

Every target runs inside Docker. Gradle (**9.7.1**) on Eclipse Temurin JDK
(**25**) is the build image; the runtime image is a JRE 25 with the assembled
distribution. Kotlin (**2.4.20**) is pinned in the root `build.gradle.kts`. The
source is baked in via `COPY`, nothing is bind-mounted, and the assembled
distribution is streamed out of the container over stdout by the `export` stage's
entrypoint (see `make build`).

There is no Gradle wrapper: the Gradle version *is* the image tag, which keeps the
pin in one place — the `Makefile`'s `GRADLE` variable — instead of in a wrapper
properties file that the build would then download over the network.

The dependency cache is warmed **in the image**, the way the Go example warms its
module cache: the base stage copies the build files alone — the root plus one per
module — and resolves everything once, so the build and test stages can run
`--offline`. That makes the cache part of the image rather than of a developer's
home directory.

## Quick start

```bash
make run        # build the runtime image and run the create -> cancel demo
make smoke      # run the demo in the runtime image and diff it against the expected lines
make test       # compile + run the tests offline, in Docker
make check      # the Kotlin counterpart of the Go example's check (= test)
make build      # assemble the distribution and extract it to ./bin
make image      # build the runtime image (JRE + distribution)

# architecture gate
make a-check        # fails on any hexagonal-layering violation
make a-check-graph  # prints the declared layers as a Mermaid flowchart
```

Run `make help` to list all targets. The Go example also has `vet`, `fmt`, `lint`
and `tidy` targets; this one has no counterpart because it carries no linter —
the gate here is the compiler, the tests and a-check.

## Where the architecture is enforced

Two mechanisms, deliberately layered:

**The build enforces the coarse direction.** One module per layer, and the
dependency is declared only in the direction that is allowed:

| Module | may depend on | the file that says so |
| --- | --- | --- |
| `:hexagon:domain` | nothing — not a module, not a library | `hexagon/domain/build.gradle.kts` |
| `:hexagon:application` | `:hexagon:domain` (as `api`) | `hexagon/application/build.gradle.kts` |
| `:adapters:*` (one per adapter) | `:hexagon:application` | each adapter's own `build.gradle.kts` |
| `:composition` | the application and all adapters | `composition/build.gradle.kts` |

The hexagon modules have no adapters on their classpath, so "the core must not
know the adapters" is a compile error rather than a review remark. And because
each adapter has its own module, an external library lands exactly where it
belongs — a Postgres adapter pulls its driver, an HTTP notifier its client, and
neither leaks into the core or into a sibling adapter. Today none of them needs
one; the build files say where it would go.

**a-check checks what a compiler cannot see.** Inside the application module, the
per-slice and per-port boundaries are a matter of paths, not of modules —
`lateral-slice`, `port-locality` and `port-direction-mismatch` are the rules that
watch them.

## Architecture gate (a-check)

The example is validated by [a-check](https://github.com/pt9912/a-check)
(**v0.20.0**). The layer and edge mapping lives in [`.a-check.yml`](.a-check.yml);
the make targets come from the generated [`a-check.mk`](a-check.mk), and the image
digest is pinned in the `Makefile`.

Three rules are **active**, each verified by injecting a violation into a copy of
the tree:

| Rule | Injection | Finding |
| --- | --- | --- |
| `lateral-slice` | a slice imports a foreign slice | `lateral-slice: 1`, exit 1 |
| `port-locality` | a slice imports a foreign slice-local port | `port-locality: 1`, exit 1 |
| `port-direction-mismatch` | a *driven* adapter imports an inbound port | `port-direction-mismatch: 1`, exit 1 |

The direction finding survives a declared `allow` edge — that is what
"categorical" means for this rule. On the clean tree the gate reports
`0 Befund(e)`.

### The Kotlin-specific half: import resolution

Kotlin imports are fully qualified dotted names, so a-check needs to be told how
they map onto paths — that is the `resolution` block, and getting it wrong is
**silent**. With the sources under `<module>/src/main/kotlin`, that is the
multi-module recipe: one `package_base`, one `roots` entry per module.

```yaml
resolution:
  kotlin:
    mode: fixed-root
    package_base: "hexslice"
    roots: ["hexagon/domain/src/main/kotlin/hexslice", …]
```

With no `package_base`, a-check keeps the FQN dotted
(`hexslice.hexagon.domain.order.Order`), it matches no layer glob, the import is
treated as repo-external and **no rule judges it** — the gate still says
`0 Befund(e)`, and it would say `0 Befund(e)` for a tree full of violations. The
tell is on **stderr**: a-check's boundary advisory reported *"0 von 18
Import-Symbolen lösen auf eine Schicht auf"*. Read the advisory, not just the
summary — this example was measured green-but-blind before the `package_base`
line existed.

Unlike the Go toolchain targets — which bake the source in via `COPY` and use no
mounts — a-check reads the tree through a **read-only** bind mount (`:ro`) and
runs **network-isolated** (`--network none`); that is a-check's own contract and
it never writes to the repository.

## CLI usage

Build the runtime image once (`make image`), then invoke subcommands against it:

```bash
# self-contained lifecycle demo (default when no subcommand is given)
docker run --rm hexslice-order-kotlin demo

# single use-case invocations
docker run --rm hexslice-order-kotlin create -customer cust-1 \
  -line BOOK-1:2:1999:EUR -line PEN-7:5:150:EUR
docker run --rm hexslice-order-kotlin cancel -id ORD-000001 -reason "changed mind"
```

`make build` extracts the assembled distribution to `./bin/orderctl` (a start
script plus jars). Running it needs a JRE 25 on the host — `make run` is the
Docker path, and the one the Go example's static binary does not need.

> The repository is in-memory, so state does not survive across separate process
> invocations. Use the `demo` subcommand to see a full create → cancel lifecycle
> in one run.

## Tests

`make test` compiles every module and runs four suites: the domain rules, both
slice handlers (with hand-written port stubs, per module), and one **end-to-end**
test in `:composition` that drives the composition root's own wiring — the real
adapters, the real slices, the real CLI — through the `demo` scenario, a `create`
over the flag surface, a domain-error path and a malformed flag. Only the process
boundary (`main`'s exit code) is out of scope.

The end-to-end test lives in `:composition`, and that is forced rather than chosen:
`:adapters:driving:cli:order` does not have the driven adapters on its classpath, so
a test over there would not even compile. The composition root is the only module
that sees both sides — which is also why `main` is four lines: the wiring it drives
lives in `application(...)`, and the test drives the same function instead of a
lookalike that could drift from it.

`make smoke` adds the check no suite can make: it starts the **runtime image**, runs
the demo and diffs the output against the expected lines (`SMOKE_EXPECTED` in the
Makefile). That is the packaging level — and the level this example got wrong once:
the `application` plugin derived the distribution name from the module, so the image
held `composition` instead of `orderctl`. Only the built image reveals that, and
`make smoke` now fails on it. It is deliberately *not* part of `make check`, which
stays a source-level gate: run `smoke` whenever you touched the Dockerfile, the
build files or the entrypoint.

## How the folders map to HexSlice

```text
settings.gradle.kts · build.gradle.kts   the module graph (see above)
hexagon/domain/                          module: business core
  src/main/kotlin/hexslice/hexagon/domain/order/    Order aggregate, value
                                                    objects, rules — depends on
                                                    nothing
hexagon/application/                     module: the application core
  src/main/kotlin/hexslice/hexagon/application/order/
    ports/outbound/                      business-area shared port (OrderRepository)
    createorder/                         vertical slice: handler, validator + ports
      ports/inbound/                     use-case interface + request/result (CreateOrder)
      ports/outbound/                    use-case-local port (IdGenerator)
    cancelorder/                         vertical slice: handler, validator + ports
      ports/inbound/                     use-case interface + request/result (CancelOrder)
      ports/outbound/                    use-case-local port (Notifier)
adapters/driving/cli/order/              module: driving adapter — the CLI
adapters/driven/memory/order/            module: driven adapter — OrderRepository
adapters/driven/id/                      module: driven adapter — IdGenerator
adapters/driven/notify/                  module: driven adapter — Notifier
composition/                             module: composition root (wires everything)
src/test/kotlin/ is per module:  hexagon/domain/… and hexagon/application/…
```

Every one of those directories carries its own `build.gradle.kts` with a
`src/main/kotlin/hexslice/…` convention layout — the module paths above are the
Gradle project paths (`:hexagon:domain`, `:adapters:driven:id`, …).

Both port directions are present, across all three port scopes:

| Direction | Scope | Port | Location (under `hexagon/application/src/main/kotlin/hexslice/hexagon/application/order/`) |
| --- | --- | --- | --- |
| Inbound | Use-case-local | `CreateOrder` | `createorder/ports/inbound` |
| Inbound | Use-case-local | `CancelOrder` | `cancelorder/ports/inbound` |
| Outbound | Use-case-local | `IdGenerator` | `createorder/ports/outbound` |
| Outbound | Use-case-local | `Notifier` | `cancelorder/ports/outbound` |
| Outbound | Business-area shared | `OrderRepository` | `ports/outbound` |

As in the Go example, the inbound port owns its contract: the interface *and* the
request/result types belong to the port package, so a driving adapter can import
the port alone and never the slice. The slice implements the port and imports its
own port for those types — inward, like every other dependency.

## What differs from the Go example

Same architecture, five differences — four from the language, one from the build:

* **The module graph.** The Go example is a single module, so only a-check guards
  its layer boundaries. Here the hexagon modules have no adapters on their
  classpath at all: the coarse rule is a compile error, and a-check adds the
  fine-grained per-slice, per-port and direction rules. That is also why this
  README has a section the Go one does not.
* **Nominal vs. structural typing.** Kotlin's type system makes the adapter *name*
  the port it implements (`class InMemoryOrderRepository : OrderRepository`), so
  the driven adapters import `ports_outbound` and the gate needs the edge
  `driven_adapters -> ports_outbound`. Go satisfies its interfaces structurally
  and needs neither the import nor the edge — that single edge is the whole
  difference in the two `.a-check.yml` files' edge lists.
* **Errors.** The Go twin returns sentinel errors (`(Order, error)`,
  `errors.Is`); here they are a sealed hierarchy — `OrderError.NoLines`,
  `CreateOrderError.NoCustomer` — thrown and asserted with `assertFailsWith`. The
  vocabulary is the same set of names.
* **Identifiers and value objects.** `OrderId`/`CustomerId` are value classes, so
  the type system carries what the Go twin's named string types carry. `Money` and
  `Line` are `data class`es with a private constructor — validated by `of` is the
  only way in — annotated with `@ConsistentCopyVisibility` so the generated
  `copy()` cannot bypass that validation either.
* **Context.** Go threads a `context.Context` through the ports and handlers for
  cancellation and deadlines; the Kotlin twin has no equivalent here, so the
  parameter is simply absent.

## Dependency direction

Dependencies always point inward. Adapters depend on the ports and on the domain;
the core never imports an adapter — and in this example it cannot, because the
module is not on its classpath. The composition root is the only module that knows
both sides, and the place where "the handler satisfies the port" is checked by the
compiler.

```text
CLI (driving)  ->  inbound ports  ->  use-case slice  ->  domain
                                            |
                                            v
                          outbound ports  <-  adapters (driven)

modules:  composition -> adapters/* -> application -> domain
          (nothing depends on composition; domain depends on nothing)
```
