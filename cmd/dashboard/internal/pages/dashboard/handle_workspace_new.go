package dashboard

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/plaenen/webx"
	"github.com/plaenen/webx/cmd/dashboard/internal/pages/dashboard/templates"
	"github.com/plaenen/webx/cmd/dashboard/internal/session"
	"github.com/starfederation/datastar-go/datastar"
)

// PlanResolver returns the limits for a given tier.
type PlanResolver func(tier webx.PlanTier) webx.PlanLimits

// WorkspaceNewPageHandler renders the create-workspace form (GET).
func WorkspaceNewPageHandler(
	providers webx.ProviderStore,
	regions webx.RegionStore,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		provList, err := providers.List()
		if err != nil {
			http.Error(w, "failed to load providers", http.StatusInternalServerError)
			return
		}

		var regList []webx.Region
		for _, p := range provList {
			regs, err := regions.ListByProvider(p.ID)
			if err != nil {
				http.Error(w, "failed to load regions", http.StatusInternalServerError)
				return
			}
			regList = append(regList, regs...)
		}

		w.Header().Set("Content-Type", "text/html")
		templates.WorkspaceNew(provList, regList).Render(r.Context(), w)
	}
}

type wsNewSignals struct {
	WsNew struct {
		Name       string `json:"name"`
		ProviderID string `json:"provider_id"`
		RegionID   string `json:"region_id"`
	} `json:"ws_new"`
}

// WorkspaceNewSubmitHandler processes the create-workspace form via SSE.
func WorkspaceNewSubmitHandler(
	workspaces webx.WorkspaceStore,
	roles webx.RoleStore,
	memberships webx.MembershipStore,
	providers webx.ProviderStore,
	regions webx.RegionStore,
	defaultTier webx.PlanTier,
	planResolver PlanResolver,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		us := session.GetUserSession(r.Context())
		if us == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		var signals wsNewSignals
		if err := datastar.ReadSignals(r, &signals); err != nil {
			writeWSError(w, r, "Failed to read form data")
			return
		}

		name := strings.TrimSpace(signals.WsNew.Name)
		if name == "" {
			writeWSError(w, r, "Name is required.")
			return
		}

		providerID := strings.TrimSpace(signals.WsNew.ProviderID)
		regionID := strings.TrimSpace(signals.WsNew.RegionID)

		if providerID == "" || regionID == "" {
			writeWSError(w, r, "Provider and region are required.")
			return
		}

		// Validate provider exists
		prov, err := providers.Get(providerID)
		if err != nil || prov == nil {
			writeWSError(w, r, "Invalid provider selected.")
			return
		}

		// Validate region exists and belongs to provider
		reg, err := regions.Get(regionID)
		if err != nil || reg == nil || reg.ProviderID != providerID {
			writeWSError(w, r, "Invalid region selected.")
			return
		}

		limits := planResolver(defaultTier)
		result, err := webx.CreateWorkspace(workspaces, roles, memberships, webx.CreateWorkspaceInput{
			Name:       name,
			OwnerID:    us.UserId,
			Tier:       defaultTier,
			Limits:     &limits,
			ProviderID: providerID,
			RegionID:   regionID,
		})
		if err != nil {
			msg := "Failed to create workspace."
			if errors.Is(err, webx.ErrWorkspaceLimitReached) {
				msg = "You have reached the maximum number of workspaces for your plan."
			}
			writeWSError(w, r, msg)
			return
		}

		sse := datastar.NewSSE(w, r)
		_ = sse.Redirect(fmt.Sprintf("/ws/%s/", result.Workspace.ID))
	}
}

func writeWSError(w http.ResponseWriter, r *http.Request, msg string) {
	sse := datastar.NewSSE(w, r)
	sse.MarshalAndPatchSignals(map[string]any{
		"ws_new": map[string]any{
			"submitting": false,
			"error":      msg,
		},
	})
}
