package pipeline

import (
	"fmt"
	"sort"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
)

// ImprimirDistribuciones muestra las distribuciones clave del dataset
// para verificación visual rápida (RF-M1-05).
func ImprimirDistribuciones(cs []domain.Comprador) {
	total := len(cs)
	var afiliados, desistidos int
	categorias := make(map[string]int)
	segmentos := make(map[string]int)
	edades := make(map[string]int)
	proyectos := make(map[string]int)

	for _, c := range cs {
		if c.Afiliado {
			afiliados++
		}
		if c.Desistio {
			desistidos++
		}
		categorias[c.Categoria]++
		segmentos[c.Segmento]++
		edades[c.RangoEdad]++
		proyectos[c.Proyecto]++
	}

	fmt.Println("═══ Criterios duros (RF-M1-01) ═══")
	fmt.Printf("  total:      %d (esperado %d)\n", total, TotalEsperado)
	fmt.Printf("  afiliados:  %d (esperado %d)\n", afiliados, AfiliadosEsperados)
	fmt.Printf("  desistidos: %d (esperado %d)\n", desistidos, DesistidosEsperados)

	fmt.Println("\n═══ Distribución de categoría ═══")
	imprimirMapa(categorias)

	fmt.Println("\n═══ Distribución de segmento ═══")
	imprimirMapa(segmentos)

	fmt.Println("\n═══ Distribución de rango_edad ═══")
	imprimirMapa(edades)

	fmt.Println("\n═══ Top 10 proyectos por compradores ═══")
	imprimirTop(proyectos, 10)
}

// imprimirMapa imprime un mapa de conteos ordenado por valor descendente.
func imprimirMapa(m map[string]int) {
	type par struct {
		clave string
		n     int
	}
	pares := make([]par, 0, len(m))
	for k, v := range m {
		pares = append(pares, par{k, v})
	}
	sort.Slice(pares, func(i, j int) bool { return pares[i].n > pares[j].n })
	for _, p := range pares {
		fmt.Printf("  %-20s %d\n", p.clave, p.n)
	}
}

// imprimirTop imprime los N primeros del mapa ordenado por valor descendente.
func imprimirTop(m map[string]int, n int) {
	type par struct {
		clave string
		n     int
	}
	pares := make([]par, 0, len(m))
	for k, v := range m {
		pares = append(pares, par{k, v})
	}
	sort.Slice(pares, func(i, j int) bool { return pares[i].n > pares[j].n })
	if len(pares) > n {
		pares = pares[:n]
	}
	for i, p := range pares {
		fmt.Printf("  %2d. %-40s %d\n", i+1, p.clave, p.n)
	}
}
