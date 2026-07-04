# HexSlice Architecture

**English** | [Deutsch](README.de.md)

This repository describes and documents the **HexSlice Architecture** as an
architectural concept. It contains no source code and is independent of any
programming language or framework.

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
* Inbound adapters call use cases.
* Outbound adapters implement ports.
* The core does not depend on technical infrastructure.

## Documentation

The architecture is documented here:

[HexSlice Architecture](docs/architecture/hexslice-architecture.md)
