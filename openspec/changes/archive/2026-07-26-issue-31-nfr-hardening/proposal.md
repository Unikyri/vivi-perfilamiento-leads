# Proposal: Local NFR Hardening — k6 Load Evidence, Per-IP Rate Limit, Repo Docs (Issue #31)

## Intent

Three audited NFR gaps have no verification: NFR-E-02 (load evidence), NFR-S-04 (public-API per-IP rate limit), NFR-M-05 (README quickstart + RF map + architecture diagrams). Today the API has no client rate limit, no reproducible load evidence, and no local setup guide, so scale claims are unverified and a curious visitor can hammer the demo API.

## Scope

### In Scope
- `tests/carga/endpoints.js`: k6 constant-arrival-rate 100 req/s for 60s over no-LLM routes (`GET /api/leads`, `GET /salud`), thresholds p95 < 300 ms and failure rate < 1%; loopback-only `BASE_URL` (reject non-loopback), never invokes an LLM/provider.
- `tests/carga/README.md`: run command plus measured p95, throughput, and machine.
- `internal/adapters/http/ratelimit.go` + `_test.go`: route-aware outer middleware, 30 requests / 60 s fixed window per client identity, `/salud` and static GETs exempt, API paths limited.
- Safe identity: canonicalized `RemoteAddr` by default; `X-Forwarded-For`/`Forwarded` honored only when the direct peer matches an explicitly configured trusted-proxy CIDR allowlist (default empty). Bounded/evicted per-identity state with injected clock.
- Rejection returns existing sanitized envelope: HTTP 429, `codigo: LIMITE_TASA`, no IP/config/header leakage.
- `README.md`: ≤5-command local quickstart, RF→package map, limiter policy and single-process limitation.
- `docs/arquitectura.md`: Mermaid C4 L1–L3 and sequence diagrams (from Wiki doc 11).

### Out of Scope
- The 20-concurrent-conversation LLM-stub scenario, CI k6 job, distributed/edge limiter, `Retry-After`.
- Contract v1.1, error catalog, DB schema, provider limiter, deployment/production probing.

## Capabilities

### New Capabilities
- `limite-tasa-http`: per-client HTTP rate limiting, route exemptions, trusted-proxy identity resolution, 429 `LIMITE_TASA` envelope.
- `evidencia-carga-local`: local-only k6 load scenario, thresholds, and documented evidence.

### Modified Capabilities
- None.

## Approach

Wrap the fully registered mux at `http.Server.Handler` in `cmd/servidor/main.go`. Reuse `writeError` via a new package sentinel mapped to `LIMITE_TASA`; do not add a new error code. Trusted-proxy CIDRs come from config with a safe empty default and fail on malformed values. Tests use `httptest` + fake clock: 31st request blocked, per-IP isolation, `/salud` exempt, envelope shape, spoofed header from untrusted peer ignored, eviction.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `internal/adapters/http/ratelimit*.go` | New | Middleware, identity, tests |
| `cmd/servidor/main.go` | Modified | Wrap handler |
| `internal/infrastructure/config/config.go` | Modified | Trusted-proxy CIDRs |
| `tests/carga/*`, `docs/arquitectura.md`, `README.md` | New/Modified | Evidence and docs |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Header spoofing evades limiter | Med | Trust forwarded headers only from configured peers |
| Unbounded identity map | Med | Capacity limit + expiry eviction |
| Limiter breaks anti-sleep ping/demo | Med | Explicit `/salud` and static exemptions, tests |
| k6 unavailable locally | Med | Documented Docker k6 command; manual evidence |
| Process-local only | High | Documented limitation, not a scaling boundary |

## Rollback Plan

Revert the `Handler` wrap line in `main.go` (limiter becomes inert), then revert the change branch commits. No schema, Contract, or data migration to undo; docs and k6 files are additive.

## Dependencies

- #23 (API), #27 (deploy) already merged; local Docker PostgreSQL and Docker/k6 for evidence.

## Success Criteria

- [ ] k6 100 req/s × 60 s meets p95 < 300 ms and <1% failures; result documented with machine
- [ ] 31st API request within a minute returns 429 `LIMITE_TASA`; other IPs unaffected
- [ ] 50 sequential `/salud` requests all succeed
- [ ] Spoofed forwarded header from untrusted peer does not change identity
- [ ] README quickstart (≤5 commands) + RF map, and `docs/arquitectura.md` Mermaid diagrams present
- [ ] `go vet ./...` and `go test ./... -count=1` pass
