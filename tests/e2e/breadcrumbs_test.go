package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestBreadcrumbsPage_RendersBreadcrumbs(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/breadcrumbs", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Should have multiple breadcrumbs containers.
	containers := page.Locator("div.breadcrumbs")
	count, err := containers.Count()
	if err != nil {
		t.Fatalf("count breadcrumbs: %v", err)
	}
	if count < 3 {
		t.Errorf("expected at least 3 breadcrumbs containers, got %d", count)
	}
}

func TestBreadcrumbsPage_ItemsWithLinks(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/breadcrumbs", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// The first breadcrumbs container should have items with links.
	first := page.Locator("div.breadcrumbs").First()
	links := first.Locator("li a")
	linkCount, err := links.Count()
	if err != nil {
		t.Fatalf("count links: %v", err)
	}
	if linkCount < 2 {
		t.Errorf("expected at least 2 links in first breadcrumbs, got %d", linkCount)
	}

	// The last item should not be a link (current page).
	items := first.Locator("li")
	itemCount, err := items.Count()
	if err != nil {
		t.Fatalf("count items: %v", err)
	}
	lastItem := items.Nth(itemCount - 1)
	lastLinks, err := lastItem.Locator("a").Count()
	if err != nil {
		t.Fatalf("count last item links: %v", err)
	}
	if lastLinks != 0 {
		t.Errorf("expected last breadcrumb item to not be a link, got %d links", lastLinks)
	}
}

func TestBreadcrumbsPage_MaxWidthConstraint(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/breadcrumbs", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// The long path example should have max-w-xs class.
	maxWidth := page.Locator("div.breadcrumbs.max-w-xs")
	count, err := maxWidth.Count()
	if err != nil {
		t.Fatalf("count max-width breadcrumbs: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 breadcrumbs with max-w-xs, got %d", count)
	}
}

func TestBreadcrumbsPage_DeepNavigation(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/breadcrumbs", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// The deep navigation example should have 5 items.
	containers := page.Locator("div.breadcrumbs")
	last := containers.Last()
	items := last.Locator("li")
	count, err := items.Count()
	if err != nil {
		t.Fatalf("count deep nav items: %v", err)
	}
	if count != 5 {
		t.Errorf("expected 5 items in deep navigation breadcrumbs, got %d", count)
	}
}
