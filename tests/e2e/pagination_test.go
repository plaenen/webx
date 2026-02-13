package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestPaginationPage_Render(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/pagination", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	containers := page.Locator("main .join")
	count, err := containers.Count()
	if err != nil {
		t.Fatalf("count pagination containers: %v", err)
	}
	if count == 0 {
		t.Error("no pagination containers found")
	}
}

func TestPaginationPage_ActiveButton(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/pagination", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	active := page.Locator("main .join button.btn-active")
	count, err := active.Count()
	if err != nil {
		t.Fatalf("count btn-active: %v", err)
	}
	if count == 0 {
		t.Error("no active pagination button found")
	}
}

func TestPaginationPage_DisabledButton(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/pagination", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	disabled := page.Locator("main .join button.btn-disabled")
	count, err := disabled.Count()
	if err != nil {
		t.Fatalf("count btn-disabled: %v", err)
	}
	if count == 0 {
		t.Error("no disabled pagination button found")
	}
}

func TestPaginationPage_JoinItemButtons(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/pagination", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	items := page.Locator("main .join button.join-item")
	count, err := items.Count()
	if err != nil {
		t.Fatalf("count join-item buttons: %v", err)
	}
	if count < 4 {
		t.Errorf("expected at least 4 join-item buttons, got %d", count)
	}
}
