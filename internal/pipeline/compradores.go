package pipeline

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
	"github.com/xuri/excelize/v2"
)

// Conteos de aceptación — verificados contra el archivo oficial v2.
const (
	TotalEsperado       = 4142
	AfiliadosEsperados  = 3020
	DesistidosEsperados = 550
)

// LeerCompradores transforma el xlsx v2 en []domain.Comprador.
func LeerCompradores(rutaXLSX string) ([]domain.Comprador, error) {
	f, err := excelize.OpenFile(rutaXLSX)
	if err != nil {
		return nil, fmt.Errorf("abriendo %s: %w", rutaXLSX, err)
	}
	defer f.Close()

	hojas := f.GetSheetList()
	if len(hojas) == 0 {
		return nil, fmt.Errorf("el archivo no tiene hojas")
	}
	filas, err := f.GetRows(hojas[0])
	if err != nil {
		return nil, fmt.Errorf("leyendo filas: %w", err)
	}
	if len(filas) < 2 {
		return nil, fmt.Errorf("el archivo no tiene datos")
	}

	out := make([]domain.Comprador, 0, len(filas)-1)
	for i, fila := range filas[1:] {
		c := domain.Comprador{ID: i + 1}
		c.Proyecto = col(fila, 0)
		c.ProyectoID = Slug(c.Proyecto)
		c.Etapa = col(fila, 1)
		c.FechaOpcion = fechaExcel(col(fila, 2))
		c.Desistio = strings.EqualFold(strings.TrimSpace(col(fila, 3)), "Si")
		c.Entidad = normalizarEntidad(col(fila, 4))
		c.Medio = col(fila, 5)
		c.ValorCOP = valorVivienda(col(fila, 6))
		c.Afiliado = strings.TrimSpace(col(fila, 8)) != "" // PERIODO_AFILIADO
		c.Segmento = mapear(segmentoMap, col(fila, 9))
		c.Categoria = mapear(categoriaMap, col(fila, 10))
		c.RangoEdad = mapear(edadMap, col(fila, 11))
		c.PersonasACargo = entero(col(fila, 12))
		c.Piramide = col(fila, 14)
		out = append(out, c)
	}
	return out, nil
}

// Validar comprueba los criterios de aceptación duros (RF-M1-01).
func Validar(cs []domain.Comprador) error {
	total := len(cs)
	var afiliados, desistidos int
	for _, c := range cs {
		if c.Afiliado {
			afiliados++
		}
		if c.Desistio {
			desistidos++
		}
	}
	if total != TotalEsperado {
		return fmt.Errorf("total = %d, esperado %d", total, TotalEsperado)
	}
	if afiliados != AfiliadosEsperados {
		return fmt.Errorf("afiliados = %d, esperado %d", afiliados, AfiliadosEsperados)
	}
	if desistidos != DesistidosEsperados {
		return fmt.Errorf("desistidos = %d, esperado %d", desistidos, DesistidosEsperados)
	}
	return nil
}

func col(fila []string, i int) string {
	if i < len(fila) {
		return strings.TrimSpace(fila[i])
	}
	return ""
}

func mapear(m map[string]string, v string) string {
	if r, ok := m[strings.TrimSpace(v)]; ok {
		return r
	}
	return "SIN_DATO"
}

// valorVivienda aplica la regla VLR_VIVIENDA / 10000 (instructivo oficial).
func valorVivienda(v string) int64 {
	n, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
	if err != nil {
		return 0
	}
	return n / 10000
}

func entero(v string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(v))
	return n
}

func normalizarEntidad(v string) string {
	s := strings.ToUpper(strings.TrimSpace(v))
	switch {
	case s == "":
		return "SIN_DATO"
	case strings.Contains(s, "COLSUBSIDIO"):
		return "COLSUBSIDIO"
	case strings.Contains(s, "CONTADO"):
		return "CONTADO"
	default:
		return "BANCO"
	}
}

// fechaExcel convierte el serial de Excel a "YYYY-MM-DD".
func fechaExcel(v string) string {
	n, err := strconv.Atoi(strings.TrimSpace(v))
	if err != nil {
		return ""
	}
	t, err := excelize.ExcelDateToTime(float64(n), false)
	if err != nil {
		return ""
	}
	return t.Format("2006-01-02")
}

var noAlfaNum = regexp.MustCompile(`[^a-z0-9]+`)

// Slug convierte un nombre de proyecto a identificador ASCII en minúsculas.
func Slug(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	for _, r := range s {
		switch r {
		case 'á', 'à', 'ä', 'â':
			b.WriteRune('a')
		case 'é', 'è', 'ë', 'ê':
			b.WriteRune('e')
		case 'í', 'ì', 'ï', 'î':
			b.WriteRune('i')
		case 'ó', 'ò', 'ö', 'ô':
			b.WriteRune('o')
		case 'ú', 'ù', 'ü', 'û':
			b.WriteRune('u')
		case 'ñ':
			b.WriteRune('n')
		default:
			if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' {
				b.WriteRune(r)
			}
		}
	}
	return strings.Trim(noAlfaNum.ReplaceAllString(b.String(), "_"), "_")
}
