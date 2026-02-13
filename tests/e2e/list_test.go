package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestListPage_ListsRender(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/list", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	lists := page.Locator(".list")
	count, err := lists.Count()
	if err != nil {
		t.Fatalf("count lists: %v", err)
	}
	if count == 0 {
		t.Error("no list components found on list page")
	}
}

func TestListPage_RowsPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/list", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	rows := page.Locator(".list-row")
	count, err := rows.Count()
	if err != nil {
		t.Fatalf("count list-row: %v", err)
	}
	if count == 0 {
		t.Error("no list-row elements found")
	}
}

func TestListPage_ColGrowPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/list", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	grow := page.Locator(".list-col-grow")
	count, err := grow.Count()
	if err != nil {
		t.Fatalf("count list-col-grow: %v", err)
	}
	if count == 0 {
		t.Error("no list-col-grow elements found")
	}
}

func TestListPage_ColWrapPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/list", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	wrap := page.Locator(".list-col-wrap")
	count, err := wrap.Count()
	if err != nil {
		t.Fatalf("count list-col-wrap: %v", err)
	}
	if count == 0 {
		t.Error("no list-col-wrap elements found")
	}
}
