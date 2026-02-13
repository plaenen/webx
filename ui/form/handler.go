package form

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/starfederation/datastar-go/datastar"
)

// FieldError represents a validation error for a specific form field.
type FieldError struct {
	// Field is the signal name (e.g. "email_error").
	Field string
	// Message is the error text shown to the user.
	Message string
}

// SubmitFunc processes a form submission.
// It receives the form ID and the raw request, and returns field errors.
// Return nil or empty slice for success.
type SubmitFunc func(formID string, r *http.Request) []FieldError

// Handler returns an http.HandlerFunc that processes form submissions via SSE.
// On validation failure, it patches field error signals.
// On success, it calls the onSuccess callback to patch success state.
//
// Mount at your form's Action path:
//
//	r.Post("/api/auth/login", form.Handler(loginHandler, loginSuccess))
func Handler(validate SubmitFunc, onSuccess func(formID string, sse *datastar.ServerSentEventGenerator)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		formID := r.URL.Query().Get("id")
		if formID == "" {
			http.Error(w, "missing id query parameter", http.StatusBadRequest)
			return
		}

		sanitizedID := strings.ReplaceAll(formID, "-", "_")
		sse := datastar.NewSSE(w, r)

		// Set submitting = true
		sse.MarshalAndPatchSignals(map[string]any{
			sanitizedID: map[string]any{
				"submitting": true,
				"error":      "",
			},
		})

		errors := validate(formID, r)

		if len(errors) > 0 {
			// Patch field errors and clear submitting
			patch := map[string]any{
				"submitting": false,
			}
			for _, e := range errors {
				patch[e.Field] = e.Message
			}
			sse.MarshalAndPatchSignals(map[string]any{
				sanitizedID: patch,
			})
			return
		}

		// Clear submitting on success
		sse.MarshalAndPatchSignals(map[string]any{
			sanitizedID: map[string]any{
				"submitting": false,
			},
		})

		if onSuccess != nil {
			onSuccess(formID, sse)
		}
	}
}

// ReadSignals reads the form's namespaced signals from the request.
// Pass a pointer to your signals struct.
//
//	type LoginSignals struct {
//	    Email    string `json:"email"`
//	    Password string `json:"password"`
//	}
//	var signals LoginSignals
//	if err := form.ReadSignals("login", r, &signals); err != nil { ... }
func ReadSignals(formID string, r *http.Request, dest any) error {
	sanitizedID := strings.ReplaceAll(formID, "-", "_")
	wrapper := map[string]any{sanitizedID: dest}
	if err := datastar.ReadSignals(r, &wrapper); err != nil {
		return fmt.Errorf("read form signals: %w", err)
	}
	return nil
}
