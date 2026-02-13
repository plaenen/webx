package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestJoinPage_JoinsRender(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/join", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	joins := page.Locator("main .join")
	count, err := joins.Count()
	if err != nil {
		t.Fatalf("count joins: %v", err)
	}
	if count == 0 {
		t.Error("no join components found on join page")
	}
}

func TestJoinPage_VerticalPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/join", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	loc := page.Locator("main .join.join-vertical")
	count, err := loc.Count()
	if err != nil {
		t.Fatalf("count join-vertical: %v", err)
	}
	if count == 0 {
		t.Error("no vertical join found")
	}
}

func TestJoinPage_JoinItemsPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/join", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	items := page.Locator("main .join-item")
	count, err := items.Count()
	if err != nil {
		t.Fatalf("count join-items: %v", err)
	}
	if count < 3 {
		t.Errorf("expected at least 3 join-items, got %d", count)
	}
}

func TestJoinPage_InputWithButton(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/join", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	input := page.Locator("main .join input.join-item[type='text']")
	count, err := input.Count()
	if err != nil {
		t.Fatalf("count input join-items: %v", err)
	}
	if count == 0 {
		t.Error("no input join-item found")
	}
}
