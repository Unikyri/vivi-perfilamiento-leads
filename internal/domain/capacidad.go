package domain

// ItemDesglose is one line of the capacity breakdown (Contract v1.1 §2.2).
type ItemDesglose struct {
	Concepto string      `json:"concepto"`
	Monto    int64       `json:"monto"`
	Regla    string      `json:"regla"`
	Fuente   FuenteCampo `json:"fuente"`
}

// Capacidad is the deterministic engine output.
type Capacidad struct {
	PresupuestoMax    int64          `json:"presupuesto_max"`
	CreditoMax        int64          `json:"credito_max"`
	SubsidioAplicable int64          `json:"subsidio_aplicable"`
	RecursosPropios   int64          `json:"recursos_propios"`
	Ratio             float64        `json:"ratio"`
	Confianza         float64        `json:"confianza"`
	Desglose          []ItemDesglose `json:"desglose"`
}

type Intencion struct {
	Nivel     Nivel    `json:"nivel"`
	Confianza Nivel    `json:"confianza"`
	Senales   []string `json:"senales"`
}

type Recomendacion struct {
	ProyectoID        string  `json:"proyecto_id"`
	Nombre            string  `json:"nombre"`
	Zona              string  `json:"zona"`
	PrecioDesde       int64   `json:"precio_desde"`
	Razon             string  `json:"razon"`
	Vecinos           int     `json:"vecinos"`
	TasaDesistimiento float64 `json:"tasa_desistimiento"`
	BrochureURL       string  `json:"brochure_url"`
	Recorrido360URL   string  `json:"recorrido_360_url"`
}

// Proyecto is a catalog entry produced by the data pipeline.
type Proyecto struct {
	ProyectoID      string `json:"proyecto_id"`
	Nombre          string `json:"nombre"`
	Zona            string `json:"zona"`
	PrecioDesde     int64  `json:"precio_desde"`
	PrecioHasta     int64  `json:"precio_hasta"`
	EsVIS           bool   `json:"es_vis"`
	BrochureURL     string `json:"brochure_url"`
	Recorrido360URL string `json:"recorrido_360_url"`
}

// Comprador is a row from data/compradores.json used by the kNN twin.
type Comprador struct {
	ID             int    `json:"id"`
	ProyectoID     string `json:"proyecto_id"`
	Proyecto       string `json:"proyecto"`
	Etapa          string `json:"etapa"`
	Afiliado       bool   `json:"afiliado"`
	Categoria      string `json:"categoria"`
	Segmento       string `json:"segmento"`
	RangoEdad      string `json:"rango_edad"`
	PersonasACargo int    `json:"personas_a_cargo"`
	Piramide       string `json:"piramide"`
	ValorCOP       int64  `json:"valor_cop"`
	Entidad        string `json:"entidad"`
	Medio          string `json:"medio"`
	Desistio       bool   `json:"desistio"`
	FechaOpcion    string `json:"fecha_opcion"`
}
