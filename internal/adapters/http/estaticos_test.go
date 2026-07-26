package http

import (
	"bytes"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"path"
	"regexp"
	"strings"
	"testing"
)

func TestRegistrarEstaticosServesSPAAssetsAndPreservesRoutes(t *testing.T) {
	index, err := fs.ReadFile(staticFiles, "estaticos/index.html")
	if err != nil {
		t.Fatalf("read embedded entry: %v", err)
	}
	assetPath, asset, err := embeddedAsset()
	if err != nil {
		t.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte("api"))
	})
	mux.HandleFunc("GET /salud", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	if err := RegistrarEstaticos(mux); err != nil {
		t.Fatalf("register static handler: %v", err)
	}

	tests := []struct {
		name        string
		path        string
		status      int
		body        []byte
		static      bool
		contentType string
	}{
		{name: "entry", path: "/", status: http.StatusOK, body: index, static: true, contentType: "text/html"},
		{name: "hashed asset", path: "/" + assetPath, status: http.StatusOK, body: asset, static: true},
		{name: "browser fallback", path: "/dashboard/leads", status: http.StatusOK, body: index, static: true, contentType: "text/html"},
		{name: "missing asset", path: "/assets/missing.js", status: http.StatusNotFound, static: true},
		{name: "api precedence", path: "/api/unknown", status: http.StatusAccepted, body: []byte("api")},
		{name: "health precedence", path: "/salud", status: http.StatusNoContent},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			mux.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, tt.path, nil))
			if recorder.Code != tt.status {
				t.Fatalf("status = %d, want %d", recorder.Code, tt.status)
			}
			if tt.body != nil && !bytes.Equal(recorder.Body.Bytes(), tt.body) {
				t.Fatalf("body does not match expected packaged bytes")
			}
			if !tt.static {
				return
			}
			if got := recorder.Header().Get("X-Content-Type-Options"); got != "nosniff" {
				t.Errorf("X-Content-Type-Options = %q, want nosniff", got)
			}
			if got := recorder.Header().Get("Content-Security-Policy"); got != staticCSP {
				t.Errorf("Content-Security-Policy = %q, want %q", got, staticCSP)
			}
			if tt.contentType != "" && recorder.Header().Get("Content-Type") == "" {
				t.Errorf("Content-Type is empty, want %s", tt.contentType)
			}
		})
	}
}

func TestPackagedIndexReferencesExistingEmbeddedFiles(t *testing.T) {
	index, err := fs.ReadFile(staticFiles, "estaticos/index.html")
	if err != nil {
		t.Fatalf("read embedded entry: %v", err)
	}

	refs := regexp.MustCompile(`(?:src|href)="(/[^\"]+)"`).FindAllSubmatch(index, -1)
	if len(refs) == 0 {
		t.Fatal("packaged entry has no local asset references")
	}
	for _, ref := range refs {
		name := strings.TrimPrefix(string(ref[1]), "/")
		if !fs.ValidPath(name) {
			t.Errorf("invalid packaged asset path %q", name)
			continue
		}
		info, err := fs.Stat(staticFiles, path.Join("estaticos", name))
		if err != nil {
			t.Errorf("index references missing packaged file %q: %v", name, err)
		} else if info.IsDir() {
			t.Errorf("index references packaged directory %q", name)
		}
	}
}

func TestRegistrarEstaticosRejectsNilMux(t *testing.T) {
	if err := RegistrarEstaticos(nil); err == nil {
		t.Fatal("RegistrarEstaticos(nil) returned nil error")
	}
}

func embeddedAsset() (string, []byte, error) {
	entries, err := fs.ReadDir(staticFiles, "estaticos/assets")
	if err != nil {
		return "", nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := "assets/" + entry.Name()
		content, err := fs.ReadFile(staticFiles, "estaticos/"+name)
		if err != nil {
			return "", nil, err
		}
		return name, content, nil
	}
	return "", nil, fs.ErrNotExist
}
