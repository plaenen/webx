# webx

Go web framework built on Chi, Templ, Tailwind CSS, Datastar, and TemplUI.

## Prerequisites

- Go 1.24+
- [Datastar Pro license](https://data-star.dev/)

## BYOL (Bring Your Own License)

Datastar Pro is a commercial product. You must purchase your own license and place the files in the `byol/` directory, which is gitignored.

```
byol/
  datastar/
    datastar-pro.js
    datastar-inspector.js
    datastar-pro-rocket.js    # optional, if included in your license
```

After purchasing, download the JS files from your Datastar account and copy them into `byol/datastar/`. The build embeds this directory and serves the files at `/assets/js/datastar/`.

Without these files, the project will not compile.

## Getting started

```bash
# Install dependencies
go tool task install:all

# Place your Datastar Pro files
mkdir -p byol/datastar
cp /path/to/your/datastar-pro.js byol/datastar/
cp /path/to/your/datastar-inspector.js byol/datastar/

# Generate templ + build Tailwind
go tool templ generate
go tool gotailwind

# Run the showcase
go run ./cmd/showcase
```

The showcase starts at [http://localhost:3000](http://localhost:3000).
