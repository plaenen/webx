// Package ds provides typed helpers for Datastar HTML attributes.
//
// Datastar's parameterized plugins use a colon separator (e.g. data-on:click),
// NOT a hyphen (data-on-click). A hyphen is silently ignored as an unknown
// plugin. This package makes that mistake impossible by construction.
package ds

import "github.com/a-h/templ"

// --- Parameterized attributes (colon syntax) ---

// On returns a data-on:<event> attribute.
//
//	ds.On("click", expr) → {"data-on:click": expr}
func On(event, expr string) templ.Attributes {
	return templ.Attributes{"data-on:" + event: expr}
}

// OnClick is shorthand for On("click", expr).
func OnClick(expr string) templ.Attributes {
	return On("click", expr)
}

// Bind returns a data-bind:<signal> attribute.
func Bind(signal string) templ.Attributes {
	return templ.Attributes{"data-bind:" + signal: ""}
}

// ClassToggle returns a data-class:<name> attribute (single class toggle).
func ClassToggle(name, expr string) templ.Attributes {
	return templ.Attributes{"data-class:" + name: expr}
}

// Attr returns a data-attr:<name> attribute.
func Attr(name, expr string) templ.Attributes {
	return templ.Attributes{"data-attr:" + name: expr}
}

// Style returns a data-style:<prop> attribute.
func Style(prop, expr string) templ.Attributes {
	return templ.Attributes{"data-style:" + prop: expr}
}

// Computed returns a data-computed:<name> attribute.
func Computed(name, expr string) templ.Attributes {
	return templ.Attributes{"data-computed:" + name: expr}
}

// Indicator returns a data-indicator:<name> attribute.
func Indicator(name string) templ.Attributes {
	return templ.Attributes{"data-indicator:" + name: ""}
}

// Ref returns a data-ref:<name> attribute.
func Ref(name string) templ.Attributes {
	return templ.Attributes{"data-ref:" + name: ""}
}

// --- Standalone attributes (no colon) ---

// Signals returns a data-signals attribute.
func Signals(value string) templ.Attributes {
	return templ.Attributes{"data-signals": value}
}

// Show returns a data-show attribute.
func Show(expr string) templ.Attributes {
	return templ.Attributes{"data-show": expr}
}

// Text returns a data-text attribute.
func Text(expr string) templ.Attributes {
	return templ.Attributes{"data-text": expr}
}

// Class returns a data-class attribute (object syntax).
func Class(value string) templ.Attributes {
	return templ.Attributes{"data-class": value}
}

// Init returns a data-init attribute.
func Init(expr string) templ.Attributes {
	return templ.Attributes{"data-init": expr}
}

// Effect returns a data-effect attribute.
func Effect(expr string) templ.Attributes {
	return templ.Attributes{"data-effect": expr}
}

// --- Combiner ---

// Merge combines multiple templ.Attributes into one.
// Later values overwrite earlier ones for the same key.
func Merge(attrs ...templ.Attributes) templ.Attributes {
	merged := templ.Attributes{}
	for _, attr := range attrs {
		for k, v := range attr {
			merged[k] = v
		}
	}
	return merged
}
