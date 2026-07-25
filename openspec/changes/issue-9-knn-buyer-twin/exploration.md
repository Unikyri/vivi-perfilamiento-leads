## Exploration: issue-9-knn-buyer-twin

### Current State
Issue #9 owns only a deterministic buyer-twin Gower kNN in `internal/domain/motor/knn.go` and `knn_test.go`; #8 owns capacity files and is excluded. Doc 13 §3 and Contract §4.1/§4.2 govern the behavior. Gower weights are category .35, zone .20, age .15, affiliation .15, dependents .15; absent dimensions are omitted and remaining weights renormalized. Results are limited to `min(k,n)` and ties break by buyer ID ascending.

`domain.Comprador` has canonical `ProyectoID`, affiliation, category, age bracket, dependents, price, and outcome, but no zone. `domain.Proyecto` owns catalog zone. `Perfil` carries lead income and desired zone. The catalog/pipeline source supplies zones separately by canonical project ID.

### Sketch Reconciliation
- Do not add/infer buyer zone. Zone comes only from a separate `map[ProyectoID]Zona` catalog input.
- Do not use project name, slug, price, or price-based category imputation as a proxy.
- Do not use independently indexed parallel slices; accept buyer records and keyed catalog metadata.
- Compute category mode and age median from affiliated buyers of the same project only. With fewer than five affiliates, omit unavailable dimensions rather than impute.
- Map lead `ingreso_hogar` to A/B/C only for kNN; absent income leaves the category absent.

### Affected Areas
- `internal/domain/motor/knn.go` — pure projection, project-local imputation, Gower distance, and deterministic neighbors.
- `internal/domain/motor/knn_test.go` — KNN-1..5, safe zone provenance, no-proxy, and missing-data tests.
- `internal/domain/motor/capacidad.go` and `capacidad_test.go` — read-only, owned by #8.
- Pipeline/data/Docs — read-only.

### Approaches
1. **Buyer records plus keyed catalog-zone map** — recommended. It preserves contracts, avoids positional misalignment, and supports deterministic project-local statistics.
2. **Add zone to Comprador** — rejected; alters contract shape and duplicates catalog data.
3. **Parallel slices** — rejected; unsafe after filtering or sorting.
4. **Name/price proxy** — rejected by Doc 13 and data authority.

### Recommendation
Use the pure motor API with buyer records, explicit lead affiliation, and catalog zones keyed by canonical project ID. Build same-project affiliate statistics internally; never mutate inputs. Resolve a deterministic representative for each age bracket and treat unknown catalog IDs as missing zone.

### Risks
- `RangoEdad` needs a fixed numeric representation, especially `55+`.
- Unknown `ProyectoID` must omit zone, never guess.
- Issue #8 boundary must remain untouched.

### Ready for Proposal
Yes — scope, safe inputs, rejected sketch drift, and remaining age-representation decision are explicit.
