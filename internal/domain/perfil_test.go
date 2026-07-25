package domain

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestEnumJSONWireValues(t *testing.T) {
	cases := []struct {
		name  string
		value any
		want  string
	}{
		{"EstadoLead", EstadoLeadEnNutricion, `"EN_NUTRICION"`}, {"Ruta", RutaRemarketing, `"REMARKETING"`},
		{"FuenteCampo", FuenteCampoVerificadoBase, `"VERIFICADO_BASE"`}, {"Nivel", NivelMedia, `"MEDIA"`},
		{"TipoMensaje", TipoMensajeAudio, `"AUDIO"`}, {"TipoContenido", TipoContenidoTarjetasProyectos, `"TARJETAS_PROYECTOS"`},
		{"AutorMensaje", AutorMensajeVivi, `"VIVI"`}, {"EstadoPlan", EstadoPlanCompletado, `"COMPLETADO"`},
		{"TipoHito", TipoHitoCesantias, `"CESANTIAS"`}, {"EstadoHito", EstadoHitoNotificado, `"NOTIFICADO"`},
		{"Semaforo", SemaforoAmbar, `"AMBAR"`},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if reflect.TypeOf(tt.value).Kind() != reflect.String {
				t.Errorf("underlying kind = %s, want string", reflect.TypeOf(tt.value).Kind())
			}
			got, err := json.Marshal(tt.value)
			if err != nil || string(got) != tt.want {
				t.Errorf("json.Marshal() = %s, %v; want %s", got, err, tt.want)
			}
		})
	}
}

func TestPerfilAccessors(t *testing.T) {
	perfil := Perfil{
		"int64":   {Valor: int64(12)},
		"int":     {Valor: 13},
		"float64": {Valor: 14.0},
		"bool":    {Valor: true},
		"texto":   {Valor: "hogar"},
	}

	enteroCases := []struct {
		name, key string
		want      int64
		ok        bool
	}{
		{"int64", "int64", 12, true},
		{"int", "int", 13, true},
		{"json float64", "float64", 14, true},
		{"absent", "missing", 0, false},
		{"incompatible", "texto", 0, false},
	}
	for _, tt := range enteroCases {
		t.Run("Entero/"+tt.name, func(t *testing.T) {
			got, ok := perfil.Entero(tt.key)
			if got != tt.want || ok != tt.ok {
				t.Errorf("Entero(%q) = (%d, %t), want (%d, %t)", tt.key, got, ok, tt.want, tt.ok)
			}
		})
	}

	boolCases := []struct {
		name, key string
		want      bool
		ok        bool
	}{
		{"valid", "bool", true, true}, {"absent", "missing", false, false}, {"incompatible", "texto", false, false},
	}
	for _, tt := range boolCases {
		t.Run("Booleano/"+tt.name, func(t *testing.T) {
			got, ok := perfil.Booleano(tt.key)
			if got != tt.want || ok != tt.ok {
				t.Errorf("Booleano(%q) = (%t, %t), want (%t, %t)", tt.key, got, ok, tt.want, tt.ok)
			}
		})
	}

	texto, ok := perfil.Texto("texto")
	if texto != "hogar" || !ok {
		t.Errorf("Texto(texto) = (%q, %t), want (hogar, true)", texto, ok)
	}
	if _, ok := perfil.Texto("missing"); ok {
		t.Error("Texto(missing) returned ok=true")
	}
}

func TestPerfilVerificationAndKeySets(t *testing.T) {
	perfil := Perfil{
		"base":     {Fuente: FuenteCampoVerificadoBase},
		"declared": {Fuente: FuenteCampoDeclarado},
		"inferred": {Fuente: FuenteCampoInferido},
	}
	if !perfil.EsVerificado("base") || perfil.EsVerificado("declared") || perfil.EsVerificado("inferred") || perfil.EsVerificado("missing") {
		t.Fatal("EsVerificado did not distinguish the three sources and absent fields")
	}
	assertKeys(t, CamposReconocidos, []string{
		"ingreso_hogar", "recursos_propios", "personas_hogar", "tipo_hogar", "edad", "zona_deseada",
		"plazo_compra_meses", "arriendo_actual", "tiene_vivienda", "recibio_subsidio", "reporte_credito",
		"situacion_laboral", "caja_externa", "hogar_con_afiliado", "cedula_familiar_afiliado", "cesantias_entidad",
		"categoria", "segmento",
	})
	assertKeys(t, CamposCriticos, []string{"ingreso_hogar", "recursos_propios", "tiene_vivienda", "recibio_subsidio"})
}

