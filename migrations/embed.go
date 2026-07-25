// Package migrations exposes the canonical DDL as an embedded string.
// The SQL is the single source of truth for the database schema (Contract §5).
package migrations

import _ "embed"

//go:embed 001_esquema_inicial.sql
var EsquemaInicial string
