package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestStepsPage_StepsRender(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/steps", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	containers := page.Locator("main .steps")
	count, err := containers.Count()
	if err != nil {
		t.Fatalf("count steps containers: %v", err)
	}
	if count == 0 {
		t.Error("no steps containers found on steps page")
	}
}

func TestStepsPage_StepItemsPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/steps", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	items := page.Locator("main .step")
	count, err := items.Count()
	if err != nil {
		t.Fatalf("count step items: %v", err)
	}
	if count < 3 {
		t.Errorf("expected at least 3 step items, got %d", count)
	}
}

func TestStepsPage_DirectionVariants(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/steps", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	directions := []string{"steps-horizontal", "steps-vertical"}
	for _, d := range directions {
		loc := page.Locator("main .steps." + d)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", d, err)
		}
		if count == 0 {
			t.Errorf("no %s steps found", d)
		}
	}
}

func TestStepsPage_ColorVariants(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/steps", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	colors := []string{"step-primary", "step-success", "step-warning", "step-error", "step-info"}
	for _, c := range colors {
		loc := page.Locator("main .step." + c)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", c, err)
		}
		if count == 0 {
			t.Errorf("no %s step found", c)
		}
	}
}
