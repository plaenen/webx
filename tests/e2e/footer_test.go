package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestFooterPage_Renders(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/footer", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	footers := page.Locator("main footer.footer")
	count, err := footers.Count()
	if err != nil {
		t.Fatalf("count footers: %v", err)
	}
	if count == 0 {
		t.Error("no footer elements found")
	}
}

func TestFooterPage_TitleElements(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/footer", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	titles := page.Locator("main footer.footer .footer-title")
	count, err := titles.Count()
	if err != nil {
		t.Fatalf("count footer-title: %v", err)
	}
	if count < 3 {
		t.Errorf("expected at least 3 footer titles, got %d", count)
	}
}

func TestFooterPage_Horizontal(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/footer", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	horiz := page.Locator("main footer.footer.footer-horizontal")
	count, err := horiz.Count()
	if err != nil {
		t.Fatalf("count footer-horizontal: %v", err)
	}
	if count == 0 {
		t.Error("no horizontal footer found")
	}
}

func TestFooterPage_Center(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/footer", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	center := page.Locator("main footer.footer.footer-center")
	count, err := center.Count()
	if err != nil {
		t.Fatalf("count footer-center: %v", err)
	}
	if count == 0 {
		t.Error("no centered footer found")
	}
}
