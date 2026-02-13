---
name: migrate-component
description: Migrate a DaisyUI component into a WebX templ component following the migration guide
---

Migrate the DaisyUI component **$ARGUMENTS** into a WebX templ component.

Follow the migration guide at `docs/component-migration-guide.md` exactly. Read it first before writing any code.

Steps:

1. Look up the component on https://daisyui.com/components/$ARGUMENTS/ to understand its HTML structure, CSS classes, modifiers, and whether it uses a hidden input (checkbox/radio) or `<details>` for interactivity.

2. Determine the tier:
   - **Tier 1 (CSS-only)** if purely visual/structural
   - **Tier 2 (Interactive)** if it toggles, shows/hides, or changes state

3. Create the package at `ui/$ARGUMENTS/$ARGUMENTS.templ` following the tier template from the migration guide.

4. Map all DaisyUI modifiers to typed Go consts (Variant, Size, Position, etc.).

5. For Tier 2 components: define a signals struct, use `utils.Signals()`, and use `ds.*` helpers for ALL Datastar attributes. Never write raw `data-on-click` or similar — always use `ds.OnClick()`, `ds.Attr()`, etc.

6. Generate and build:
   ```sh
   go tool templ generate
   go build ./cmd/showcase
   ```

7. If the component introduces new DaisyUI classes not yet in the codebase, rebuild CSS:
   ```sh
   go tool gotailwind -i static/css/input.css -o static/css/output.css
   ```

8. Create a showcase page at `cmd/showcase/internal/pages/$ARGUMENTS.templ` demonstrating the component with multiple variants/sizes. Register the route in `cmd/showcase/main.go` and add a sidebar link in `layouts/showcase.templ`.

9. Write an E2E test at `tests/e2e/$ARGUMENTS_test.go`. For CSS-only components, verify elements render with expected classes. For interactive components, verify toggle/state behavior using Playwright.

10. Run E2E tests to verify:
    ```sh
    go test ./tests/e2e/... -run $ARGUMENTS -v
    ```
