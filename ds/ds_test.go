package ds_test

import (
	"testing"

	"github.com/plaenen/webx/ds"
)

func TestOn(t *testing.T) {
	attrs := ds.On("click", "$open = !$open")
	assertAttr(t, attrs, "data-on:click", "$open = !$open")
}

func TestOnClick(t *testing.T) {
	attrs := ds.OnClick("$open = !$open")
	assertAttr(t, attrs, "data-on:click", "$open = !$open")
}

func TestBind(t *testing.T) {
	attrs := ds.Bind("value")
	assertKey(t, attrs, "data-bind:value")
}

func TestClassToggle(t *testing.T) {
	attrs := ds.ClassToggle("drawer-open", "$drawer.open")
	assertAttr(t, attrs, "data-class:drawer-open", "$drawer.open")
}

func TestAttr(t *testing.T) {
	attrs := ds.Attr("disabled", "$loading")
	assertAttr(t, attrs, "data-attr:disabled", "$loading")
}

func TestStyle(t *testing.T) {
	attrs := ds.Style("opacity", "$visible ? 1 : 0")
	assertAttr(t, attrs, "data-style:opacity", "$visible ? 1 : 0")
}

func TestComputed(t *testing.T) {
	attrs := ds.Computed("fullName", "$first + ' ' + $last")
	assertAttr(t, attrs, "data-computed:fullName", "$first + ' ' + $last")
}

func TestIndicator(t *testing.T) {
	attrs := ds.Indicator("loading")
	assertKey(t, attrs, "data-indicator:loading")
}

func TestRef(t *testing.T) {
	attrs := ds.Ref("input")
	assertKey(t, attrs, "data-ref:input")
}

func TestSignals(t *testing.T) {
	attrs := ds.Signals(`{"open": false}`)
	assertAttr(t, attrs, "data-signals", `{"open": false}`)
}

func TestShow(t *testing.T) {
	attrs := ds.Show("$visible")
	assertAttr(t, attrs, "data-show", "$visible")
}

func TestText(t *testing.T) {
	attrs := ds.Text("$message")
	assertAttr(t, attrs, "data-text", "$message")
}

func TestClass(t *testing.T) {
	attrs := ds.Class("{'active': $active}")
	assertAttr(t, attrs, "data-class", "{'active': $active}")
}

func TestInit(t *testing.T) {
	attrs := ds.Init("console.log('init')")
	assertAttr(t, attrs, "data-init", "console.log('init')")
}

func TestEffect(t *testing.T) {
	attrs := ds.Effect("console.log($count)")
	assertAttr(t, attrs, "data-effect", "console.log($count)")
}

func TestMerge(t *testing.T) {
	merged := ds.Merge(
		ds.OnClick("$open = !$open"),
		ds.Attr("disabled", "$loading"),
	)
	assertAttr(t, merged, "data-on:click", "$open = !$open")
	assertAttr(t, merged, "data-attr:disabled", "$loading")
}

func TestMergeOverwrite(t *testing.T) {
	merged := ds.Merge(
		ds.OnClick("first"),
		ds.OnClick("second"),
	)
	assertAttr(t, merged, "data-on:click", "second")
}

func TestNoHyphenInParameterizedAttrs(t *testing.T) {
	// This test ensures we never accidentally produce data-on-click etc.
	tests := []struct {
		name  string
		attrs map[string]any
	}{
		{"On", ds.On("click", "x")},
		{"OnClick", ds.OnClick("x")},
		{"Bind", ds.Bind("val")},
		{"ClassToggle", ds.ClassToggle("active", "x")},
		{"Attr", ds.Attr("disabled", "x")},
		{"Style", ds.Style("color", "x")},
		{"Computed", ds.Computed("name", "x")},
		{"Indicator", ds.Indicator("load")},
		{"Ref", ds.Ref("el")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key := range tt.attrs {
				// Must start with "data-"
				if len(key) < 6 || key[:5] != "data-" {
					t.Fatalf("key %q does not start with data-", key)
				}
				// After "data-<plugin>", the parameter must be separated by ":"
				// A hyphen instead of colon would be silently ignored by Datastar.
				rest := key[5:] // e.g. "on:click", "bind:value"
				hasColon := false
				for _, c := range rest {
					if c == ':' {
						hasColon = true
						break
					}
				}
				if !hasColon {
					t.Errorf("key %q has no colon separator — would be silently ignored by Datastar", key)
				}
			}
		})
	}
}

// --- helpers ---

func assertAttr(t *testing.T, attrs map[string]any, key string, want any) {
	t.Helper()
	got, ok := attrs[key]
	if !ok {
		t.Fatalf("key %q not found in attributes: %v", key, attrs)
	}
	if got != want {
		t.Errorf("attrs[%q] = %v, want %v", key, got, want)
	}
}

func assertKey(t *testing.T, attrs map[string]any, key string) {
	t.Helper()
	if _, ok := attrs[key]; !ok {
		t.Fatalf("key %q not found in attributes: %v", key, attrs)
	}
}
