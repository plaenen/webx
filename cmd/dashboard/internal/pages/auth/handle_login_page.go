package auth

import (
	"net/http"

	"github.com/plaenen/webx/cmd/dashboard/internal/pages/auth/templates"
	"github.com/plaenen/webx/cmd/dashboard/internal/session"
)

func LoginPage(w http.ResponseWriter, r *http.Request) {
	if us := session.GetUserSession(r.Context()); us != nil && us.IsLoggedIn() {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	templates.LoginPage().Render(r.Context(), w)
}
