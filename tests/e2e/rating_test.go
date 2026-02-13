package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestRatingPage_Renders(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/rating", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	containers := page.Locator("main div.rating")
	count, err := containers.Count()
	if err != nil {
		t.Fatalf("count rating containers: %v", err)
	}
	if count == 0 {
		t.Error("no rating containers found")
	}
}

func TestRatingPage_MaskVariants(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/rating", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	masks := []string{"mask-star", "mask-star-2", "mask-heart"}
	for _, m := range masks {
		loc := page.Locator("main .rating input." + m)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", m, err)
		}
		if count == 0 {
			t.Errorf("no %s rating inputs found", m)
		}
	}
}

func TestRatingPage_SizeVariants(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/rating", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	sizes := []string{"rating-xs", "rating-sm", "rating-lg", "rating-xl"}
	for _, s := range sizes {
		loc := page.Locator("main div.rating." + s)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", s, err)
		}
		if count == 0 {
			t.Errorf("no %s rating found", s)
		}
	}
}

func TestRatingPage_HalfStars(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/rating", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	half := page.Locator("main div.rating.rating-half")
	count, err := half.Count()
	if err != nil {
		t.Fatalf("count rating-half: %v", err)
	}
	if count == 0 {
		t.Error("no half-star rating found")
	}

	halves := page.Locator("main div.rating.rating-half input.mask-half-1")
	hcount, err := halves.Count()
	if err != nil {
		t.Fatalf("count mask-half-1: %v", err)
	}
	if hcount < 5 {
		t.Errorf("expected at least 5 half-1 inputs, got %d", hcount)
	}
}
