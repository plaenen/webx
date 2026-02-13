package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/plaenen/webx/ui/validator"
	"github.com/plaenen/webx/validators"
)

type validateHandlers struct{}

func newValidateHandlers() *validateHandlers {
	return &validateHandlers{}
}

func (v *validateHandlers) register(r chi.Router) {
	r.Get("/api/validate/email", validator.Handler(func(value string) validator.Result {
		res := validators.Email(value, false)
		return validator.Result{Valid: res.Valid, Error: res.Error}
	}))
	r.Get("/api/validate/email-mx", validator.Handler(func(value string) validator.Result {
		res := validators.Email(value, true)
		return validator.Result{Valid: res.Valid, Error: res.Error}
	}))
}
