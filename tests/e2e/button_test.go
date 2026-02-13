package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestButtonPage_VariantsRender(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/button", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	variants := []string{"btn-primary", "btn-secondary", "btn-accent"}
	for _, v := range variants {
		loc := page.Locator("." + v)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", v, err)
		}
		if count == 0 {
			t.Errorf("no button with class %s found", v)
		}
	}
}

func TestButtonPage_SizesRender(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/button", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	sizes := []string{"btn-sm", "btn-lg"}
	for _, s := range sizes {
		loc := page.Locator("." + s)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", s, err)
		}
		if count == 0 {
			t.Errorf("no button with class %s found", s)
		}
	}
}
