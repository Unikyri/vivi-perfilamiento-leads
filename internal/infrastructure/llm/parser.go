package llm

import (
	"bytes"
	"encoding/json"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
)

type campoWire struct {
	Campo     *string         `json:"campo"`
	Valor     json.RawMessage `json:"valor"`
	Fuente    *string         `json:"fuente"`
	Confianza *float64        `json:"confianza"`
	Confirm   *bool           `json:"requiere_confirmacion"`
}
type salidaWire struct {
	Campos    *[]campoWire     `json:"campos_extraidos"`
	Intencion *json.RawMessage `json:"intencion"`
	Respuesta *string          `json:"respuesta"`
	Accion    *string          `json:"accion"`
}

func invalid(k ErrorKind) (usecase.SalidaTurno, error) {
	return usecase.SalidaTurno{}, providerError(k, nil)
}
func ParseSalida(data []byte, motor map[string]int64) (usecase.SalidaTurno, error) {
	var w salidaWire
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if err := d.Decode(&w); err != nil {
		if json.Valid(data) {
			return invalid(KindInvalidOutput)
		}
		return invalid(KindMalformed)
	}
	var extra any
	if d.Decode(&extra) != io.EOF {
		return invalid(KindMalformed)
	}
	if w.Campos == nil || w.Intencion == nil || w.Respuesta == nil || w.Accion == nil {
		return invalid(KindInvalidOutput)
	}
	out := usecase.SalidaTurno{Respuesta: *w.Respuesta, Accion: *w.Accion, CamposExtraidos: make([]usecase.CampoExtraido, len(*w.Campos))}
	for i, c := range *w.Campos {
		if c.Campo == nil || c.Fuente == nil || c.Confianza == nil || c.Confirm == nil || len(c.Valor) == 0 || !domain.CamposReconocidos[*c.Campo] || (*c.Fuente != "DECLARADO" && *c.Fuente != "INFERIDO") || *c.Confianza < 0 || *c.Confianza > 1 {
			return invalid(KindInvalidOutput)
		}
		var value any
		json.Unmarshal(c.Valor, &value)
		out.CamposExtraidos[i] = usecase.CampoExtraido{Campo: *c.Campo, Valor: value, Fuente: domain.FuenteCampo(*c.Fuente), Confianza: *c.Confianza, RequiereConfirmacion: *c.Confirm}
	}
	var intent domain.Intencion
	x := json.NewDecoder(bytes.NewReader(*w.Intencion))
	x.DisallowUnknownFields()
	if x.Decode(&intent) != nil || x.Decode(&extra) != io.EOF || !validLevel(string(intent.Nivel)) || !validLevel(string(intent.Confianza)) || intent.Senales == nil {
		return invalid(KindInvalidOutput)
	}
	out.Intencion = intent
	if !validAction(out.Accion) || !validResponse(out.Respuesta, motor) {
		return invalid(KindInvalidOutput)
	}
	return out, nil
}
func validLevel(v string) bool { return v == "ALTA" || v == "MEDIA" || v == "BAJA" }
func validAction(v string) bool {
	switch v {
	case usecase.AccionContinuar, usecase.AccionPerfilCompleto, usecase.AccionConsentimientoSi, usecase.AccionConsentimientoNo, usecase.AccionPausarContacto, usecase.AccionFueraDeDominio, usecase.AccionAudioInint:
		return true
	}
	return false
}

var currencyPattern = regexp.MustCompile(`(?i)(?:\$\s*|\b(?:cop|pesos?)\s*)(-?\d(?:[\d.,]*\d)?)\s*(mill(?:on|ón)(?:es)?)?|(-?\d(?:[\d.,]*\d)?)\s*(cop|pesos?|mill(?:on|ón)(?:es)?)`)
var bareAmountPattern = regexp.MustCompile(`-?\d(?:[\d.,]*\d)?`)
var monetaryContextPattern = regexp.MustCompile(`(?i)\b(?:mensual(?:es)?|cuota(?:s)?|ingreso(?:s)?|salario(?:s)?|presupuesto(?:s)?|precio(?:s)?|valor(?:es)?)\b`)
var lexicalWordPattern = regexp.MustCompile(`[\p{L}\p{N}]+`)

const maxContextGapWords = 3

func validResponse(response string, motor map[string]int64) bool {
	allowed := map[int64]bool{}
	for _, n := range motor {
		allowed[n] = true
	}
	for _, match := range currencyPattern.FindAllStringSubmatch(response, -1) {
		token, unit := match[1], match[2]
		if token == "" {
			token, unit = match[3], match[4]
		}
		n, ok := parseCurrencyAmount(token, unit)
		if !ok || !allowed[n] {
			return false
		}
	}
	for _, indexes := range bareAmountPattern.FindAllStringIndex(response, -1) {
		start, end := indexes[0], indexes[1]
		if hasNumericBoundary(response, start, end) || !isContextualBareAmount(response[start:end]) || !hasMonetaryContext(response, start, end) {
			continue
		}
		n, ok := parseCurrencyAmount(response[start:end], "")
		if !ok || !allowed[n] {
			return false
		}
	}
	return true
}

func hasNumericBoundary(response string, start, end int) bool {
	if start > 0 && isNumericCharacter(response[start-1]) {
		return true
	}
	if start > 1 && (response[start-1] == '.' || response[start-1] == ',') && isNumericCharacter(response[start-2]) {
		return true
	}
	if start > 0 && (response[start-1] == '-' || response[start-1] == '/') {
		return true
	}
	if end < len(response) && isNumericCharacter(response[end]) {
		return true
	}
	if end+1 < len(response) && (response[end] == '.' || response[end] == ',') && isNumericCharacter(response[end+1]) {
		return true
	}
	if end < len(response) && (response[end] == '-' || response[end] == '/') {
		return true
	}
	return false
}

