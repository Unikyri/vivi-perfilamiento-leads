# Limite Tasa HTTP Specification

## Purpose

Protect public API routes with a bounded, process-local per-client request limit without trusting spoofable identity headers.

## Requirements

### Requirement: Route-aware fixed-window limit

Every `/api` request MUST consume one of 30 requests in a 60-second fixed window for its client identity. `GET /salud` and static `GET` requests MUST be exempt. The limiter MUST remain process-local and MUST NOT add `Retry-After`, a distributed store, provider limiting, or an API-contract change.

#### Scenario: Thirty-first API request
- GIVEN one client sends 30 `/api` requests in its current window
- WHEN it sends request 31
- THEN the response is HTTP 429 with `error.codigo=LIMITE_TASA`

#### Scenario: Exempt health route
- GIVEN one client sends 50 sequential `GET /salud` requests
- WHEN they are handled
- THEN all succeed without consuming API quota

### Requirement: Safe, bounded client identity

The default identity MUST be the canonical host from `RemoteAddr`. `X-Forwarded-For` and `Forwarded` MUST be ignored unless the direct peer belongs to an explicitly configured trusted-proxy CIDR allowlist, which defaults empty; malformed configured CIDRs MUST fail configuration loading. Per-identity state MUST expire after its window and MUST remain bounded with eviction.

#### Scenario: Untrusted header is ignored
- GIVEN an untrusted direct peer with a spoofed forwarding header
- WHEN it requests an API route
- THEN quota is attributed to its canonical remote address

#### Scenario: Distinct client and expiration
- GIVEN one exhausted identity and a second identity, then an expired first window
- WHEN each sends an API request
- THEN the second succeeds and the first receives a fresh quota after expiration

### Requirement: Sanitized rejection

A rate-limit rejection MUST use the existing error envelope and MUST expose no IP address, trusted-proxy configuration, forwarding header, or limiter internals.

#### Scenario: Envelope
- GIVEN an exhausted API identity
- WHEN the middleware rejects its request
- THEN status is 429 and `codigo` is `LIMITE_TASA`
- AND details are empty and response text contains no identity/configuration value
