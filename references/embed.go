// Package references exposes canonical normative JSON files as raw bytes
// via go:embed. No parsing, no types, no init().
package references

import _ "embed"

//go:embed constantes.json
var ConstantesJSON []byte

//go:embed subsidios_2026.json
var SubsidiosJSON []byte

//go:embed calendario_economico.json
var CalendarioJSON []byte
