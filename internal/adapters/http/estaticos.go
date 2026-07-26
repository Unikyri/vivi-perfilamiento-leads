package http

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"path"
	"strings"
	"time"
)

const staticCSP = "default-src 'self'; base-uri 'self'; object-src 'none'; frame-ancestors 'none'; script-src 'self'; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; font-src 'self' https://fonts.gstatic.com; img-src 'self' data:; connect-src 'self'"

// staticFiles contains the output produced by the Vite build. CI packages it
// before compiling the Go binary; the entry check below prevents a bad binary.
//
//go:embed all:estaticos
var staticFiles embed.FS

// RegistrarEstaticos registers the SPA after the API and health routes. The
// ServeMux gives those more-specific routes precedence over the root subtree.
func RegistrarEstaticos(mux *http.ServeMux) error {
	if mux == nil {
		return errors.New("static routes require a mux")
	}
	root, err := fs.Sub(staticFiles, "estaticos")
	if err != nil {
		return fmt.Errorf("static root: %w", err)
	}
	if info, err := fs.Stat(root, "index.html"); err != nil || info.IsDir() {
		return fmt.Errorf("static entry index.html is missing")
	}
	mux.Handle("/", staticHandler(root))
	return nil
}

func staticHandler(root fs.FS) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", staticCSP)
		w.Header().Set("X-Content-Type-Options", "nosniff")
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		requestPath := strings.TrimPrefix(r.URL.Path, "/")
		if requestPath == "" {
			requestPath = "index.html"
		}
		if fs.ValidPath(requestPath) {
			if info, err := fs.Stat(root, requestPath); err == nil && !info.IsDir() {
				if data, err := fs.ReadFile(root, requestPath); err == nil {
					http.ServeContent(w, r, path.Base(requestPath), info.ModTime(), bytes.NewReader(data))
					return
				}
			}
		}

		// An asset miss must not silently become a successful HTML response.
		// Extensionless paths outside /assets are client-side browser routes.
		if !strings.HasPrefix(requestPath, "assets/") && path.Ext(requestPath) == "" {
			if data, err := fs.ReadFile(root, "index.html"); err == nil {
				http.ServeContent(w, r, "index.html", time.Time{}, bytes.NewReader(data))
				return
			}
		}
		http.NotFound(w, r)
	})
}
