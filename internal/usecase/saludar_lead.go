package usecase

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
)

const URLPolitica = "https://www.colsubsidio.com/politica-tratamiento-datos"

const despedidaConsentimiento = "Gracias por contármelo. Respeto tu decisión y no continuaré con el proceso. Si deseas recibir orientación en otro momento, puedes escribirnos."

var (
	greetingAmountPattern    = regexp.MustCompile(`(?i)(\$\s*\d|\b\d[\d.,]*\s*(m|millones?|cop|pesos?))`)
	greetingForbiddenPattern = regexp.MustCompile(`(?i)\b(ingreso|salario|sueldo|personas?\s+(a\s+)?cargo|hogar|familia|cu[aá]ntas?\s+personas?)\b`)
)

// SaludarLead persists the first Vivi message for a new lead after validating
// provider output against a deterministic, audience-specific copy contract.
type SaludarLead struct {
	Leads LeadRepository
	LLM   LLMProvider
	IDs   GeneradorID
	Reloj Reloj
}

func (uc *SaludarLead) Ejecutar(ctx context.Context, event Evento) error {
	if uc == nil || uc.Leads == nil || uc.IDs == nil || uc.Reloj == nil {
		return fmt.Errorf("%w: dependencias del saludo requeridas", ErrValidacion)
	}
	if strings.TrimSpace(event.LeadID) == "" {
		return fmt.Errorf("%w: lead_id vacio", ErrValidacion)
	}
	lead, err := uc.Leads.PorID(ctx, event.LeadID)
	if err != nil {
		return err
	}
	if lead == nil {
		return ErrNoEncontrado
	}

	texto, amount := saludoDeterminista(lead)
	if uc.LLM != nil {
		turn := EntradaTurno{
			LeadID: lead.LeadID, Nombre: lead.Nombre, Perfil: copiarPerfil(lead.Perfil),
			Capacidad: copiarCapacidad(lead.Capacidad), NumerosDelMotor: numerosDelMotor(lead.Capacidad),
			EsAfiliado: lead.Afiliado, MensajeUsuario: "Redacta el saludo inicial y solicita consentimiento.",
		}
		if salida, providerErr := uc.LLM.GenerarTurno(ctx, turn); providerErr == nil &&
			ValidarSaludo(salida.Respuesta, lead.Nombre, lead.Afiliado, amount) {
			texto = salida.Respuesta
		}
	}
	return uc.Leads.AgregarMensaje(ctx, &domain.Mensaje{
		MensajeID: uc.IDs.Nuevo(), LeadID: lead.LeadID, Autor: domain.AutorMensajeVivi,
		TipoContenido: domain.TipoContenidoTexto, Texto: texto, CreadoEn: uc.Reloj.Ahora().UTC(),
	})
}

func saludoDeterminista(lead *domain.Lead) (string, string) {
	amount := ""
	if lead.Afiliado && lead.Capacidad != nil && lead.Capacidad.SubsidioAplicable > 0 {
		amount = formatoSubsidio(lead.Capacidad.SubsidioAplicable)
		return fmt.Sprintf("¡Hola %s! 👋 Como afiliada a Colsubsidio, el motor identifica un subsidio aplicable de hasta %s. Consulta la política de tratamiento de datos: %s. Al continuar autorizas el tratamiento de tus datos. ¿Qué sueñas con comprar este año?", lead.Nombre, amount, URLPolitica), amount
	}
	return fmt.Sprintf("¡Hola %s! 👋 Estoy aquí para orientarte en tu camino hacia la vivienda. Consulta la política de tratamiento de datos: %s. Al continuar autorizas el tratamiento de tus datos. ¿Cómo está tu situación laboral para acompañarte mejor?", lead.Nombre, URLPolitica), amount
}

func formatoSubsidio(value int64) string {
	millions, remainder := value/1_000_000, value%1_000_000
	return fmt.Sprintf("$%d,%dM", millions, remainder/100_000)
}

func ValidarSaludo(text, name string, affiliate bool, amount string) bool {
	text = strings.TrimSpace(text)
	if text == "" || !strings.Contains(text, URLPolitica) ||
		!strings.Contains(text, "Al continuar autorizas el tratamiento de tus datos") ||
		strings.Count(text, "?") != 1 || (name != "" && !strings.Contains(text, name)) {
		return false
	}
	if greetingForbiddenPattern.MatchString(text) {
		return false
	}
	if affiliate && amount != "" {
		if !strings.Contains(text, amount) || !strings.Contains(strings.ToLower(text), "sueñ") {
			return false
		}
		withoutAmount := strings.Replace(text, amount, "", 1)
		return !greetingAmountPattern.MatchString(withoutAmount)
	}
	lower := strings.ToLower(text)
	hasJobMarker := strings.Contains(lower, "emple") || strings.Contains(lower, "trabaj") || strings.Contains(lower, "labor")
	return hasJobMarker && !greetingAmountPattern.MatchString(text)
}

func (uc *SaludarLead) RechazarConsentimiento(ctx context.Context, lead *domain.Lead, entrada EntradaMensaje) error {
	if uc == nil || uc.Leads == nil || uc.IDs == nil || uc.Reloj == nil || lead == nil {
		return fmt.Errorf("%w: dependencias del rechazo requeridas", ErrValidacion)
	}
	if lead.Estado != domain.EstadoLeadPerfilando {
		return fmt.Errorf("%w: lead %q en estado %s", ErrLeadNoPerfilando, lead.LeadID, lead.Estado)
	}
	received := entrada.RecibidoEn
	if received.IsZero() {
		received = uc.Reloj.Ahora()
	}
	messageID := entrada.MensajeID
	if messageID == "" {
		messageID = uc.IDs.Nuevo()
	}
	inbound := &domain.Mensaje{
		MensajeID: messageID, LeadID: lead.LeadID, Autor: domain.AutorMensajeLead,
		TipoContenido: domain.TipoContenidoTexto, Texto: entrada.Texto, CreadoEn: received.UTC(),
	}
	if entrada.Tipo == domain.TipoMensajeAudio && entrada.Audio != nil {
		inbound.Adjunto = map[string]any{"audio_original": true}
	}
	if err := uc.Leads.AgregarMensaje(ctx, inbound); err != nil {
		return fmt.Errorf("guardar mensaje de consentimiento: %w", err)
	}
	if err := lead.Transicionar(domain.EstadoLeadCalificado); err != nil {
		return fmt.Errorf("transicionar lead: %w", err)
	}
	if err := lead.Transicionar(domain.EstadoLeadDespedido); err != nil {
		return fmt.Errorf("transicionar lead: %w", err)
	}
	lead.Ruta = domain.RutaDespedida
	lead.ActualizadoEn = uc.Reloj.Ahora()
	if err := uc.Leads.Guardar(ctx, lead); err != nil {
		return fmt.Errorf("guardar lead despedido: %w", err)
	}
	if err := uc.Leads.AgregarMensaje(ctx, &domain.Mensaje{
		MensajeID: uc.IDs.Nuevo(), LeadID: lead.LeadID, Autor: domain.AutorMensajeVivi,
		TipoContenido: domain.TipoContenidoTexto, Texto: despedidaConsentimiento, CreadoEn: uc.Reloj.Ahora().UTC(),
	}); err != nil {
		return fmt.Errorf("guardar despedida: %w", err)
	}
	return nil
}
