# Exploration: issue-23-http-api

## Current State
`cmd/servidor/main.go` currently wires configuration, PostgreSQL migrations, LLM provider, and only `GET /salud`. Existing completed use cases from Issues #18–#22 are callable, including profiling, message processing, qualification, ficha generation, nutrition milestones, the in-memory event bus, and deterministic coordinator. No ADK is in scope.

The HTTP adapter lacks controllers, JSON/error presentation, API route registration, and composition wiring. Frontend models expect the Contract §3 API: lead creation, messages, conversation, queue, ficha, buyer persona, demo time, and reset.

## Risks and Decisions Needed
- Contract v1.1 and architecture documents must govern exact DTOs, errors, privacy, queue, simulated-time, and reset behavior.
- Queue, buyer-persona, production clock, and reset may require missing application capabilities; handlers must not invent business logic.
- HTTP can expose sensitive lead data; logs and errors must remain privacy-safe.
- Scope is complex and needs bounded slices of at most 400 authored runtime/test lines.

## Provisional Delivery
1. HTTP foundation: error/JSON helpers, router seam, creation, messages, conversation.
2. Queue and ficha presentation.
3. Buyer persona and simulated time after contract validation.
4. Demo reset only if a safe Contract-defined boundary exists.

## Explicit Exclusion
ADK is explicitly excluded by the user. Integration uses the existing synchronous plain-Go use cases, event bus, and coordinator only.
