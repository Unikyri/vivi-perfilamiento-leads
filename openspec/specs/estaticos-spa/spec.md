# estaticos-spa Specification

## Requirements

### Requirement: Packaged SPA
Production builds MUST package Vite output before Go compilation. `GET /` MUST return the SPA entry; hashed-asset paths MUST return packaged bytes.

#### Scenario: Entry and asset
- GIVEN packaged production assets
- WHEN `/` and a referenced asset are requested
- THEN each returns its matching content.

#### Scenario: Missing assets
- GIVEN Vite output is absent
- WHEN the artifact is built
- THEN the build fails without a runnable artifact.

### Requirement: Safe fallback
Unknown browser paths MUST return the entry. `/api/*` and `/salud` MUST retain their handlers. Static responses MUST set `nosniff` and an NFR-S-03-compliant CSP.

#### Scenario: Precedence
- GIVEN registered API and health handlers
- WHEN an API or health path is requested
- THEN static fallback is not used.
