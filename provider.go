package webx

import "time"

// ProviderType distinguishes managed cloud providers from self-hosted.
// Apps define their own constants (e.g. "managed", "selfhosted").
type ProviderType string

// Provider represents a deployment target for workspaces.
type Provider struct {
	ID          string
	Name        string       // e.g. "AWS", "Google Cloud", "My Home Server"
	Type        ProviderType // app-defined, e.g. "managed" or "selfhosted"
	OwnerID     string       // self-hosted: user who created it (empty for managed)
	PairingCode string       // self-hosted: short code for pairing
	PairedAt    *time.Time   // self-hosted: nil = pending pairing
	CreatedAt   time.Time
}

// Region represents a deployment region within a provider.
type Region struct {
	ID         string
	ProviderID string
	Name       string // e.g. "us-east-1"
	Label      string // e.g. "US East (Virginia)"
	CreatedAt  time.Time
}

// ProviderStore persists providers.
type ProviderStore interface {
	Create(p Provider) error
	Get(id string) (*Provider, error)
	List() ([]Provider, error)
	ListByType(t ProviderType) ([]Provider, error)
	ListByOwner(ownerID string) ([]Provider, error)
	Update(p Provider) error
	Delete(id string) error
}

// RegionStore persists regions.
type RegionStore interface {
	Create(r Region) error
	Get(id string) (*Region, error)
	ListByProvider(providerID string) ([]Region, error)
	Delete(id string) error
}
