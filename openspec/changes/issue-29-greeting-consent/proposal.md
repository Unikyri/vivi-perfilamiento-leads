# Proposal: Personalized Greeting and Data Consent (Issue #29)

## Scope
Implement a first-message use case driven by `LeadNuevo`. Affiliate greetings use the actual non-zero `Capacidad.SubsidioAplicable`, name, policy URL, and one dream question; they never ask income or household. Non-affiliate greetings are warm, amount-free, include the policy URL, and ask one employment question. Consent wording is declarative, leaving exactly one question mark.

The greeting makes one LLM attempt and validates it; an error, empty, or non-compliant draft falls back to a deterministic template. `CONSENTIMIENTO_NO` is processed before field merge/capacity recalculation: retain the inbound denial, append one dignified farewell, use existing lifecycle edges to finish `DESPEDIDO` with `RutaDespedida`, and emit no completion event.

## Constraints
Keep the existing `LeadNuevo` subscription and Contract vocabulary. No migration, consent field, new lifecycle edge, coordinator expansion, or nutrition-plan policy. Idempotent replay is deferred.

## Evidence and verification
Table-driven tests cover affiliate/non-affiliate content, single-question rule, fallback, automatic subscription, and denial side effects. Validate focused use cases, full Go test/build/vet. Forecast: 200–260 authored lines.

## Rollback
Revert the single PR; there is no schema, contract, or persistent-data migration.
