package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestStatusPage_RendersColorVariants(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/status", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	for _, variant := range []string{
		"status-neutral",
		"status-primary",
		"status-secondary",
		"status-accent",
		"status-info",
		"status-success",
		"status-warning",
		"status-error",
	} {
		loc := page.Locator("div.status." + variant)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", variant, err)
		}
		if count < 1 {
			t.Errorf("expected at least 1 element with class %s, got %d", variant, count)
		}
	}
}

func TestStatusPage_RendersSizes(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/status", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	for _, size := range []string{
		"status-xs",
		"status-sm",
		"status-md",
		"status-lg",
		"status-xl",
	} {
		loc := page.Locator("div.status." + size)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", size, err)
		}
		if count < 1 {
			t.Errorf("expected at least 1 element with class %s, got %d", size, count)
		}
	}
}

func TestStatusPage_PingAnimationUsesGrid(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/status", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Ping animation should be wrapped in an inline-grid container with two status elements.
	grid := page.Locator("div.inline-grid")
	count, err := grid.Count()
	if err != nil {
		t.Fatalf("count grid wrappers: %v", err)
	}
	if count < 1 {
		t.Fatalf("expected at least 1 inline-grid wrapper for ping, got %d", count)
	}

	// The first grid should contain an animate-ping status.
	pingStatus := grid.First().Locator("div.status.animate-ping")
	pingCount, err := pingStatus.Count()
	if err != nil {
		t.Fatalf("count ping status: %v", err)
	}
	if pingCount != 1 {
		t.Errorf("expected 1 animate-ping status inside grid, got %d", pingCount)
	}
}

func TestStatusPage_BounceAnimation(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/status", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	bounce := page.Locator("div.status.animate-bounce")
	count, err := bounce.Count()
	if err != nil {
		t.Fatalf("count bounce status: %v", err)
	}
	if count < 1 {
		t.Errorf("expected at least 1 animate-bounce status, got %d", count)
	}
}
