package icon

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

func renderComponent(t *testing.T, c templ.Component) string {
	t.Helper()
	var buf bytes.Buffer
	err := c.Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	return buf.String()
}

func TestIconTypeIsCallable(t *testing.T) {
	var i IconType = Icon("circle")
	c := i()
	result := renderComponent(t, c)
	if !strings.Contains(result, "<svg") {
		t.Fatalf("expected SVG output, got: %s", result)
	}
}

func TestIconTypeWithProps(t *testing.T) {
	var i IconType = Icon("circle")
	c := i(Props{Size: 16, Class: "my-icon"})
	result := renderComponent(t, c)

	if !strings.Contains(result, `width="16"`) {
		t.Errorf("expected width=16, got: %s", result)
	}
	if !strings.Contains(result, `height="16"`) {
		t.Errorf("expected height=16, got: %s", result)
	}
	if !strings.Contains(result, `class="my-icon"`) {
		t.Errorf("expected class my-icon, got: %s", result)
	}
}

func TestIconTypeAssignableFromDef(t *testing.T) {
	// Verify that generated icon defs are assignable to IconType.
	var i IconType = Circle
	c := i()
	result := renderComponent(t, c)
	if !strings.Contains(result, "<svg") {
		t.Fatalf("expected SVG output, got: %s", result)
	}
}

func TestIconTypeInStruct(t *testing.T) {
	// Simulates a component Props struct using IconType.
	type ComponentProps struct {
		Icon IconType
	}
	p := ComponentProps{Icon: Circle}
	c := p.Icon(Props{Size: 20})
	result := renderComponent(t, c)
	if !strings.Contains(result, `width="20"`) {
		t.Errorf("expected width=20, got: %s", result)
	}
}

func TestIconReturnsTemplComponent(t *testing.T) {
	// The existing pattern: call icon def to get templ.Component.
	var c templ.Component = Circle()
	result := renderComponent(t, c)
	if !strings.Contains(result, "<svg") {
		t.Fatalf("expected SVG output, got: %s", result)
	}
}

func TestIconDefaultProps(t *testing.T) {
	result := renderComponent(t, Circle())

	if !strings.Contains(result, `width="24"`) {
		t.Errorf("expected default width=24, got: %s", result)
	}
	if !strings.Contains(result, `fill="none"`) {
		t.Errorf("expected default fill=none, got: %s", result)
	}
	if !strings.Contains(result, `stroke="currentColor"`) {
		t.Errorf("expected default stroke=currentColor, got: %s", result)
	}
	if !strings.Contains(result, `stroke-width="2"`) {
		t.Errorf("expected default stroke-width=2, got: %s", result)
	}
}

func TestIconCustomStroke(t *testing.T) {
	result := renderComponent(t, Circle(Props{Stroke: "red", StrokeWidth: "3"}))

	if !strings.Contains(result, `stroke="red"`) {
		t.Errorf("expected stroke=red, got: %s", result)
	}
	if !strings.Contains(result, `stroke-width="3"`) {
		t.Errorf("expected stroke-width=3, got: %s", result)
	}
}

func TestIconColorFallsBackToStroke(t *testing.T) {
	result := renderComponent(t, Circle(Props{Color: "blue"}))

	if !strings.Contains(result, `stroke="blue"`) {
		t.Errorf("expected stroke=blue from Color fallback, got: %s", result)
	}
}

func TestIconStrokeOverridesColor(t *testing.T) {
	result := renderComponent(t, Circle(Props{Color: "blue", Stroke: "green"}))

	if !strings.Contains(result, `stroke="green"`) {
		t.Errorf("expected stroke=green (Stroke overrides Color), got: %s", result)
	}
}

func TestIconCaching(t *testing.T) {
	// Clear cache for this test.
	iconMutex.Lock()
	iconContents = make(map[string]string)
	iconMutex.Unlock()

	r1 := renderComponent(t, Circle())
	r2 := renderComponent(t, Circle())
	if r1 != r2 {
		t.Errorf("cached result differs:\nfirst:  %s\nsecond: %s", r1, r2)
	}

	// Different props should produce different output.
	r3 := renderComponent(t, Circle(Props{Size: 16}))
	if r1 == r3 {
		t.Error("expected different output for different props")
	}
}

func TestIconNotFound(t *testing.T) {
	i := Icon("nonexistent-icon-xyz")
	c := i()
	var buf bytes.Buffer
	err := c.Render(context.Background(), &buf)
	if err == nil {
		t.Fatal("expected error for nonexistent icon")
	}
	if !strings.Contains(err.Error(), "nonexistent-icon-xyz") {
		t.Errorf("error should mention icon name, got: %v", err)
	}
}
