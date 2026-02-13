# DaisyUI → WebX Component Migration Guide

> For agents converting DaisyUI HTML patterns into WebX templ components with Datastar interactivity.

## Project layout

```
ui/<component>/         — one directory per component
  <component>.templ     — templ source (you write this)
  <component>_templ.go  — generated (run `go tool templ generate`)
ds/ds.go                — typed Datastar attribute helpers (ALWAYS use for data-on:*, data-attr:*, etc.)
utils/utils.go          — TwMerge, If, IfElse, RandomID, MergeAttributes
utils/signals.go        — SignalManager for namespaced Datastar signals
utils/data_class.go     — DataClass builder for data-class object syntax
```

## Imports

```go
import "github.com/plaenen/webx/utils"       // TwMerge, If, RandomID, SignalManager
import "github.com/plaenen/webx/ds"           // typed Datastar attributes
```

Only import `ds` when the component uses Datastar interactivity.

---

## Component tiers

Every DaisyUI component falls into one of two tiers:

### Tier 1 — CSS-only (no state)

Components that are purely structural/visual. No signals, no Datastar attributes.

**Examples:** card, navbar, menu, badge, alert, breadcrumb, footer, hero, stat, table, timeline, avatar, separator

**Pattern:** Props struct + `utils.TwMerge` + `{ children... }`

### Tier 2 — Interactive (stateful)

Components that show/hide, toggle, or otherwise change state. These need Datastar signals.

**Examples:** drawer, dropdown, modal/dialog, collapse/accordion, tabs, tooltip, popover, swap

**Pattern:** Props struct + SignalManager + `ds.*` helpers + DaisyUI CSS mechanism (checkbox, radio, or CSS classes)

---

## Tier 1: CSS-only component template

Use this for any component that has no interactive state.

### Single-element component

```go
package badge

import "github.com/plaenen/webx/utils"

type Variant string

const (
	VariantDefault   Variant = ""
	VariantPrimary   Variant = "badge-primary"
	VariantSecondary Variant = "badge-secondary"
	// ... map each DaisyUI modifier to a const
)

type Size string

const (
	SizeDefault Size = ""
	SizeLg      Size = "badge-lg"
	SizeSm      Size = "badge-sm"
)

type Props struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	Variant    Variant
	Size       Size
}

templ Badge(props ...Props) {
	{{ var p Props }}
	if len(props) > 0 {
		{{ p = props[0] }}
	}
	<span
		if p.ID != "" {
			id={ p.ID }
		}
		class={ utils.TwMerge("badge", string(p.Variant), string(p.Size), p.Class) }
		{ p.Attributes... }
	>
		{ children... }
	</span>
}
```

### Multi-part component (compound)

Some DaisyUI components have sub-elements. Create a separate templ for each part. They share the same package.

```go
package card

import "github.com/plaenen/webx/utils"

type Props struct {
	ID         string
	Class      string
	Attributes templ.Attributes
}

templ Card(props ...Props) {
	// ... renders <div class="card ...">{ children... }</div>
}

templ Body(props ...Props) {
	// ... renders <div class="card-body ...">{ children... }</div>
}

templ Title(props ...Props) {
	// ... renders <h2 class="card-title ...">{ children... }</h2>
}

templ Actions(props ...Props) {
	// ... renders <div class="card-actions justify-end ...">{ children... }</div>
}
```

### Rules for CSS-only components

1. Package name = component name (lowercase, singular): `badge`, `alert`, `stat`
2. Props use variadic `...Props` so callers can omit them: `@badge.Badge()` works
3. First line of templ body: `{{ var p Props }}` then `if len(props) > 0 { {{ p = props[0] }} }`
4. DaisyUI base class is always the first arg to `TwMerge`: `utils.TwMerge("badge", ...)`
5. `p.Class` is always the **last** arg to `TwMerge` so callers can override defaults
6. `{ p.Attributes... }` is spread on the root element to allow arbitrary HTML attributes
7. `{ children... }` renders slot content
8. ID is only set when non-empty (avoids empty `id=""` in HTML)
9. Map every DaisyUI modifier to a typed const: `badge-primary` → `VariantPrimary`
10. Include sensible defaults in `TwMerge` (e.g. `card bg-base-200 shadow-sm`)

