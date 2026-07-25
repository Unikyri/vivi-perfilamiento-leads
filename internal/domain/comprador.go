package domain

// Comprador representa un registro del dataset de compradores históricos (Excel v2).
// Los campos siguen el esquema definido en Contrato v1.1 §4.1.
type Comprador struct {
	ID             int    `json:"id"`
	Proyecto       string `json:"proyecto"`
	ProyectoID     string `json:"proyecto_id"`
	Etapa          string `json:"etapa"`
	FechaOpcion    string `json:"fecha_opcion"`
	Desistio       bool   `json:"desistio"`
	Entidad        string `json:"entidad"`
	Medio          string `json:"medio"`
	ValorCOP       int64  `json:"valor_cop"`
	Afiliado       bool   `json:"afiliado"`
	Segmento       string `json:"segmento"`
	Categoria      string `json:"categoria"`
	RangoEdad      string `json:"rango_edad"`
	PersonasACargo int    `json:"personas_a_cargo"`
	Piramide       string `json:"piramide"`
}
