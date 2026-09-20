// The driven id adapter: implements the create-order slice's IdGenerator port.
//
// Its own build file for the same reason as its siblings: if this generator ever
// came from a library — a ULID or UUID implementation — that dependency would be
// declared here and nowhere else. Today an AtomicLong is enough.
plugins {
    id("org.jetbrains.kotlin.jvm")
}

dependencies {
    implementation(project(":hexagon:application"))

    // Every module declares how it is tested. It sits in testImplementation, so
    // no production classpath gains anything from it.
    testImplementation(kotlin("test"))
}
