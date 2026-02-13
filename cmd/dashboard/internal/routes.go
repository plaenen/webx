package internal

import (
	"io/fs"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/plaenen/webx"
	"github.com/plaenen/webx/cmd/dashboard/internal/fragments"
	"github.com/plaenen/webx/cmd/dashboard/internal/pages"
	"github.com/plaenen/webx/cmd/dashboard/internal/session"
	"github.com/plaenen/webx/cmd/dashboard/internal/static"
	"github.com/plaenen/webx/session/memory"
	tokenmem "github.com/plaenen/webx/token/memory"
	"github.com/plaenen/webx/ui"
	"github.com/plaenen/webx/ui/validator"
	"github.com/plaenen/webx/validators"
)

func SetupRoutes(r chi.Router, pro bool) {

	sessionStore := memory.New()
	userSessionStore := session.NewMemSessionProvider()
	tokenStore := tokenmem.NewMemTokenStore(15 * time.Minute)

	// Session + CSRF middleware
	r.Use(webx.SessionMiddleware(sessionStore))
	r.Use(webx.SecurityHeadersMiddleware())
	r.Use(session.Middleware(userSessionStore))

	// Set dev-mode flag and base path on every request
	const basePath = "/"
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wctx := webx.FromContext(r.Context())
			wctx.DevMode = true
			wctx.DatastarPro = pro
			wctx.ShowDatastarInspector = pro
			wctx.BasePath = basePath
			next.ServeHTTP(w, r.WithContext(wctx.WithContext(r.Context())))
		})
	})

	// Serve static files (css, js) at /assets/
	staticFS, _ := fs.Sub(static.Static, "static")
	r.Handle("/assets/*", cacheControl(http.StripPrefix("/assets/", http.FileServerFS(staticFS))))

	// Serve byol files (datastar pro) at /assets/js/datastar/
	if pro {
		byolFS, _ := fs.Sub(static.Byol, "byol/datastar")
		r.Handle("/assets/js/datastar/*", cacheControl(http.StripPrefix("/assets/js/datastar/", http.FileServerFS(byolFS))))
	}

	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	// Public routes (login, verify, validation, shared UI)
	pages.RegisterRoutes(r, tokenStore, userSessionStore)

	ui.RegisterRoutes(r)

	r.Get("/api/validate/email", validator.Handler(func(value string) validator.Result {
		res := validators.Email(value, false)
		return validator.Result{Valid: res.Valid, Error: res.Error}
	}))

	fragments.RegisterRoutes(r, tokenStore)

	// Protected routes — require authenticated session
	r.Group(func(r chi.Router) {
		r.Use(session.RequireAuth)

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		})

		r.Get("/dashboard", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte("dashboard (coming soon)"))
		})
	})
}
