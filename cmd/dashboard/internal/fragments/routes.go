package fragments

import (
	"github.com/go-chi/chi/v5"
	"github.com/plaenen/webx"
	"github.com/plaenen/webx/cmd/dashboard/internal/fragments/login"
)

func RegisterRoutes(r chi.Router, tokens webx.TokenStore, sendMagicLink login.SendMagicLinkFunc) {
	loginHandler := login.NewHandler(tokens, sendMagicLink)
	r.Post("/api/auth/login", loginHandler.SubmitHandler())
}
