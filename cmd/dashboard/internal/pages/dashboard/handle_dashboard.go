package dashboard

import (
	"net/http"

	"github.com/plaenen/webx"
	"github.com/plaenen/webx/cmd/dashboard/internal/pages/dashboard/templates"
	"github.com/plaenen/webx/cmd/dashboard/internal/session"
)

// WorkspacePickerHandler lists the user's workspaces so they can select one.
func WorkspacePickerHandler(memberships webx.MembershipStore, workspaces webx.WorkspaceStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		us := session.GetUserSession(r.Context())
		if us == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		mbs, err := memberships.ListByUser(us.UserId)
		if err != nil {
			http.Error(w, "failed to list memberships", http.StatusInternalServerError)
			return
		}

		var wsList []webx.Workspace
		for _, m := range mbs {
			ws, err := workspaces.Get(m.WorkspaceID)
			if err != nil || ws == nil {
				continue
			}
			wsList = append(wsList, *ws)
		}

		// If the user has exactly one workspace, redirect directly.
		if len(wsList) == 1 {
			http.Redirect(w, r, "/ws/"+wsList[0].ID+"/", http.StatusSeeOther)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		templates.WorkspacePicker(wsList).Render(r.Context(), w)
	}
}
