# Delta for Procesar Mensaje

## ADDED Requirements

### Requirement: Consent denial precedes normal turn mutation

When a valid provider result has `Accion=CONSENTIMIENTO_NO`, `ProcesarMensaje` MUST persist the inbound denial evidence and invoke `SaludarLead.RechazarConsentimiento` before field normalization, `aplicarCampos`, capacity calculation, intention assignment, or normal response/completion handling. It MUST preserve the pre-denial profile, capacity, and intention, append exactly the refusal farewell, and publish no event. It SHALL retain existing validation and provider-dispatch behavior before this classified action is available.

#### Scenario: Denial with extracted fields
- GIVEN a `PERFILANDO` lead and a result with `CONSENTIMIENTO_NO` plus extracted fields
- WHEN the text turn is processed
- THEN no extracted field, capacity, or intention is saved
- AND only the inbound denial, farewell, final route, and final state are persisted

#### Scenario: Denial does not complete a profile
- GIVEN a lead whose profile would otherwise be complete
- WHEN the result is `CONSENTIMIENTO_NO`
- THEN no `PerfilCompleto` event is published
- AND the final state is `DESPEDIDO`
