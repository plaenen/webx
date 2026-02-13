package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestFieldsetPage_FieldsetRender(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/fieldset", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	containers := page.Locator("main fieldset.fieldset")
	count, err := containers.Count()
	if err != nil {
		t.Fatalf("count fieldset elements: %v", err)
	}
	if count == 0 {
		t.Error("no fieldset elements found on fieldset page")
	}
}

func TestFieldsetPage_LegendPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/fieldset", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	legends := page.Locator("main fieldset.fieldset legend.fieldset-legend")
	count, err := legends.Count()
	if err != nil {
		t.Fatalf("count fieldset-legend: %v", err)
	}
	if count == 0 {
		t.Error("no fieldset-legend elements found")
	}
}

func TestFieldsetPage_LabelPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/fieldset", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	labels := page.Locator("main fieldset.fieldset p.label")
	count, err := labels.Count()
	if err != nil {
		t.Fatalf("count label: %v", err)
	}
	if count == 0 {
		t.Error("no label elements found inside fieldset")
	}
}
