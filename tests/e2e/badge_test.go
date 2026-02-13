package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestBadgePage_BadgesRender(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/badge", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	badges := page.Locator(".badge")
	count, err := badges.Count()
	if err != nil {
		t.Fatalf("count badges: %v", err)
	}
	if count == 0 {
		t.Error("no badge components found on badge page")
	}
}

func TestBadgePage_VariantsPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/badge", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	variants := []string{
		"badge-primary",
		"badge-secondary",
		"badge-accent",
		"badge-info",
		"badge-success",
		"badge-warning",
		"badge-error",
		"badge-neutral",
	}
	for _, v := range variants {
		loc := page.Locator("." + v)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", v, err)
		}
		if count == 0 {
			t.Errorf("no %s badge found", v)
		}
	}
}

func TestBadgePage_SizesPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/badge", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	sizes := []string{
		"badge-xs",
		"badge-sm",
		"badge-md",
		"badge-lg",
		"badge-xl",
	}
	for _, s := range sizes {
		loc := page.Locator("." + s)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", s, err)
		}
		if count == 0 {
			t.Errorf("no %s badge found", s)
		}
	}
}

func TestBadgePage_StylesPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/badge", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	styles := []string{
		"badge-outline",
		"badge-dash",
		"badge-soft",
		"badge-ghost",
	}
	for _, s := range styles {
		loc := page.Locator("." + s)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", s, err)
		}
		if count == 0 {
			t.Errorf("no %s badge found", s)
		}
	}
}
