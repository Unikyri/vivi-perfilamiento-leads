# Exploration: Load Evidence, Client-IP Rate Limit, and Repository Docs (Issue #31)

## Current state
There is no HTTP client rate-limit middleware. `LIMITE_TASA` already maps to a sanitized 429 error envelope; the existing limiter protects LLM providers only. The server assigns a bare mux to `http.Server.Handler`. README lacks a local quickstart and k6 evidence. No production endpoint was contacted.

## Decision
Use a route-aware outer HTTP middleware. Limit `/api` only; exempt `/salud` and static assets. Apply a bounded 30-request/one-minute fixed-window policy per canonical client identity. Default identity is `RemoteAddr`; forwarded headers are accepted only from configured trusted peers, so clients cannot spoof an IP. Keep implementation process-local and document that limitation. The k6 script must reject non-loopback targets and never invoke an LLM/provider.

## Evidence
Issue #31 body, current HTTP composition/error tests, and persisted Engram explorations #630 and state #629.
