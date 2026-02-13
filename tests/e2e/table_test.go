package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestTablePage_TablesRender(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/table", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	tables := page.Locator("table.table")
	count, err := tables.Count()
	if err != nil {
		t.Fatalf("count tables: %v", err)
	}
	if count == 0 {
		t.Error("no table components found on table page")
	}
}

func TestTablePage_ZebraPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/table", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	zebra := page.Locator(".table-zebra")
	count, err := zebra.Count()
	if err != nil {
		t.Fatalf("count table-zebra: %v", err)
	}
	if count == 0 {
		t.Error("no table-zebra table found")
	}
}

func TestTablePage_SizesPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/table", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	sizes := []string{"table-xs", "table-sm", "table-lg"}
	for _, s := range sizes {
		loc := page.Locator("." + s)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", s, err)
		}
		if count == 0 {
			t.Errorf("no %s table found", s)
		}
	}
}

func TestTablePage_PinRowsPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/table", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	pinRows := page.Locator(".table-pin-rows")
	count, err := pinRows.Count()
	if err != nil {
		t.Fatalf("count table-pin-rows: %v", err)
	}
	if count == 0 {
		t.Error("no table-pin-rows table found")
	}
}
