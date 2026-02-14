package webx

import "time"

// Subscription tracks a workspace's billing subscription.
type Subscription struct {
	ID             string
	WorkspaceID    string
	CustomerID     string // billing provider customer ID
	SubscriptionID string // billing provider subscription ID
	Tier           PlanTier
	Active         bool
	CreatedAt      time.Time
}

// SubscriptionStore persists subscriptions.
type SubscriptionStore interface {
	Create(s Subscription) error
	Get(id string) (*Subscription, error)
	GetByWorkspace(workspaceID string) (*Subscription, error)
	Update(s Subscription) error
	Delete(id string) error
}
