# HexSlice Architecture

[English](README.md) | **Deutsch**

Dieses Repository beschreibt und dokumentiert die **HexSlice Architecture** als
Architekturkonzept. Es enthält keinen Quellcode und ist unabhängig von
Programmiersprache und Framework.

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
* Inbound Adapter rufen Use Cases auf.
* Outbound Adapter implementieren Ports.
* Der Core hängt nicht von technischer Infrastruktur ab.

## Dokumentation

Die Architektur ist hier dokumentiert:

[HexSlice Architecture](docs/architecture/hexslice-architecture.de.md)
