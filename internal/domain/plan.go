package domain

import "time"

type Hito struct {
	HitoID      string     `json:"hito_id"`
	Tipo        TipoHito   `json:"tipo"`
	Fecha       string     `json:"fecha"`
	Monto       *int64     `json:"monto"`
	Descripcion string     `json:"descripcion"`
	Estado      EstadoHito `json:"estado"`
}

type PlanNutricion struct {
	PlanID           string     `json:"plan_id"`
	LeadID           string     `json:"lead_id"`
	Estado           EstadoPlan `json:"estado"`
	ConsentimientoEn *time.Time `json:"consentimiento_en"`
	Frecuencia       string     `json:"frecuencia"`
	MetaMonto        int64      `json:"meta_monto"`
	MetaDescripcion  string     `json:"meta_descripcion"`
	Hitos            []Hito     `json:"hitos"`
}
