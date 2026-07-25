package llm

import (
	"context"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
	"regexp"
	"strings"
)

const (
	genericGuardrailReply = "Puedo ayudarte únicamente con vivienda, subsidios y compra de vivienda."
	privacyGuardrailReply = "Por privacidad, no puedo consultar datos de otras personas. Puedo ayudarte con tu proceso de vivienda."
)

type Guardarrailes struct{ next usecase.LLMProvider }

func ConGuardarrailes(next usecase.LLMProvider) *Guardarrailes { return &Guardarrailes{next: next} }

var (
	jailbreakPattern = regexp.MustCompile(`(?i)\b(?:ignora (?:tus )?instrucciones|olvida tus reglas|jailbreak|modo desarrollador)\b`)
	promptPattern    = regexp.MustCompile(`(?i)\b(?:system prompt|prompt del sistema|instrucciones internas|configuraci[oó]n interna|archivo interno|repite tus instrucciones|lista los archivos skill[.]md)\b`)
	rolePattern      = regexp.MustCompile(`(?i)\b(?:eres ahora|ahora eres|act[uú]a como|finge ser|pretend[ae] ser|modo desarrollador)\b`)
	privacyPattern   = regexp.MustCompile(`(?i)\b(?:c[eé]dula\s+\d{6,12}|c[eé]dula\s+de\s+(?:otra|otro|tercera?)\s+persona|(?:otro|otra|tercero)\s+lead)\b`)
	outsidePattern   = regexp.MustCompile(`(?i)\b(?:escr[ií]beme\s+c[oó]digo\s+en\s+python|hazme un ensayo sobre la revoluci[oó]n francesa|cu[aá]l es la capital de francia|por qui[eé]n votar|explosivo|hackear|contrase[nñ]a)\b`)
	leakPattern      = regexp.MustCompile(`(?i)(system prompt|prompt del sistema|instrucciones internas|skill|config(?:uraci[oó]n)? interna|archivo interno|internal[_ -]?file|\.env|go\.mod|\.kiro|credential|credencial|secreto)`)
	leadPattern      = regexp.MustCompile(`(?i)\blead[_ -]?id\s*["']?\s*[:=]\s*["']?([a-z0-9_-]+)`)
)

func (g *Guardarrailes) Nombre() string {
	if g == nil || g.next == nil {
		return "unconfigured"
	}
	return g.next.Nombre()
}
func (g *Guardarrailes) GenerarTurno(ctx context.Context, in usecase.EntradaTurno) (usecase.SalidaTurno, error) {
	return g.run(in, func() (usecase.SalidaTurno, error) { return g.next.GenerarTurno(ctx, in) })
}
func (g *Guardarrailes) ProcesarAudio(ctx context.Context, audio usecase.Audio, in usecase.EntradaTurno) (usecase.SalidaTurno, error) {
	return g.run(in, func() (usecase.SalidaTurno, error) { return g.next.ProcesarAudio(ctx, audio, in) })
}
func (g *Guardarrailes) run(in usecase.EntradaTurno, call func() (usecase.SalidaTurno, error)) (usecase.SalidaTurno, error) {
	if blocked, privacy := blockedInput(in.MensajeUsuario); blocked {
		return safeTurn(privacy), nil
	}
	if g == nil || g.next == nil {
		return usecase.SalidaTurno{}, providerError(KindConfig, nil)
	}
	out, err := call()
	return g.validate(out, err, in)
}
func (g *Guardarrailes) validate(out usecase.SalidaTurno, err error, in usecase.EntradaTurno) (usecase.SalidaTurno, error) {
	if err != nil {
		return out, err
	}
	if leakPattern.MatchString(out.Respuesta) || foreignLead(out.Respuesta, in.LeadID) || !validResponse(out.Respuesta, in.NumerosDelMotor) {
		return safeTurn(false), nil
	}
	return out, nil
}
func (g *Guardarrailes) CircuitBreakerState() BreakerState {
	if owner, ok := g.next.(breakerStateOwner); ok {
		return owner.CircuitBreakerState()
	}
	return BreakerClosed
}
func blockedInput(text string) (bool, bool) {
	if privacyPattern.MatchString(text) {
		return true, true
	}
	return jailbreakPattern.MatchString(text) || promptPattern.MatchString(text) || rolePattern.MatchString(text) || outsidePattern.MatchString(text), false
}
func safeTurn(privacy bool) usecase.SalidaTurno {
	reply := genericGuardrailReply
	if privacy {
		reply = privacyGuardrailReply
	}
	return usecase.SalidaTurno{CamposExtraidos: []usecase.CampoExtraido{}, Intencion: domain.Intencion{Nivel: domain.NivelBaja, Confianza: domain.NivelBaja, Senales: []string{}}, Respuesta: reply, Accion: usecase.AccionFueraDeDominio}
}
func foreignLead(response, own string) bool {
	for _, match := range leadPattern.FindAllStringSubmatch(response, -1) {
		if len(match) < 2 || !strings.EqualFold(match[1], own) {
			return true
		}
	}
	return false
}

type breakerStateOwner interface{ CircuitBreakerState() BreakerState }
