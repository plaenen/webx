package login

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/plaenen/webx"
	"github.com/plaenen/webx/validators"
	"github.com/starfederation/datastar-go/datastar"
)

// SendMagicLinkFunc sends a magic link to the given email.
// In dev mode this can be a no-op since the link is shown on screen.
type SendMagicLinkFunc func(ctx context.Context, email, link string) error

// Handler handles login form submissions.
type Handler struct {
	tokens        webx.TokenStore
	sendMagicLink SendMagicLinkFunc
}

// NewHandler creates a login handler backed by the given token store.
// The sendMagicLink function is called to deliver the link in non-dev mode.
func NewHandler(tokens webx.TokenStore, sendMagicLink SendMagicLinkFunc) *Handler {
	return &Handler{tokens: tokens, sendMagicLink: sendMagicLink}
}

// loginSubmitSignals matches the full signal payload from the login form.
// The validator.Input component nests its state under "login_email".
type loginSubmitSignals struct {
	LoginEmail struct {
		Value string `json:"value"`
		Valid bool   `json:"valid"`
	} `json:"login_email"`
}

// SubmitHandler returns an http.HandlerFunc for the login form SSE endpoint.
// On valid email, it creates a token in the store and sends the magic link.
func (h *Handler) SubmitHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		formID := r.URL.Query().Get("id")
		if formID == "" {
			http.Error(w, "missing id query parameter", http.StatusBadRequest)
			return
		}

		var signals loginSubmitSignals
		if err := datastar.ReadSignals(r, &signals); err != nil {
			writeFormError(w, r, formID, "Failed to read form data")
			return
		}

		email := signals.LoginEmail.Value
		if email == "" {
			writeFormError(w, r, formID, "Email is required")
			return
		}

		res := validators.Email(email, false)
		if !res.Valid {
			writeFormError(w, r, formID, res.Error)
			return
		}

		// Create a login token for this email
		token, err := h.tokens.Create(email)
		if err != nil {
			writeFormError(w, r, formID, "Something went wrong, please try again")
			return
		}

		link := fmt.Sprintf("/verify?token=%s", token.Value)

		// Send the magic link (no-op in dev mode — the link is shown on screen).
		wctx := webx.FromContext(r.Context())
		if !wctx.DevMode && h.sendMagicLink != nil {
			if err := h.sendMagicLink(r.Context(), email, link); err != nil {
				writeFormError(w, r, formID, "Failed to send login email, please try again")
				return
			}
		}

		sse := datastar.NewSSE(w, r)
		sse.PatchElementTempl(MagicLinkSent(link))
	}
}

func writeFormError(w http.ResponseWriter, r *http.Request, formID, msg string) {
	sanitizedID := strings.ReplaceAll(formID, "-", "_")
	sse := datastar.NewSSE(w, r)
	sse.MarshalAndPatchSignals(map[string]any{
		sanitizedID: map[string]any{
			"submitting": false,
			"error":      msg,
		},
	})
}
