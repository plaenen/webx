package webx

// BillingProvider abstracts a billing/payment backend (e.g. Stripe).
type BillingProvider interface {
	// CreateCustomer creates a billing customer for a workspace.
	CreateCustomer(workspaceID, email string) (customerID string, err error)
	// CreateSubscription starts a subscription for the given tier.
	CreateSubscription(customerID string, tier PlanTier) (subscriptionID string, err error)
	// CancelSubscription cancels the workspace's subscription.
	CancelSubscription(subscriptionID string) error
	// GetPortalURL returns a URL where the user can manage billing.
	GetPortalURL(customerID string) (string, error)
	// ChangePlan upgrades or downgrades to a new tier.
	ChangePlan(subscriptionID string, tier PlanTier) error
}

// NoopBillingProvider is a no-op implementation for dev/testing.
type NoopBillingProvider struct{}

func (NoopBillingProvider) CreateCustomer(_, _ string) (string, error)              { return "noop", nil }
func (NoopBillingProvider) CreateSubscription(_ string, _ PlanTier) (string, error) { return "noop", nil }
func (NoopBillingProvider) CancelSubscription(_ string) error                       { return nil }
func (NoopBillingProvider) GetPortalURL(_ string) (string, error)                   { return "", nil }
func (NoopBillingProvider) ChangePlan(_ string, _ PlanTier) error                   { return nil }
