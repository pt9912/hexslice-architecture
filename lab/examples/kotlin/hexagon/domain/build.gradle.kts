// The domain: business rules, technology-independent.
//
// There is deliberately no production dependency here — not on another module,
// not on a library. Kotlin's stdlib is the language runtime and the only thing
// this module compiles against. The day someone adds one, the architecture is
// broken, and this file is where it becomes a compile error instead of a
// review remark.
plugins {
    id("org.jetbrains.kotlin.jvm")
}

dependencies {
    // Tests are not part of the production dependency graph.
    testImplementation(kotlin("test"))
}
