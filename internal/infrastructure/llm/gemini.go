package llm

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
)

const (
	DefaultGeminiBaseURL = "https://generativelanguage.googleapis.com"
	DefaultGeminiModel   = "gemini-2.0-flash"
)

type GeminiOption func(*GeminiProvider)

func WithGeminiHTTPDoer(client HTTPDoer) GeminiOption {
	return func(p *GeminiProvider) { p.client = client }
}
func WithGeminiBaseURL(baseURL string) GeminiOption {
	return func(p *GeminiProvider) { p.baseURL = strings.TrimRight(baseURL, "/") }
}
func WithGeminiModel(model string) GeminiOption { return func(p *GeminiProvider) { p.model = model } }

type GeminiProvider struct {
	apiKey  string
	baseURL string
	model   string
	client  HTTPDoer
}

func NewGeminiProvider(apiKey string, options ...GeminiOption) *GeminiProvider {
	p := &GeminiProvider{apiKey: apiKey, baseURL: DefaultGeminiBaseURL, model: DefaultGeminiModel, client: http.DefaultClient}
	for _, option := range options {
		option(p)
	}
	return p
}

func (p *GeminiProvider) Nombre() string { return "gemini" }

func (p *GeminiProvider) GenerarTurno(ctx context.Context, input usecase.EntradaTurno) (usecase.SalidaTurno, error) {
	return p.turn(ctx, input, nil)
}

func (p *GeminiProvider) ProcesarAudio(ctx context.Context, audio usecase.Audio, input usecase.EntradaTurno) (usecase.SalidaTurno, error) {
	decoded, err := base64.StdEncoding.DecodeString(audio.Base64)
	if err != nil {
		return invalid(KindConfig)
	}
	return p.turn(ctx, input, &geminiAudio{MIME: audio.MIME, Data: decoded})
}

type geminiAudio struct {
	MIME string `json:"mimeType"`
	Data []byte `json:"data"`
}
type geminiPart struct {
	Text       string       `json:"text,omitempty"`
	InlineData *geminiAudio `json:"inlineData,omitempty"`
}
type geminiContent struct {
	Role  string       `json:"role,omitempty"`
	Parts []geminiPart `json:"parts"`
}
type geminiRequest struct {
	SystemInstruction geminiContent   `json:"systemInstruction"`
	Contents          []geminiContent `json:"contents"`
	GenerationConfig  struct {
		ResponseMIMEType string `json:"responseMimeType"`
	} `json:"generationConfig"`
}
type geminiResponse struct {
	Candidates []struct {
		Content geminiContent `json:"content"`
	} `json:"candidates"`
}

func (p *GeminiProvider) turn(ctx context.Context, input usecase.EntradaTurno, audio *geminiAudio) (usecase.SalidaTurno, error) {
	if p == nil || p.apiKey == "" || p.client == nil || p.model == "" || p.baseURL == "" {
		return invalid(KindConfig)
	}
	prompt := BuildPrompt(input)
	body := geminiRequest{SystemInstruction: geminiContent{Parts: []geminiPart{{Text: prompt.Instruction}}}, Contents: []geminiContent{{Role: "user", Parts: []geminiPart{{Text: prompt.UserData}}}}, GenerationConfig: struct {
		ResponseMIMEType string `json:"responseMimeType"`
	}{ResponseMIMEType: "application/json"}}
	if audio != nil {
		body.Contents[0].Parts = append(body.Contents[0].Parts, geminiPart{InlineData: audio})
	}
	encoded, err := json.Marshal(body)
	if err != nil {
		return invalid(KindConfig)
	}
	return p.call(ctx, encoded, input.NumerosDelMotor)
}

func (p *GeminiProvider) call(ctx context.Context, body []byte, motor map[string]int64) (usecase.SalidaTurno, error) {
	endpoint := strings.TrimRight(p.baseURL, "/") + "/v1beta/models/" + url.PathEscape(p.model) + ":generateContent"
	data, err := DoRequest(ctx, p.client, http.MethodPost, endpoint, body, http.Header{"Content-Type": {"application/json"}, "X-Goog-Api-Key": {p.apiKey}})
	if err != nil {
		return usecase.SalidaTurno{}, err
	}
	text, err := extractGemini(data)
	if err != nil {
		return usecase.SalidaTurno{}, err
	}
	out, err := ParseSalida([]byte(text), motor)
	if ErrorKindOf(err) != KindMalformed {
		return out, err
	}
	data, err = DoRequest(ctx, p.client, http.MethodPost, endpoint, body, http.Header{"Content-Type": {"application/json"}, "X-Goog-Api-Key": {p.apiKey}})
	if err != nil {
		return usecase.SalidaTurno{}, err
	}
	text, err = extractGemini(data)
	if err != nil {
		return usecase.SalidaTurno{}, err
	}
	return ParseSalida([]byte(text), motor)
}

func extractGemini(data []byte) (string, error) {
	var response geminiResponse
	if json.Unmarshal(data, &response) != nil || len(response.Candidates) == 0 {
		return "", providerError(KindMalformed, nil)
	}
	var text strings.Builder
	for _, part := range response.Candidates[0].Content.Parts {
		text.WriteString(part.Text)
	}
	if text.Len() == 0 {
		return "", providerError(KindMalformed, nil)
	}
	return text.String(), nil
}
