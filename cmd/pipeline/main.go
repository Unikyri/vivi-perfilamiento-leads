package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/pipeline"
)

func main() {
	validar := flag.Bool("validar", false, "Solo imprime distribuciones sin escribir archivos (RF-M1-05)")
	flag.Parse()

	const entrada = "data/hackathon_VIVIENDAv2.xlsx"
	const salida = "data/compradores.json"

	cs, err := pipeline.LeerCompradores(entrada)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR leyendo el dataset: %v\n", err)
		os.Exit(1)
	}
	if err := pipeline.Validar(cs); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR de validación: %v\n", err)
		fmt.Fprintln(os.Stderr, "Los conteos no cuadran. Revisar reglas de decodificación (RF-M1-01).")
		os.Exit(1)
	}

	if *validar {
		pipeline.ImprimirDistribuciones(cs)
		fmt.Println("\n✔ Todos los criterios duros OK.")
		return
	}

	sort.Slice(cs, func(i, j int) bool { return cs[i].ID < cs[j].ID }) // determinismo

	f, err := os.Create(salida)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR creando %s: %v\n", salida, err)
		os.Exit(1)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(cs); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR escribiendo JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("OK: %d compradores → %s\n", len(cs), salida)
}