---

## Tier 2: Interactive component template

Use this when the DaisyUI component requires JavaScript-like behavior. DaisyUI components use hidden checkboxes, radio buttons, or CSS-only mechanisms. We keep those DOM elements but control them with Datastar signals.

### Signal struct

Define a struct for the component's reactive state:

```go
type DrawerSignals struct {
	Open bool `json:"open"`
}

type TabsSignals struct {
	Active string `json:"active"`
}

type AccordionSignals struct {
	Open bool `json:"open"`
}
```

Rules:
- Field names must be exported (uppercase first letter)
- Must have `json:"..."` tags (lowercase, used in Datastar expressions)
- Keep it minimal — only state that changes at runtime

### SignalManager

Create a `SignalManager` from the component's ID and signal struct:

```go
signals := utils.Signals(id, DrawerSignals{Open: false})
```

This produces:
- `signals.DataSignals` → `{"my_component":{"open":false}}` (for `data-signals` attribute)
- `signals.Signal("open")` → `$my_component.open` (reference in expressions)
- `signals.Toggle("open")` → `$my_component.open = !$my_component.open`
- `signals.Set("open", "false")` → `$my_component.open = false`
- `signals.Conditional("open", "'yes'", "'no'")` → `$my_component.open ? 'yes' : 'no'`

The ID is sanitized: hyphens become underscores (`my-drawer` → `my_drawer`).

### The `ds` package (MANDATORY for Datastar attributes)

NEVER write Datastar attributes as raw strings. Use `ds.*` helpers:

```go
// WRONG — hyphen instead of colon, silently ignored by Datastar:
templ.Attributes{"data-on-click": expr}

// WRONG — raw string, typo-prone:
data-on:click={ expr }

// CORRECT — typed helper, colon guaranteed:
{ ds.OnClick(expr)... }
```

Available helpers:

| Helper | Produces | Use for |
|--------|----------|---------|
| `ds.OnClick(expr)` | `data-on:click` | Click handlers |
| `ds.On("change", expr)` | `data-on:change` | Any DOM event |
| `ds.Attr("checked", expr)` | `data-attr:checked` | Reactive HTML attributes |
| `ds.Attr("disabled", expr)` | `data-attr:disabled` | Reactive disabled state |
| `ds.Show(expr)` | `data-show` | Conditional visibility |
| `ds.Text(expr)` | `data-text` | Reactive text content |
| `ds.Class(value)` | `data-class` | Reactive class object |
| `ds.ClassToggle(name, expr)` | `data-class:name` | Toggle single class |
| `ds.Bind(signal)` | `data-bind:signal` | Two-way form binding |
| `ds.Style(prop, expr)` | `data-style:prop` | Reactive inline style |
| `ds.Merge(a, b, ...)` | Combined attributes | Multiple ds.* on one element |

Spread syntax on elements: `{ ds.OnClick(expr)... }`
Through Props: `Attributes: ds.OnClick(expr)`
Multiple through Props: `Attributes: ds.Merge(ds.OnClick(a), ds.Attr("disabled", b))`

### DaisyUI's checkbox/radio mechanism

Many DaisyUI interactive components use a hidden `<input>` as a CSS toggle. The CSS uses `:checked` pseudo-selectors to show/hide siblings. **You must keep these inputs in the DOM** — DaisyUI's CSS rules depend on them via sibling selectors (`~`).

Control the input's checked state via Datastar instead of letting the user click it directly:

```go
<input
	type="checkbox"
	class="drawer-toggle"
	aria-hidden="true"
	tabindex="-1"
	{ ds.Attr("checked", signals.Signal("open"))... }
/>
```

This pattern applies to: drawer, modal, collapse, accordion, dropdown (when using checkbox).

