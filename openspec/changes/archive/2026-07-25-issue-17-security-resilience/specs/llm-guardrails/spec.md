# Delta for LLM Guardrails

## ADDED Requirements

### Requirement: Deterministic Input Containment

The system MUST classify adversarial text before either provider is invoked and return a Contract-valid redirection template with `FUERA_DE_DOMINIO`. The fixture MUST contain all 15 entries: jailbreak (1, 2, 14), prompt extraction (3, 4, 13), role change (5, 6, 7), third-party data (8, 9, 15), and out-of-domain (10, 11, 12). Recognized third-party requests MUST use the privacy template. **Decision:** the three out-of-domain entries are contained guardrail cases, not pass-through requests; each therefore receives the redirection template and zero calls. Novel wording outside the fixture MAY rely on the existing hardened prompt until a deterministic pattern is added.

#### Scenario: Fixture containment and zero calls
- GIVEN a counting primary and fallback provider and each fixture entry
- WHEN the text turn is processed
- THEN every one of the 15 results is a safe template with `FUERA_DE_DOMINIO`
- AND both provider call counts remain zero

#### Scenario: Permitted housing conversation
- GIVEN a non-adversarial housing or subsidy question
- WHEN its text turn is processed
- THEN it MUST reach the existing provider composition normally

### Requirement: Audio Guardrail Parity

The system MUST apply the same input containment and output validation to `ProcesarAudio` using its available transcribed/associated turn text. A blocked audio turn MUST return the same safe result without a provider call; accepted audio SHALL preserve existing provider capability behavior.

#### Scenario: Blocked audio turn
- GIVEN audio with associated text matching an adversarial fixture case
- WHEN `ProcesarAudio` is processed
- THEN it returns the template and invokes neither provider

### Requirement: Safe Provider Output

The system MUST locally suppress a provider response that exposes a system prompt, skill/config/internal-file marker, explicit foreign `lead_id`, or a monetary amount absent from `NumerosDelMotor`. It MUST return a Contract-valid safe result without a retry, regeneration, fallback, or extra provider call. A permitted motor amount MUST remain available in a normal response.

#### Scenario: Leakage and amount validation
- GIVEN a successful provider response with a leak or unauthorized peso amount
- WHEN output validation runs
- THEN the unsafe content is not returned and no additional provider call occurs

### Requirement: Privacy-Safe Observability

LLM metrics MUST emit structured JSON sufficient for NFR-D-03, including `lead_id`, event, latency, provider, and outcome kind. They MUST NOT log message text, audio, cédulas, prompts, raw responses, provider bodies, credentials, or secrets.

#### Scenario: Sanitized event
- GIVEN accepted, rejected, failed, and breaker-open turns
- WHEN metrics are recorded
- THEN each record contains allowed operational fields only

### Requirement: Decorated Composition and Resilience Preservation

The configuration factory MUST compose `Metricas(Guardarrailes(FallbackProvider(primary, fallback)))`; every wrapper MUST implement and delegate the existing `LLMProvider` port. Health MUST obtain the live breaker state through wrappers. Wrappers MUST NOT add a breaker, timeout, limiter, fallback path, or change eligibility: three consecutive eligible failures SHALL open the existing breaker for 60 seconds, and guardrail outcomes MUST NOT count as failures.

#### Scenario: Decorated breaker health
- GIVEN the decorated composition receives three eligible primary failures
- WHEN health is queried before 60 seconds elapse
- THEN it reports `ABIERTO` and the existing fallback behavior remains intact

### Requirement: Stable Contract and Configuration Boundary

This change MUST NOT alter Contract v1.1, `LLMProvider`, `EntradaTurno`, `SalidaTurno`, configuration signatures, environment-variable names, or public HTTP behavior.

#### Scenario: Existing configuration compatibility
- GIVEN a valid pre-change provider configuration
- WHEN the factory creates the decorated provider
- THEN selection and the unchanged port/config contracts remain valid
