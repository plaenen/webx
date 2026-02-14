package dashboard

import (
	"net/http"

	"github.com/plaenen/webx"
	"github.com/plaenen/webx/cmd/dashboard/internal/pages/dashboard/templates"
	"github.com/plaenen/webx/cmd/dashboard/internal/session"
	"github.com/plaenen/webx/layouts"
)

// WorkspaceDashboardHandler renders the workspace dashboard with the full layout.
func WorkspaceDashboardHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wsc := webx.GetWorkspaceContext(r.Context())
		if wsc == nil {
			http.Error(w, "workspace context missing", http.StatusInternalServerError)
			return
		}

		us := session.GetUserSession(r.Context())

		user := layouts.UserInfo{}
		if us != nil {
			user.Name = us.DisplayName
			user.Email = us.Email
		}

		w.Header().Set("Content-Type", "text/html")
		templates.WorkspaceDashboard(wsc.Workspace, wsc, user, r.URL.Path).Render(r.Context(), w)
	}
}
