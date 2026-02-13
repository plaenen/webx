package fragments

import (
	"github.com/go-chi/chi/v5"
	"github.com/plaenen/webx"
	"github.com/plaenen/webx/cmd/dashboard/internal/fragments/login"
)

func RegisterRoutes(r chi.Router, tokens webx.TokenStore) {
	loginHandler := login.NewHandler(tokens)
	r.Post("/api/auth/login", loginHandler.SubmitHandler())
}