func isNumericCharacter(value byte) bool {
	return value >= '0' && value <= '9'
}

func isContextualBareAmount(token string) bool {
	if strings.HasPrefix(token, "-") {
		return false
	}
	whole, _, ok := splitNumericAmount(token)
	if !ok || len(strings.TrimLeft(whole, "0")) < 4 {
		return false
	}
	// A bare year or compact YYYYMMDD date is not a monetary amount.
	digits := strings.ReplaceAll(whole, ".", "")
	if (len(digits) == 4 && digits >= "1900" && digits <= "2100") || (len(digits) == 8 && digits[0:4] >= "1900" && digits[0:4] <= "2100") {
		return false
	}
	return true
}

func hasMonetaryContext(response string, start, end int) bool {
	for _, context := range monetaryContextPattern.FindAllStringIndex(response, -1) {
		if context[1] <= start && contextPairIsAdjacent(response, context[1], start) {
			return true
		}
		if context[0] >= end && contextPairIsAdjacent(response, end, context[0]) {
			return true
		}
	}
	return false
}

func contextPairIsAdjacent(response string, left, right int) bool {
	between := response[left:right]
	if strings.ContainsAny(between, ".!?;\n\r,") || containsNumericToken(between) {
		return false
	}
	return len(lexicalWordPattern.FindAllStringIndex(between, -1)) <= maxContextGapWords
}

func containsNumericToken(value string) bool {
	for _, match := range bareAmountPattern.FindAllString(value, -1) {
		if strings.IndexFunc(match, func(r rune) bool { return r >= '0' && r <= '9' }) >= 0 {
			return true
		}
	}
	return false
}

const maxInt64Value = int64(^uint64(0) >> 1)

func parseCurrencyAmount(token, unit string) (int64, bool) {
	if strings.Contains(strings.ToLower(unit), "millon") || strings.Contains(strings.ToLower(unit), "millón") {
		return parseMillions(token)
	}
	whole, fraction, ok := splitNumericAmount(token)
	if !ok || strings.Trim(fraction, "0") != "" {
		return 0, false
	}
	return parseWholeAmount(whole)
}

func parseMillions(token string) (int64, bool) {
	whole, fraction, ok := splitNumericAmount(token)
	if !ok {
		return 0, false
	}
	if len(fraction) > 6 {
		if strings.Trim(fraction[6:], "0") != "" {
			return 0, false
		}
		fraction = fraction[:6]
	}
	wholeValue, ok := parseWholeAmount(whole)
	if !ok || wholeValue > maxInt64Value/1_000_000 {
		return 0, false
	}
	result := wholeValue * 1_000_000
	if fraction == "" {
		return result, true
	}
	fractionValue, err := strconv.ParseInt(fraction, 10, 64)
	if err != nil {
		return 0, false
	}
	for len(fraction) < 6 {
		fractionValue *= 10
		fraction += "0"
	}
	if fractionValue > maxInt64Value-result {
		return 0, false
	}
	return result + fractionValue, true
}

func splitNumericAmount(token string) (whole, fraction string, ok bool) {
	token = strings.TrimSpace(token)
	if token == "" || strings.HasPrefix(token, "-") || strings.ContainsAny(token, " +") {
		return "", "", false
	}
	commas, dots := strings.Count(token, ","), strings.Count(token, ".")
	decimalSeparator := byte(0)
	switch {
	case commas > 0 && dots > 0:
		// Colombian notation uses dots for grouping and comma for decimals.
		if strings.LastIndex(token, ",") < strings.LastIndex(token, ".") {
			return "", "", false
		}
		decimalSeparator = ','
	case commas == 1 && dots == 0:
		if len(token)-strings.LastIndex(token, ",")-1 <= 2 {
			decimalSeparator = ','
		}
	case dots == 1 && commas == 0:
		if len(token)-strings.LastIndex(token, ".")-1 <= 2 {
			decimalSeparator = '.'
		}
	}
	if decimalSeparator != 0 {
		parts := strings.Split(token, string(decimalSeparator))
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" || !allDigits(parts[1]) {
			return "", "", false
		}
		grouped := parts[0]
		groupSeparator := byte(0)
		if strings.Contains(grouped, ".") {
			groupSeparator = '.'
		} else if strings.Contains(grouped, ",") {
			groupSeparator = ','
		}
		if !validGroupedInteger(grouped, groupSeparator) {
			return "", "", false
		}
		return strings.NewReplacer(".", "", ",", "").Replace(grouped), parts[1], true
	}
	separator := byte(0)
	if commas > 0 {
		separator = ','
	} else if dots > 0 {
		separator = '.'
	}
	if !validGroupedInteger(token, separator) {
		return "", "", false
	}
	return strings.NewReplacer(".", "", ",", "").Replace(token), "", true
}

func validGroupedInteger(value string, separator byte) bool {
	if separator == 0 {
		return allDigits(value)
	}
	parts := strings.Split(value, string(separator))
	if len(parts) < 2 || len(parts[0]) < 1 || len(parts[0]) > 3 || !allDigits(parts[0]) {
		return false
	}
	for _, part := range parts[1:] {
		if len(part) != 3 || !allDigits(part) {
			return false
		}
	}
	return true
}

func allDigits(value string) bool {
	if value == "" {
		return false
	}
	for _, char := range value {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

func parseWholeAmount(value string) (int64, bool) {
	if !allDigits(value) {
		return 0, false
	}
	n, err := strconv.ParseInt(value, 10, 64)
	return n, err == nil
}
