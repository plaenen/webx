package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestMenuPage_MenuRender(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/menu", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	containers := page.Locator("main ul.menu")
	count, err := containers.Count()
	if err != nil {
		t.Fatalf("count menu containers: %v", err)
	}
	if count == 0 {
		t.Error("no menu containers found on menu page")
	}
}

func TestMenuPage_DirectionVariants(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/menu", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	loc := page.Locator("main ul.menu.menu-horizontal")
	count, err := loc.Count()
	if err != nil {
		t.Fatalf("count menu-horizontal: %v", err)
	}
	if count == 0 {
		t.Error("no horizontal menu found")
	}
}

func TestMenuPage_ActiveItem(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/menu", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	active := page.Locator("main ul.menu li.menu-active")
	count, err := active.Count()
	if err != nil {
		t.Fatalf("count menu-active: %v", err)
	}
	if count == 0 {
		t.Error("no active menu item found")
	}
}

func TestMenuPage_SizeVariants(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/menu", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	sizes := []string{"menu-xs", "menu-sm", "menu-lg"}
	for _, s := range sizes {
		loc := page.Locator("main ul.menu." + s)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", s, err)
		}
		if count == 0 {
			t.Errorf("no %s menu found", s)
		}
	}
}

func TestMenuPage_MenuTitle(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/menu", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	title := page.Locator("main ul.menu li.menu-title")
	count, err := title.Count()
	if err != nil {
		t.Fatalf("count menu-title: %v", err)
	}
	if count == 0 {
		t.Error("no menu-title found")
	}
}
