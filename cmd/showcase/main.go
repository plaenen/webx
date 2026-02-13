package main

import (
	"fmt"
	"io/fs"
	"log/slog"
	"net"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/plaenen/webx"
	"github.com/plaenen/webx/cmd/showcase/internal/pages"
	"github.com/plaenen/webx/cmd/showcase/internal/static"
	"github.com/plaenen/webx/session/memory"
)

func main() {
	r := chi.NewRouter()

	// Session + CSRF middleware
	store := memory.New()
	r.Use(webx.SessionMiddleware(store))

	// Set dev-mode flag on every request
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wctx := webx.FromContext(r.Context())
			wctx.DevMode = true
			next.ServeHTTP(w, r.WithContext(wctx.WithContext(r.Context())))
		})
	})

	// Serve static files (css, js) at /assets/
	staticFS, _ := fs.Sub(static.Static, "static")
	r.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServerFS(staticFS)))

	// Serve byol files (datastar) at /assets/js/datastar/
	byolFS, _ := fs.Sub(static.Byol, "byol/datastar")
	r.Handle("/assets/js/datastar/*", http.StripPrefix("/assets/js/datastar/", http.FileServerFS(byolFS)))

	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	r.Get("/", templ.Handler(pages.Home()).ServeHTTP)

	ln, err := net.Listen("tcp", ":3000")
	if err != nil {
		ln, err = net.Listen("tcp", ":0")
		if err != nil {
			slog.Error("failed to listen", "error", err)
			return
		}
	}

	slog.Info("server started", "address", fmt.Sprintf("http://localhost:%d", ln.Addr().(*net.TCPAddr).Port))

	if err := http.Serve(ln, r); err != nil {
		slog.Error("server failed", "error", err)
	}
}
