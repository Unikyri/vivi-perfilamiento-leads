# Tasks: Issue 31 NFR Hardening

## Review Workload Forecast

| Field | Value |
|---|---|
| Estimated changed lines | 620–790 total; 300–360 / 180–250 / 140–180 per slice |
| 400-line budget risk | High (total); Low per slice |
| Chained PRs recommended | Yes |
| Suggested split | PR1 limiter; PR2 local harness; PR3 docs/evidence |
| Delivery strategy | ask-on-risk |
| Chain strategy | pending |

Decision needed before apply: Yes
Chained PRs recommended: Yes
Chain strategy: pending
400-line budget risk: High

### Suggested Work Units

| Unit | Goal | Likely PR | Focused test command | Runtime harness | Rollback boundary |
|---|---|---|---|---|---|
| 1 | Safe API limiter | PR1 | `go test ./internal/adapters/http ./internal/infrastructure/config -count=1` | `go test -race ./internal/adapters/http -count=1` | limiter/config/main wrap |
| 2 | Loopback-only load evidence | PR2 | k6 script validation | local stub server + both k6 scenarios | `tests/carga/`, test-only server |
| 3 | Reader docs/evidence | PR3 | Mermaid/doc review | record local results only | README/docs only |

## Phase 1: Limiter and Safe Identity

- [x] 1.1 Add failing table/race tests in `internal/adapters/http/ratelimit_test.go`: `/api` 30/31, exact 60s reset, identity isolation, `/salud`/static bypass, downstream not called, bounded expiry/capacity, mapped IPv4.
- [x] 1.2 Add failing trusted-proxy cases: untrusted spoof ignored; loopback-trusted `Forwarded`/`X-Forwarded-For` chain selection; malformed/ambiguous header falls back to peer; rejection envelope leaks no identity/configuration.
- [x] 1.3 Create `internal/adapters/http/ratelimit.go` with mutex-protected fixed-window state, canonical `RemoteAddr`, trusted-prefix resolver, and 429 sentinel; map it through existing `LIMITE_TASA` in `errores.go`.
- [x] 1.4 Extend `internal/infrastructure/config/config.go` and `_test.go` for empty/default, IPv4/IPv6 allowlists, and malformed `TRUSTED_PROXY_CIDRS` startup failure; document the empty production-safe default in `.env.example`.
- [x] 1.5 Wrap the fully registered mux once in `cmd/servidor/main.go`; preserve all routes, no `Retry-After`, provider limit, Contract, database, or deployment change.

## Phase 2: Test-Only Local Load Harness

- [x] 2.1 Create a test-only local server command under `tests/carga/` that uses an in-process LLM stub (no keys/provider/network), exposes no-LLM API/health and seeded conversation flow, and trusts only `127.0.0.1/32` and `::1/128` for this harness.
- [x] 2.2 Create `tests/carga/endpoints.js`: reject credentials/non-loopback `BASE_URL` before requests; run 100 req/s for 60s over only `GET /salud` and `GET /api/leads`, unique virtual identities in `X-Forwarded-For`, p95 <300ms and unexpected failures <1% (429 expected).
- [x] 2.3 Create `tests/carga/conversations.js`: 20 local LLM-stub conversations with unique forwarded identities and a p95-overhead threshold; never call a provider or public endpoint.
- [x] 2.4 Add `tests/carga/README.md` commands that start/await/stop the loopback server, optionally use pinned Docker k6, and record machine, throughput, p95, and pass/fail results.

## Phase 3: Documentation and Evidence

- [x] 3.1 Update `README.md` with ≤5-command quickstart, RF-to-package map from repository sources, 30/min per-identity policy, process-local limitation, and local-harness-only proxy trust.
- [x] 3.2 Create `docs/arquitectura.md` by transcribing—not inferring—Wiki doc 11 C4 L1/L2/L3 and sequences 4.1–4.3 from `/tmp/vivi-wiki/11‐Arquitectura-de-Software-—-Vivi-·-Fase-0.md`.

## Phase 4: Validation

- [x] 4.1 Run focused/race tests, then `go test ./... -count=1` and `go vet ./...`; fix only Issue 31 regressions.
- [x] 4.2 Run both documented k6 scenarios against the started loopback/stub server, persist measured local evidence, verify non-loopback refusal, and stop the harness; do not add k6 or provider calls to CI.
