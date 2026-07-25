package domain

// Proyecto representa un proyecto inmobiliario del catálogo de Vivi.
// Solo entran los 16 proyectos con brochure + recorrido 360.
// Campos según Contrato v1.1 §4.2.
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
