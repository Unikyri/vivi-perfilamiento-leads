package llm

type ErrorKind string

const (
	KindConfig        ErrorKind = "config"
	KindCanceled      ErrorKind = "canceled"
	KindTimeout       ErrorKind = "timeout"
	KindRateLimit     ErrorKind = "rate_limit"
	KindUnavailable   ErrorKind = "unavailable"
	KindHTTP          ErrorKind = "http"
	KindMalformed     ErrorKind = "malformed"
	KindInvalidOutput ErrorKind = "invalid_output"
	KindExhausted     ErrorKind = "exhausted"
	KindGuardrail     ErrorKind = "guardrail"
)

type ProviderError struct{ Kind ErrorKind }

const KindCapability, KindCircuitOpen ErrorKind = "capability", "circuit_open"

func (e *ProviderError) Error() string                      { return "llm provider: " + string(e.Kind) }
func providerError(kind ErrorKind, _ error, _ ...int) error { return &ProviderError{Kind: kind} }
func ErrorKindOf(err error) ErrorKind {
	if e, ok := err.(*ProviderError); ok {
		return e.Kind
	}
	return ""
}
func statusError(status int) error {
	if status == 429 {
		return providerError(KindRateLimit, nil)
	}
	if status == 408 || status == 500 || status == 502 || status == 503 || status == 504 {
		return providerError(KindUnavailable, nil)
	}
	return providerError(KindHTTP, nil)
}

func Eligible(err error) bool {
	switch ErrorKindOf(err) {
	case KindMalformed, KindTimeout, KindRateLimit, KindUnavailable, KindCircuitOpen:
		return true
	default:
		return false
	}
}
