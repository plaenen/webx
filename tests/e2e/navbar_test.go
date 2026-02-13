package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestNavbarPage_NavbarsRender(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/navbar", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Count navbars on the showcase page (excluding the layout navbar)
	navbars := page.Locator("main .navbar")
	count, err := navbars.Count()
	if err != nil {
		t.Fatalf("count navbars: %v", err)
	}
	if count == 0 {
		t.Error("no navbar components found on navbar page")
	}
}

func TestNavbarPage_SectionsPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/navbar", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	sections := []string{"navbar-start", "navbar-center", "navbar-end"}
	for _, s := range sections {
		loc := page.Locator("main ." + s)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", s, err)
		}
		if count == 0 {
			t.Errorf("no %s section found", s)
		}
	}
}

func TestNavbarPage_ColoredPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/navbar", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	colors := []string{"bg-neutral", "bg-primary"}
	for _, c := range colors {
		loc := page.Locator("main .navbar." + c)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", c, err)
		}
		if count == 0 {
			t.Errorf("no %s navbar found", c)
		}
	}
}
