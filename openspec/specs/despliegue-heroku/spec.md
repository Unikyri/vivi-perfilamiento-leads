# despliegue-heroku Specification

## Requirements

### Requirement: Reproducible packaged build
Heroku MUST run the Node buildpack before Go so Vite assets are packaged before Go compilation. CI and Heroku MUST use `go.mod`'s Go version and MUST fail when assets are absent.

#### Scenario: Ordered build
- GIVEN web dependencies and source
- WHEN the build runs
- THEN Node packages assets before Go compiles on the declared version.

#### Scenario: Missing assets
- GIVEN no Vite output
- WHEN the build runs
- THEN it fails without a deployable artifact.

### Requirement: Inert secret-safe configuration
Configuration MUST default `DEMO_SEED` to false, MUST NOT contain credential values, and MUST remain inert in this change. Future promotion SHALL require successful CI and operator-provided secret references.

#### Scenario: No promotion
- GIVEN configuration without credentials
- WHEN it is evaluated
- THEN no secret is exposed and no deployment occurs.

### Requirement: Real acceptance and health guidance
Acceptance MUST use the real API flow; `?mock=1` MAY be a labeled visual fallback only. Docs MUST cover `/salud` verification and MAY offer an external `/salud` ping.

#### Scenario: Acceptance and ping
- GIVEN real and mock flows
- WHEN acceptance or anti-sleep guidance is assessed
- THEN only real API passes; docs list health verification and optional ping.
