package usecase

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
)

var (
	ErrValidacion          = errors.New("VALIDACION")
	ErrAudioInvalido       = errors.New("AUDIO_INVALIDO")
	ErrLeadNoPerfilando    = errors.New("TRANSICION_INVALIDA")
	ErrLimiteTurnos        = fmt.Errorf("%w: LIMITE_TURNOS", ErrValidacion)
	ErrSalidaTurnoInvalida = errors.New("SALIDA_TURNO_INVALIDA")
)

type EntradaMensaje struct {
	LeadID string
	Tipo   domain.TipoMensaje
	Texto  string
	Audio  *Audio
}

type ProcesarMensaje struct {
	Leads LeadRepository
	LLM   LLMProvider
	IDs   GeneradorID
	Bus   BusEventos
	Reloj Reloj
}

func (uc *ProcesarMensaje) Ejecutar(ctx context.Context, entrada EntradaMensaje) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := validarEntrada(entrada); err != nil {
		return err
	}
	lead, err := uc.Leads.PorID(ctx, entrada.LeadID)
	if err != nil {
		return err
	}
	if lead.Estado != domain.EstadoLeadPerfilando {
		return fmt.Errorf("%w: lead %q en estado %s", ErrLeadNoPerfilando, lead.LeadID, lead.Estado)
	}
	conversation, err := uc.Leads.Conversacion(ctx, lead.LeadID)
	if err != nil {
		return err
	}
	limit := 6
	if lead.Afiliado {
		limit = 4
	}
	turns := 0
	for _, message := range conversation {
		if message.Autor == domain.AutorMensajeLead {
			turns++
		}
	}
	if turns >= limit {
		return ErrLimiteTurnos
	}

	current := domain.Mensaje{LeadID: lead.LeadID, Texto: entrada.Texto, TipoContenido: domain.TipoContenidoTexto}
	history := historialReciente(conversation, current)
	turn := EntradaTurno{
		LeadID:            lead.LeadID,
		Nombre:            lead.Nombre,
		MensajeUsuario:    entrada.Texto,
		Perfil:            copiarPerfil(lead.Perfil),
		Capacidad:         copiarCapacidad(lead.Capacidad),
		NumerosDelMotor:   numerosDelMotor(lead.Capacidad),
		HistorialReciente: history,
		EsAfiliado:        lead.Afiliado,
	}
	var salida SalidaTurno
	if entrada.Tipo == domain.TipoMensajeAudio {
		_, decodeErr := base64.StdEncoding.Strict().DecodeString(entrada.Audio.Base64)
		if decodeErr != nil {
			return fmt.Errorf("%w: %v", ErrAudioInvalido, decodeErr)
		}
		audio := *entrada.Audio
		salida, err = uc.LLM.ProcesarAudio(ctx, audio, turn)
	} else {
		salida, err = uc.LLM.GenerarTurno(ctx, turn)
	}
	if err != nil {
		return err
	}
	if err := validarSalida(salida); err != nil {
		return err
	}
	aplicarCamposBasicos(lead, salida.CamposExtraidos)
	lead.Intencion = &salida.Intencion
	lead.ActualizadoEn = uc.Reloj.Ahora()

	inbound := &domain.Mensaje{
		MensajeID: uc.IDs.Nuevo(), LeadID: lead.LeadID, Autor: domain.AutorMensajeLead,
		TipoContenido: domain.TipoContenidoTexto, Texto: entrada.Texto, CreadoEn: lead.ActualizadoEn,
	}
	if entrada.Tipo == domain.TipoMensajeAudio {
		inbound.Adjunto = map[string]any{"audio_original": true}
	}
	if err := uc.Leads.AgregarMensaje(ctx, inbound); err != nil {
		return fmt.Errorf("guardar mensaje entrante: %w", err)
	}
	if err := uc.Leads.Guardar(ctx, lead); err != nil {
		return fmt.Errorf("guardar lead: %w", err)
	}
	response := &domain.Mensaje{
		MensajeID: uc.IDs.Nuevo(), LeadID: lead.LeadID, Autor: domain.AutorMensajeVivi,
		TipoContenido: domain.TipoContenidoTexto, Texto: salida.Respuesta, CreadoEn: uc.Reloj.Ahora(),
	}
	if err := uc.Leads.AgregarMensaje(ctx, response); err != nil {
		return fmt.Errorf("guardar respuesta: %w", err)
	}
	return nil
}

