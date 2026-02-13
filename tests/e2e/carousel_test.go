package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestCarouselPage_CarouselRender(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/carousel", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	containers := page.Locator("main .carousel")
	count, err := containers.Count()
	if err != nil {
		t.Fatalf("count carousel containers: %v", err)
	}
	if count == 0 {
		t.Error("no carousel containers found on carousel page")
	}
}

func TestCarouselPage_CarouselItemsPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/carousel", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	items := page.Locator("main .carousel-item")
	count, err := items.Count()
	if err != nil {
		t.Fatalf("count carousel-item: %v", err)
	}
	if count < 3 {
		t.Errorf("expected at least 3 carousel items, got %d", count)
	}
}

func TestCarouselPage_SnapVariants(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/carousel", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	snaps := []string{"carousel-start", "carousel-center", "carousel-end"}
	for _, s := range snaps {
		loc := page.Locator("main .carousel." + s)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", s, err)
		}
		if count == 0 {
			t.Errorf("no %s carousel found", s)
		}
	}
}

func TestCarouselPage_VerticalCarousel(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/carousel", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	loc := page.Locator("main .carousel.carousel-vertical")
	count, err := loc.Count()
	if err != nil {
		t.Fatalf("count carousel-vertical: %v", err)
	}
	if count == 0 {
		t.Error("no vertical carousel found")
	}
}
