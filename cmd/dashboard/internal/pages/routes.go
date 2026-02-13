package pages

import (
	"github.com/go-chi/chi/v5"
	"github.com/plaenen/webx"
	"github.com/plaenen/webx/cmd/dashboard/internal/pages/auth"
	"github.com/plaenen/webx/cmd/dashboard/internal/session"
)

func RegisterRoutes(r chi.Router, tokens webx.TokenStore, sessions session.SessionStore) {
	// Public routes (no auth required)
	r.Get("/login", auth.LoginPage)
	r.Get("/verify", auth.VerifyHandler(tokens, sessions))
}