func validarEntrada(entrada EntradaMensaje) error {
	if strings.TrimSpace(entrada.LeadID) == "" {
		return fmt.Errorf("%w: lead_id vacio", ErrValidacion)
	}
	switch entrada.Tipo {
	case domain.TipoMensajeTexto:
		if !utf8.ValidString(entrada.Texto) || strings.TrimSpace(entrada.Texto) == "" || utf8.RuneCountInString(entrada.Texto) > 2000 {
			return ErrValidacion
		}
	case domain.TipoMensajeAudio:
		if entrada.Audio == nil || entrada.Audio.DuracionS < 1 || entrada.Audio.DuracionS > 60 {
			return ErrAudioInvalido
		}
		if entrada.Audio.MIME != "audio/webm" && entrada.Audio.MIME != "audio/ogg" && entrada.Audio.MIME != "audio/mpeg" {
			return ErrAudioInvalido
		}
		decoded, err := base64.StdEncoding.Strict().DecodeString(entrada.Audio.Base64)
		if err != nil || len(decoded) == 0 || len(decoded) > 2*1024*1024 {
			return ErrAudioInvalido
		}
	default:
		return fmt.Errorf("%w: tipo de mensaje invalido", ErrValidacion)
	}
	return nil
}

func historialReciente(conversacion []domain.Mensaje, actual domain.Mensaje) []domain.Mensaje {
	start := 0
	if len(conversacion) > 6 {
		start = len(conversacion) - 6
	}
	history := make([]domain.Mensaje, 0, len(conversacion)-start+1)
	history = append(history, conversacion[start:]...)
	return append(history, actual)
}

func numerosDelMotor(capacidad *domain.Capacidad) map[string]int64 {
	numeros := map[string]int64{"presupuesto_max": 0, "credito_max": 0, "subsidio_aplicable": 0, "recursos_propios": 0}
	if capacidad == nil {
		return numeros
	}
	numeros["presupuesto_max"] = capacidad.PresupuestoMax
	numeros["credito_max"] = capacidad.CreditoMax
	numeros["subsidio_aplicable"] = capacidad.SubsidioAplicable
	numeros["recursos_propios"] = capacidad.RecursosPropios
	return numeros
}

func validarSalida(salida SalidaTurno) error {
	seen := make(map[string]bool)
	for _, campo := range salida.CamposExtraidos {
		if !domain.CamposReconocidos[campo.Campo] {
			continue
		}
		if seen[campo.Campo] || campo.Confianza < 0 || campo.Confianza > 1 {
			return ErrSalidaTurnoInvalida
		}
		switch campo.Fuente {
		case domain.FuenteCampoDeclarado, domain.FuenteCampoInferido, domain.FuenteCampoVerificadoBase:
		default:
			return ErrSalidaTurnoInvalida
		}
		seen[campo.Campo] = true
	}
	return nil
}

func aplicarCamposBasicos(lead *domain.Lead, campos []CampoExtraido) {
	if lead.Perfil == nil {
		lead.Perfil = domain.Perfil{}
	}
	updated := copiarPerfil(lead.Perfil)
	for _, campo := range campos {
		if !domain.CamposReconocidos[campo.Campo] {
			continue
		}
		if existente, ok := updated[campo.Campo]; ok && existente.Fuente == domain.FuenteCampoVerificadoBase {
			continue
		}
		fuente := campo.Fuente
		if fuente == domain.FuenteCampoVerificadoBase {
			fuente = domain.FuenteCampoDeclarado
		}
		updated[campo.Campo] = domain.CampoPerfil{Valor: campo.Valor, Fuente: fuente, Confianza: campo.Confianza, RequiereConfirmacion: campo.RequiereConfirmacion}
	}
	lead.Perfil = updated
}

func copiarPerfil(perfil domain.Perfil) domain.Perfil {
	copia := make(domain.Perfil, len(perfil))
	for clave, campo := range perfil {
		copia[clave] = campo
	}
	return copia
}

func copiarCapacidad(capacidad *domain.Capacidad) *domain.Capacidad {
	if capacidad == nil {
		return nil
	}
	copia := *capacidad
	copia.Desglose = append([]domain.ItemDesglose(nil), capacidad.Desglose...)
	return &copia
}
