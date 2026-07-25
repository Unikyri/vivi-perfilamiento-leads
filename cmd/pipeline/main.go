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
	const salidaCompradores = "data/compradores.json"
	const salidaProyectos = "data/proyectos.json"
	const rutaMapa = "data/mapa_proyectos.json"

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

	// Cargar mapa de proyectos y asignar IDs canónicos antes de escribir compradores.json
	mapa, err := pipeline.CargarMapa(rutaMapa)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR cargando mapa: %v\n", err)
		os.Exit(1)
	}
	cs = pipeline.AsignarProyectoID(mapa, cs)

	sort.Slice(cs, func(i, j int) bool { return cs[i].ID < cs[j].ID }) // determinismo

	if err := escribirJSON(salidaCompradores, cs); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR escribiendo %s: %v\n", salidaCompradores, err)
		os.Exit(1)
	}
	fmt.Printf("OK: %d compradores → %s\n", len(cs), salidaCompradores)

	// Construir catálogo de proyectos cruzando mapa con compradores
	proys, err := pipeline.ConstruirProyectos(mapa, cs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR construyendo proyectos: %v\n", err)
		os.Exit(1)
	}
	if err := escribirJSON(salidaProyectos, proys); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR escribiendo %s: %v\n", salidaProyectos, err)
		os.Exit(1)
	}
	fmt.Printf("OK: %d proyectos → %s\n", len(proys), salidaProyectos)
}

// escribirJSON crea un archivo con JSON indentado y determinista.
func escribirJSON(ruta string, v any) error {
	f, err := os.Create(ruta)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
