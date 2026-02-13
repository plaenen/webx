package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestProgressPage_ProgressRender(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/progress", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	containers := page.Locator("main progress.progress")
	count, err := containers.Count()
	if err != nil {
		t.Fatalf("count progress elements: %v", err)
	}
	if count == 0 {
		t.Error("no progress elements found on progress page")
	}
}

func TestProgressPage_ColorVariants(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/progress", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	colors := []string{
		"progress-primary", "progress-secondary", "progress-accent",
		"progress-neutral", "progress-info", "progress-success",
		"progress-warning", "progress-error",
	}
	for _, c := range colors {
		loc := page.Locator("main progress." + c)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", c, err)
		}
		if count == 0 {
			t.Errorf("no %s progress found", c)
		}
	}
}

func TestProgressPage_HasValueAttribute(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/progress", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	withValue := page.Locator("main progress.progress[value]")
	count, err := withValue.Count()
	if err != nil {
		t.Fatalf("count progress with value: %v", err)
	}
	if count == 0 {
		t.Error("no progress elements with value attribute found")
	}
}
