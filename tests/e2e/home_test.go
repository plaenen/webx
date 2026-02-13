package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestHomePage_LoadsWithTitle(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL); err != nil {
		t.Fatalf("goto: %v", err)
	}
	title, err := page.Title()
	if err != nil {
		t.Fatalf("title: %v", err)
	}
	if title == "" {
		t.Error("page title is empty")
	}
}

func TestHomePage_NavigationLinksVisible(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL, pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// The sidebar menu should have navigation links.
	links := page.Locator("aside a, .drawer-side a, nav a")
	count, err := links.Count()
	if err != nil {
		t.Fatalf("count links: %v", err)
	}
	if count == 0 {
		t.Error("no navigation links found")
	}
}

func TestHomePage_ReadmeRendered(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL, pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	prose := page.Locator(".prose")
	count, err := prose.Count()
	if err != nil {
		t.Fatalf("count prose: %v", err)
	}
	if count == 0 {
		t.Error("no rendered markdown (prose) found on home page")
	}
}
