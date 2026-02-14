package webx

import "time"

// PlanTier identifies a pricing tier. Apps define their own tier constants.
type PlanTier string

// PlanLimits defines the resource limits for a plan tier.
// A zero value means unlimited for numeric fields.
type PlanLimits struct {
	MaxWorkspaces   int      // per user, 0 = unlimited
	MaxMembers      int      // per workspace, 0 = unlimited
	MaxStorageBytes int64    // per workspace, 0 = unlimited
	MaxTokens       int      // per workspace, 0 = unlimited
	AllowSelfHost   bool     // whether self-hosted providers are available
	AllowedRegions  []string // region IDs available; empty = all
}

// Plan is a named plan with limits.
type Plan struct {
	ID        string
	Tier      PlanTier
	Name      string // display name, e.g. "Pro Plan"
	Limits    PlanLimits
	CreatedAt time.Time
}
