# CLAUDE.md

## Tech stack

Go + Chi (routing) + Templ (templating) + Tailwind CSS (styling) + Datastar (frontend interactivity) + TemplUI (component library)

Module: `github.com/plaenen/webx`

## Commands

- `go tool templ generate` — generate Go code from .templ files
- `go build ./cmd/...` — build the application
- `go run ./cmd/...` — run the application
- `go tool task install:all` — install all dependencies (mod tidy + templui init + add all components)
- `go tool task templui:add` — add/update all templui components
- `go tool gotailwind` — build Tailwind CSS

## Project structure

```
cmd/            — main entry point (create when needed)
internal/       — internal packages (create when needed)
ui/             — templui components (one dir per component, e.g. ui/button/button.templ)
utils/          — shared templ utilities (TwMerge, If, RandomID, etc.)
assets/css/     — Tailwind CSS
assets/js/      — JS for interactive components (managed by templui)
docs/           — reference documentation
.templui.json   — templui config
Taskfile.yaml   — task runner config
```

## Component pattern

TemplUI components follow this pattern — match it when creating or modifying components:

```go
package mycomponent

type Props struct {
    ID         string
    Class      string
    Attributes templ.Attributes
    // component-specific fields
}

templ MyComponent(props ...Props) {
    {{ var p Props }}
    if len(props) > 0 {
        {{ p = props[0] }}
    }
    // use utils.TwMerge() to combine Tailwind classes
}
```

Import components as: `"github.com/plaenen/webx/ui/button"`
Import utils as: `"github.com/plaenen/webx/utils"`

## Rules

- **No Co-Author lines** in commits
- **No custom CSS/JS** — use TemplUI + Datastar only
- **Wrap errors**: `fmt.Errorf("context: %w", err)`
- **Fix root cause** — when fixing a bug, fix the root cause, not the symptom
- **Use templui** for all UI — components, fragments, pages live in `./ui`
- **Use Datastar for all frontend interactivity** — no raw JS
  - Read [Datastar Go reference](./docs/datastar-go-reference.md) before writing Go-side Datastar code
  - Read [Datastar HTML elements reference](./docs/html-datastar-elements-reference.md) before writing templ-side Datastar attributes
- **Run `go tool templ generate`** after editing any `.templ` file
- **Ask if backward compatibility is required** before making changes
