# Local load harness

This harness is test-only. It binds the stub server to `127.0.0.1`, uses deterministic in-process responses, and never calls PostgreSQL, an LLM provider, credentials, or a public endpoint.

## Quick path (local k6)

From the repository root:

```bash
go run ./tests/carga/servidor-local > /tmp/vivi-carga.log 2>&1 &
SERVER_PID=$!
until curl --fail --silent http://127.0.0.1:8080/salud >/dev/null; do sleep 0.1; done
BASE_URL=http://127.0.0.1:8080 k6 run tests/carga/endpoints.js
BASE_URL=http://127.0.0.1:8080 k6 run tests/carga/conversations.js
kill "$SERVER_PID"; wait "$SERVER_PID" 2>/dev/null || true
```

The server can also be stopped with `kill <SERVER_PID>` or Ctrl-C in its foreground terminal. The readiness probe and both scripts use loopback only.

## Docker alternative (optional)

Docker is not required. On Linux, keep the server running and use the pinned image:

```bash
docker run --rm --network host -e BASE_URL=http://127.0.0.1:8080 -v "$PWD/tests/carga:/scripts:ro" grafana/k6:0.52.0 run /scripts/endpoints.js
docker run --rm --network host -e BASE_URL=http://127.0.0.1:8080 -v "$PWD/tests/carga:/scripts:ro" grafana/k6:0.52.0 run /scripts/conversations.js
```

`BASE_URL` must be an http URL with no credentials and hostname exactly `127.0.0.1`, `localhost`, or `::1`. Both scripts reject other targets before making a request. The endpoint scenario runs 100 requests/second for 60 seconds over only `GET /salud` and `GET /api/leads`; its deterministic pool of 300 virtual client identities keeps every API identity below the production 30 req/min limit, so every response must be 200. The conversation scenario starts 20 one-iteration VUs against the local LLM stub.

## Results

Measurements below are local-only; do not substitute production or public-endpoint measurements.

| Scenario | Requests / throughput | p95 | Machine | Result |
|---|---:|---:|---|---|
| endpoints.js | 6,000 / 99.998 req/s | 0.888 ms | Linux WSL2, local Go stub, k6 v0.52.0 | PASS — 0% failures |
| conversations.js (20 VUs) | 20 / 1,526 req/s | 13 ms own overhead | Linux WSL2, local Go stub, k6 v0.52.0 | PASS — 0% failures |
