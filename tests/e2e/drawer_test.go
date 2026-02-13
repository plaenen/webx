package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestShowcaseLayout_HamburgerTogglesSidebar(t *testing.T) {
	// Use a mobile viewport so the hamburger is visible (hidden on lg:).
	page, err := browser.NewPage(pw.BrowserNewPageOptions{
		Viewport: &pw.Size{Width: 375, Height: 812},
	})
	if err != nil {
		t.Fatalf("new page: %v", err)
	}
	t.Cleanup(func() { page.Close() })

	if _, err := page.Goto(baseURL, pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Sidebar should be hidden on mobile.
	sidebar := page.Locator(".drawer-side aside")
	visible, err := sidebar.IsVisible()
	if err != nil {
		t.Fatalf("sidebar visible check: %v", err)
	}
	if visible {
		t.Fatal("sidebar should be hidden on mobile before toggle")
	}

	// Click the hamburger button.
	hamburger := page.Locator("button[aria-label='Toggle menu']")
	if err := hamburger.Click(); err != nil {
		t.Fatalf("click hamburger: %v", err)
	}

	// Sidebar should become visible.
	if err := sidebar.WaitFor(pw.LocatorWaitForOptions{
		State: pw.WaitForSelectorStateVisible,
	}); err != nil {
		t.Fatalf("sidebar did not open after hamburger click: %v", err)
	}

	// Click the overlay to close. The overlay sits behind the aside in a CSS
	// grid, so click to the right of the aside where the overlay is exposed.
	overlay := page.Locator(".drawer-overlay")
	asideBox, err := sidebar.BoundingBox()
	if err != nil {
		t.Fatalf("aside bounding box: %v", err)
	}
	if err := overlay.Click(pw.LocatorClickOptions{
		Position: &pw.Position{X: asideBox.Width + 20, Y: asideBox.Height / 2},
	}); err != nil {
		t.Fatalf("click overlay: %v", err)
	}

	// Sidebar should hide again.
	if err := sidebar.WaitFor(pw.LocatorWaitForOptions{
		State: pw.WaitForSelectorStateHidden,
	}); err != nil {
		t.Fatalf("sidebar did not close after overlay click: %v", err)
	}
}

func TestDrawerPage_ToggleOpensPanel(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/drawer", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Click the "Toggle Drawer" button in the left-side inline demo.
	toggleBtn := page.Locator("button", pw.PageLocatorOptions{
		HasText: "Toggle Drawer",
	}).First()
	if err := toggleBtn.Click(); err != nil {
		t.Fatalf("click toggle: %v", err)
	}

	// The left slide-in panel should become visible.
	panel := page.Locator(".border-r.absolute")
	if err := panel.WaitFor(pw.LocatorWaitForOptions{
		State: pw.WaitForSelectorStateVisible,
	}); err != nil {
		t.Fatalf("panel did not become visible: %v", err)
	}
}

func TestDrawerPage_OverlayClosesPanel(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/drawer", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Open the left drawer first.
	toggleBtn := page.Locator("button", pw.PageLocatorOptions{
		HasText: "Toggle Drawer",
	}).First()
	if err := toggleBtn.Click(); err != nil {
		t.Fatalf("click toggle: %v", err)
	}

	// Wait for overlay to appear (scope to the first demo container).
	overlay := page.Locator(".bg-black\\/30").First()
	if err := overlay.WaitFor(pw.LocatorWaitForOptions{
		State: pw.WaitForSelectorStateVisible,
	}); err != nil {
		t.Fatalf("overlay not visible: %v", err)
	}

	// Click the overlay to close.
	if err := overlay.Click(); err != nil {
		t.Fatalf("click overlay: %v", err)
	}

	// Overlay should become hidden.
	if err := overlay.WaitFor(pw.LocatorWaitForOptions{
		State: pw.WaitForSelectorStateHidden,
	}); err != nil {
		t.Fatalf("overlay did not hide after click: %v", err)
	}
}

func TestDrawerPage_RightSideDrawerToggles(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/drawer", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	toggleBtn := page.Locator("button", pw.PageLocatorOptions{
		HasText: "Toggle Right Drawer",
	})
	if err := toggleBtn.Click(); err != nil {
		t.Fatalf("click toggle: %v", err)
	}

	panel := page.Locator(".border-l.absolute")
	if err := panel.WaitFor(pw.LocatorWaitForOptions{
		State: pw.WaitForSelectorStateVisible,
	}); err != nil {
		t.Fatalf("right panel did not become visible: %v", err)
	}
}

func TestDrawerPage_StateIndicatorUpdates(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/drawer", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// The first state badge should initially show "closed".
	badge := page.Locator("code.badge").First()
	text, err := badge.TextContent()
	if err != nil {
		t.Fatalf("badge text: %v", err)
	}
	if text != "closed" {
		t.Errorf("initial badge text = %q, want %q", text, "closed")
	}

	// Click toggle - badge should show "open".
	toggleBtn := page.Locator("button", pw.PageLocatorOptions{
		HasText: "Toggle Drawer",
	}).First()
	if err := toggleBtn.Click(); err != nil {
		t.Fatalf("click toggle: %v", err)
	}
	if err := page.Locator("code.badge", pw.PageLocatorOptions{
		HasText: "open",
	}).First().WaitFor(pw.LocatorWaitForOptions{
		State: pw.WaitForSelectorStateVisible,
	}); err != nil {
		t.Fatalf("badge did not change to 'open': %v", err)
	}
}
