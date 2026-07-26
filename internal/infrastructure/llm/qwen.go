package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
)

const (
	DefaultQwenBaseURL = "https://dashscope-intl.aliyuncs.com/compatible-mode/v1"
	DefaultQwenModel   = "qwen3.7-plus"
)

type QwenOption func(*QwenProvider)

func WithQwenHTTPDoer(client HTTPDoer) QwenOption { return func(p *QwenProvider) { p.client = client } }
func WithQwenModel(model string) QwenOption       { return func(p *QwenProvider) { p.model = model } }

// QwenProvider speaks the OpenAI-compatible DashScope chat-completions protocol.
type QwenProvider struct {
	apiKey  string
	baseURL string
	model   string
	client  HTTPDoer
}

func NewQwenProvider(apiKey, baseURL string, options ...QwenOption) *QwenProvider {
	if baseURL == "" {
		baseURL = DefaultQwenBaseURL
	}
	p := &QwenProvider{apiKey: apiKey, baseURL: strings.TrimRight(baseURL, "/"), model: DefaultQwenModel, client: http.DefaultClient}
	for _, option := range options {
		option(p)
	}
	return p
}

func (p *QwenProvider) Nombre() string { return "qwen" }

func (p *QwenProvider) GenerarTurno(ctx context.Context, input usecase.EntradaTurno) (usecase.SalidaTurno, error) {
	if p == nil || p.apiKey == "" || p.client == nil || p.baseURL == "" || p.model == "" {
		return invalid(KindConfig)
	}
	prompt := BuildPrompt(input)
	body := qwenRequest{Model: p.model, Messages: []qwenMessage{{Role: "system", Content: prompt.Instruction}, {Role: "user", Content: prompt.UserData}}, ResponseFormat: qwenResponseFormat{Type: "json_object"}}
	encoded, err := json.Marshal(body)
	if err != nil {
		return invalid(KindConfig)
	}
	return p.call(ctx, encoded, input.NumerosDelMotor)
}

func (p *QwenProvider) ProcesarAudio(context.Context, usecase.Audio, usecase.EntradaTurno) (usecase.SalidaTurno, error) {
	return invalid(KindCapability)
}

type qwenMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type qwenResponseFormat struct {
	Type string `json:"type"`
}
type qwenRequest struct {
	Model          string             `json:"model"`
	Messages       []qwenMessage      `json:"messages"`
	ResponseFormat qwenResponseFormat `json:"response_format"`
}
type qwenResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (p *QwenProvider) call(ctx context.Context, body []byte, motor map[string]int64) (usecase.SalidaTurno, error) {
	endpoint := strings.TrimRight(p.baseURL, "/") + "/chat/completions"
	headers := http.Header{"Content-Type": {"application/json"}, "Authorization": {"Bearer " + p.apiKey}}
	data, err := DoRequest(ctx, p.client, http.MethodPost, endpoint, body, headers)
	if err != nil {
		return usecase.SalidaTurno{}, err
	}
	text, err := extractQwen(data)
	if err != nil {
		return usecase.SalidaTurno{}, err
	}
	out, err := ParseSalida([]byte(text), motor)
	if ErrorKindOf(err) != KindMalformed {
		return out, err
	}
	data, err = DoRequest(ctx, p.client, http.MethodPost, endpoint, body, headers)
	if err != nil {
		return usecase.SalidaTurno{}, err
	}
	text, err = extractQwen(data)
	if err != nil {
		return usecase.SalidaTurno{}, err
	}
	return ParseSalida([]byte(text), motor)
}

func extractQwen(data []byte) (string, error) {
	var response qwenResponse
	if json.Unmarshal(data, &response) != nil || len(response.Choices) == 0 || response.Choices[0].Message.Content == "" {
		return "", providerError(KindMalformed, nil)
	}
	return response.Choices[0].Message.Content, nil
}
