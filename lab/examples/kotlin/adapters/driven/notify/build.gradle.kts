// The driven notifier adapter: implements the cancel-order slice's Notifier port.
//
// This is the adapter that would grow a dependency first — an SMTP client, an HTTP
// client, a message-broker library. It is declared here, in the module that uses
// it, so the hexagon and the sibling adapters stay free of it. Today an Appendable
// is enough.
plugins {
    id("org.jetbrains.kotlin.jvm")
}

dependencies {
    implementation(project(":hexagon:application"))

    // Every module declares how it is tested. It sits in testImplementation, so
    // no production classpath gains anything from it.
    testImplementation(kotlin("test"))
}
