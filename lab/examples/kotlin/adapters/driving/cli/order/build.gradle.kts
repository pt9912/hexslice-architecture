// The driving adapter: the CLI.
//
// It speaks the slices' inbound ports and knows nothing else about the core — not
// the domain, not a slice. Its own build file is where a library would go if this
// adapter needed one (an argument parser, a terminal library); today it uses the
// Kotlin stdlib.
plugins {
    id("org.jetbrains.kotlin.jvm")
}

dependencies {
    implementation(project(":hexagon:application"))

    // Every module declares how it is tested. It sits in testImplementation, so
    // no production classpath gains anything from it.
    testImplementation(kotlin("test"))
}
