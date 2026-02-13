package handlers

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/plaenen/webx/ui/form"
	"github.com/plaenen/webx/validators"
	"github.com/starfederation/datastar-go/datastar"
)

type formHandlers struct{}

func newFormHandlers() *formHandlers {
	return &formHandlers{}
}

func (f *formHandlers) register(r chi.Router) {
	r.Post("/api/form/login", f.login())
	r.Post("/api/form/contact", f.contact())
}

type loginFormSignals struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (f *formHandlers) login() http.HandlerFunc {
	return form.Handler(
		func(formID string, r *http.Request) []form.FieldError {
			var signals loginFormSignals
			if err := form.ReadSignals(formID, r, &signals); err != nil {
				return []form.FieldError{{Field: "error", Message: "Failed to read form data"}}
			}

			var errs []form.FieldError
			if signals.Email == "" {
				errs = append(errs, form.FieldError{Field: "email_error", Message: "Email is required"})
			} else {
				res := validators.Email(signals.Email, false)
				if !res.Valid {
					errs = append(errs, form.FieldError{Field: "email_error", Message: res.Error})
				}
			}
			if signals.Password == "" {
				errs = append(errs, form.FieldError{Field: "password_error", Message: "Password is required"})
			} else if len(signals.Password) < 8 {
				errs = append(errs, form.FieldError{Field: "password_error", Message: "Password must be at least 8 characters"})
			}
			return errs
		},
		func(formID string, sse *datastar.ServerSentEventGenerator) {
			sanitizedID := strings.ReplaceAll(formID, "-", "_")
			sse.MarshalAndPatchSignals(map[string]any{
				sanitizedID: map[string]any{
					"success": "Login successful!",
				},
			})
		},
	)
}

type contactFormSignals struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

func (f *formHandlers) contact() http.HandlerFunc {
	return form.Handler(
		func(formID string, r *http.Request) []form.FieldError {
			var signals contactFormSignals
			if err := form.ReadSignals(formID, r, &signals); err != nil {
				return []form.FieldError{{Field: "error", Message: "Failed to read form data"}}
			}

			var errs []form.FieldError
			if signals.Name == "" {
				errs = append(errs, form.FieldError{Field: "name_error", Message: "Name is required"})
			}
			if signals.Email == "" {
				errs = append(errs, form.FieldError{Field: "email_error", Message: "Email is required"})
			} else {
				res := validators.Email(signals.Email, false)
				if !res.Valid {
					errs = append(errs, form.FieldError{Field: "email_error", Message: res.Error})
				}
			}
			if signals.Message == "" {
				errs = append(errs, form.FieldError{Field: "message_error", Message: "Message is required"})
			}
			return errs
		},
		func(formID string, sse *datastar.ServerSentEventGenerator) {
			sanitizedID := strings.ReplaceAll(formID, "-", "_")
			sse.MarshalAndPatchSignals(map[string]any{
				sanitizedID: map[string]any{
					"success": "Message sent successfully!",
				},
			})
		},
	)
}
