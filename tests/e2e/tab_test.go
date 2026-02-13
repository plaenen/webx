package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestTabPage_TabsRender(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/tab", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	containers := page.Locator("main .tabs")
	count, err := containers.Count()
	if err != nil {
		t.Fatalf("count tabs containers: %v", err)
	}
	if count == 0 {
		t.Error("no tabs containers found on tab page")
	}
}

func TestTabPage_VariantStyles(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/tab", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	variants := []string{"tabs-border", "tabs-lift", "tabs-box"}
	for _, v := range variants {
		loc := page.Locator("main .tabs." + v)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", v, err)
		}
		if count == 0 {
			t.Errorf("no %s tabs found", v)
		}
	}
}

func TestTabPage_ActiveTab(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/tab", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	active := page.Locator("main .tabs .tab.tab-active")
	count, err := active.Count()
	if err != nil {
		t.Fatalf("count tab-active: %v", err)
	}
	if count == 0 {
		t.Error("no active tab found")
	}
}

func TestTabPage_RadioTabsWithContent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/tab", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	radios := page.Locator("main .tabs input[type='radio'].tab")
	count, err := radios.Count()
	if err != nil {
		t.Fatalf("count radio tabs: %v", err)
	}
	if count < 3 {
		t.Errorf("expected at least 3 radio tabs, got %d", count)
	}

	content := page.Locator("main .tabs .tab-content")
	ccount, err := content.Count()
	if err != nil {
		t.Fatalf("count tab-content: %v", err)
	}
	if ccount < 3 {
		t.Errorf("expected at least 3 tab-content panels, got %d", ccount)
	}
}

func TestTabPage_SizeVariants(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/tab", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	sizes := []string{"tabs-xs", "tabs-lg"}
	for _, s := range sizes {
		loc := page.Locator("main .tabs." + s)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", s, err)
		}
		if count == 0 {
			t.Errorf("no %s tabs found", s)
		}
	}
}
