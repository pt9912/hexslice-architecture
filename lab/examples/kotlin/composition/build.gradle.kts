// The composition root: the only module that knows both sides.
//
// It wires the driven adapters to the slices, hands the handlers to the driving
// adapter, and owns the executable — hence the `application` plugin and the
// mainClass. Nothing depends on this module; it is the top of the graph.
plugins {
    id("org.jetbrains.kotlin.jvm")
    application
}

dependencies {
    implementation(project(":hexagon:application"))
    implementation(project(":adapters:driving:cli:order"))
    implementation(project(":adapters:driven:memory:order"))
    implementation(project(":adapters:driven:id"))
    implementation(project(":adapters:driven:notify"))

    // Every module declares how it is tested. It sits in testImplementation, so
    // no production classpath gains anything from it.
    testImplementation(kotlin("test"))
}

application {
    mainClass = "hexslice.composition.OrderCtlKt"

    // The distribution is named after the executable, not after the module: the
    // `application` plugin derives that name from the *project* name, which is
    // `composition` here — the Dockerfile, the Makefile and the README all speak
    // about `orderctl`.
    applicationName = "orderctl"
}
