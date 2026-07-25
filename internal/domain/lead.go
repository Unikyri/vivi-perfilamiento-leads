package domain

import "time"

// Lead is the central system entity.
type Lead struct {
	LeadID        string     `json:"lead_id"`
	Nombre        string     `json:"nombre"`
	Telefono      string     `json:"telefono"`
	Cedula        string     `json:"cedula"`
	Fuente        string     `json:"fuente"`
	Estado        EstadoLead `json:"estado"`
	Ruta          Ruta       `json:"ruta"`
	Afiliado      bool       `json:"afiliado"`
	Prioridad     float64    `json:"prioridad"`
	ConsumeCupo10 bool       `json:"consume_cupo_10"`
	Perfil        Perfil     `json:"perfil"`
	Capacidad     *Capacidad `json:"capacidad"`
	Intencion     *Intencion `json:"intencion"`
	Version       int        `json:"-"`
	CreadoEn      time.Time  `json:"creado_en"`
	ActualizadoEn time.Time  `json:"actualizado_en"`
}

type Mensaje struct {
	MensajeID     string         `json:"mensaje_id"`
	LeadID        string         `json:"-"`
	Autor         AutorMensaje   `json:"autor"`
	TipoContenido TipoContenido  `json:"tipo_contenido"`
	Texto         string         `json:"texto"`
	CreadoEn      time.Time      `json:"creado_en"`
	Adjunto       map[string]any `json:"adjunto"`
}
