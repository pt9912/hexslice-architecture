// Shared build configuration. What each module may depend on is declared in that
// module's own build.gradle.kts — those four files are where this example's
// architecture is enforced, so they read top to bottom: domain, application,
// adapters, composition.
import org.jetbrains.kotlin.gradle.dsl.JvmTarget

plugins {
    kotlin("jvm") version "2.4.20" apply false
}

subprojects {
    repositories {
        mavenCentral()
    }

    plugins.withId("org.jetbrains.kotlin.jvm") {
        // Unique archive names, derived from the module path. Two modules end in
        // `order` — the memory adapter and the CLI — so their jars would both be
        // `order.jar`, and `installDist` refuses two identical entries in lib/.
        val moduleName = path.removePrefix(":").replace(':', '-')
        extensions.configure<org.gradle.api.plugins.BasePluginExtension>("base") {
            archivesName.set(moduleName)
        }

        extensions.configure<org.jetbrains.kotlin.gradle.dsl.KotlinJvmProjectExtension>("kotlin") {
            compilerOptions {
                // No toolchain declaration on purpose: the build image *is* JDK 25,
                // and a declared toolchain would make Gradle look for (and try to
                // download) a JDK instead of using the one it runs on — which
                // breaks the offline build.
                jvmTarget.set(JvmTarget.JVM_25)
            }
        }
    }

    tasks.withType<org.gradle.api.tasks.testing.Test>().configureEach {
        useJUnitPlatform()
    }
}
