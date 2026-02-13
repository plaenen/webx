package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestFabPage_FabsRender(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/fab", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	fabs := page.Locator("main .fab")
	count, err := fabs.Count()
	if err != nil {
		t.Fatalf("count fabs: %v", err)
	}
	if count == 0 {
		t.Error("no fab components found on fab page")
	}
}

func TestFabPage_FlowerPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/fab", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	loc := page.Locator("main .fab.fab-flower")
	count, err := loc.Count()
	if err != nil {
		t.Fatalf("count fab-flower: %v", err)
	}
	if count == 0 {
		t.Error("no flower fab found")
	}
}

func TestFabPage_CloseButtonPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/fab", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	loc := page.Locator("main .fab-close")
	count, err := loc.Count()
	if err != nil {
		t.Fatalf("count fab-close: %v", err)
	}
	if count == 0 {
		t.Error("no fab-close element found")
	}
}

func TestFabPage_SpeedDialButtons(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/fab", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Each fab should contain multiple buttons
	buttons := page.Locator("main .fab .btn")
	count, err := buttons.Count()
	if err != nil {
		t.Fatalf("count fab buttons: %v", err)
	}
	if count < 4 {
		t.Errorf("expected at least 4 buttons in fabs, got %d", count)
	}
}
