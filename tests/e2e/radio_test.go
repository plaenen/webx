package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestRadioPage_Renders(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/radio", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	radios := page.Locator("main input.radio")
	count, err := radios.Count()
	if err != nil {
		t.Fatalf("count radios: %v", err)
	}
	if count == 0 {
		t.Error("no radio elements found")
	}
}

func TestRadioPage_VariantStyles(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/radio", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	variants := []string{
		"radio-primary",
		"radio-secondary",
		"radio-accent",
		"radio-neutral",
		"radio-info",
		"radio-success",
		"radio-warning",
		"radio-error",
	}
	for _, v := range variants {
		loc := page.Locator("main input.radio." + v)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", v, err)
		}
		if count == 0 {
			t.Errorf("no %s radio found", v)
		}
	}
}

func TestRadioPage_SizeVariants(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/radio", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	sizes := []string{"radio-xs", "radio-sm", "radio-md", "radio-lg", "radio-xl"}
	for _, s := range sizes {
		loc := page.Locator("main input.radio." + s)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", s, err)
		}
		if count == 0 {
			t.Errorf("no %s radio found", s)
		}
	}
}

func TestRadioPage_Disabled(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/radio", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	disabled := page.Locator("main input.radio[disabled]")
	count, err := disabled.Count()
	if err != nil {
		t.Fatalf("count disabled: %v", err)
	}
	if count < 2 {
		t.Errorf("expected at least 2 disabled radios, got %d", count)
	}
}
