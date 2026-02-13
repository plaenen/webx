# Datastar + Go Backend: Comprehensive Reference for Coding Agents

> **Version**: Datastar v1.0.0-RC.7 | Go SDK v1.0.3 (`github.com/starfederation/datastar-go`)
> **Last Updated**: February 2026
> **Go Requirement**: Go 1.24+

---

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Installation & Setup](#installation--setup)
3. [Frontend: Datastar Attributes Reference](#frontend-datastar-attributes-reference)
4. [Frontend: Actions Reference](#frontend-actions-reference)
5. [Frontend: Datastar Expressions](#frontend-datastar-expressions)
6. [Backend: Go SDK API Reference](#backend-go-sdk-api-reference)
7. [SSE Events Protocol](#sse-events-protocol)
8. [Complete Patterns & Examples](#complete-patterns--examples)
9. [Security Guidelines](#security-guidelines)
10. [Common Patterns & How-Tos](#common-patterns--how-tos)
11. [Northstar Best Practices (Production Boilerplate)](#northstar-best-practices-production-boilerplate)

---

## Architecture Overview

Datastar is a **hypermedia-driven reactive framework** (~11 KiB) that combines backend-driven UI (like htmx) with frontend reactivity (like Alpine.js) in a single library. The core principles are:

- **Backend drives the frontend**: The server sends HTML fragments and state updates over Server-Sent Events (SSE). There is no separate REST API layer.
- **No JavaScript build step**: All frontend reactivity is declared via `data-*` HTML attributes.
- **Signals are the state model**: Reactive state lives in "signals" (prefixed with `$` in expressions). The backend can read, modify, and patch signals.
- **DOM morphing**: Datastar uses Idiomorph to intelligently morph only changed parts of the DOM, preserving state and event listeners.

### Request/Response Flow

```
Browser                          Go Backend
  │                                  │
  │ ── @get('/endpoint') ──────────► │  (sends all signals as query param or JSON body)
  │                                  │
  │ ◄── text/event-stream ───────── │  (streams SSE events: patch-elements, patch-signals)
  │                                  │
  │  [Datastar morphs DOM &          │
  │   updates signals reactively]    │
```

### Content Types the Backend Can Return

| Content-Type | Behavior |
|---|---|
| `text/event-stream` | Standard SSE with Datastar events (recommended) |
| `text/html` | HTML elements patched into DOM by ID |
| `application/json` | JSON signals patched into state |
| `text/javascript` | JavaScript executed in browser |

---

## Installation & Setup

### Frontend (HTML)

Include Datastar via CDN script tag (version-locked):

```html
<script type="module" src="https://cdn.jsdelivr.net/gh/starfederation/[email protected]/bundles/datastar.js"></script>
```

Or self-host by downloading the script file and serving it:

```html
<script type="module" src="/static/js/datastar.js"></script>
```

### Backend (Go)

Install the Go SDK:

```bash
go get github.com/starfederation/datastar-go
```

Import in your Go code:

```go
import (
    "net/http"
    "github.com/starfederation/datastar-go/datastar"
)
```

### Minimal Complete Example

```go
package main

import (
    "fmt"
    "net/http"
    "time"

    "github.com/starfederation/datastar-go/datastar"
)

func main() {
    mux := http.NewServeMux()

    // Serve the HTML page
    mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/html")
        fmt.Fprint(w, `<!DOCTYPE html>
<html>
<head>
    <script type="module" src="https://cdn.jsdelivr.net/gh/starfederation/[email protected]/bundles/datastar.js"></script>
</head>
<body>
    <div data-signals:count="0">
        <div id="output">Count: <span data-text="$count"></span></div>
        <button data-on:click="@post('/increment')">Increment</button>
    </div>
</body>
</html>`)
    })

    // Handle the increment action
    mux.HandleFunc("POST /increment", func(w http.ResponseWriter, r *http.Request) {
        // Read signals from client
        store := &struct {
            Count int `json:"count"`
        }{}
        if err := datastar.ReadSignals(r, store); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        store.Count++

        // Create SSE writer and send updates
        sse := datastar.NewSSE(w, r)
        sse.MarshalAndPatchSignals(store)
        sse.PatchElements(fmt.Sprintf(
            `<div id="output">Count: <span data-text="$count">%d</span></div>`,
            store.Count,
        ))
    })

    fmt.Println("Server running at http://localhost:8080")
    http.ListenAndServe(":8080", mux)
}
```

---

## Frontend: Datastar Attributes Reference

All Datastar functionality is declared via `data-*` HTML attributes. Attributes are evaluated in DOM order (depth-first). Hyphenated keys are auto-converted to camelCase for signals.

### Signal & State Attributes

#### `data-signals`

Patches (adds, updates, removes) one or more signals. This is how you declare reactive state.

```html
<!-- Single signal -->
<div data-signals:count="0"></div>

<!-- Multiple signals (object syntax) -->
<div data-signals="{name: 'Alice', age: 30, active: true}"></div>

<!-- Nested signals (dot notation) -->
<div data-signals:user.name="'Bob'"></div>

<!-- Nested signals (object syntax) -->
<div data-signals="{user: {name: 'Bob', email: 'bob@example.com'}}"></div>

<!-- Remove a signal by setting to null -->
<div data-signals="{obsoleteSignal: null}"></div>
```

**Modifiers:**
- `__ifmissing` — only set if signal doesn't already exist: `data-signals:foo__ifmissing="1"`
- `__case` — `.camel` (default), `.kebab`, `.snake`, `.pascal`

**Rules:**
- Signals beginning with `_` are NOT sent to the backend by default.
- Signal names cannot contain `__` (reserved for modifier delimiter).
- Keys in `data-signals:*` are auto-converted to camelCase.

#### `data-bind`

Two-way data binding between a signal and an input/select/textarea element.

```html
<!-- Creates signal $name and binds to input value -->
<input data-bind:name />

<!-- Alternative value syntax -->
<input data-bind="name" />

<!-- With initial HTML value (if no signal predefined) -->
<input data-bind:username value="defaultUser" />

<!-- Predefined signal type preserved -->
<div data-signals:count="0">
    <select data-bind:count>
        <option value="10">10</option>  <!-- $count will be number 10, not string "10" -->
    </select>
</div>

<!-- Checkbox array binding -->
<div data-signals:colors="[]">
    <input data-bind:colors type="checkbox" value="red" />
    <input data-bind:colors type="checkbox" value="blue" />
</div>

<!-- File upload (auto base64 encoded) -->
<input type="file" data-bind:files multiple />
<!-- Signal format: { name: string, contents: string, mime: string }[] -->
```

#### `data-computed`

Creates a read-only derived signal that auto-updates when dependencies change.

```html
<div data-computed:full-name="$firstName + ' ' + $lastName"></div>
<div data-computed:total="$price * $quantity"></div>

<!-- Object syntax for multiple computeds -->
<div data-computed="{fullName: () => $firstName + ' ' + $lastName}"></div>
```

#### `data-ref`

Creates a signal that references a DOM element.

```html
<canvas data-ref:my-canvas></canvas>
<div data-text="$myCanvas.tagName"></div>  <!-- outputs "CANVAS" -->
```

### Display & Rendering Attributes

#### `data-text`

Binds element text content to an expression.

```html
<span data-text="$count"></span>
<span data-text="$name.toUpperCase()"></span>
<span data-text="`Hello, ${$name}!`"></span>
```

#### `data-show`

Shows/hides an element based on expression truthiness.

```html
<div data-show="$isLoggedIn">Welcome back!</div>

<!-- Prevent FOUC with initial display:none -->
<div data-show="$loaded" style="display: none">Content</div>
```

#### `data-class`

Conditionally adds/removes CSS classes.

```html
<!-- Single class -->
<div data-class:font-bold="$isImportant"></div>
<div data-class:hidden="!$visible"></div>

<!-- Multiple classes (object syntax) -->
<div data-class="{active: $isActive, 'text-red-500': $hasError}"></div>
```

#### `data-attr`

Sets any HTML attribute reactively.

```html
<button data-attr:disabled="$isSubmitting"></button>
<a data-attr:href="`/users/${$userId}`"></a>
<img data-attr:src="$imageUrl" />

<!-- Multiple attributes (object syntax) -->
<div data-attr="{'aria-label': $label, disabled: $isDisabled}"></div>
```

#### `data-style`

Sets inline CSS styles reactively.

```html
<div data-style:color="$isError ? 'red' : 'green'"></div>
<div data-style:display="$hidden && 'none'"></div>

<!-- Multiple styles (object syntax) -->
<div data-style="{color: $textColor, 'font-size': `${$size}px`}"></div>
```

### Event Attributes

#### `data-on`

Attaches event listeners. The `evt` variable is available in the expression.

```html
<button data-on:click="@post('/submit')">Submit</button>
<button data-on:click="$count++">Increment</button>
<input data-on:input="$search = evt.target.value" />
<div data-on:my-custom-event="$data = evt.detail"></div>

<!-- Form submit (auto-prevents default) -->
<form data-on:submit="@post('/login')">...</form>
```

**Modifiers (chainable with `__`):**

| Modifier | Description |
|---|---|
| `__once` | Fire only once |
| `__debounce.500ms` | Debounce (accepts ms/s, `.leading`, `.notrailing`) |
| `__throttle.1s` | Throttle (accepts ms/s, `.noleading`, `.trailing`) |
| `__delay.500ms` | Delay execution |
| `__window` | Listen on window |
| `__outside` | Trigger when event is outside element |
| `__prevent` | Call `preventDefault()` |
| `__stop` | Call `stopPropagation()` |
| `__passive` | Passive event listener |
| `__capture` | Capture phase |
| `__viewtransition` | Wrap in View Transition API |

```html
<button data-on:click__debounce.300ms="@post('/search')">Search</button>
<input data-on:keydown__throttle.100ms="handleKey()" />
<div data-on:click__outside="$menuOpen = false"></div>
```

#### `data-on-intersect`

Fires when element enters/exits viewport.

```html
<div data-on-intersect="@get('/load-more')">Loading...</div>
<div data-on-intersect__once__full="$seen = true"></div>
<div data-on-intersect__half="$visible = true"></div>
```

#### `data-on-interval`

Runs expression at regular intervals (default: 1 second).

```html
<div data-on-interval="@get('/poll')"></div>
<div data-on-interval__duration.5s="@get('/check-status')"></div>
<div data-on-interval__duration.100ms="$elapsed++"></div>
```

#### `data-on-signal-patch`

Fires when any signal is patched.

```html
<div data-on-signal-patch="console.log('Signal changed:', patch)"></div>

<!-- With filter -->
<div data-on-signal-patch-filter="{include: /^user/}"
     data-on-signal-patch="@post('/sync')"></div>
```

### Lifecycle & Control Attributes

#### `data-init`

Runs expression when element is initialized (page load, DOM patch, attribute change).

```html
<div data-init="@get('/initial-data')"></div>
<div data-init__delay.1s="$ready = true"></div>
```

#### `data-effect`

Runs expression on load AND whenever referenced signals change (side effects).

```html
<div data-effect="document.title = `Count: ${$count}`"></div>
<div data-effect="$total = $price * $quantity"></div>
```

#### `data-indicator`

Creates a boolean signal that is `true` while a fetch request is in flight.

```html
<button data-on:click="@get('/data')" data-indicator:loading>
    Load Data
</button>
<div data-show="$loading">Loading...</div>
<button data-attr:disabled="$loading">Submit</button>
```

**Important:** Place `data-indicator` BEFORE `data-init` when using together:

```html
<div data-indicator:fetching data-init="@get('/endpoint')"></div>
```

#### `data-ignore`

Tells Datastar to skip processing an element and its descendants.

```html
<div data-ignore>
    <div data-some-third-party-lib="">...</div>
</div>

<!-- __self: only ignore element itself, not children -->
<div data-ignore__self>
    <div data-text="$stillProcessed"></div>
</div>
```

#### `data-ignore-morph`

Skips morphing during PatchElements for this element and children.

```html
<div data-ignore-morph>
    <video><!-- preserves playback state --></video>
</div>
```

#### `data-preserve-attr`

Preserves specified attributes during morphing.

```html
<details open data-preserve-attr="open">
    <summary>Collapsible</summary>
    <p>Content preserved across morphs</p>
</details>
```

### Debugging Attributes

#### `data-json-signals`

Displays all signals as formatted JSON (great for debugging).

```html
<pre data-json-signals></pre>

<!-- With filter -->
<pre data-json-signals="{include: /user/}"></pre>
<pre data-json-signals="{exclude: /password/}"></pre>

<!-- Compact format -->
<pre data-json-signals__terse="{include: /count/}"></pre>
```

---

## Frontend: Actions Reference

Actions are helper functions prefixed with `@` that can be used in Datastar expressions.

### Utility Actions

#### `@peek(callable)`

Access signal values without subscribing to changes.

```html
<!-- $foo changes trigger re-evaluation, $bar changes do NOT -->
<div data-text="$foo + @peek(() => $bar)"></div>
```

#### `@setAll(value, filter?)`

Set multiple signals at once.

```html
<button data-on:click="@setAll(false, {include: /^is/})">Reset All</button>
<button data-on:click="@setAll('', {include: /^form\./, exclude: /_temp$/})">Clear Form</button>
```

#### `@toggleAll(filter?)`

Toggle boolean signals.

```html
<button data-on:click="@toggleAll({include: /^settings\./})">Toggle Settings</button>
```

### Backend Actions

All backend actions send an HTTP request. By default, ALL signals (except `_`-prefixed local ones) are included. For GET requests, signals go as a `datastar` query parameter. For other methods, signals are sent as a JSON body.

#### `@get(uri, options?)`

```html
<button data-on:click="@get('/api/data')">Load</button>
<div data-init="@get('/initial-state')"></div>
```

#### `@post(uri, options?)`

```html
<button data-on:click="@post('/api/submit')">Submit</button>
```

#### `@put(uri, options?)` / `@patch(uri, options?)` / `@delete(uri, options?)`

```html
<button data-on:click="@put('/api/users/1')">Update</button>
<button data-on:click="@delete('/api/users/1')">Delete</button>
```

#### Backend Action Options

```html
<button data-on:click="@post('/endpoint', {
    filterSignals: {include: /^form\./, exclude: /_temp$/},
    headers: {'X-CSRF-Token': 'abc123'},
    openWhenHidden: true,
    contentType: 'json',
    retry: 'auto',
    retryInterval: 1000,
    retryScaler: 2,
    retryMaxWaitMs: 30000,
    retryMaxCount: 10,
    requestCancellation: 'auto'
})">Submit</button>
```

| Option | Type | Default | Description |
|---|---|---|---|
| `contentType` | `'json'` \| `'form'` | `'json'` | `'form'` finds closest form, validates, sends as form data |
| `filterSignals` | `{include: RegExp, exclude?: RegExp}` | all non-`_` signals | Filter which signals to send |
| `selector` | `string \| null` | `null` | CSS selector for form when `contentType: 'form'` |
| `headers` | `object` | `{}` | Additional request headers |
| `openWhenHidden` | `boolean` | `false` (GET), `true` (others) | Keep SSE open when tab hidden |
| `payload` | `object` | - | Override fetch payload |
| `retry` | `'auto'` \| `'error'` \| `'always'` \| `'never'` | `'auto'` | Retry strategy |
| `retryInterval` | `number` | `1000` | Base retry interval (ms) |
| `retryScaler` | `number` | `2` | Exponential backoff multiplier |
| `retryMaxWaitMs` | `number` | `30000` | Max wait between retries |
| `retryMaxCount` | `number` | `10` | Max retry attempts |
| `requestCancellation` | `'auto'` \| `'disabled'` \| `AbortController` | `'auto'` | `'auto'` cancels previous request on same element |

#### Form Submission with File Upload

```html
<form enctype="multipart/form-data">
    <input type="file" name="document" />
    <button data-on:click="@post('/upload', {contentType: 'form'})">Upload</button>
</form>
```

### Request Headers Sent by Datastar

Every request includes:
- `Datastar-Request: true`
- Signals in `datastar` query param (GET) or JSON body (other methods)

---

## Frontend: Datastar Expressions

Datastar expressions are strings evaluated in a sandboxed context. They support JavaScript syntax but with important differences:

### Signal Access

Signals are accessed with `$` prefix:

```html
<div data-text="$count"></div>
<div data-text="$user.name"></div>
<div data-text="$items.length"></div>
```

### JavaScript in Expressions

```html
<div data-text="$name.toUpperCase()"></div>
<div data-text="`Hello, ${$name}! You have ${$count} items.`"></div>
<div data-text="$count > 0 ? 'Has items' : 'Empty'"></div>
<div data-text="$items.map(i => i.name).join(', ')"></div>
```

### Available Variables

- `$signalName` — reactive signal value
- `el` — the current DOM element
- `evt` — event object (in `data-on` expressions)
- `patch` — signal patch details (in `data-on-signal-patch`)

### Attribute Casing Rules

HTML attributes are case-insensitive. Datastar converts:

- **Signal-defining attributes** (bind, signals, computed, ref, indicator): hyphen → camelCase
  - `data-signals:my-signal` → signal `$mySignal`
  - `data-bind:user-name` → signal `$userName`
- **Other attributes** (class, on, attr, style): hyphen → kebab-case (default)
  - `data-class:text-blue-700` → class `text-blue-700`
  - `data-on:my-event` → event `my-event`

Use `__case` modifier to override: `.camel`, `.kebab`, `.snake`, `.pascal`

```html
<div data-on:widget-loaded__case.camel="handleLoad()"></div>
```

---

## Backend: Go SDK API Reference

### Package Import

```go
import "github.com/starfederation/datastar-go/datastar"
```

### Reading Signals from Requests

```go
func datastar.ReadSignals(r *http.Request, signals any) error
```

Extracts signals from the request and unmarshals into a struct pointer.

- **GET requests**: reads from `datastar` query parameter (URL-encoded JSON)
- **Other methods**: reads from JSON request body

```go
// Define a struct matching your signal shape
type Store struct {
    Count   int    `json:"count"`
    Name    string `json:"name"`
    Active  bool   `json:"active"`
}

// Nested signals
type FormStore struct {
    User struct {
        Name  string `json:"name"`
        Email string `json:"email"`
    } `json:"user"`
    Preferences struct {
        Theme    string `json:"theme"`
        Language string `json:"language"`
    } `json:"preferences"`
}

func handler(w http.ResponseWriter, r *http.Request) {
    store := &Store{}
    if err := datastar.ReadSignals(r, store); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    // Use store.Count, store.Name, etc.
}
```

### Creating an SSE Writer

```go
func datastar.NewSSE(w http.ResponseWriter, r *http.Request, opts ...SSEOption) *ServerSentEventGenerator
```

Creates a new SSE writer that automatically sets required headers:
- `Content-Type: text/event-stream`
- `Cache-Control: no-cache`
- `Connection: keep-alive`

```go
sse := datastar.NewSSE(w, r)

// With options
sse := datastar.NewSSE(w, r,
    datastar.WithCompression(
        datastar.WithGzip(),
        datastar.WithBrotli(),
    ),
    datastar.WithContext(ctx),
)
```

**SSE Options:**

| Option | Description |
|---|---|
| `WithCompression(opts...)` | Enable SSE compression (Gzip, Brotli, Deflate, Zstd) |
| `WithContext(ctx)` | Custom context for cancellation |

### ServerSentEventGenerator Methods

#### Patching Elements (DOM Updates)

```go
// Patch HTML elements into DOM (morphs by default, matches by ID)
func (sse *ServerSentEventGenerator) PatchElements(elements string, opts ...PatchElementOption) error

// Formatted patch
func (sse *ServerSentEventGenerator) PatchElementf(format string, args ...any) error

// Patch using templ component (github.com/a-h/templ)
func (sse *ServerSentEventGenerator) PatchElementTempl(c TemplComponent, opts ...PatchElementOption) error

// Patch using GoStar renderer
func (sse *ServerSentEventGenerator) PatchElementGostar(child GoStarElementRenderer, opts ...PatchElementOption) error
```

**PatchElement Options:**

```go
// Patch mode
datastar.WithModeOuter()      // default — morphs outer HTML
datastar.WithModeInner()      // morphs inner HTML
datastar.WithModeReplace()    // replaces outer HTML (no morph)
datastar.WithModePrepend()    // prepend to target children
datastar.WithModeAppend()     // append to target children
datastar.WithModeBefore()     // insert before target
datastar.WithModeAfter()      // insert after target
datastar.WithModeRemove()     // remove target
datastar.WithMode(mode)       // custom ElementPatchMode

// Target selector
datastar.WithSelector("#my-element")
datastar.WithSelectorID("my-element")      // shorthand for "#my-element"
datastar.WithSelectorf("#item-%d", itemID)

// View transitions
datastar.WithViewTransitions()
datastar.WithoutViewTransitions()

// Event metadata
datastar.WithPatchElementsEventID("event-123")
datastar.WithRetryDuration(2 * time.Second)
```

**Usage examples:**

```go
sse := datastar.NewSSE(w, r)

// Basic element patch (morphs into element with matching ID)
sse.PatchElements(`<div id="output">Hello World!</div>`)

// Formatted patch
sse.PatchElementf(`<div id="user-%d">%s</div>`, userID, userName)

// Append to a list
sse.PatchElements(
    `<li>New Item</li>`,
    datastar.WithSelector("#item-list"),
    datastar.WithModeAppend(),
)

// Replace inner HTML
sse.PatchElements(
    `<p>New content</p>`,
    datastar.WithSelector("#container"),
    datastar.WithModeInner(),
)

// With view transitions
sse.PatchElements(
    `<div id="page">New page content</div>`,
    datastar.WithViewTransitions(),
)
```

#### Removing Elements

```go
func (sse *ServerSentEventGenerator) RemoveElement(selector string, opts ...PatchElementOption) error
func (sse *ServerSentEventGenerator) RemoveElementByID(id string) error
func (sse *ServerSentEventGenerator) RemoveElementf(selectorFormat string, args ...any) error
```

```go
sse.RemoveElement("#temporary-message")
sse.RemoveElementByID("notification-5")
sse.RemoveElementf("#item-%d", itemID)
```

#### Patching Signals (State Updates)

```go
// Patch with raw JSON bytes
func (sse *ServerSentEventGenerator) PatchSignals(signalsContents []byte, opts ...PatchSignalsOption) error

// Marshal a Go value and patch
func (sse *ServerSentEventGenerator) MarshalAndPatchSignals(signals any, opts ...PatchSignalsOption) error

// Marshal and patch only if signals don't exist
func (sse *ServerSentEventGenerator) MarshalAndPatchSignalsIfMissing(signals any, opts ...PatchSignalsOption) error

// Raw JSON string, only if missing
func (sse *ServerSentEventGenerator) PatchSignalsIfMissingRaw(signalsJSON string) error
```

**PatchSignals Options:**

```go
datastar.WithOnlyIfMissing(true)              // only patch if signal doesn't exist
datastar.WithPatchSignalsEventID("evt-123")
datastar.WithPatchSignalsRetryDuration(2 * time.Second)
```

**Usage examples:**

```go
sse := datastar.NewSSE(w, r)

// Patch with struct (auto-marshaled to JSON)
sse.MarshalAndPatchSignals(map[string]any{
    "count":   42,
    "message": "Updated!",
    "user": map[string]any{
        "name":  "Alice",
        "email": "alice@example.com",
    },
})

// Patch with typed struct
type Signals struct {
    Count   int    `json:"count"`
    Message string `json:"message"`
}
sse.MarshalAndPatchSignals(Signals{Count: 42, Message: "Updated!"})

// Raw JSON bytes
sse.PatchSignals([]byte(`{"count": 42, "message": "Updated!"}`))

// Only set if missing (useful for defaults)
sse.MarshalAndPatchSignalsIfMissing(map[string]any{
    "theme": "dark",
})

// Remove a signal (set to null in JSON)
sse.PatchSignals([]byte(`{"obsoleteSignal": null}`))
```

#### Executing JavaScript

```go
func (sse *ServerSentEventGenerator) ExecuteScript(scriptContents string, opts ...ExecuteScriptOption) error
```

**ExecuteScript Options:**

```go
datastar.WithExecuteScriptAutoRemove(true)     // auto-remove script tag after execution (default: true)
datastar.WithExecuteScriptAttributes("type", "module")
datastar.WithExecuteScriptEventID("evt-123")
datastar.WithExecuteScriptRetryDuration(2 * time.Second)
```

```go
sse.ExecuteScript(`console.log("Hello from server!")`)
sse.ExecuteScript(`alert("Operation complete!")`)
```

#### Redirecting

```go
func (sse *ServerSentEventGenerator) Redirect(url string, opts ...ExecuteScriptOption) error
func (sse *ServerSentEventGenerator) Redirectf(format string, args ...any) error
```

```go
sse.Redirect("/dashboard")
sse.Redirectf("/users/%d/profile", userID)
```

**Note:** In Firefox, wrap redirects in `setTimeout` to avoid URL replacement. The SDK's `Redirect` method handles this automatically.

#### Replacing Browser URL

```go
func (sse *ServerSentEventGenerator) ReplaceURL(u url.URL, opts ...ExecuteScriptOption) error
func (sse *ServerSentEventGenerator) ReplaceURLQuerystring(r *http.Request, values url.Values, opts ...ExecuteScriptOption) error
```

#### Console Logging

```go
func (sse *ServerSentEventGenerator) ConsoleLog(msg string, opts ...ExecuteScriptOption) error
func (sse *ServerSentEventGenerator) ConsoleLogf(format string, args ...any) error
func (sse *ServerSentEventGenerator) ConsoleError(err error, opts ...ExecuteScriptOption) error
```

#### Dispatching Custom Events

```go
func (sse *ServerSentEventGenerator) DispatchCustomEvent(eventName string, detail any, opts ...DispatchCustomEventOption) error
```

```go
sse.DispatchCustomEvent("notification", map[string]any{
    "title": "New message",
    "body":  "You have a new message",
}, datastar.WithDispatchCustomEventSelector("#app"))
```

#### Prefetching URLs

```go
func (sse *ServerSentEventGenerator) Prefetch(urls ...string) error
```

#### Connection State

```go
func (sse *ServerSentEventGenerator) IsClosed() bool
func (sse *ServerSentEventGenerator) Context() context.Context
```

#### Low-Level Send

```go
func (sse *ServerSentEventGenerator) Send(eventType EventType, dataLines []string, opts ...SSEEventOption) error
```

### Convenience Functions for Templates

Generate `data-on:click` attribute values for use in Go templates:

```go
datastar.GetSSE("/endpoint")      // returns: @get('/endpoint')
datastar.PostSSE("/endpoint")     // returns: @post('/endpoint')
datastar.PutSSE("/endpoint")      // returns: @put('/endpoint')
datastar.PatchSSE("/endpoint")    // returns: @patch('/endpoint')
datastar.DeleteSSE("/endpoint")   // returns: @delete('/endpoint')

// With format args
datastar.GetSSE("/users/%d", userID)   // returns: @get('/users/42')
```

Useful in Go HTML templates:

```go
// In templ:
<button data-on:click={datastar.PostSSE("/increment")}>+1</button>

// In html/template:
<button data-on:click="{{postSSE "/api/submit"}}">Submit</button>
```

### Event Types (Constants)

```go
const (
    EventTypePatchElements = "datastar-patch-elements"
    EventTypePatchSignals  = "datastar-patch-signals"
)
```

### Element Patch Modes

```go
const (
    ElementPatchModeOuter   // "outer"   — morph outer HTML (default, recommended)
    ElementPatchModeInner   // "inner"   — morph inner HTML
    ElementPatchModeRemove  // "remove"  — remove element
    ElementPatchModePrepend // "prepend" — prepend to children
    ElementPatchModeAppend  // "append"  — append to children
    ElementPatchModeBefore  // "before"  — insert before as sibling
    ElementPatchModeAfter   // "after"   — insert after as sibling
    ElementPatchModeReplace // "replace" — replace outer HTML (no morph)
)
```

---

## SSE Events Protocol

When you need to manually construct SSE events (without the SDK), here is the wire format:

### `datastar-patch-elements`

```
event: datastar-patch-elements
data: elements <div id="foo">Hello world!</div>

```

Multi-line elements:

```
event: datastar-patch-elements
data: selector #target
data: mode inner
data: useViewTransition true
data: elements <div>
data: elements     Hello world!
data: elements </div>

```

Remove elements:

```
event: datastar-patch-elements
data: selector #foo
data: mode remove

```

SVG namespace:

```
event: datastar-patch-elements
data: namespace svg
data: elements <circle id="c1" cx="100" r="50" cy="75"></circle>

```

| Data Line | Values | Default |
|---|---|---|
| `data: selector #css-selector` | Any CSS selector | (matches by element ID) |
| `data: mode <mode>` | `outer`, `inner`, `replace`, `prepend`, `append`, `before`, `after`, `remove` | `outer` |
| `data: namespace <ns>` | `svg`, `mathml` | (HTML) |
| `data: useViewTransition true` | `true`/`false` | `false` |
| `data: elements <html>` | HTML content | (required unless mode is remove) |

### `datastar-patch-signals`

```
event: datastar-patch-signals
data: signals {count: 42, name: "Alice"}

```

Only-if-missing:

```
event: datastar-patch-signals
data: onlyIfMissing true
data: signals {theme: "dark", lang: "en"}

```

Remove signals (set to null):

```
event: datastar-patch-signals
data: signals {obsoleteSignal: null}

```

**Important:** Every SSE event MUST end with a double newline (`\n\n`).

---

## Complete Patterns & Examples

### Pattern 1: CRUD Application

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "sync"

    "github.com/starfederation/datastar-go/datastar"
)

type Todo struct {
    ID   int    `json:"id"`
    Text string `json:"text"`
    Done bool   `json:"done"`
}

var (
    todos   = []Todo{}
    todosMu sync.RWMutex
    nextID  = 1
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /", indexHandler)
    mux.HandleFunc("POST /todos", addTodoHandler)
    mux.HandleFunc("PUT /todos/toggle", toggleTodoHandler)
    mux.HandleFunc("DELETE /todos", deleteTodoHandler)
    http.ListenAndServe(":8080", mux)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")
    fmt.Fprint(w, `<!DOCTYPE html>
<html>
<head>
    <script type="module" src="https://cdn.jsdelivr.net/gh/starfederation/[email protected]/bundles/datastar.js"></script>
</head>
<body>
    <div data-signals="{newTodo: ''}">
        <h1>Todo App</h1>
        <input data-bind:new-todo placeholder="Add a todo..." />
        <button data-on:click="@post('/todos')" data-attr:disabled="$newTodo === ''">
            Add
        </button>
        <div id="todo-list" data-init="@get('/todos')"></div>
    </div>
</body>
</html>`)
}

type AddTodoSignals struct {
    NewTodo string `json:"newTodo"`
}

func addTodoHandler(w http.ResponseWriter, r *http.Request) {
    signals := &AddTodoSignals{}
    if err := datastar.ReadSignals(r, signals); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    todosMu.Lock()
    todo := Todo{ID: nextID, Text: signals.NewTodo, Done: false}
    todos = append(todos, todo)
    nextID++
    todosMu.Unlock()

    sse := datastar.NewSSE(w, r)
    sse.PatchElements(renderTodoList())
    sse.MarshalAndPatchSignals(map[string]any{"newTodo": ""})
}

func toggleTodoHandler(w http.ResponseWriter, r *http.Request) {
    signals := &struct {
        ToggleID int `json:"toggleId"`
    }{}
    if err := datastar.ReadSignals(r, signals); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    todosMu.Lock()
    for i := range todos {
        if todos[i].ID == signals.ToggleID {
            todos[i].Done = !todos[i].Done
        }
    }
    todosMu.Unlock()

    sse := datastar.NewSSE(w, r)
    sse.PatchElements(renderTodoList())
}

func deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
    signals := &struct {
        DeleteID int `json:"deleteId"`
    }{}
    if err := datastar.ReadSignals(r, signals); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    todosMu.Lock()
    for i := range todos {
        if todos[i].ID == signals.DeleteID {
            todos = append(todos[:i], todos[i+1:]...)
            break
        }
    }
    todosMu.Unlock()

    sse := datastar.NewSSE(w, r)
    sse.PatchElements(renderTodoList())
}

func renderTodoList() string {
    todosMu.RLock()
    defer todosMu.RUnlock()

    html := `<div id="todo-list">`
    for _, t := range todos {
        checked := ""
        style := ""
        if t.Done {
            checked = "checked"
            style = "text-decoration: line-through;"
        }
        html += fmt.Sprintf(`
            <div id="todo-%d" style="%s">
                <input type="checkbox" %s
                    data-on:click="$toggleId = %d; @put('/todos/toggle')" />
                <span>%s</span>
                <button data-on:click="$deleteId = %d; @delete('/todos')">×</button>
            </div>`,
            t.ID, style, checked, t.ID, t.Text, t.ID)
    }
    html += `</div>`
    return html
}
```

### Pattern 2: Real-Time Updates (Long-Lived SSE)

```go
func clockHandler(w http.ResponseWriter, r *http.Request) {
    sse := datastar.NewSSE(w, r)
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-r.Context().Done():
            return
        case t := <-ticker.C:
            if sse.IsClosed() {
                return
            }
            sse.PatchElements(fmt.Sprintf(
                `<div id="clock">%s</div>`,
                t.Format("15:04:05"),
            ))
        }
    }
}
```

**Frontend for real-time:**

```html
<div id="clock" data-init="@get('/clock')">--:--:--</div>
```

### Pattern 3: Loading Indicators

```html
<div data-signals="{query: ''}">
    <input data-bind:query placeholder="Search..." />
    <button
        data-on:click="@post('/search')"
        data-indicator:searching
        data-attr:disabled="$searching"
    >
        <span data-show="!$searching">Search</span>
        <span data-show="$searching">Searching...</span>
    </button>
    <div id="results"></div>
</div>
```

### Pattern 4: Using with templ (Go Templating)

```go
// signals.go
type CounterSignals struct {
    Count int `json:"count"`
}

// handler.go
func incrementHandler(w http.ResponseWriter, r *http.Request) {
    signals := &CounterSignals{}
    if err := datastar.ReadSignals(r, signals); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    signals.Count++

    sse := datastar.NewSSE(w, r)
    sse.PatchElementTempl(counterComponent(*signals))
    sse.MarshalAndPatchSignals(signals)
}
```

```go
// counter.templ
templ counterComponent(signals CounterSignals) {
    <div id="counter" data-signals={ templ.JSONString(signals) }>
        <span data-text="$count"></span>
        <button data-on:click={ datastar.PostSSE("/increment") }>+1</button>
    </div>
}
```

### Pattern 5: Form Handling with Validation

```html
<form data-signals="{form: {email: '', password: ''}}">
    <input data-bind:form.email type="email" placeholder="Email" />
    <input data-bind:form.password type="password" placeholder="Password" />
    <button
        data-on:click="@post('/login')"
        data-indicator:submitting
        data-attr:disabled="$submitting || $form.email === '' || $form.password === ''"
    >
        Login
    </button>
    <div id="errors"></div>
</form>
```

```go
type LoginForm struct {
    Form struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    } `json:"form"`
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
    signals := &LoginForm{}
    if err := datastar.ReadSignals(r, signals); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    sse := datastar.NewSSE(w, r)

    if signals.Form.Email == "" || signals.Form.Password == "" {
        sse.PatchElements(`<div id="errors" class="error">All fields required</div>`)
        return
    }

    // Validate credentials...
    if !authenticate(signals.Form.Email, signals.Form.Password) {
        sse.PatchElements(`<div id="errors" class="error">Invalid credentials</div>`)
        return
    }

    sse.Redirect("/dashboard")
}
```

### Pattern 6: Infinite Scroll / Load More

```html
<div id="feed">
    <!-- items here -->
</div>
<div
    data-signals:page="1"
    data-on-intersect__once="$page++; @get('/feed')"
    id="load-trigger"
>
    Loading more...
</div>
```

```go
func feedHandler(w http.ResponseWriter, r *http.Request) {
    signals := &struct {
        Page int `json:"page"`
    }{}
    datastar.ReadSignals(r, signals)

    items := fetchItems(signals.Page, 20)

    sse := datastar.NewSSE(w, r)

    for _, item := range items {
        sse.PatchElements(
            fmt.Sprintf(`<div id="item-%d">%s</div>`, item.ID, item.Title),
            datastar.WithSelector("#feed"),
            datastar.WithModeAppend(),
        )
    }

    // Reset trigger for next page
    if len(items) == 20 {
        sse.PatchElements(
            fmt.Sprintf(`<div
                data-signals:page="%d"
                data-on-intersect__once="$page++; @get('/feed')"
                id="load-trigger"
            >Loading more...</div>`, signals.Page+1),
        )
    } else {
        sse.RemoveElementByID("load-trigger")
    }
}
```

---

## Security Guidelines

### Always Escape User Input

Signal values are visible in source code and can be modified by users. Never trust signal values; always validate on the backend.

```go
// ALWAYS validate signals server-side
func handler(w http.ResponseWriter, r *http.Request) {
    store := &Store{}
    if err := datastar.ReadSignals(r, store); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Validate and sanitize
    if store.Count < 0 || store.Count > 1000 {
        http.Error(w, "Invalid count", http.StatusBadRequest)
        return
    }
}
```

### Escape HTML in Templates

When injecting user content into HTML that will be patched, always escape it to prevent XSS:

```go
import "html"

sse.PatchElements(fmt.Sprintf(
    `<div id="msg">%s</div>`,
    html.EscapeString(userInput),
))
```

### Use `data-ignore` for Untrusted Content

```html
<div data-ignore>
    <!-- User-generated content that should not be processed by Datastar -->
</div>
```

### Signals Starting with `_` Are Local

Signals prefixed with `_` are not sent to the backend and serve as local/private state:

```html
<div data-signals:_menuOpen="false">
    <!-- $\_menuOpen is never sent to server -->
</div>
```

---

## Common Patterns & How-Tos

### Redirect from Backend

```go
sse := datastar.NewSSE(w, r)
sse.Redirect("/new-page")
// Or with format:
sse.Redirectf("/users/%d", userID)
```

### Debounced Search

```html
<input
    data-bind:query
    data-on:input__debounce.300ms="@get('/search')"
    placeholder="Search..."
/>
<div id="results"></div>
```

### Polling / Auto-Refresh

```html
<div data-on-interval__duration.5s="@get('/status')">
    <div id="status">Loading...</div>
</div>
```

### Conditional Backend Requests

```html
<button
    data-on:click="$count > 0 && @post('/checkout')"
    data-attr:disabled="$count === 0"
>
    Checkout
</button>
```

### Multiple Elements in One Response

```go
sse := datastar.NewSSE(w, r)
sse.PatchElements(`<div id="header">Updated Header</div>`)
sse.PatchElements(`<div id="sidebar">Updated Sidebar</div>`)
sse.PatchElements(`<div id="content">Updated Content</div>`)
sse.MarshalAndPatchSignals(map[string]any{
    "lastUpdated": time.Now().Format(time.RFC3339),
})
```

### Using with Chi Router

```go
import "github.com/go-chi/chi/v5"

r := chi.NewRouter()
r.Get("/", indexHandler)
r.Post("/api/submit", submitHandler)
r.Get("/api/data", dataHandler)
http.ListenAndServe(":8080", r)
```

### Using with Gorilla Sessions for Per-User State

```go
import "github.com/gorilla/sessions"

var store = sessions.NewCookieStore([]byte("secret-key"))

func handler(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "session")
    count, _ := session.Values["count"].(int)
    count++
    session.Values["count"] = count
    session.Save(r, w)

    sse := datastar.NewSSE(w, r)
    sse.MarshalAndPatchSignals(map[string]any{"count": count})
}
```

### SSE Compression

```go
sse := datastar.NewSSE(w, r,
    datastar.WithCompression(
        datastar.WithGzip(datastar.WithGzipLevel(6)),
        datastar.WithBrotli(datastar.WithBrotliLevel(4)),
        datastar.WithClientPriority(), // respect client's Accept-Encoding preference
    ),
)
```

### Check If SSE Connection Is Still Open

```go
func streamHandler(w http.ResponseWriter, r *http.Request) {
    sse := datastar.NewSSE(w, r)

    for i := 0; i < 100; i++ {
        if sse.IsClosed() {
            return // Client disconnected
        }
        sse.PatchElementf(`<div id="progress">%d%%</div>`, i)
        time.Sleep(100 * time.Millisecond)
    }
}
```

---

## Northstar Best Practices (Production Boilerplate)

> Source: [github.com/zangster300/northstar](https://github.com/zangster300/northstar) — a production-grade boilerplate for real-time Hypermedia applications with Datastar, Go, NATS, Templ, and Tailwind CSS. Star count: 200+. Referenced in the official templ documentation.

### Project Structure

Northstar uses a **feature-based architecture** where each feature is a self-contained module with its own routes, handlers, services, and templ templates. This is the recommended pattern for Datastar + Go applications beyond trivial size.

```
northstar/
├── cmd/web/
│   ├── main.go                     # Entry point: signal handling, server setup, graceful shutdown
│   ├── build/main.go               # esbuild bundler for web components (TypeScript → JS)
│   └── downloader/main.go          # CLI tool to fetch Datastar + DaisyUI from CDN
├── config/
│   ├── config.go                   # Shared config struct + env loading
│   ├── config_dev.go               # //go:build dev — sets Environment = Dev
│   └── config_prod.go              # //go:build !dev — sets Environment = Prod
├── features/                       # ← EACH FEATURE IS A SELF-CONTAINED MODULE
│   ├── common/
│   │   ├── components/             # Shared templ components (icons, indicators, nav)
│   │   └── layouts/                # Base HTML layout (<!DOCTYPE>, head, Datastar script)
│   ├── counter/                    # Feature: counter (signals + sessions)
│   │   ├── routes.go               # Route registration only
│   │   ├── handlers.go             # HTTP handlers (create SSE, patch signals/elements)
│   │   └── pages/
│   │       └── counter.templ       # Signal struct + templ components
│   ├── index/                      # Feature: todo MVC (NATS KV + real-time watchers)
│   │   ├── routes.go
│   │   ├── handlers.go
│   │   ├── services/
│   │   │   └── todo_service.go     # Business logic + NATS KV persistence
│   │   ├── components/
│   │   │   └── todo.templ          # Complex UI components
│   │   └── pages/
│   │       └── index.templ         # Page-level template
│   └── monitor/                    # Feature: real-time system monitoring (streaming SSE)
│       ├── routes.go
│       ├── handlers.go
│       └── pages/
│           └── monitor.templ
├── nats/
│   └── nats.go                     # Embedded NATS server setup with JetStream
├── router/
│   └── router.go                   # Mounts all features, static files, dev hot-reload
├── web/
│   ├── resources/
│   │   ├── static/                 # Built assets (Datastar JS, CSS, bundled components)
│   │   └── styles/                 # TailwindCSS source
│   └── libs/                       # Web component source (TypeScript)
│       ├── web-components/         # Vanilla custom elements
│       └── lit/                    # Lit-based web components
├── Dockerfile                      # Multi-stage build → scratch container
├── Taskfile.yml                    # Task runner: live reload, build, debug
└── go.mod
```

**Key structural rules:**
- Each feature has its own Go package under `features/`.
- Route registration is in `routes.go` — returns `error`, takes a `chi.Router` and any shared dependencies (session store, NATS, etc.).
- Handlers live in `handlers.go` — a struct with methods, constructed via `NewHandlers(...)`.
- Services (if needed) handle business logic and persistence — handlers call services, never access storage directly.
- Templ templates are colocated in `pages/` and `components/` subdirs.
- Signal structs (used for JSON serialization) are defined in the templ file alongside the templates that use them.

### Application Startup Pattern

Northstar uses `errgroup` for concurrent lifecycle management with graceful shutdown:

```go
package main

import (
    "context"
    "fmt"
    "log/slog"
    "net"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "myapp/config"
    "myapp/nats"
    "myapp/router"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/gorilla/sessions"
    "golang.org/x/sync/errgroup"
)

func main() {
    ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer cancel()

    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: config.Global.LogLevel,
    }))
    slog.SetDefault(logger)

    if err := run(ctx); err != nil && err != http.ErrServerClosed {
        slog.Error("error running server", "error", err)
        os.Exit(1)
    }
}

func run(ctx context.Context) error {
    addr := fmt.Sprintf("%s:%s", config.Global.Host, config.Global.Port)
    slog.Info("server started", "addr", addr)
    defer slog.Info("server shutdown complete")

    eg, egctx := errgroup.WithContext(ctx)

    r := chi.NewMux()
    r.Use(middleware.Logger, middleware.Recoverer)

    // Session store for per-user state
    sessionStore := sessions.NewCookieStore([]byte(config.Global.SessionSecret))
    sessionStore.MaxAge(86400 * 30)
    sessionStore.Options.Path = "/"
    sessionStore.Options.HttpOnly = true
    sessionStore.Options.Secure = false // set true in production with HTTPS
    sessionStore.Options.SameSite = http.SameSiteLaxMode

    // Start embedded NATS (optional — for real-time features)
    ns, err := nats.SetupNATS(ctx)
    if err != nil {
        return err
    }

    // Mount all feature routes
    if err := router.SetupRoutes(egctx, r, sessionStore, ns); err != nil {
        return fmt.Errorf("error setting up routes: %w", err)
    }

    srv := &http.Server{
        Addr:    addr,
        Handler: r,
        BaseContext: func(l net.Listener) context.Context {
            return egctx // propagate cancellation to all handlers
        },
        ErrorLog: slog.NewLogLogger(slog.Default().Handler(), slog.LevelError),
    }

    // Start server
    eg.Go(func() error {
        err := srv.ListenAndServe()
        if err != nil && err != http.ErrServerClosed {
            return fmt.Errorf("server error: %w", err)
        }
        return nil
    })

    // Graceful shutdown
    eg.Go(func() error {
        <-egctx.Done()
        shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        slog.Debug("shutting down server...")
        return srv.Shutdown(shutdownCtx)
    })

    return eg.Wait()
}
```

**Key patterns:**
- `signal.NotifyContext` listens for SIGINT/SIGTERM for clean shutdown.
- `errgroup` manages concurrent goroutines (server + shutdown watcher).
- `srv.BaseContext` propagates the errgroup context to all request handlers — when the server shuts down, all SSE connections close cleanly via `r.Context().Done()`.
- `slog` with JSON handler is used throughout for structured logging.

### Feature Module Pattern (Routes → Handlers → Services)

Each feature follows this three-layer pattern:

**routes.go** — Pure route registration, no logic:

```go
package counter

import (
    "github.com/go-chi/chi/v5"
    "github.com/gorilla/sessions"
)

func SetupRoutes(router chi.Router, sessionStore sessions.Store) error {
    handlers := NewHandlers(sessionStore)

    router.Get("/counter", handlers.CounterPage)
    router.Get("/counter/data", handlers.CounterData)

    router.Route("/counter/increment", func(r chi.Router) {
        r.Post("/global", handlers.IncrementGlobal)
        r.Post("/user", handlers.IncrementUser)
    })

    return nil
}
```

**handlers.go** — HTTP layer, creates SSE, delegates to service:

```go
package counter

import (
    "net/http"
    "sync/atomic"

    "myapp/features/counter/pages"

    "github.com/Jeffail/gabs/v2"
    "github.com/gorilla/sessions"
    "github.com/starfederation/datastar-go/datastar"
)

type Handlers struct {
    globalCounter atomic.Uint32
    sessionStore  sessions.Store
}

func NewHandlers(sessionStore sessions.Store) *Handlers {
    return &Handlers{sessionStore: sessionStore}
}

// Page render — plain templ, no SSE
func (h *Handlers) CounterPage(w http.ResponseWriter, r *http.Request) {
    if err := pages.CounterPage().Render(r.Context(), w); err != nil {
        http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    }
}

// SSE data endpoint — sends initial fragment via PatchElementTempl
func (h *Handlers) CounterData(w http.ResponseWriter, r *http.Request) {
    userCount, _, err := h.getUserValue(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    store := pages.CounterSignals{Global: h.globalCounter.Load(), User: userCount}

    if err := datastar.NewSSE(w, r).PatchElementTempl(pages.Counter(store)); err != nil {
        http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    }
}

// Action endpoint — reads signals, mutates state, patches signals back
func (h *Handlers) IncrementUser(w http.ResponseWriter, r *http.Request) {
    val, sess, err := h.getUserValue(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    val++
    sess.Values["count"] = val
    if err := sess.Save(r, w); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    update := gabs.New()
    update.Set(h.globalCounter.Add(1), "global")
    update.Set(val, "user")

    if err := datastar.NewSSE(w, r).MarshalAndPatchSignals(update); err != nil {
        http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    }
}
```

**Key handler patterns:**
- Page handlers render templ directly (no SSE).
- Data/init handlers use `PatchElementTempl` to send the initial fragment.
- Action handlers (POST/PUT/DELETE) use `MarshalAndPatchSignals` to update reactive state without replacing DOM.
- Use `gabs.New()` from `Jeffail/gabs/v2` for ad-hoc JSON construction when you don't need a full struct.
- `datastar.ReadSignals(r, &store)` to deserialize incoming signals into a typed struct.

### Router Composition Pattern

Mount all features in a central router with `errors.Join` for combined error handling:

```go
package router

import (
    "errors"
    "fmt"
    "net/http"

    "myapp/config"
    counterFeature "myapp/features/counter"
    indexFeature "myapp/features/index"
    monitorFeature "myapp/features/monitor"
    "myapp/web/resources"

    "github.com/go-chi/chi/v5"
    "github.com/gorilla/sessions"
    "github.com/delaneyj/toolbelt/embeddednats"
)

func SetupRoutes(ctx context.Context, router chi.Router, sessionStore *sessions.CookieStore, ns *embeddednats.Server) error {
    if config.Global.Environment == config.Dev {
        setupReload(router)
    }

    router.Handle("/static/*", resources.Handler())

    if err := errors.Join(
        indexFeature.SetupRoutes(router, sessionStore, ns),
        counterFeature.SetupRoutes(router, sessionStore),
        monitorFeature.SetupRoutes(router),
    ); err != nil {
        return fmt.Errorf("error setting up routes: %w", err)
    }

    return nil
}
```

### RESTful Route Naming for Datastar Endpoints

Northstar uses **nested chi routes** with RESTful verbs matching the Datastar actions:

```go
router.Route("/api", func(apiRouter chi.Router) {
    apiRouter.Route("/todos", func(todosRouter chi.Router) {
        todosRouter.Get("/", handlers.TodosSSE)          // @get — long-lived SSE stream
        todosRouter.Put("/reset", handlers.ResetTodos)    // @put
        todosRouter.Put("/cancel", handlers.CancelEdit)   // @put
        todosRouter.Put("/mode/{mode}", handlers.SetMode) // @put with path param

        todosRouter.Route("/{idx}", func(todoRouter chi.Router) {
            todoRouter.Post("/toggle", handlers.ToggleTodo) // @post
            todoRouter.Get("/edit", handlers.StartEdit)     // @get
            todoRouter.Put("/edit", handlers.SaveEdit)      // @put
            todoRouter.Delete("/", handlers.DeleteTodo)     // @delete
        })
    })
})
```

In templ, the corresponding Datastar actions use the SDK's helper functions:

```go
// templ file
data-on:click={ datastar.PostSSE("/api/todos/%d/toggle", i) }
data-on:click={ datastar.PutSSE("/api/todos/mode/%d", modeIdx) }
data-on:click={ datastar.DeleteSSE("/api/todos/%d", i) }
data-on:click={ datastar.GetSSE("/api/todos/%d/edit", i) }
```

### Templ + Datastar Integration Patterns

#### Signal Structs Defined Alongside Templates

Define signal structs in the same templ file that uses them:

```go
// pages/counter.templ
package pages

import "github.com/starfederation/datastar-go/datastar"

type CounterSignals struct {
    Global uint32 `json:"global"`
    User   uint32 `json:"user"`
}

templ Counter(signals CounterSignals) {
    <div
        id="container"
        data-signals={ templ.JSONString(signals) }
        class="flex flex-col gap-4"
    >
        @CounterButtons()
        @CounterCounts()
    </div>
}
```

**Key technique:** Use `templ.JSONString(signals)` to serialize Go structs into `data-signals`. This is type-safe and avoids manual JSON construction.

#### Page Shell + data-init Pattern

The page renders an empty container shell, then immediately loads data via SSE:

```go
templ CounterPage() {
    @layouts.Base("Counter") {
        @components.Navigation(components.PageCounter)
        <article class="prose mx-auto m-2">
            <div id="container" data-init={ datastar.GetSSE("/counter/data") }></div>
        </article>
    }
}
```

The `data-init` fires `@get('/counter/data')` on page load, which returns an SSE stream that patches the actual content into `#container` via `PatchElementTempl`.

#### Long-Lived SSE with Disabled Request Cancellation

For real-time features where the SSE stream should persist across user interactions (not be cancelled by subsequent Datastar requests):

```go
templ IndexPage(title string) {
    @layouts.Base(title) {
        <div class="flex flex-col w-full min-h-screen">
            @components.Navigation(components.PageIndex)
            <div id="todos-container"
                 data-init="@get('/api/todos',{requestCancellation: 'disabled'})"></div>
        </div>
    }
}
```

**Critical:** `{requestCancellation: 'disabled'}` prevents the SSE connection from being closed when other `@get`/`@post` calls are made. Use this for any feature with a persistent watcher (real-time updates, notifications, live data).

#### Loading Indicators with data-indicator

Northstar uses a reusable indicator component pattern:

```go
// Shared component
templ SseIndicator(signalName string) {
    <div class="loading-dots text-primary"
         data-class={ fmt.Sprintf("{'loading ml-4': $%s}", signalName) }></div>
}

// Usage in a todo row — each row gets its own indicator signal
templ TodoRow(mode TodoViewMode, todo *Todo, i int, isEditing bool) {
    {{
        fetchingSignalName := fmt.Sprintf("fetching%d", i)
    }}
    <li class="flex items-center gap-8 p-2 group">
        <label
            data-on:click={ datastar.PostSSE("/api/todos/%d/toggle", i) }
            data-indicator={ fetchingSignalName }
        >
            <!-- checkbox icon -->
        </label>
        @SseIndicator(fetchingSignalName)
        <button
            data-on:click={ datastar.DeleteSSE("/api/todos/%d", i) }
            data-indicator={ fetchingSignalName }
            data-attrs-disabled={ fetchingSignalName + "" }
        >
            <!-- delete icon -->
        </button>
    </li>
}
```

**Pattern:** Each interactive element references the same `data-indicator` signal name. While any request from that element is in flight, the signal is `true`, driving both loading spinners and disabled states.

#### Inline Form Handling with data-bind + Keyboard Events

```go
templ TodoInput(i int) {
    <input
        id="todoInput"
        class="flex-1 w-full italic input input-bordered input-lg"
        placeholder="What needs to be done?"
        data-bind:input
        data-on:keydown={ fmt.Sprintf(`
            if (evt.key !== 'Enter' || !$input.trim().length) return;
            %s;
            $input = '';
        `, datastar.PutSSE("/api/todos/%d/edit", i)) }
    />
}
```

**Pattern:** Bind an input to a local signal with `data-bind:input`, then handle Enter key inline — execute the SSE action and clear the signal. No form element needed.

#### Click-Outside to Cancel

```go
if i >= 0 {
    data-on:click__outside={ datastar.PutSSE("/api/todos/cancel") }
}
```

Use `data-on:click__outside` to cancel editing when clicking away from an input.

### Real-Time Streaming with NATS KV Watchers

The most powerful Northstar pattern uses NATS KV watchers to push updates to all connected clients in real time.

**Service layer — NATS KV setup:**

```go
package services

import (
    "context"
    "encoding/json"
    "time"

    "github.com/delaneyj/toolbelt/embeddednats"
    "github.com/nats-io/nats.go/jetstream"
)

type TodoService struct {
    kv    jetstream.KeyValue
    store sessions.Store
}

func NewTodoService(ns *embeddednats.Server, store sessions.Store) (*TodoService, error) {
    nc, err := ns.Client()
    if err != nil {
        return nil, fmt.Errorf("error creating nats client: %w", err)
    }

    js, err := jetstream.New(nc)
    if err != nil {
        return nil, fmt.Errorf("error creating jetstream client: %w", err)
    }

    kv, err := js.CreateOrUpdateKeyValue(context.Background(), jetstream.KeyValueConfig{
        Bucket:      "todos",
        Description: "Datastar Todos",
        Compression: true,
        TTL:         time.Hour,
        MaxBytes:    16 * 1024 * 1024,
    })
    if err != nil {
        return nil, fmt.Errorf("error creating key value: %w", err)
    }

    return &TodoService{kv: kv, store: store}, nil
}

func (s *TodoService) WatchUpdates(ctx context.Context, sessionID string) (jetstream.KeyWatcher, error) {
    return s.kv.Watch(ctx, sessionID)
}

func (s *TodoService) SaveMVC(ctx context.Context, sessionID string, mvc *TodoMVC) error {
    b, err := json.Marshal(mvc)
    if err != nil {
        return fmt.Errorf("failed to marshal mvc: %w", err)
    }
    _, err = s.kv.Put(ctx, sessionID, b)
    return err
}
```

**Handler — SSE watcher loop:**

```go
func (h *Handlers) TodosSSE(w http.ResponseWriter, r *http.Request) {
    sessionID, mvc, err := h.todoService.GetSessionMVC(w, r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    sse := datastar.NewSSE(w, r)

    ctx := r.Context()
    watcher, err := h.todoService.WatchUpdates(ctx, sessionID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer watcher.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case entry := <-watcher.Updates():
            if entry == nil {
                continue
            }
            if err := json.Unmarshal(entry.Value(), mvc); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            c := components.TodosMVCView(mvc)
            if err := sse.PatchElementTempl(c); err != nil {
                if err := sse.ConsoleError(err); err != nil {
                    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
                }
                return
            }
        }
    }
}
```

**How it works:** Action handlers (toggle, edit, delete) only write to NATS KV via `SaveMVC`. The watcher loop in `TodosSSE` picks up the change and streams the updated HTML fragment to the client. This decouples mutations from rendering and enables multi-client sync.

#### Embedded NATS Server

Northstar embeds a NATS server directly into the application binary — no separate NATS deployment needed for development or simple deployments:

```go
package nats

import (
    "context"
    "fmt"
    "log/slog"

    "github.com/delaneyj/toolbelt"
    "github.com/delaneyj/toolbelt/embeddednats"
    natsserver "github.com/nats-io/nats-server/v2/server"
)

func SetupNATS(ctx context.Context) (*embeddednats.Server, error) {
    natsPort, err := getFreeNatsPort()
    if err != nil {
        return nil, fmt.Errorf("error obtaining NATS port: %w", err)
    }

    ns, err := embeddednats.New(ctx, embeddednats.WithNATSServerOptions(&natsserver.Options{
        JetStream: true,
        NoSigs:    true,
        Port:      natsPort,
        StoreDir:  "data/nats",
    }))
    if err != nil {
        return nil, fmt.Errorf("error creating embedded nats server: %w", err)
    }

    ns.WaitForServer()
    slog.Info("NATS started", "port", natsPort)
    return ns, nil
}
```

Use `NATS_PORT` env var to pin the port, or let it find a free port automatically. Data persists in `data/nats`.

### Real-Time Signal Streaming (Without NATS)

For simpler real-time use cases without NATS, use a ticker-based streaming pattern:

```go
func (h *Handlers) MonitorEvents(w http.ResponseWriter, r *http.Request) {
    memT := time.NewTicker(time.Second)
    defer memT.Stop()

    cpuT := time.NewTicker(time.Second)
    defer cpuT.Stop()

    sse := datastar.NewSSE(w, r)
    for {
        select {
        case <-r.Context().Done():
            slog.Debug("client disconnected")
            return

        case <-memT.C:
            vm, _ := mem.VirtualMemory()
            memStats := pages.SystemMonitorSignals{
                MemTotal:       humanize.Bytes(vm.Total),
                MemUsed:        humanize.Bytes(vm.Used),
                MemUsedPercent: fmt.Sprintf("%.2f%%", vm.UsedPercent),
            }
            if err := sse.MarshalAndPatchSignals(memStats); err != nil {
                return
            }

        case <-cpuT.C:
            cpuTimes, _ := cpu.Times(false)
            cpuStats := pages.SystemMonitorSignals{
                CpuUser:   formatDuration(cpuTimes[0].User),
                CpuSystem: formatDuration(cpuTimes[0].System),
                CpuIdle:   formatDuration(cpuTimes[0].Idle),
            }
            if err := sse.MarshalAndPatchSignals(cpuStats); err != nil {
                return
            }
        }
    }
}
```

**Corresponding templ — signals initialized inline, updated by SSE:**

```go
templ MonitorPage() {
    @layouts.Base("System Monitoring") {
        <div
            id="container"
            data-init={ datastar.GetSSE("/monitor/events") }
            data-signals="{memTotal:'', memUsed:'', memUsedPercent:'', cpuUser:'', cpuSystem:'', cpuIdle:''}"
        >
            <div id="mem">
                <h1>Memory</h1>
                <p>Total: <span data-text="$memTotal"></span></p>
                <p>Used: <span data-text="$memUsed"></span></p>
                <p>Used (%): <span data-text="$memUsedPercent"></span></p>
            </div>
            <div id="cpu">
                <h1>CPU</h1>
                <p>User: <span data-text="$cpuUser"></span></p>
                <p>System: <span data-text="$cpuSystem"></span></p>
                <p>Idle: <span data-text="$cpuIdle"></span></p>
            </div>
        </div>
    }
}
```

**Key pattern:** Use `omitempty` JSON tags on the signal struct so partial updates only patch the signals that changed:

```go
type SystemMonitorSignals struct {
    MemTotal       string `json:"memTotal,omitempty"`
    MemUsed        string `json:"memUsed,omitempty"`
    MemUsedPercent string `json:"memUsedPercent,omitempty"`
    CpuUser        string `json:"cpuUser,omitempty"`
    CpuSystem      string `json:"cpuSystem,omitempty"`
    CpuIdle        string `json:"cpuIdle,omitempty"`
}
```

### Build Tag Configuration (Dev vs Prod)

Northstar uses Go build tags for environment-specific behavior:

```go
// config/config.go — shared
type Environment string

const (
    Dev  Environment = "dev"
    Prod Environment = "prod"
)

type Config struct {
    Environment   Environment
    Host          string
    Port          string
    LogLevel      slog.Level
    SessionSecret string
}

var Global *Config

func init() {
    sync.OnceFunc(func() { Global = Load() })()
}

func loadBase() *Config {
    godotenv.Load()
    return &Config{
        Host:          getEnv("HOST", "0.0.0.0"),
        Port:          getEnv("PORT", "8080"),
        SessionSecret: getEnv("SESSION_SECRET", "session-secret"),
    }
}
```

```go
// config/config_dev.go
//go:build dev
package config

func Load() *Config {
    cfg := loadBase()
    cfg.Environment = Dev
    return cfg
}
```

```go
// config/config_prod.go
//go:build !dev
package config

func Load() *Config {
    cfg := loadBase()
    cfg.Environment = Prod
    return cfg
}
```

Build with tags: `go build -tags=dev` for development, `go build -tags=prod` (or no tag) for production.

### Static Asset Handling (Dev vs Embedded)

Northstar uses `hashfs` for cache-busting in production and direct file serving in dev:

```go
// web/resources/static_prod.go
//go:build !dev
package resources

import (
    "embed"
    "net/http"
    "github.com/benbjohnson/hashfs"
)

//go:embed static
var StaticDirectory embed.FS
var StaticSys = hashfs.NewFS(StaticDirectory)

func Handler() http.Handler {
    return hashfs.FileServer(StaticSys)
}

func StaticPath(path string) string {
    return "/" + StaticSys.HashName("static/" + path)
}
```

```go
// web/resources/static_dev.go
//go:build dev
package resources

import (
    "net/http"
    "os"
)

func Handler() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Cache-Control", "no-store")
        http.StripPrefix("/static/", http.FileServerFS(os.DirFS(StaticDirectoryPath))).ServeHTTP(w, r)
    })
}

func StaticPath(path string) string {
    return "/static/" + path
}
```

In templ templates, always use `resources.StaticPath(...)` for asset URLs:

```go
<script defer type="module" src={ resources.StaticPath("datastar/datastar.js") }></script>
<link href={ resources.StaticPath("index.css") } rel="stylesheet" type="text/css"/>
```

### Self-Hosting Datastar (Recommended)

Northstar self-hosts the Datastar JS bundle rather than using a CDN. A downloader utility fetches the latest version:

```go
files := map[string]string{
    "https://raw.githubusercontent.com/starfederation/datastar/develop/bundles/datastar.js":
        resources.StaticDirectoryPath + "/datastar/datastar.js",
    "https://raw.githubusercontent.com/starfederation/datastar/develop/bundles/datastar.js.map":
        resources.StaticDirectoryPath + "/datastar/datastar.js.map",
}
```

This ensures version pinning, offline development, and no third-party CDN dependency. The JS is then embedded into the binary at build time via `//go:embed static`.

### Hot Reload in Development

Northstar implements SSE-based hot reload via a dedicated endpoint:

```go
func setupReload(router chi.Router) {
    reloadChan := make(chan struct{}, 1)
    var hotReloadOnce sync.Once

    // SSE endpoint — browser connects and waits for reload signal
    router.Get("/reload", func(w http.ResponseWriter, r *http.Request) {
        sse := datastar.NewSSE(w, r)
        reload := func() { sse.ExecuteScript("window.location.reload()") }
        hotReloadOnce.Do(reload)
        select {
        case <-reloadChan:
            reload()
        case <-r.Context().Done():
        }
    })

    // Trigger endpoint — called by esbuild or Air on file change
    router.Get("/hotreload", func(w http.ResponseWriter, r *http.Request) {
        select {
        case reloadChan <- struct{}{}:
        default:
        }
        w.WriteHeader(http.StatusOK)
    })
}
```

In the base layout, conditionally include the reload connection in dev mode:

```go
templ Base(title string) {
    <!DOCTYPE html>
    <html lang="en">
        <head>...</head>
        <body>
            if config.Global.Environment == config.Dev {
                <div data-init="@get('/reload', {retryMaxCount: 1000, retryInterval:20, retryMaxWaitMs:200})"></div>
            }
            { children... }
        </body>
    </html>
}
```

### Web Components + Datastar Integration

Northstar shows how to bridge custom elements (Vanilla or Lit) with Datastar signals:

```go
// templ — pass signal values as attributes, listen for custom events
<reverse-component
    data-on:reverse="$_reversed = evt.detail.value"
    data-attr:name="$_name">
</reverse-component>
```

```typescript
// TypeScript — vanilla custom element
class ReverseComponent extends HTMLElement {
  static get observedAttributes() { return ["name"]; }

  attributeChangedCallback(name: string, oldValue: string, newValue: string) {
    const reversed = [...newValue].reverse().join("");
    this.dispatchEvent(new CustomEvent("reverse", { detail: { value: reversed } }));
  }
}
customElements.define("reverse-component", ReverseComponent);
```

**Pattern:** Use `data-attr:propname="$signal"` to flow Datastar signals into web component attributes. Use `data-on:customevent="$signal = evt.detail.value"` to flow custom events back into signals. Note the `_` prefix on `$_name` and `$_reversed` — these are **local-only signals** that don't get sent to the server.

### Error Handling with ConsoleError

Use `sse.ConsoleError(err)` to send errors to the browser console without breaking the SSE stream:

```go
if err := sse.PatchElementTempl(c); err != nil {
    if err := sse.ConsoleError(err); err != nil {
        http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    }
    return
}
```

### Docker Deployment

Multi-stage build with `scratch` base and UPX compression for minimal image size:

```dockerfile
FROM docker.io/golang:1.25.0-alpine AS build
RUN apk add --no-cache upx
WORKDIR /src
COPY . ./
RUN go mod download
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -ldflags="-s" -o /bin/main ./cmd/web
RUN upx -9 -k /bin/main

FROM scratch
ENV PORT=9001
COPY --from=build /bin/main /
ENTRYPOINT ["/main"]
```

### Key Dependencies (Northstar Stack)

| Package | Purpose |
|---|---|
| `github.com/go-chi/chi/v5` | HTTP router with nested routes and middleware |
| `github.com/starfederation/datastar-go` | Datastar Go SDK (SSE, signals, templ helpers) |
| `github.com/a-h/templ` | Type-safe HTML templating with Go |
| `github.com/gorilla/sessions` | Cookie-based session management |
| `github.com/delaneyj/toolbelt/embeddednats` | Embedded NATS server |
| `github.com/nats-io/nats.go/jetstream` | NATS JetStream client (KV store, watchers) |
| `github.com/Jeffail/gabs/v2` | Dynamic JSON construction |
| `github.com/benbjohnson/hashfs` | Content-hash based static file serving |
| `github.com/evanw/esbuild` | TypeScript/JS bundling for web components |
| `github.com/joho/godotenv` | .env file loading |
| `golang.org/x/sync/errgroup` | Concurrent goroutine management |

### Northstar Taskfile (Development Workflow)

The recommended development workflow uses [Taskfile](https://taskfile.dev/) to orchestrate concurrent processes:

```yaml
version: "3"

env:
  NATS_PORT: 4222
  STATIC_DIR: "web/resources/static"

tasks:
  live:templ:
    cmds: [go tool templ generate -watch]

  live:web:styles:
    cmds: [go tool gotailwind -i web/resources/styles/styles.css -o $STATIC_DIR/index.css -w]

  live:web:bundle:
    cmds: [go run cmd/web/build/main.go -watch]

  live:server:
    cmds:
      - |
        go tool air \
          -build.cmd "go build -tags=dev -o tmp/bin/main ./cmd/web" \
          -build.bin "tmp/bin/main" \
          -build.exclude_dir "data,node_modules" \
          -build.include_ext "go,templ" \
          -misc.clean_on_exit "true"

  live:
    deps: [live:templ, live:web:styles, live:web:bundle, live:server]
```

Run `go tool task live` to start all watchers concurrently: templ generation, TailwindCSS, esbuild, and Air (Go live reload). Changes to `.templ` files trigger templ → Go → Air rebuild → hot reload in browser automatically.

---

## Quick Reference Card

### Frontend Cheatsheet

```
SIGNALS:        data-signals:name="'value'"     data-signals="{a: 1, b: 2}"
BIND:           data-bind:name                  (two-way binding to input)
COMPUTED:       data-computed:full="$a + $b"    (read-only derived signal)
TEXT:           data-text="$name"               (bind text content)
SHOW:           data-show="$visible"            (toggle visibility)
CLASS:          data-class:active="$isActive"   (toggle CSS class)
ATTR:           data-attr:disabled="$loading"   (set any HTML attribute)
STYLE:          data-style:color="$color"       (set inline style)
ON EVENT:       data-on:click="@post('/url')"   (event listener)
INIT:           data-init="@get('/data')"       (run on load)
EFFECT:         data-effect="$c = $a + $b"      (reactive side effect)
INDICATOR:      data-indicator:loading           (true during fetch)
REF:            data-ref:el                     (DOM element reference)
INTERVAL:       data-on-interval="@get('/poll')" (periodic execution)
INTERSECT:      data-on-intersect="@get('/more')"(viewport intersection)
DEBUG:          data-json-signals                (show all signals as JSON)
```

### Go Backend Cheatsheet

```go
// Read signals
store := &MyStore{}
datastar.ReadSignals(r, store)

// Create SSE writer
sse := datastar.NewSSE(w, r)

// Patch DOM elements
sse.PatchElements(`<div id="x">content</div>`)
sse.PatchElementf(`<div id="item-%d">%s</div>`, id, text)
sse.PatchElements(html, datastar.WithSelector("#target"), datastar.WithModeAppend())

// Remove elements
sse.RemoveElement("#selector")
sse.RemoveElementByID("element-id")

// Patch signals
sse.MarshalAndPatchSignals(map[string]any{"key": "value"})
sse.PatchSignals([]byte(`{"key": "value"}`))

// Execute JS
sse.ExecuteScript(`console.log("hello")`)

// Redirect
sse.Redirect("/new-url")

// Check connection
sse.IsClosed()

// Template helpers
datastar.GetSSE("/url")   // → @get('/url')
datastar.PostSSE("/url")  // → @post('/url')
```
