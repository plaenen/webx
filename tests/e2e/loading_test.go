package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestLoadingPage_LoadingRender(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/loading", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	containers := page.Locator("main .loading")
	count, err := containers.Count()
	if err != nil {
		t.Fatalf("count loading elements: %v", err)
	}
	if count == 0 {
		t.Error("no loading elements found on loading page")
	}
}

func TestLoadingPage_TypeVariants(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/loading", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	types := []string{"loading-spinner", "loading-dots", "loading-ring", "loading-ball", "loading-bars", "loading-infinity"}
	for _, typ := range types {
		loc := page.Locator("main .loading." + typ)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", typ, err)
		}
		if count == 0 {
			t.Errorf("no %s loading found", typ)
		}
	}
}

func TestLoadingPage_SizeVariants(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/loading", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	sizes := []string{"loading-xs", "loading-sm", "loading-md", "loading-lg", "loading-xl"}
	for _, s := range sizes {
		loc := page.Locator("main .loading." + s)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", s, err)
		}
		if count == 0 {
			t.Errorf("no %s loading found", s)
		}
	}
}
