rootProject.name = "orderctl"

// One Gradle module per layer — and one per adapter, because each adapter is where
// its *own* external libraries would land: a Postgres adapter pulls its driver, an
// HTTP notifier its client, and neither leaks into the core or into a sibling
// adapter.
//
// This module graph is the dependency rule made compilable: the hexagon modules
// cannot see the adapters, because no such dependency exists on their classpath.
// The compiler objects long before the architecture gate does.
include(
    ":hexagon:domain",
    ":hexagon:application",
    ":adapters:driving:cli:order",
    ":adapters:driven:memory:order",
    ":adapters:driven:id",
    ":adapters:driven:notify",
    ":composition",
)
