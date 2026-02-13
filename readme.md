# webx

Beautiful, accessible UI components built with Go, Templ, Tailwind CSS, DaisyUI, and Datastar.

## Philosophy

webx is inspired by [daisyui](https://daisyui.com). Like daisyui, this is **not** a component library you install as a dependency. It's a collection of reusable components that you copy into your project and own.

**Why copy/paste instead of packaging?**

- **You own the code.** Need to change something? Just edit it.
- **No versioning conflicts** or breaking changes to worry about.
- **Zero configuration complexity.** Customize by changing the code.
- **Server-side rendering** with Go/templ — no Node.js build step for components.
- **~15KB of JavaScript** total (the Datastar library) for all frontend interactivity.

## Component structure

Each component lives in its own directory under `ui/` and follows a consistent pattern:

```
ui/button/
  button.templ       # Component template (Props struct, variants, rendering)
  button_templ.go    # Generated Go code (do not edit)
```

Interactive components that need a server-side handler add extra files:

```
ui/form/
  form.templ         # Component template
  form_templ.go      # Generated Go code
  handler.go         # SSE handler for Datastar interactivity
```

Components use DaisyUI classes for styling and Datastar attributes for interactivity:

```go
package button

type Props struct {
    ID         string
    Class      string
    Attributes templ.Attributes
    Variant    Variant           // btn-primary, btn-secondary, etc.
    Size       Size              // btn-lg, btn-sm, btn-xs
    Disabled   bool
    OnClick    string            // Datastar action expression
}

templ Button(props ...Props) {
    // Optional variadic props — sensible zero-value defaults
    // DaisyUI classes merged with utils.TwMerge()
    // Datastar attributes spread via ds.* helpers
}
```

Import and use in your templ files:

```go
import "github.com/plaenen/webx/ui/button"

@button.Button(button.Props{Variant: button.VariantPrimary}) {
    Click me
}
```

## Prerequisites

- Go 1.24+

## Getting started

```bash
# Install dependencies (downloads DaisyUI, Datastar, tidies modules)
go tool task install:all

# Generate templ + build Tailwind
go tool templ generate
go tool gotailwind

# Run the showcase (uses open-source Datastar)
go run ./cmd/showcase serve
```

The showcase starts at [http://localhost:3000](http://localhost:3000).

## Datastar Pro (optional)

The showcase supports both open-source Datastar and Datastar Pro. Open-source is the default and requires no license.

To use Datastar Pro, purchase a license from [data-star.dev](https://data-star.dev/) and place the files in the `byol/` directory (gitignored):

```
byol/
  datastar/
    datastar-pro.js
    datastar-inspector.js
    datastar-pro-rocket.js    # optional
```

Then start the showcase with the `--pro` flag:

```bash
go run ./cmd/showcase serve --pro
```

Pro mode enables the Datastar inspector and loads from your local BYOL files.
