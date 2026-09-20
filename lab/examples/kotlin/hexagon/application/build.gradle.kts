// The application core: the use-case slices and their ports.
//
// It may depend on the domain — inward, and on nothing else. The ports are
// defined here, so an adapter that implements one depends on *this* module; the
// dependency is never declared in the other direction, and it cannot be, because
// this module does not know that adapters exist. No library either: the core
// stays free of infrastructure concerns.
plugins {
    id("org.jetbrains.kotlin.jvm")
}

dependencies {
    // `api`, not `implementation`: the port signatures are written in domain
    // types, so whoever implements a port needs the domain on its classpath too.
    api(project(":hexagon:domain"))

    testImplementation(kotlin("test"))
}
