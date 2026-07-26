package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
)

const (
	demoClockKey     = "fecha_simulada"
	approvedDemoDate = "2026-07-26T00:00:00Z"
)

type DemoRepository struct{ pool pgxPool }

func NuevoDemoRepository(pool pgxPool) *DemoRepository { return &DemoRepository{pool: pool} }

var _ usecase.DemoRepository = (*DemoRepository)(nil)
var _ usecase.DemoResetRepository = (*DemoRepository)(nil)
var _ usecase.DemoSeedRepository = (*DemoRepository)(nil)
var _ usecase.DemoResetSeedRepository = (*DemoRepository)(nil)

const demoLeadUpsertSQL = `INSERT INTO leads (lead_id,nombre,telefono,cedula,fuente,estado,ruta,afiliado,prioridad,consume_cupo_10,perfil,capacidad,intencion,version,creado_en,actualizado_en) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16) ON CONFLICT (lead_id) DO UPDATE SET nombre=EXCLUDED.nombre,telefono=EXCLUDED.telefono,cedula=EXCLUDED.cedula,fuente=EXCLUDED.fuente,estado=EXCLUDED.estado,ruta=EXCLUDED.ruta,afiliado=EXCLUDED.afiliado,prioridad=EXCLUDED.prioridad,consume_cupo_10=EXCLUDED.consume_cupo_10,perfil=EXCLUDED.perfil,capacidad=EXCLUDED.capacidad,intencion=EXCLUDED.intencion,version=EXCLUDED.version,creado_en=EXCLUDED.creado_en,actualizado_en=EXCLUDED.actualizado_en`

func (r *DemoRepository) FechaSimulada(ctx context.Context) (time.Time, error) {
	var value string
	err := r.pool.QueryRow(ctx, `SELECT valor FROM demo WHERE clave=$1`, demoClockKey).Scan(&value)
	if err != nil {
		if err == pgx.ErrNoRows {
			return time.Time{}, nil
		}
		return time.Time{}, repositoryError("demo", demoClockKey, err)
	}
	if value == "" {
		return time.Time{}, nil
	}
	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		parsed, err = time.Parse("2006-01-02", value)
	}
	if err != nil {
		return time.Time{}, repositoryError("demo", demoClockKey, err)
	}
	return parsed.UTC(), nil
}

func (r *DemoRepository) GuardarFechaSimulada(ctx context.Context, value time.Time) error {
	if value.IsZero() {
		return usecase.ErrValidacion
	}
	_, err := r.pool.Exec(ctx, `INSERT INTO demo (clave,valor) VALUES ($1,$2) ON CONFLICT (clave) DO UPDATE SET valor=EXCLUDED.valor`, demoClockKey, value.UTC().Format(time.RFC3339Nano))
	return repositoryError("demo", demoClockKey, err)
}

