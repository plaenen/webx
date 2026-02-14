package internal

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/plaenen/webx"
	"github.com/plaenen/webx/cmd/dashboard/internal/session"
	"github.com/plaenen/webx/idgen"
)

// SeedResult holds the output of a dev seed so the caller can set up
// auto-login or log the session details.
type SeedResult struct {
	SessionID string // cookie value for auto-login
	UserID    string
	Email     string
}

// SeedDevData populates stores with test data for local development.
func SeedDevData(
	workspaces webx.WorkspaceStore,
	memberships webx.MembershipStore,
	roles webx.RoleStore,
	sessions session.SessionStore,
	providers webx.ProviderStore,
	regions webx.RegionStore,
) (*SeedResult, error) {
	now := time.Now()
	userID := idgen.Generate("user")
	email := "dev@example.com"

	// --- Managed providers + regions ---

	aws := webx.Provider{
		ID:        "prov_aws",
		Name:      "AWS",
		Type:      ProviderManaged,
		CreatedAt: now,
	}
	gcp := webx.Provider{
		ID:        "prov_gcp",
		Name:      "Google Cloud",
		Type:      ProviderManaged,
		CreatedAt: now,
	}
	for _, p := range []webx.Provider{aws, gcp} {
		if err := providers.Create(p); err != nil {
			return nil, fmt.Errorf("seed provider %q: %w", p.Name, err)
		}
	}

	awsRegions := []webx.Region{
		{ID: "reg_aws_use1", ProviderID: "prov_aws", Name: "us-east-1", Label: "US East (Virginia)", CreatedAt: now},
		{ID: "reg_aws_euw1", ProviderID: "prov_aws", Name: "eu-west-1", Label: "EU West (Ireland)", CreatedAt: now},
	}
	gcpRegions := []webx.Region{
		{ID: "reg_gcp_usc1", ProviderID: "prov_gcp", Name: "us-central1", Label: "US Central (Iowa)", CreatedAt: now},
	}
	for _, r := range append(awsRegions, gcpRegions...) {
		if err := regions.Create(r); err != nil {
			return nil, fmt.Errorf("seed region %q: %w", r.Name, err)
		}
	}

	// --- Workspaces (owner role + membership created automatically) ---

	acme, err := webx.CreateWorkspace(workspaces, roles, memberships, webx.CreateWorkspaceInput{
		Name:       "Acme Corp",
		OwnerID:    userID,
		Tier:       PlanPro,
		ProviderID: "prov_aws",
		RegionID:   "reg_aws_use1",
	})
	if err != nil {
		return nil, fmt.Errorf("seed acme: %w", err)
	}

	personal, err := webx.CreateWorkspace(workspaces, roles, memberships, webx.CreateWorkspaceInput{
		Name:       "Personal",
		OwnerID:    userID,
		Tier:       PlanFree,
		ProviderID: "prov_aws",
		RegionID:   "reg_aws_use1",
	})
	if err != nil {
		return nil, fmt.Errorf("seed personal: %w", err)
	}

	// --- Extra roles for Acme (beyond the auto-created owner) ---

	editorRole := webx.Role{
		ID:          idgen.Generate("role"),
		WorkspaceID: acme.Workspace.ID,
		Name:        "editor",
		Capabilities: []string{
			"dashboard:read",
			"members:read",
			"invoices:read",
			"invoices:create",
			"invoices:update",
			"reports:read",
		},
		CreatedAt: now,
	}
	viewerRole := webx.Role{
		ID:          idgen.Generate("role"),
		WorkspaceID: acme.Workspace.ID,
		Name:        "viewer",
		Capabilities: []string{
			"dashboard:read",
			"invoices:read",
			"reports:read",
		},
		CreatedAt: now,
	}
	for _, r := range []webx.Role{editorRole, viewerRole} {
		if err := roles.Create(r); err != nil {
			return nil, fmt.Errorf("seed role %q: %w", r.Name, err)
		}
	}

	// --- User session (pre-logged-in) ---

	sessionID, err := sessions.Create(session.UserSession{
		UserId:      userID,
		Email:       email,
		DisplayName: "Dev User",
		CreatedAt:   now,
	})
	if err != nil {
		return nil, fmt.Errorf("seed user session: %w", err)
	}

	slog.Info("dev seed complete",
		"email", email,
		"user_id", userID,
		"session_id", sessionID,
		"workspaces", []string{acme.Workspace.ID, personal.Workspace.ID},
	)

	return &SeedResult{
		SessionID: sessionID,
		UserID:    userID,
		Email:     email,
	}, nil
}