func TestPerfilSchemaAndJSON(t *testing.T) {
	tipo := reflect.TypeOf(CampoPerfil{})
	want := []struct {
		name, tag string
		typ       reflect.Type
	}{
		{"Valor", `json:"valor"`, reflect.TypeOf((*any)(nil)).Elem()},
		{"Fuente", `json:"fuente"`, reflect.TypeOf(FuenteCampo(""))},
		{"Confianza", `json:"confianza"`, reflect.TypeOf(float64(0))},
		{"RequiereConfirmacion", `json:"requiere_confirmacion"`, reflect.TypeOf(false)},
		{"ActualizadoEn", `json:"actualizado_en"`, reflect.TypeOf(time.Time{})},
	}
	if tipo.NumField() != len(want) {
		t.Fatalf("CampoPerfil has %d fields, want %d", tipo.NumField(), len(want))
	}
	for _, tt := range want {
		field, ok := tipo.FieldByName(tt.name)
		if !ok || field.Type != tt.typ || string(field.Tag) != tt.tag {
			t.Errorf("CampoPerfil.%s metadata = (%v, %q), want (%v, %q)", tt.name, field.Type, field.Tag, tt.typ, tt.tag)
		}
	}

	stamp := time.Date(2026, time.July, 25, 0, 0, 0, 0, time.UTC)
	got, err := json.Marshal(CampoPerfil{Valor: int64(7), Fuente: FuenteCampoVerificadoBase, Confianza: 0.9, RequiereConfirmacion: false, ActualizadoEn: stamp})
	if err != nil {
		t.Fatal(err)
	}
	wantJSON := `{"valor":7,"fuente":"VERIFICADO_BASE","confianza":0.9,"requiere_confirmacion":false,"actualizado_en":"2026-07-25T00:00:00Z"}`
	if string(got) != wantJSON {
		t.Errorf("CampoPerfil JSON = %s, want %s", got, wantJSON)
	}
	if reflect.TypeOf(Perfil{}).NumMethod() != 4 {
		t.Errorf("Perfil exported method count = %d, want 4", reflect.TypeOf(Perfil{}).NumMethod())
	}
}

func assertKeys(t *testing.T, got map[string]bool, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("key set length = %d, want %d", len(got), len(want))
	}
	for _, key := range want {
		if !got[key] {
			t.Errorf("key set missing %q or value is false", key)
		}
	}
}

func TestCapacityContractTypes(t *testing.T) {
	assertJSONFields(t, ItemDesglose{}, []string{"Concepto", "Monto", "Regla", "Fuente"})
	assertJSONFields(t, Capacidad{}, []string{"PresupuestoMax", "CreditoMax", "SubsidioAplicable", "RecursosPropios", "Ratio", "Confianza", "Desglose"})
	assertJSONFields(t, Intencion{}, []string{"Nivel", "Confianza", "Senales"})
	assertJSONFields(t, Recomendacion{}, []string{"ProyectoID", "Nombre", "Zona", "PrecioDesde", "Razon", "Vecinos", "TasaDesistimiento", "BrochureURL", "Recorrido360URL"})
	assertJSONFields(t, Proyecto{}, []string{"ProyectoID", "Nombre", "Zona", "PrecioDesde", "PrecioHasta", "EsVIS", "BrochureURL", "Recorrido360URL"})
	assertJSONFields(t, Comprador{}, []string{"ID", "ProyectoID", "Proyecto", "Etapa", "Afiliado", "Categoria", "Segmento", "RangoEdad", "PersonasACargo", "Piramide", "ValorCOP", "Entidad", "Medio", "Desistio", "FechaOpcion"})
}

func assertJSONFields(t *testing.T, value any, names []string) {
	t.Helper()
	typ := reflect.TypeOf(value)
	if typ.NumField() != len(names) {
		t.Fatalf("%s fields = %d, want %d", typ.Name(), typ.NumField(), len(names))
	}
	for _, name := range names {
		field, ok := typ.FieldByName(name)
		if !ok || field.Tag.Get("json") == "" {
			t.Errorf("%s.%s missing JSON metadata", typ.Name(), name)
		}
	}
}

func TestLeadAndPlanContractTypes(t *testing.T) {
	assertJSONFields(t, Lead{}, []string{"LeadID", "Nombre", "Telefono", "Cedula", "Fuente", "Estado", "Ruta", "Afiliado", "Prioridad", "ConsumeCupo10", "Perfil", "Capacidad", "Intencion", "Version", "CreadoEn", "ActualizadoEn"})
	assertJSONFields(t, Mensaje{}, []string{"MensajeID", "LeadID", "Autor", "TipoContenido", "Texto", "CreadoEn", "Adjunto"})
	assertJSONFields(t, Hito{}, []string{"HitoID", "Tipo", "Fecha", "Monto", "Descripcion", "Estado"})
	assertJSONFields(t, PlanNutricion{}, []string{"PlanID", "LeadID", "Estado", "ConsentimientoEn", "Frecuencia", "MetaMonto", "MetaDescripcion", "Hitos"})
}

func TestFichaContractTypes(t *testing.T) {
	assertJSONFields(t, AlertaDesistimiento{}, []string{"Activa", "TasaVecinos", "Detalle"})
	assertJSONFields(t, Identificacion{}, []string{"Nombre", "Afiliada", "Categoria", "Telefono"})
	assertJSONFields(t, Ficha{}, []string{"FichaID", "LeadID", "GeneradaEn", "ConfianzaPerfil", "BandaAdvertencia", "Identificacion", "Capacidad", "Perfil", "Intencion", "Recomendaciones", "Beneficios", "ArgumentosVenta", "AlertaDesistimiento", "ConsumeCupo10"})
}
