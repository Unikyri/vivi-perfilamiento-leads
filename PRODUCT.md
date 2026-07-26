# Product

<!-- impeccable:product-schema 1 -->

## Platform

web

## Users

- **Asesor de vivienda (primary, dashboard):** Colsubsidio housing advisor who receives a lead already qualified by Vivi and must close a sale. Works from the advisor panel (`web/src/views/cola.ts`, `ficha.ts`, `chat.ts`) to triage a lead queue, open a lead's WhatsApp conversation, and read its generated "ficha" (profile card) before calling.
- **Lead / prospect (WhatsApp channel):** a Colsubsidio affiliate exploring housing options who chats with "Vivi," the WhatsApp assistant, before ever talking to a human advisor.
- **Gerencia view (`web/src/views/gerencia.ts`) is a hackathon demo extra, not a confirmed product audience** — do not treat its "buyer persona vivo" metrics as a real stakeholder need without further confirmation.

## Product Purpose

Vivi turns inbound leads into qualified, ready-to-call prospects ("Vivi convierte leads en vecinos"). It is a WhatsApp housing advisor that already knows the prospect before greeting them: it cross-references the lead against Colsubsidio affiliate data, computes buying capacity with a deterministic and auditable engine, recommends housing projects via a k-nearest-neighbors "buyer twin" model, and hands the human advisor a ready-made "ficha" with score, subsidies, and priority. Success is a shorter, better-informed handoff from bot to advisor.

## Positioning

Deterministic, auditable capacity/eligibility engine (not LLM-generated numbers) combined with automatic affiliate-data cross-referencing — the advisor never has to ask the prospect for information Colsubsidio already has on file. Recommendations come from a similarity-based "buyer twin" (kNN over past buyers), not generic catalog filters.

## Operating Context

- Built for Hackathon Colsubsidio × 30X; two-block hackathon git flow (`feature/bloque-a` decision core, `feature/bloque-b` data/experience), integration only through the HTTP API and the Contract (doc 10 in the project wiki).
- Advisor works the dashboard alongside the lead's live WhatsApp chat, moving leads through a queue (`cola`) into a qualified "ficha" state before handoff (`ENTREGADO`).
- Domain vocabulary from the specs: `CALIFICADO`, `ASESOR`, `ficha`, `capacidad`, `GemeloKNN`, `cupo`, badges for data provenance (`verificado` / `declarado` / `inferido`), confidence levels (`ALTA`/`MEDIA`/`BAJA`), withdrawal-rate alert threshold at `0.20`.

## Capabilities and Constraints

- Backend: Go 1.24+, ADK Go 2.0, Clean Architecture. Frontend: TypeScript, MVC pattern (`web/src/models|views|controllers`). DB: PostgreSQL. LLM: Gemini primary, switchable to Qwen/Anthropic. Deploy: Heroku.
- The lead-qualification, capacity, and recommendation engine is deterministic and provider-free; ficha generation explicitly must not call an LLM (`openspec/specs/generar-ficha`).
- The doc 10 Contract governs data shapes across the two hackathon blocks and overrides any conflict; doc 13 "Criterios del Motor" fixes exact numbers the code must reproduce.
- Accessibility: focus-visible outlines and `prefers-reduced-motion` support are already implemented per "Doc 12 §5" of the project wiki.

## Brand Commitments

- Product/assistant name: "Vivi." Tagline: "Vivi convierte leads en vecinos."
- Colsubsidio brand system already established in `web/src/estilos/marca.css`: primary blue `#003DA6`, accent yellow `#FFC700`, semaphore colors for lead status (green/amber/gray), Colsubsidio logo in the top bar. Treat this as binding — not open for reinterpretation without cause.

## Evidence on Hand

- All data currently powering the dashboard is synthetic/mocked for the hackathon: `data/afiliados_mock.json`, `data/proyectos.json`, `data/compradores.json`, `data/mapa_proyectos.json`, and `web/src/mock/servidor-mock.ts`. There is no live connection to real Colsubsidio affiliate or catalog systems today.
- Project brochures exist as real content under `data/brochures/*.md` (e.g., `inari.md`, `versalles.md`, `zarzal.md`) — usable as real project copy/evidence, not fabricated.
- Do not invent real testimonials, live production metrics, or claims of a live Colsubsidio data connection; the demo runs entirely on the mock/fixture data above.

## Product Principles

1. Determinism and auditability over LLM guesswork wherever a number reaches the advisor (capacity, score, recommendations).
2. Advisor time is the scarce resource: surface only what changes the call — status, capacity, top recommendations, priority — never raw data the advisor would have to interpret.
3. The bot already knows what Colsubsidio knows: never make the prospect re-declare data the affiliate record already holds.
4. Contract (doc 10) is the source of truth across blocks; UI and engine must reproduce its exact field order and thresholds, not approximate them.

## Accessibility & Inclusion

Doc 12 §5 of the project wiki already establishes a standard: visible focus outlines on interactive elements and `prefers-reduced-motion` support, both already implemented in `marca.css`. No further requirement confirmed beyond this.
