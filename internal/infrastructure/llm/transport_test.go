package llm

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

type fakeDoer func(*http.Request) (*http.Response, error)

func (f fakeDoer) Do(r *http.Request) (*http.Response, error) { return f(r) }
func response(status int, body string) *http.Response {
	return &http.Response{StatusCode: status, Body: io.NopCloser(strings.NewReader(body))}
}
func TestTransportTable(t *testing.T) {
	canceled, stop := context.WithCancel(context.Background())
	stop()
	cases := []struct {
		name string
		ctx  context.Context
		do   fakeDoer
		want ErrorKind
	}{{"ok", context.Background(), func(r *http.Request) (*http.Response, error) {
		if _, ok := r.Context().Deadline(); !ok {
			t.Error("no deadline")
		}
		return response(200, "ok"), nil
	}, ""}, {"rate", context.Background(), func(*http.Request) (*http.Response, error) { return response(429, "secret-body"), nil }, KindRateLimit}, {"timeout", context.Background(), func(r *http.Request) (*http.Response, error) {
		if d, ok := r.Context().Deadline(); !ok || time.Until(d) > RequestTimeout {
			t.Error("missing cap")
		}
		return nil, context.DeadlineExceeded
	}, KindTimeout}, {"canceled", canceled, func(*http.Request) (*http.Response, error) { return nil, context.Canceled }, KindCanceled}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			body, err := DoRequest(tc.ctx, tc.do, http.MethodPost, "http://example.invalid", []byte("secret"), nil)
			if tc.want == "" && (err != nil || string(body) != "ok") {
				t.Fatalf("body=%q err=%v", body, err)
			}
			if tc.want != "" && (ErrorKindOf(err) != tc.want || strings.Contains(err.Error(), "secret")) {
				t.Fatalf("kind=%q err=%v", ErrorKindOf(err), err)
			}
		})
	}
}
