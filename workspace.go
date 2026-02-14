package webx

import "time"

// Workspace represents a logical tenant.
type Workspace struct {
	ID         string
	Name       string
	PlanTier   PlanTier
	ProviderID string
	RegionID   string
	CreatedAt  time.Time
}

// Membership binds a user to a workspace with a specific role.
type Membership struct {
	ID          string
	WorkspaceID string
	UserID      string
	RoleID      string
	CreatedAt   time.Time
}

// Role defines a named set of capabilities within a workspace.
type Role struct {
	ID           string
	WorkspaceID  string   // empty = global/system role
	Name         string   // e.g. "owner", "editor", "viewer"
	Capabilities []string // e.g. ["invoices:read", "invoices:write"]
	System       bool     // built-in, cannot be deleted
	CreatedAt    time.Time
}

// WorkspaceStore persists workspaces.
type WorkspaceStore interface {
	Create(ws Workspace) error
	Get(id string) (*Workspace, error)
	List() ([]Workspace, error)
	Update(ws Workspace) error
	Delete(id string) error
}

// MembershipStore persists workspace memberships.
type MembershipStore interface {
	Create(m Membership) error
	Get(id string) (*Membership, error)
	GetByUserAndWorkspace(userID, workspaceID string) (*Membership, error)
	ListByWorkspace(workspaceID string) ([]Membership, error)
	ListByUser(userID string) ([]Membership, error)
	Delete(id string) error
}

// RoleStore persists roles.
type RoleStore interface {
	Create(r Role) error
	Get(id string) (*Role, error)
	GetByNameAndWorkspace(name, workspaceID string) (*Role, error)
	ListByWorkspace(workspaceID string) ([]Role, error)
	Update(r Role) error
	Delete(id string) error
}
