package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestFilterPage_Renders(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/filter", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	containers := page.Locator("main form.filter")
	count, err := containers.Count()
	if err != nil {
		t.Fatalf("count filter containers: %v", err)
	}
	if count == 0 {
		t.Error("no filter containers found")
	}
}

func TestFilterPage_ResetButtons(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/filter", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	resets := page.Locator("main form.filter input[type='reset']")
	count, err := resets.Count()
	if err != nil {
		t.Fatalf("count reset buttons: %v", err)
	}
	if count < 3 {
		t.Errorf("expected at least 3 reset buttons, got %d", count)
	}
}

func TestFilterPage_RadioOptions(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/filter", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	radios := page.Locator("main form.filter input[type='radio'].btn")
	count, err := radios.Count()
	if err != nil {
		t.Fatalf("count radio options: %v", err)
	}
	if count < 4 {
		t.Errorf("expected at least 4 radio options, got %d", count)
	}
}

func TestFilterPage_CheckedDefault(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/filter", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	checked := page.Locator("main form.filter input[type='radio'][checked]")
	count, err := checked.Count()
	if err != nil {
		t.Fatalf("count checked radios: %v", err)
	}
	if count == 0 {
		t.Error("no default-checked radio found")
	}
}
