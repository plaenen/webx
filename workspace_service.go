package webx

import (
	"errors"
	"fmt"
	"time"

	"github.com/plaenen/webx/idgen"
)

var (
	// ErrWorkspaceLimitReached is returned when a user has reached the
	// maximum number of workspaces allowed by their plan.
	ErrWorkspaceLimitReached = errors.New("workspace limit reached for current plan")

	// ErrMemberLimitReached is returned when a workspace has reached the
	// maximum number of members allowed by its plan.
	ErrMemberLimitReached = errors.New("member limit reached for current plan")
)

// CreateWorkspaceInput is the input for creating a new workspace.
type CreateWorkspaceInput struct {
	Name       string
	OwnerID    string
	Tier       PlanTier    // stored on the workspace
	Limits     *PlanLimits // if non-nil, enforce MaxWorkspaces before creation
	ProviderID string
	RegionID   string
}

// CreateWorkspaceResult is the output of workspace creation.
type CreateWorkspaceResult struct {
	Workspace  Workspace
	OwnerRole  Role
	Membership Membership
}

// CreateWorkspace creates a workspace with a system "owner" role and assigns
// the creating user as owner. This enforces the invariant that every workspace
// always has at least one owner. When Limits is provided, it checks that the
// user has not exceeded the workspace cap before creating.
func CreateWorkspace(
	workspaces WorkspaceStore,
	roles RoleStore,
	memberships MembershipStore,
	input CreateWorkspaceInput,
) (*CreateWorkspaceResult, error) {
	if input.Limits != nil && input.Limits.MaxWorkspaces > 0 {
		existing, err := memberships.ListByUser(input.OwnerID)
		if err != nil {
			return nil, fmt.Errorf("checking workspace limit: %w", err)
		}
		if len(existing) >= input.Limits.MaxWorkspaces {
			return nil, ErrWorkspaceLimitReached
		}
	}

	now := time.Now()

	ws := Workspace{
		ID:         idgen.Generate("ws"),
		Name:       input.Name,
		PlanTier:   input.Tier,
		ProviderID: input.ProviderID,
		RegionID:   input.RegionID,
		CreatedAt:  now,
	}
	if err := workspaces.Create(ws); err != nil {
		return nil, fmt.Errorf("creating workspace: %w", err)
	}

	ownerRole := Role{
		ID:           idgen.Generate("role"),
		WorkspaceID:  ws.ID,
		Name:         "owner",
		Capabilities: []string{"*"},
		System:       true,
		CreatedAt:    now,
	}
	if err := roles.Create(ownerRole); err != nil {
		return nil, fmt.Errorf("creating owner role: %w", err)
	}

	mbr := Membership{
		ID:          idgen.Generate("mbr"),
		WorkspaceID: ws.ID,
		UserID:      input.OwnerID,
		RoleID:      ownerRole.ID,
		CreatedAt:   now,
	}
	if err := memberships.Create(mbr); err != nil {
		return nil, fmt.Errorf("creating owner membership: %w", err)
	}

	return &CreateWorkspaceResult{
		Workspace:  ws,
		OwnerRole:  ownerRole,
		Membership: mbr,
	}, nil
}

// AddMemberInput is the input for adding a member to a workspace.
type AddMemberInput struct {
	WorkspaceID string
	UserID      string
	RoleID      string
	Limits      *PlanLimits // if non-nil, enforce MaxMembers before adding
}

// AddMember adds a user to a workspace with the given role. When Limits is
// provided, it checks that the workspace has not exceeded the member cap.
func AddMember(
	memberships MembershipStore,
	input AddMemberInput,
) (*Membership, error) {
	if input.Limits != nil && input.Limits.MaxMembers > 0 {
		existing, err := memberships.ListByWorkspace(input.WorkspaceID)
		if err != nil {
			return nil, fmt.Errorf("checking member limit: %w", err)
		}
		if len(existing) >= input.Limits.MaxMembers {
			return nil, ErrMemberLimitReached
		}
	}

	mbr := Membership{
		ID:          idgen.Generate("mbr"),
		WorkspaceID: input.WorkspaceID,
		UserID:      input.UserID,
		RoleID:      input.RoleID,
		CreatedAt:   time.Now(),
	}
	if err := memberships.Create(mbr); err != nil {
		return nil, fmt.Errorf("creating membership: %w", err)
	}

	return &mbr, nil
}