### CSS class rules and `drawer-toggle`

DaisyUI components that use hidden inputs need those CSS classes in the output. Tailwind CSS v4 only generates CSS for classes it detects in source files. If you add a new component that uses a DaisyUI class for the first time (e.g., `modal-toggle`, `collapse-toggle`), you must **rebuild Tailwind CSS** so the rules are generated:

```sh
go tool gotailwind -i static/css/input.css -o static/css/output.css
```

If a component's toggle doesn't work, check `static/css/output.css` for the relevant class rules. Missing rules = Tailwind didn't detect the class in source.

### Interactive component example (drawer)

```go
package drawer

import (
	"github.com/plaenen/webx/ds"
	"github.com/plaenen/webx/utils"
)

type DrawerSignals struct {
	Open bool `json:"open"`
}

type Props struct {
	ID    string
	Class string
}

templ Drawer(props Props) {
	{{
		id := props.ID
		if id == "" {
			id = utils.RandomID()
		}
		signals := utils.Signals(id, DrawerSignals{Open: false})
	}}
	<div
		data-signals={ signals.DataSignals }
		class={ utils.TwMerge("drawer lg:drawer-open", props.Class) }
	>
		<input
			type="checkbox"
			class="drawer-toggle"
			aria-hidden="true"
			tabindex="-1"
			{ ds.Attr("checked", signals.Signal("open"))... }
		/>
		{ children... }
	</div>
}
```

Key points:
- `data-signals` declares the reactive state on the container element
- `data-attr:checked` syncs the hidden checkbox with the signal
- Child components reference the same signal by creating a `SignalManager` with the same ID
- The `data-signals` attribute is only on the **outermost** container (once per component instance)

### Explicit Datastar props

For commonly used Datastar interactions, add explicit props instead of requiring callers to pass raw `Attributes`. The component merges them internally:

```go
type Props struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	Variant    Variant
	Size       Size
	Disabled   bool
	OnClick    string   // explicit Datastar on:click
}

templ Button(props ...Props) {
	// ...
	{{
		attrs := p.Attributes
		if p.OnClick != "" {
			attrs = ds.Merge(attrs, ds.OnClick(p.OnClick))
		}
	}}
	<button ... { attrs... }>
		{ children... }
	</button>
}
```

Add explicit props for interactions that are **common for the component type**:
- Button → `OnClick`
- Input → `OnChange`, `OnInput`
- Form → `OnSubmit`
- Anything toggleable → no explicit prop needed (use SignalManager internally)

---

## Step-by-step migration process

### 1. Identify the DaisyUI component

