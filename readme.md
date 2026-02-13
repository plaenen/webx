# webx

Go web framework built on Chi, Templ, Tailwind CSS, DaisyUI, and Datastar.

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
