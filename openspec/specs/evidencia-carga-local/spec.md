# Evidencia Carga Local Specification

## Purpose

Provide reproducible, local-only NFR evidence and concise repository onboarding documentation without provider or deployment activity.

## Requirements

### Requirement: Loopback-only no-provider k6 scenario

`tests/carga/endpoints.js` MUST reject a `BASE_URL` whose host is not loopback before requests begin. It MUST run constant-arrival-rate 100 requests/second for 60 seconds only against `GET /api/leads` and `GET /salud`, MUST invoke no LLM/provider route, and MUST declare p95 below 300 ms and failure rate below 1% thresholds.

#### Scenario: Non-loopback target is refused
- GIVEN `BASE_URL` names a non-loopback host
- WHEN the script starts
- THEN it fails before making any HTTP request

#### Scenario: Local-only execution
- GIVEN a local loopback server and k6 installed locally or through Docker
- WHEN the documented command runs
- THEN requests target only loopback no-LLM routes and report threshold results

### Requirement: Local evidence and reviewable docs

`tests/carga/README.md` MUST document the local command and record measured p95, throughput, and machine. Root `README.md` MUST provide a five-command-or-fewer local quickstart, RF-to-package map, limiter policy, and single-process limitation. `docs/arquitectura.md` MUST contain Mermaid C4 L1–L3 and sequence diagrams sourced from the existing architecture context.

#### Scenario: Documentation review
- GIVEN a repository reader
- WHEN they follow the quickstart and load-evidence guide
- THEN they can identify the local-only boundary, commands, RF owners, limiter limitations, and diagrams