Go to [daisyui.com/components](https://daisyui.com/components/) and read the component's HTML structure. Note:
- What CSS classes it uses (base class + modifiers)
- Whether it needs a hidden input (checkbox/radio) for interactivity
- What sub-parts it has (body, title, actions, content, etc.)
- What modifiers exist (variants, sizes, positions)

### 2. Determine the tier

- Does it toggle, show/hide, or change state? → **Tier 2** (interactive)
- Is it purely structural/visual? → **Tier 1** (CSS-only)

### 3. Create the package

```sh
mkdir -p ui/<component>
```

Create `ui/<component>/<component>.templ` following the appropriate tier template above.

### 4. Map DaisyUI modifiers to Go types

For each group of mutually exclusive modifiers, create a typed string:

```go
// DaisyUI: alert-info, alert-success, alert-warning, alert-error
type Variant string
const (
	VariantDefault Variant = ""
	VariantInfo    Variant = "alert-info"
	VariantSuccess Variant = "alert-success"
	VariantWarning Variant = "alert-warning"
	VariantError   Variant = "alert-error"
)
```

Naming conventions:
- `Variant` for visual style (primary, secondary, success, error, ghost, etc.)
- `Size` for size modifiers (xs, sm, md, lg, xl)
- `Position` for placement (top, bottom, start, end)
- Use the exact DaisyUI class as the const value

### 5. Handle DaisyUI's hidden inputs (Tier 2 only)

Check the DaisyUI docs for the component. If the HTML example includes:

```html
<input type="checkbox" class="XXX-toggle" />
```

Then you need:
1. A signals struct with the toggle state
2. The hidden input with `ds.Attr("checked", signal)`
3. The toggle class (e.g., `modal-toggle`, `collapse-toggle`, `drawer-toggle`)

If DaisyUI uses a `<details>` element instead of a checkbox:

```html
<details class="collapse">
  <summary>...</summary>
  <div class="collapse-content">...</div>
</details>
```

You have two options:
- Keep `<details>` and use `ds.Attr("open", signal)` to control it
- Replace with a div + checkbox approach if more control is needed

### 6. Generate and build

After writing the `.templ` file:

```sh
go tool templ generate          # generates _templ.go
go build ./cmd/showcase         # verify it compiles
```

If the component introduces a new DaisyUI class (first use in the codebase):

```sh
go tool gotailwind -i static/css/input.css -o static/css/output.css
```

### 7. Create a showcase page

Add a page in `cmd/showcase/internal/pages/<component>.templ` demonstrating the component. Register it in `cmd/showcase/main.go` and add a sidebar link in `layouts/showcase.templ`.

### 8. Write an E2E test

Add `tests/e2e/<component>_test.go`. For CSS-only components, verify elements render with the expected classes. For interactive components, verify the toggle/state behavior.

---

## DaisyUI CSS patterns reference

### Checkbox toggle (drawer, modal, collapse)

DaisyUI CSS rule structure:
```css
.component-toggle:checked ~ .component-side { /* visible state */ }
.component-open > .component-toggle ~ .component-side { /* always-open state */ }
```

The `~` is a sibling selector. The hidden input MUST be a sibling before the target element.

### Radio toggle (tabs, accordion)

Multiple items share a `name` attribute. Only one can be checked at a time.

```html
<input type="radio" name="my-tabs" class="tab" />
```

Use `ds.Attr("checked", signals.Equals("active", "'tab1'"))` for each radio.

### Details/summary (collapse)

DaisyUI's simpler collapse uses native `<details>`:

```html
<details class="collapse"><summary>Title</summary><div class="collapse-content">...</div></details>
```

Control with `ds.Attr("open", signal)`.

### Pure class toggle

Some components (swap, dropdown with CSS focus) toggle via CSS classes or pseudo-selectors rather than inputs. Use `data-class` to toggle the activating class:

```go
data-class={ utils.NewDataClass().Add("swap-active", signals.Signal("swapped")).Build() }
```

---

## Common mistakes to avoid

1. **Using `data-on-click` instead of `data-on:click`** — Datastar silently ignores hyphen syntax. Always use `ds.OnClick()` or `ds.On()`.

2. **Forgetting the hidden input** — DaisyUI's CSS uses sibling selectors (`~`) that require the toggle input in the DOM. Without it, the CSS rules don't match and nothing happens.

3. **Missing CSS rules** — If a DaisyUI class is used for the first time, Tailwind may not have generated its CSS. Rebuild with `go tool gotailwind -i static/css/input.css -o static/css/output.css`.

4. **Using `drawer-open` class for mobile toggle** — The `drawer-open` class means "always open" (desktop). It hides the overlay. For mobile toggle, use the checkbox `:checked` mechanism via `ds.Attr("checked", signal)`.

5. **Declaring `data-signals` in multiple places** — Only declare `data-signals` once on the outermost container. Child components create a `SignalManager` with the same ID but don't emit `data-signals` again.

6. **Signal ID mismatch** — All sub-components that share state must use the exact same ID string. The `SignalManager` sanitizes hyphens to underscores, so `"my-drawer"` and `"my-drawer"` match, but `"myDrawer"` does not.

7. **Writing Datastar attributes as raw strings in templ** — Always use `ds.*` helpers. Raw strings bypass the colon enforcement and are typo-prone.
