# Lab — Examples

Runnable, language-specific example projects that put the
[HexSlice Architecture](../docs/architecture/hexslice-architecture.md) into
practice. Each example is self-contained and independent of the others.

## Examples

| Language | Use case(s) | Highlights |
| --- | --- | --- |
| [Go](examples/go/) | Create Order, Cancel Order | Two vertical slices, both port directions across all three port scopes, CLI driving adapter, in-memory/id/notify driven adapters, `make` + Docker, [a-check](https://github.com/pt9912/a-check) architecture gate |
| [Kotlin](examples/kotlin/) | Create Order, Cancel Order | The same architecture in a second language: two vertical slices, both port directions across all three port scopes, CLI driving adapter, in-memory/id/notify driven adapters, Gradle multi-module (one module per layer and per adapter) on JDK 25 in Docker, a-check gate with the JVM import resolution |

Each example maps the architecture's folder layout onto the idioms of its
language and ships with its own `README.md` explaining how to build and run it.
