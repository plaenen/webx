package auth

import (
	"net/http"
	"time"

	"github.com/plaenen/webx"
	"github.com/plaenen/webx/cmd/dashboard/internal/pages/auth/templates"
	"github.com/plaenen/webx/cmd/dashboard/internal/session"
	"github.com/plaenen/webx/idgen"
)

// VerifyHandler returns a handler that validates magic link tokens
// and creates a user session on success.
func VerifyHandler(tokens webx.TokenStore, sessions session.SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenValue := r.URL.Query().Get("token")
		if tokenValue == "" {
			renderVerifyFailed(w, r)
			return
		}

		tok, err := tokens.Get(tokenValue)
		if err != nil || tok == nil || tok.Expired() {
			renderVerifyFailed(w, r)
			return
		}

		// Mark token as verified and delete it so it can't be reused.
		tokens.MarkVerified(tokenValue)
		tokens.Delete(tokenValue)

		// Create a user session for this email.
		sessionID, err := sessions.Create(session.UserSession{
			UserId:    idgen.Generate("user"),
			Email:     tok.Email,
			CreatedAt: time.Now(),
		})
		if err != nil {
			http.Error(w, "failed to create session", http.StatusInternalServerError)
			return
		}

		session.SetSessionCookie(w, sessionID)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func renderVerifyFailed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusUnauthorized)
	templates.VerifyFailed().Render(r.Context(), w)
}
