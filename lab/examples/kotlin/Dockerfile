# syntax=docker/dockerfile:1

ARG GRADLE_VERSION=9.7.1
ARG JDK_VERSION=25

# ---- base: Gradle on JDK 25 + warmed dependency cache ---------------------
# Same shape as the Go example's `COPY go.mod` + `go mod download`: the build files
# alone — the root plus one per module — are enough to resolve every dependency
# into this layer, so the stages below can run --offline. The warm-up happens in a
# throwaway directory, so no scratch source ever ends up in the image that
# compiles the real tree.
FROM gradle:${GRADLE_VERSION}-jdk${JDK_VERSION} AS base
WORKDIR /warm
COPY settings.gradle.kts build.gradle.kts ./
COPY hexagon/domain/build.gradle.kts hexagon/domain/
COPY hexagon/application/build.gradle.kts hexagon/application/
COPY adapters/driving/cli/order/build.gradle.kts adapters/driving/cli/order/
COPY adapters/driven/memory/order/build.gradle.kts adapters/driven/memory/order/
COPY adapters/driven/id/build.gradle.kts adapters/driven/id/
COPY adapters/driven/notify/build.gradle.kts adapters/driven/notify/
COPY composition/build.gradle.kts composition/
# The warm-up needs a source *and* a test source in every module: with no test
# file Gradle marks compileTestKotlin NO-SOURCE and never resolves the test
# classpath, so kotlin-test and the JUnit platform would be missing in the offline
# stages below. Each module gets its own warm package, so the throwaway classes
# cannot collide across the module classpath.
RUN set -eux; \
	for module in \
		hexagon/domain hexagon/application \
		adapters/driving/cli/order adapters/driven/memory/order adapters/driven/id adapters/driven/notify \
		composition; do \
		pkg=$(echo "$module" | tr '/' '.'); \
		mkdir -p "$module/src/main/kotlin/warm/$pkg" "$module/src/test/kotlin/warm/$pkg"; \
		printf 'package warm.%s\n\nclass Warm\n' "$pkg" > "$module/src/main/kotlin/warm/$pkg/Warm.kt"; \
		printf 'package warm.%s\n\nimport kotlin.test.Test\nimport kotlin.test.assertTrue\n\nclass WarmTest {\n    @Test\n    fun warm() = assertTrue(true)\n}\n' "$pkg" > "$module/src/test/kotlin/warm/$pkg/WarmTest.kt"; \
	done; \
	gradle --no-daemon --console=plain test :composition:installDist

# ---- dev: full source; runs the tests offline (no bind mounts) ------------
FROM base AS dev
WORKDIR /src
COPY . .

# ---- build: assemble the distribution (start script + jars) --------------
FROM dev AS build
RUN gradle --no-daemon --console=plain --offline :composition:installDist

# ---- export: hand the distribution out via `cat` (no bind mount) ---------
# `docker run --rm <export-image> > orderctl.tar` streams the artifact to the
# host over stdout — nothing is mounted into the container.
FROM build AS export
RUN tar -C composition/build/install -cf /orderctl.tar orderctl
ENTRYPOINT ["cat", "/orderctl.tar"]

# ---- run: JRE only, no build tool ----------------------------------------
FROM eclipse-temurin:${JDK_VERSION}-jre-noble AS run
COPY --from=build /src/composition/build/install/orderctl /opt/orderctl
ENTRYPOINT ["/opt/orderctl/bin/orderctl"]
CMD ["demo"]
