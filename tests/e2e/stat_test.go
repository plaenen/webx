package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestStatPage_StatsRender(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/stat", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	stats := page.Locator(".stats")
	count, err := stats.Count()
	if err != nil {
		t.Fatalf("count stats: %v", err)
	}
	if count == 0 {
		t.Error("no stats containers found on stat page")
	}
}

func TestStatPage_StatPartsPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/stat", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	parts := []string{"stat-title", "stat-value", "stat-desc"}
	for _, p := range parts {
		loc := page.Locator("." + p)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", p, err)
		}
		if count == 0 {
			t.Errorf("no %s elements found", p)
		}
	}
}

func TestStatPage_FigurePresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/stat", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	fig := page.Locator(".stat-figure")
	count, err := fig.Count()
	if err != nil {
		t.Fatalf("count stat-figure: %v", err)
	}
	if count == 0 {
		t.Error("no stat-figure elements found")
	}
}

func TestStatPage_ActionsPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/stat", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	actions := page.Locator(".stat-actions")
	count, err := actions.Count()
	if err != nil {
		t.Fatalf("count stat-actions: %v", err)
	}
	if count == 0 {
		t.Error("no stat-actions elements found")
	}
}

func TestStatPage_VerticalPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/stat", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	vertical := page.Locator(".stats-vertical")
	count, err := vertical.Count()
	if err != nil {
		t.Fatalf("count stats-vertical: %v", err)
	}
	if count == 0 {
		t.Error("no stats-vertical container found")
	}
}
