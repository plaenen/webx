package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestDockPage_DockRender(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/dock", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	containers := page.Locator("main .dock")
	count, err := containers.Count()
	if err != nil {
		t.Fatalf("count dock elements: %v", err)
	}
	if count == 0 {
		t.Error("no dock elements found on dock page")
	}
}

func TestDockPage_ActiveItem(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/dock", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	active := page.Locator("main .dock .dock-active")
	count, err := active.Count()
	if err != nil {
		t.Fatalf("count dock-active: %v", err)
	}
	if count == 0 {
		t.Error("no dock-active items found")
	}
}

func TestDockPage_SizeVariants(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/dock", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	sizes := []string{"dock-xs", "dock-md", "dock-lg"}
	for _, s := range sizes {
		loc := page.Locator("main .dock." + s)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", s, err)
		}
		if count == 0 {
			t.Errorf("no %s dock found", s)
		}
	}
}

func TestDockPage_DockLabels(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/dock", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	labels := page.Locator("main .dock .dock-label")
	count, err := labels.Count()
	if err != nil {
		t.Fatalf("count dock-label: %v", err)
	}
	if count < 3 {
		t.Errorf("expected at least 3 dock labels, got %d", count)
	}
}
