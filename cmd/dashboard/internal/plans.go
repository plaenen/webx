package internal

import "github.com/plaenen/webx"

// App-specific plan tiers.
const (
	PlanFree       webx.PlanTier = "free"
	PlanPro        webx.PlanTier = "pro"
	PlanEnterprise webx.PlanTier = "enterprise"
)

// App-specific provider types.
const (
	ProviderManaged    webx.ProviderType = "managed"
	ProviderSelfHosted webx.ProviderType = "selfhosted"
)

// DefaultRegionID is the region used when none is specified.
const DefaultRegionID = "reg_aws_use1"

// PlanLimits maps each tier to its resource limits.
var PlanLimits = map[webx.PlanTier]webx.PlanLimits{
	PlanFree: {
		MaxWorkspaces:  1,
		MaxMembers:     5,
		MaxStorageBytes: 1 << 30,
		MaxTokens:      100,
		AllowSelfHost:  false,
		AllowedRegions: []string{"reg_aws_use1"},
	},
	PlanPro: {
		MaxWorkspaces:  5,
		MaxMembers:     25,
		MaxStorageBytes: 10 << 30,
		MaxTokens:      1000,
		AllowSelfHost:  false,
		AllowedRegions: []string{"reg_aws_use1", "reg_aws_euw1", "reg_gcp_usc1"},
	},
	PlanEnterprise: {
		MaxWorkspaces:  0,
		MaxMembers:     0,
		MaxStorageBytes: 0,
		MaxTokens:      0,
		AllowSelfHost:  true,
		AllowedRegions: nil, // nil = all regions
	},
}

// LimitsForTier returns the limits for the given tier.
// Falls back to Free limits for unknown tiers.
func LimitsForTier(tier webx.PlanTier) webx.PlanLimits {
	if l, ok := PlanLimits[tier]; ok {
		return l
	}
	return PlanLimits[PlanFree]
}
