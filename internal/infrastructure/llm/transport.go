package llm

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"
)

const RequestTimeout = 8 * time.Second

type HTTPDoer interface {
	Do(*http.Request) (*http.Response, error)
}

func DoRequest(ctx context.Context, client HTTPDoer, method, url string, body []byte, headers http.Header) ([]byte, error) {
	if client == nil {
		return nil, providerError(KindConfig, nil)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	if err != nil {
		return nil, providerError(KindHTTP, err)
	}
	req.Header = headers.Clone()
	attempt, cancel := context.WithTimeout(ctx, RequestTimeout)
	defer cancel()
	resp, err := client.Do(req.WithContext(attempt))
	if err != nil {
		if ctx.Err() == context.Canceled {
			return nil, providerError(KindCanceled, nil)
		}
		if ctx.Err() == context.DeadlineExceeded || err == context.DeadlineExceeded || attempt.Err() == context.DeadlineExceeded {
			return nil, providerError(KindTimeout, nil)
		}
		return nil, providerError(KindUnavailable, nil)
	}
	if resp == nil {
		return nil, providerError(KindUnavailable, nil)
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, statusError(resp.StatusCode)
	}
	return data, nil
}
