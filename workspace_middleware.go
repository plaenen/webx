package webx

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/plaenen/webx/capability"
)

// WorkspaceMiddleware resolves the workspace from the URL parameter "wsID",
// verifies the user is a member, derives capabilities from the role, and
// attaches a WorkspaceContext to the request context.
//
// getUserID extracts the current user's ID from the request. This decouples
// the middleware from any specific session implementation.
func WorkspaceMiddleware(
	workspaces WorkspaceStore,
	memberships MembershipStore,
	roles RoleStore,
	getUserID func(r *http.Request) string,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wsID := chi.URLParam(r, "wsID")
			if wsID == "" {
				http.Error(w, "workspace not found", http.StatusNotFound)
				return
			}

			ws, err := workspaces.Get(wsID)
			if err != nil {
				http.Error(w, "workspace lookup error", http.StatusInternalServerError)
				return
			}
			if ws == nil {
				http.Error(w, "workspace not found", http.StatusNotFound)
				return
			}

			userID := getUserID(r)
			if userID == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			mbr, err := memberships.GetByUserAndWorkspace(userID, ws.ID)
			if err != nil {
				http.Error(w, "membership lookup error", http.StatusInternalServerError)
				return
			}
			if mbr == nil {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			role, err := roles.Get(mbr.RoleID)
			if err != nil {
				http.Error(w, "role lookup error", http.StatusInternalServerError)
				return
			}
			if role == nil {
				http.Error(w, "role not found", http.StatusInternalServerError)
				return
			}

			wsc := &WorkspaceContext{
				Workspace:    ws,
				Membership:   mbr,
				Role:         role,
				Capabilities: capability.Set(role.Capabilities),
			}

			ctx := WithWorkspaceContext(r.Context(), wsc)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireCapability returns middleware that checks the user has the required
// capability in the current workspace context. Returns 403 if not.
func RequireCapability(required string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wsc := GetWorkspaceContext(r.Context())
			if wsc == nil || !wsc.Capabilities.Can(required) {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
