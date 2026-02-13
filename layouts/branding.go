package layouts

// AppBranding configures application branding used across layouts.
// This is shared between Dashboard, Auth pages, and other layouts.
type AppBranding struct {
	// Name is the application name shown in headers/sidebars
	Name string
	// LogoFullUrl is the URL to the full logo image
	// e.g. https://example.com/logo-full.png
	LogoFullUrl string
	// Href is the logo click destination (defaults to "/")
	Href string
}

// DefaultLogoUrl returns the logo URL to use, defaulting to "/assets/logo.png" if not set.
func (b AppBranding) DefaultLogoUrl() string {
	if b.LogoFullUrl == "" {
		return "/assets/logo.png"
	}
	return b.LogoFullUrl
}

// DefaultHref returns the href to use, defaulting to "/" if not set.
func (b AppBranding) DefaultHref() string {
	if b.Href == "" {
		return "/"
	}
	return b.Href
}

// DefaultName returns the name to use, defaulting to "App" if not set.
func (b AppBranding) DefaultName() string {
	if b.Name == "" {
		return "App"
	}
	return b.Name
}
