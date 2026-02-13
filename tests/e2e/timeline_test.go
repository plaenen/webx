package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestTimelinePage_TimelinesRender(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/timeline", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	timelines := page.Locator(".timeline")
	count, err := timelines.Count()
	if err != nil {
		t.Fatalf("count timelines: %v", err)
	}
	if count == 0 {
		t.Error("no timeline components found on timeline page")
	}
}

func TestTimelinePage_PartsPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/timeline", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	parts := []string{"timeline-start", "timeline-middle", "timeline-end"}
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

func TestTimelinePage_BoxPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/timeline", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	box := page.Locator(".timeline-box")
	count, err := box.Count()
	if err != nil {
		t.Fatalf("count timeline-box: %v", err)
	}
	if count == 0 {
		t.Error("no timeline-box elements found")
	}
}

func TestTimelinePage_VerticalPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/timeline", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	vertical := page.Locator(".timeline-vertical")
	count, err := vertical.Count()
	if err != nil {
		t.Fatalf("count timeline-vertical: %v", err)
	}
	if count == 0 {
		t.Error("no timeline-vertical container found")
	}
}

func TestTimelinePage_CompactPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/timeline", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	compact := page.Locator(".timeline-compact")
	count, err := compact.Count()
	if err != nil {
		t.Fatalf("count timeline-compact: %v", err)
	}
	if count == 0 {
		t.Error("no timeline-compact container found")
	}
}
