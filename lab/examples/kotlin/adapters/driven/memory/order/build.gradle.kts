// The driven in-memory adapter: implements the order-area OrderRepository port.
//
// A real persistence adapter is where a database driver would be declared — that is
// the point of giving each adapter its own build file: the library lands in the
// module that needs it, and in no other. This one needs none: a ConcurrentHashMap
// and the JDK are enough.
plugins {
    id("org.jetbrains.kotlin.jvm")
}

dependencies {
    // Brings the domain along (`api` in the application module): the adapter maps
    // to and from domain objects.
    implementation(project(":hexagon:application"))

    // Every module declares how it is tested. It sits in testImplementation, so
    // no production classpath gains anything from it.
    testImplementation(kotlin("test"))
}
