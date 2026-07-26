# Reloj Produccion Specification

## Purpose

Separate auditable operational time from the persisted simulated date used by the demo.

## Requirements

### Requirement: Wall time and simulated time are distinct

The production `Reloj.Ahora()` MUST return current UTC wall time. `FechaSimulada()` MUST return the adapter's persisted simulated UTC date, and `Avanzar(at)` MUST change only that simulated value. The `usecase.Reloj` interface and its call sites SHALL remain unchanged.

#### Scenario: Advance does not move audit time
- GIVEN a loaded simulated date and a later demo advance
- WHEN `Avanzar` completes
- THEN `FechaSimulada()` equals the advanced date
- AND `Ahora()` remains current UTC wall time rather than the advanced date

### Requirement: Persisted demo behavior is retained

Zero-value initialization MUST persist a UTC fallback date. `AvanzarDemo` MUST persist a forward date before updating the adapter, MUST reject a backward date with `ErrTiempoSimuladoAtras`, and MUST retain reset, `TickReloj`, milestone comparison, and `/salud.fecha_simulada` behavior.

#### Scenario: Restart and non-regression
- GIVEN a persisted simulated date
- WHEN a new production clock is constructed and a backward advance is requested
- THEN the loaded date is returned and the backward request is rejected

#### Scenario: Health reports demo time
- GIVEN demo time has advanced
- WHEN `/salud` is queried
- THEN `fecha_simulada` reports the simulated date, not operational wall time
