# Saludar Lead Specification

## Purpose

Produce and persist Vivi's safe first message after `LeadNuevo` without changing Contract vocabulary or the coordinator table.

## Requirements

### Requirement: Validated deterministic greeting

On `LeadNuevo`, the system MUST make at most one `GenerarTurno` attempt for the greeting and MUST persist one `VIVI` text message. It MUST validate the returned text against the rules below; provider error, nil/empty text, or a failed validation MUST select the deterministic template. No provider retry is permitted.

#### Scenario: Provider draft is compliant
- GIVEN a profiled lead and one compliant provider response
- WHEN `LeadNuevo` is handled
- THEN exactly one greeting is persisted using the validated response
- AND no second provider call occurs

#### Scenario: Provider is unavailable or invalid
- GIVEN a provider error, empty response, or non-compliant response
- WHEN `LeadNuevo` is handled
- THEN the deterministic greeting is persisted
- AND the handler returns no unpersisted greeting

### Requirement: Audience-specific, single-question copy

For an affiliate with non-zero `Capacidad.SubsidioAplicable`, the greeting MUST contain the lead name, that motor-derived subsidy formatted as `$X,YM`, `URLPolitica`, declarative consent wording, and exactly one dream question. It MUST NOT ask income or household questions. For a non-affiliate, or an affiliate with zero/absent applicable subsidy, it MUST be warm, contain `URLPolitica`, contain no COP/peso amount, and ask exactly one job-situation question. Every selected greeting MUST contain exactly one `?`.

#### Scenario: Affiliate copy
- GIVEN an affiliate with a non-zero applicable subsidy
- WHEN the greeting is selected
- THEN it contains only that motor subsidy and one dream question
- AND it contains neither income nor household question

#### Scenario: Non-affiliate copy
- GIVEN a non-affiliate or zero-subsidy affiliate
- WHEN the greeting is selected
- THEN it contains no monetary amount and one job-situation question

### Requirement: Consent refusal is terminal and profile-safe

`RechazarConsentimiento` MUST append a dignified farewell, set `Ruta=DESPEDIDA`, and reach `DESPEDIDO` through existing `PERFILANDO → CALIFICADO → DESPEDIDO` edges. It MUST NOT mutate profile, capacity, or intention and MUST NOT publish `PerfilCompleto`.

#### Scenario: Refusal farewell
- GIVEN a `PERFILANDO` lead that refuses consent
- WHEN refusal is rejected
- THEN its stored route and final state are `DESPEDIDA` and `DESPEDIDO`
- AND the conversation contains one farewell and no completion event
