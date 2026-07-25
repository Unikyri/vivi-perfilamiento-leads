package domain

import "time"

type AlertaDesistimiento struct {
	Activa      bool    `json:"activa"`
	TasaVecinos float64 `json:"tasa_vecinos"`
	Detalle     *string `json:"detalle"`
}

type Identificacion struct {
	Nombre    string `json:"nombre"`
	Afiliada  bool   `json:"afiliada"`
	Categoria string `json:"categoria"`
	Telefono  string `json:"telefono"`
}

// Ficha is the advisor-facing Contract v1.1 output.
type Ficha struct {
	FichaID             string              `json:"ficha_id"`
	LeadID              string              `json:"lead_id"`
	GeneradaEn          time.Time           `json:"generada_en"`
	ConfianzaPerfil     float64             `json:"confianza_perfil"`
	BandaAdvertencia    *string             `json:"banda_advertencia"`
	Identificacion      Identificacion      `json:"identificacion"`
	Capacidad           Capacidad           `json:"capacidad"`
	Perfil              Perfil              `json:"perfil"`
	Intencion           Intencion           `json:"intencion"`
	Recomendaciones     []Recomendacion     `json:"recomendaciones"`
	Beneficios          []string            `json:"beneficios"`
	ArgumentosVenta     []string            `json:"argumentos_venta"`
	AlertaDesistimiento AlertaDesistimiento `json:"alerta_desistimiento"`
	ConsumeCupo10       bool                `json:"consume_cupo_10"`
}
