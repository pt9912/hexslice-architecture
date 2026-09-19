# HexSlice Architecture

**English** | [Deutsch](README.de.md)

This repository describes and documents the **HexSlice Architecture** as an
architectural concept, independent of any programming language or framework. The
concept itself is language-agnostic; runnable, language-specific example
projects live under [`lab/`](lab/).

## Architecture

HexSlice Architecture combines two architectural ideas:

* **Hexagonal Architecture** defines the system boundaries, ports, adapters, and dependency direction.
* **Vertical Slice Architecture** organizes the application core by use cases.

In short:

> HexSlice Architecture means vertical use-case slices inside a hexagonal application core, with adapters outside and ports owned by the application.

## Core Principles

* The domain contains business rules and remains technology-independent.
* The application layer contains use cases.
* Use cases are organized as vertical slices.
* Ports are defined by the application core.
* Ports live as locally as possible and as shared as necessary.
* Ports carry a direction: inbound ports are offered by the core, outbound ports are needed by it.
* Driving adapters call use cases through their inbound ports.
* Driven adapters implement outbound ports.
* The core does not depend on technical infrastructure.

## Documentation

The architecture is documented here:

[HexSlice Architecture](docs/architecture/hexslice-architecture.md)

## Examples

Runnable, language-specific examples live under [`lab/`](lab/README.md):

* [Go](lab/examples/go/) — `order` business area with *create order* and
  *cancel order* slices, a CLI driving adapter, and `make` + Docker tooling.
