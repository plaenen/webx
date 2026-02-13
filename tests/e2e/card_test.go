package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestCardPage_CardsRender(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/card", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	cards := page.Locator(".card")
	count, err := cards.Count()
	if err != nil {
		t.Fatalf("count cards: %v", err)
	}
	if count == 0 {
		t.Error("no card components found on card page")
	}
}

func TestCardPage_TitlesPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/card", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	titles := page.Locator(".card-title")
	count, err := titles.Count()
	if err != nil {
		t.Fatalf("count titles: %v", err)
	}
	if count == 0 {
		t.Error("no card titles found")
	}
}

func TestCardPage_ActionsPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/card", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	actions := page.Locator(".card-actions")
	count, err := actions.Count()
	if err != nil {
		t.Fatalf("count actions: %v", err)
	}
	if count == 0 {
		t.Error("no card actions found")
	}
}
