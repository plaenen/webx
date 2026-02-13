package internal

import (
	"context"
	"io/fs"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/plaenen/webx"
	"github.com/plaenen/webx/cmd/dashboard/internal/fragments"
	"github.com/plaenen/webx/cmd/dashboard/internal/fragments/login"
	"github.com/plaenen/webx/cmd/dashboard/internal/pages"
	"github.com/plaenen/webx/cmd/dashboard/internal/session"
	"github.com/plaenen/webx/cmd/dashboard/internal/static"
	sessionmem "github.com/plaenen/webx/session/memory"
	tokenmem "github.com/plaenen/webx/token/memory"
	"github.com/plaenen/webx/ui"
	"github.com/plaenen/webx/ui/validator"
	"github.com/plaenen/webx/validators"
)

// Option configures SetupRoutes.
type Option func(*config)

type config struct {
	sessionStore     webx.SessionStore
	userSessionStore session.SessionStore
	tokenStore       webx.TokenStore
	sendMagicLink    login.SendMagicLinkFunc
	pro              bool
	devMode          bool
}

func defaults() *config {
	return &config{
		sessionStore:     sessionmem.New(),
		userSessionStore: session.NewMemSessionProvider(),
		tokenStore:       tokenmem.NewMemTokenStore(15 * time.Minute),
		sendMagicLink: func(_ context.Context, email, link string) error {
			slog.Info("magic link", "to", email, "link", link)
			return nil
		},
		devMode: true,
	}
}

// WithSessionStore sets the framework-level session store (CSRF tokens, etc.).
func WithSessionStore(s webx.SessionStore) Option {
	return func(c *config) { c.sessionStore = s }
}

// WithUserSessionStore sets the user session store.
func WithUserSessionStore(s session.SessionStore) Option {
	return func(c *config) { c.userSessionStore = s }
}

// WithTokenStore sets the magic link token store.
func WithTokenStore(s webx.TokenStore) Option {
	return func(c *config) { c.tokenStore = s }
}

// WithSendMagicLink sets the function used to deliver magic links.
func WithSendMagicLink(fn login.SendMagicLinkFunc) Option {
	return func(c *config) { c.sendMagicLink = fn }
}

// WithPro enables Datastar Pro mode.
func WithPro(pro bool) Option {
	return func(c *config) { c.pro = pro }
}

// WithDevMode sets the dev-mode flag (defaults to true).
func WithDevMode(dev bool) Option {
	return func(c *config) { c.devMode = dev }
}

func SetupRoutes(r chi.Router, opts ...Option) {
	cfg := defaults()
	for _, opt := range opts {
		opt(cfg)
	}

	// Session + CSRF middleware
	r.Use(webx.SessionMiddleware(cfg.sessionStore))
	r.Use(webx.SecurityHeadersMiddleware())
	r.Use(session.Middleware(cfg.userSessionStore))

	// Set dev-mode flag and base path on every request
	const basePath = "/"
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wctx := webx.FromContext(r.Context())
			wctx.DevMode = cfg.devMode
			wctx.DatastarPro = cfg.pro
			wctx.ShowDatastarInspector = cfg.pro
			wctx.BasePath = basePath
			next.ServeHTTP(w, r.WithContext(wctx.WithContext(r.Context())))
		})
	})

	// Serve static files (css, js) at /assets/
	staticFS, _ := fs.Sub(static.Static, "static")
	r.Handle("/assets/*", cacheControl(http.StripPrefix("/assets/", http.FileServerFS(staticFS))))

	// Serve byol files (datastar pro) at /assets/js/datastar/
	if cfg.pro {
		byolFS, _ := fs.Sub(static.Byol, "byol/datastar")
		r.Handle("/assets/js/datastar/*", cacheControl(http.StripPrefix("/assets/js/datastar/", http.FileServerFS(byolFS))))
	}

	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	// Public routes (login, verify, validation, shared UI)
	pages.RegisterRoutes(r, cfg.tokenStore, cfg.userSessionStore)

	ui.RegisterRoutes(r)

	r.Get("/api/validate/email", validator.Handler(func(value string) validator.Result {
		res := validators.Email(value, false)
		return validator.Result{Valid: res.Valid, Error: res.Error}
	}))

	fragments.RegisterRoutes(r, cfg.tokenStore, cfg.sendMagicLink)

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
