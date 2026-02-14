package webx

import "fmt"

// LimitStatus reports usage against a limit.
type LimitStatus struct {
	Limit     int // 0 = unlimited
	Used      int
	Remaining int // -1 = unlimited
}

// CheckWorkspaceLimit returns how many workspaces a user has vs. the given limits.
func CheckWorkspaceLimit(memberships MembershipStore, userID string, limits PlanLimits) (*LimitStatus, error) {
	mbs, err := memberships.ListByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("listing user memberships: %w", err)
	}

	used := len(mbs)

	remaining := -1
	if limits.MaxWorkspaces > 0 {
		remaining = max(limits.MaxWorkspaces-used, 0)
	}

	return &LimitStatus{
		Limit:     limits.MaxWorkspaces,
		Used:      used,
		Remaining: remaining,
	}, nil
}

// CheckMemberLimit returns how many members a workspace has vs. the given limits.
func CheckMemberLimit(memberships MembershipStore, workspaceID string, limits PlanLimits) (*LimitStatus, error) {
	mbs, err := memberships.ListByWorkspace(workspaceID)
	if err != nil {
		return nil, fmt.Errorf("listing workspace memberships: %w", err)
	}

	used := len(mbs)

	remaining := -1
	if limits.MaxMembers > 0 {
		remaining = max(limits.MaxMembers-used, 0)
	}

	return &LimitStatus{
		Limit:     limits.MaxMembers,
		Used:      used,
		Remaining: remaining,
	}, nil
}
