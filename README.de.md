# HexSlice Architecture

[English](README.md) | **Deutsch**

Dieses Repository beschreibt und dokumentiert die **HexSlice Architecture** als
Architekturkonzept, unabhängig von Programmiersprache und Framework. Das Konzept
selbst ist sprachunabhängig; lauffähige, sprachspezifische Beispielprojekte
liegen unter [`lab/`](lab/).

## Architektur

HexSlice Architecture kombiniert zwei Architekturansätze:

* **Hexagonale Architektur** definiert Systemgrenzen, Ports, Adapter und Abhängigkeitsrichtungen.
* **Vertical Slice Architecture** organisiert den Application Core nach fachlichen Use Cases.

Kurz gesagt:

> HexSlice Architecture bedeutet: vertikale Use-Case-Slices innerhalb eines hexagonalen Application Cores, mit Adaptern außen und Ports im Besitz der Application.

## Grundprinzipien

* Die Domain enthält fachliche Regeln und bleibt technologieunabhängig.
* Die Application-Schicht enthält die Use Cases.
* Use Cases werden als Vertical Slices organisiert.
* Ports werden vom Application Core definiert.
* Ports leben so lokal wie möglich und so gemeinsam wie nötig.
* Ports tragen eine Richtung: Inbound Ports bietet der Core an, Outbound Ports benötigt er.
* Driving Adapter rufen Use Cases über deren Inbound Ports auf.
* Driven Adapter implementieren Outbound Ports.
* Der Core hängt nicht von technischer Infrastruktur ab.

## Dokumentation

Die Architektur ist hier dokumentiert:

[HexSlice Architecture](docs/architecture/hexslice-architecture.de.md)

## Beispiele

Lauffähige, sprachspezifische Beispiele liegen unter [`lab/`](lab/README.md):

* [Go](lab/examples/go/) — Business-Area `order` mit den Slices *create order*
  und *cancel order*, einem CLI-Driving-Adapter sowie `make`- und Docker-Tooling.
* [Kotlin](lab/examples/kotlin/) — dasselbe Beispiel in Kotlin (Gradle + JDK 25,
  nur Docker), via a-check über die JVM-Import-Auflösung gegatet.
