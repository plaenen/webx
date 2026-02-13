package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestLinkPage_LinksRender(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/link", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	links := page.Locator("main a.link")
	count, err := links.Count()
	if err != nil {
		t.Fatalf("count links: %v", err)
	}
	if count == 0 {
		t.Error("no link components found on link page")
	}
}

func TestLinkPage_ColorVariants(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/link", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	colors := []string{
		"link-neutral", "link-primary", "link-secondary", "link-accent",
		"link-success", "link-info", "link-warning", "link-error",
	}
	for _, c := range colors {
		loc := page.Locator("main a." + c)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", c, err)
		}
		if count == 0 {
			t.Errorf("no %s link found", c)
		}
	}
}

func TestLinkPage_HoverVariant(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/link", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	loc := page.Locator("main a.link-hover")
	count, err := loc.Count()
	if err != nil {
		t.Fatalf("count link-hover: %v", err)
	}
	if count == 0 {
		t.Error("no link-hover element found")
	}
}