func (r *DemoRepository) Sembrar(ctx context.Context, leads []domain.Lead) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return repositoryError("demo seed", "", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	for i := range leads {
		if err := insertDemoLead(ctx, tx, leads[i]); err != nil {
			return err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return repositoryError("demo seed", "commit", err)
	}
	return nil
}

func (r *DemoRepository) Reiniciar(ctx context.Context) (time.Time, error) {
	return r.reiniciar(ctx, nil)
}

func (r *DemoRepository) ReiniciarConSeed(ctx context.Context, leads []domain.Lead) (time.Time, error) {
	return r.reiniciar(ctx, leads)
}

func (r *DemoRepository) reiniciar(ctx context.Context, leads []domain.Lead) (time.Time, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return time.Time{}, repositoryError("demo reset", "", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	for _, table := range []string{"fichas", "hitos", "planes", "mensajes", "leads"} {
		if _, err := tx.Exec(ctx, "DELETE FROM "+table); err != nil {
			return time.Time{}, repositoryError("demo reset", table, err)
		}
	}
	if _, err := tx.Exec(ctx, `INSERT INTO demo (clave,valor) VALUES ($1,$2) ON CONFLICT (clave) DO UPDATE SET valor=EXCLUDED.valor`, demoClockKey, approvedDemoDate); err != nil {
		return time.Time{}, repositoryError("demo reset", demoClockKey, err)
	}
	for i := range leads {
		if err := insertDemoLead(ctx, tx, leads[i]); err != nil {
			return time.Time{}, err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return time.Time{}, repositoryError("demo reset", "commit", err)
	}
	date, err := time.Parse(time.RFC3339, approvedDemoDate)
	if err != nil {
		return time.Time{}, err
	}
	return date, nil
}

func insertDemoLead(ctx context.Context, tx pgx.Tx, lead domain.Lead) error {
	perfil, err := encodeJSONB(lead.Perfil)
	if err != nil {
		return err
	}
	capacidad, err := encodeJSONB(lead.Capacidad)
	if err != nil {
		return err
	}
	intencion, err := encodeJSONB(lead.Intencion)
	if err != nil {
		return err
	}
	version := lead.Version
	if version == 0 {
		version = 1
	}
	if _, err := tx.Exec(ctx, demoLeadUpsertSQL, lead.LeadID, lead.Nombre, lead.Telefono, lead.Cedula, lead.Fuente, lead.Estado, lead.Ruta, lead.Afiliado, lead.Prioridad, lead.ConsumeCupo10, perfil, capacidad, intencion, version, lead.CreadoEn, lead.ActualizadoEn); err != nil {
		return repositoryError("demo lead", lead.LeadID, err)
	}

	// Seed initial Vivi greeting message for seed leads
	var greeting string
	if lead.Afiliado {
		greeting = "¡Hola " + lead.Nombre + "! 👋 Como afiliada a Colsubsidio, el motor identifica un subsidio aplicable de hasta $52,5M. Consulta la política de tratamiento de datos: https://www.colsubsidio.com/politica-tratamiento-datos. Al continuar autorizas el tratamiento de tus datos. ¿Qué sueñas con comprar este año?"
	} else {
		greeting = "¡Hola " + lead.Nombre + "! 👋 Estoy aquí para orientarte en tu camino hacia la vivienda. Consulta la política de tratamiento de datos: https://www.colsubsidio.com/politica-tratamiento-datos. Al continuar autorizas el tratamiento de tus datos. ¿Cómo está tu situación laboral para acompañarte mejor?"
	}
	msgID := "msg-seed-" + lead.LeadID
	emptyAdj, _ := encodeJSONB(map[string]any{})
	_, _ = tx.Exec(ctx, `INSERT INTO mensajes (mensaje_id,lead_id,autor,tipo_contenido,texto,adjunto,creado_en) VALUES ($1,$2,$3,$4,$5,$6,$7) ON CONFLICT (mensaje_id) DO NOTHING`, msgID, lead.LeadID, domain.AutorMensajeVivi, domain.TipoContenidoTexto, greeting, emptyAdj, lead.CreadoEn)

	// Seed initial commercial Ficha for seed leads
	fichaID := "ficha-seed-" + lead.LeadID
	recs := []domain.Recomendacion{
		{ProyectoID: "mongui", Nombre: "Monguí", Zona: "Ciudadela Maiporé - Soacha", PrecioDesde: 156470000, Razon: "Tu presupuesto cubre el 100% de la cuota inicial", Vecinos: 622, TasaDesistimiento: 0.12, BrochureURL: "https://heyzine.com/flip-book/866af8f6a6.html", Recorrido360URL: "https://storage.net-fs.com/hosting/7532170/19/"},
	}
	if !lead.Afiliado {
		recs[0] = domain.Recomendacion{ProyectoID: "versalles", Nombre: "Versalles", Zona: "Ciudadela Maiporé - Soacha", PrecioDesde: 195200000, Razon: "Certificación EDGE, ahorro en servicios", Vecinos: 174, TasaDesistimiento: 0.15, BrochureURL: "https://heyzine.com/flip-book/be784b0d5c.html", Recorrido360URL: "https://shape.com.co/360/COLSUBSIDIO-Versalles_APTOA"}
	}
	var bandaAdv *string
	if !lead.Afiliado {
		adv := "No afiliado a Colsubsidio — consume cupo del 10% regulatorio"
		bandaAdv = &adv
	}

	capacidadObj := domain.Capacidad{
		PresupuestoMax: 180000000, CreditoMax: 117500000, SubsidioAplicable: 52500000,
		RecursosPropios: 10000000, Ratio: 0.29, Confianza: 0.94,
		Desglose: []domain.ItemDesglose{
			{Concepto: "Subsidio Mi Casa Ya / Caja", Monto: 52500000, Regla: "Afiliado Cat A", Fuente: domain.FuenteCampoVerificadoBase},
			{Concepto: "Preaprobado Bancolombia", Monto: 117500000, Regla: "Capacidad de endeudamiento 30%", Fuente: domain.FuenteCampoInferido},
			{Concepto: "Ahorro Declarado", Monto: 10000000, Regla: "Declarado en chat", Fuente: domain.FuenteCampoDeclarado},
		},
	}
	if !lead.Afiliado {
		capacidadObj.SubsidioAplicable = 0
		capacidadObj.CreditoMax = 170000000
		capacidadObj.Desglose = []domain.ItemDesglose{
			{Concepto: "Crédito Preaprobado", Monto: 170000000, Regla: "Capacidad de endeudamiento 30%", Fuente: domain.FuenteCampoInferido},
			{Concepto: "Ahorro Declarado", Monto: 10000000, Regla: "Declarado en chat", Fuente: domain.FuenteCampoDeclarado},
		}
	}

	fichaObj := domain.Ficha{
		FichaID: fichaID, LeadID: lead.LeadID, GeneradaEn: lead.CreadoEn, ConfianzaPerfil: 0.94,
		BandaAdvertencia: bandaAdv,
		Identificacion:   domain.Identificacion{Nombre: lead.Nombre, Afiliada: lead.Afiliado, Categoria: "A", Telefono: lead.Telefono},
		Capacidad:        capacidadObj, Perfil: lead.Perfil,
		Intencion:        domain.Intencion{Nivel: domain.NivelAlta, Confianza: domain.NivelAlta, Senales: []string{"Busca comprar antes de 6 meses"}},
		Recomendaciones:  recs,
		Beneficios:       []string{"Subsidio de vivienda Colsubsidio hasta $52,5M", "Tasa preferencial crédito hipotecario"},
		ArgumentosVenta:  []string{"Cuota estimada mensual ($1,4M) es adecuada para tu nivel de ingresos"},
		AlertaDesistimiento: domain.AlertaDesistimiento{Activa: false, TasaVecinos: 0.12},
		ConsumeCupo10:    lead.ConsumeCupo10,
	}

	fichaContent, errFicha := encodeJSONB(fichaObj)
	if errFicha == nil {
		_, _ = tx.Exec(ctx, `INSERT INTO fichas (ficha_id,lead_id,contenido,generada_en) VALUES ($1,$2,$3,$4) ON CONFLICT (lead_id) DO NOTHING`, fichaID, lead.LeadID, fichaContent, lead.CreadoEn)
	}

	return nil
}
