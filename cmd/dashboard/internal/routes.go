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
	"github.com/plaenen/webx/cmd/dashboard/internal/pages/dashboard"
	"github.com/plaenen/webx/cmd/dashboard/internal/session"
	"github.com/plaenen/webx/cmd/dashboard/internal/static"
	membershipmem "github.com/plaenen/webx/providers/membership/memory"
	rolemem "github.com/plaenen/webx/providers/role/memory"
	sessionmem "github.com/plaenen/webx/providers/session/memory"
	providermem "github.com/plaenen/webx/providers/provider/memory"
	regionmem "github.com/plaenen/webx/providers/region/memory"
	submem "github.com/plaenen/webx/providers/subscription/memory"
	tokenmem "github.com/plaenen/webx/providers/token/memory"
	"github.com/plaenen/webx/ui"
	"github.com/plaenen/webx/ui/validator"
	"github.com/plaenen/webx/validators"
	workspacemem "github.com/plaenen/webx/providers/workspace/memory"
)

// Option configures SetupRoutes.
type Option func(*config)

type config struct {
	sessionStore      webx.SessionStore
	userSessionStore  session.SessionStore
	tokenStore        webx.TokenStore
	workspaceStore    webx.WorkspaceStore
	membershipStore   webx.MembershipStore
	roleStore         webx.RoleStore
	subscriptionStore webx.SubscriptionStore
	providerStore     webx.ProviderStore
	regionStore       webx.RegionStore
	billingProvider   webx.BillingProvider
	sendMagicLink     login.SendMagicLinkFunc
	pro               bool
	devMode           bool
}

func defaults() *config {
	return &config{
		sessionStore:      sessionmem.New(),
		userSessionStore: session.NewMemSessionProvider(),
		tokenStore:        tokenmem.NewMemTokenStore(15 * time.Minute),
		workspaceStore:    workspacemem.New(),
		membershipStore:   membershipmem.New(),
		roleStore:         rolemem.New(),
		subscriptionStore: submem.New(),
		providerStore:     providermem.New(),
		regionStore:       regionmem.New(),
		billingProvider:   webx.NoopBillingProvider{},
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

// WithWorkspaceStore sets the workspace store.
func WithWorkspaceStore(s webx.WorkspaceStore) Option {
	return func(c *config) { c.workspaceStore = s }
}

// WithMembershipStore sets the membership store.
func WithMembershipStore(s webx.MembershipStore) Option {
	return func(c *config) { c.membershipStore = s }
}

// WithRoleStore sets the role store.
func WithRoleStore(s webx.RoleStore) Option {
	return func(c *config) { c.roleStore = s }
}

// WithSubscriptionStore sets the subscription store.
func WithSubscriptionStore(s webx.SubscriptionStore) Option {
	return func(c *config) { c.subscriptionStore = s }
}

// WithBillingProvider sets the billing provider.
func WithBillingProvider(b webx.BillingProvider) Option {
	return func(c *config) { c.billingProvider = b }
}

// WithProviderStore sets the provider store.
func WithProviderStore(s webx.ProviderStore) Option {
	return func(c *config) { c.providerStore = s }
}

// WithRegionStore sets the region store.
func WithRegionStore(s webx.RegionStore) Option {
	return func(c *config) { c.regionStore = s }
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

	// In dev mode, seed test data and enable auto-login.
	if cfg.devMode {
		seed, err := SeedDevData(cfg.workspaceStore, cfg.membershipStore, cfg.roleStore, cfg.userSessionStore, cfg.providerStore, cfg.regionStore)
		if err != nil {
			slog.Error("failed to seed dev data", "error", err)
		} else {
			// Auto-login: if no session cookie is present, set the dev session cookie.
			r.Use(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if _, err := r.Cookie("dashboard_session"); err != nil {
						session.SetSessionCookie(w, seed.SessionID)
					}
					next.ServeHTTP(w, r)
				})
			})
		}
	}

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

		r.Get("/dashboard", dashboard.WorkspacePickerHandler(cfg.membershipStore, cfg.workspaceStore))

		r.Get("/workspaces/new", dashboard.WorkspaceNewPageHandler(cfg.providerStore, cfg.regionStore))
		r.Post("/api/workspaces", dashboard.WorkspaceNewSubmitHandler(
			cfg.workspaceStore, cfg.roleStore, cfg.membershipStore, cfg.providerStore, cfg.regionStore, PlanFree, LimitsForTier,
		))

		r.Get("/providers/selfhost/new", dashboard.SelfHostNewPageHandler())
		r.Post("/api/providers/selfhost", dashboard.SelfHostSubmitHandler(cfg.providerStore, ProviderSelfHosted))

		// Workspace-scoped routes
		r.Route("/ws/{wsID}", func(r chi.Router) {
			r.Use(webx.WorkspaceMiddleware(
				cfg.workspaceStore,
				cfg.membershipStore,
				cfg.roleStore,
				func(r *http.Request) string {
					if us := session.GetUserSession(r.Context()); us != nil {
						return us.UserId
					}
					return ""
				},
			))

			r.Get("/", dashboard.WorkspaceDashboardHandler())
		})
	})
}
