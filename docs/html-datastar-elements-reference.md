# HTML Elements & Datastar Interaction Reference for Coding Agents

> **HTML Standard**: WHATWG HTML Living Standard (continuously updated, last updated Feb 2026). There is no "HTML6" — the spec formerly known as "HTML5" is now a versionless living standard. All modern browsers implement the latest living standard.
> **Datastar Version**: v1.0.x
> **Purpose**: Every HTML element a coding agent encounters, what it does, how users interact with it, what events it fires, and how to wire it up with Datastar's `data-*` attributes.

---

## Table of Contents

1. [How Datastar Binds to HTML Elements](#how-datastar-binds-to-html-elements)
2. [Form Input Elements](#form-input-elements)
3. [Button Elements](#button-elements)
4. [Select & Option Elements](#select--option-elements)
5. [Textarea](#textarea)
6. [Form Element](#form-element)
7. [Container Elements (div, section, article, main, aside, header, footer, nav)](#container-elements)
8. [Text Content Elements (span, p, h1–h6, label, strong, em, etc.)](#text-content-elements)
9. [List Elements (ul, ol, li)](#list-elements)
10. [Table Elements (table, thead, tbody, tr, th, td)](#table-elements)
11. [Link & Navigation (a, nav)](#link--navigation)
12. [Image & Media (img, video, audio, picture, source, canvas, svg)](#image--media)
13. [Dialog & Disclosure (dialog, details, summary)](#dialog--disclosure)
14. [Progress & Meter](#progress--meter)
15. [Embedded Content (iframe, object, embed)](#embedded-content)
16. [Template & Slot (template, slot)](#template--slot)
17. [Web Components (Custom Elements)](#web-components)
18. [Head & Metadata (head, meta, title, link, script, style)](#head--metadata)
19. [Sectioning & Semantic (figure, figcaption, blockquote, pre, code, hr, br)](#sectioning--semantic)
20. [Output Element](#output-element)
21. [Fieldset & Legend](#fieldset--legend)
22. [Datalist](#datalist)
23. [Global Events Reference (All Elements)](#global-events-reference)
24. [Datastar Attribute Quick Reference](#datastar-attribute-quick-reference)

---

## How Datastar Binds to HTML Elements

Every visible HTML element can use Datastar's `data-*` attributes. The core bindings are:

| Datastar Attribute | What It Does | Applies To |
|---|---|---|
| `data-signals` | Declares reactive state (signals) on this element | Any element (typically a container) |
| `data-bind:ATTR` | Two-way binds a signal to a form element's value | `<input>`, `<select>`, `<textarea>` |
| `data-on:EVENT` | Listens for a DOM event, runs an expression or action | Any element |
| `data-text` | Sets `textContent` reactively | Any element |
| `data-attr:ATTR` | Sets an HTML attribute reactively | Any element |
| `data-class:NAME` | Toggles a CSS class reactively | Any element |
| `data-style:PROP` | Sets an inline style property reactively | Any element |
| `data-show` | Toggles `display` based on a signal | Any element |
| `data-init` | Runs an expression/action when element mounts | Any element |
| `data-effect` | Runs a reactive side-effect when dependencies change | Any element |
| `data-computed:NAME` | Declares a read-only derived signal | Any element |
| `data-indicator:NAME` | Auto-true while an SSE request is in-flight | Elements with `data-on:click` actions etc. |
| `data-ref:NAME` | Stores a reference to this DOM element as a signal | Any element |
| `data-on-intersect` | Fires when element enters the viewport | Any element |
| `data-on-interval` | Fires on a repeating timer | Any element |

**Critical rule**: `data-bind` only works on form elements (`<input>`, `<select>`, `<textarea>`). For everything else, use `data-text`, `data-attr`, `data-class`, `data-style`, or `data-show`.

---

## Form Input Elements

### `<input>` — The Most Versatile Element

The `<input>` element changes behavior entirely based on its `type` attribute. Each type has different events, values, and Datastar bindings.

#### `<input type="text">` (and `search`, `url`, `tel`, `email`, `password`)

Single-line text fields. All behave similarly for Datastar purposes.

**Native behavior**: User types text. Value is a string.

**Key events**:
| Event | When It Fires | Datastar Usage |
|---|---|---|
| `input` | Every keystroke / paste / autocomplete | `data-on:input` — real-time reactivity |
| `change` | When field loses focus AND value changed | `data-on:change` — commit-on-blur |
| `focus` | Field receives focus | `data-on:focus` |
| `blur` | Field loses focus | `data-on:blur` |
| `keydown` | A key is pressed (before character appears) | `data-on:keydown` — handle Enter, Escape, shortcuts |
| `keyup` | A key is released | `data-on:keyup` |
| `paste` | Content pasted from clipboard | `data-on:paste` |
| `select` | Text selection changes within the field | `data-on:select` |

**Datastar patterns**:

```html
<!-- Two-way bind to signal named "query" -->
<input type="text" data-bind:query />

<!-- Same but with explicit signal name different from attribute -->
<input type="text" data-bind:searchTerm />

<!-- React to every keystroke -->
<input type="text" data-bind:q data-on:input="@get('/search')" />

<!-- Submit on Enter key -->
<input type="text" data-bind:input
       data-on:keydown="if(evt.key==='Enter') @post('/submit')" />

<!-- Debounced input (250ms after last keystroke) -->
<input type="text" data-bind:q
       data-on:input__debounce.250ms="@get('/search')" />

<!-- Throttled input (at most once per 500ms) -->
<input type="text" data-bind:q
       data-on:input__throttle.500ms="@get('/search')" />

<!-- Disabled based on signal -->
<input type="text" data-attr:disabled="$loading" />

<!-- Placeholder from signal -->
<input type="text" data-attr:placeholder="$placeholderText" />

<!-- Readonly based on condition -->
<input type="text" data-attr:readonly="!$isEditing" />

<!-- Set value reactively (one-way, no user input capture) -->
<input type="text" data-attr:value="$computedValue" />
```

**Important attributes**: `placeholder`, `maxlength`, `minlength`, `pattern` (regex validation), `required`, `autofocus`, `autocomplete`, `readonly`, `disabled`, `size`, `list` (datalist reference), `spellcheck`.

#### `<input type="number">`

**Native behavior**: Accepts numeric input. Has increment/decrement spinner buttons. Value is still a string in the DOM but represents a number.

**Key events**: Same as text, plus the spinner fires `input` and `change`.

```html
<input type="number" data-bind:quantity min="0" max="100" step="1" />

<!-- Reactive min/max -->
<input type="number" data-bind:price
       data-attr:min="$minPrice" data-attr:max="$maxPrice" />
```

**Important attributes**: `min`, `max`, `step`, `placeholder`, `required`, `readonly`, `disabled`.

**Gotcha**: The bound signal value is a **string**. In Datastar expressions, use `+$quantity` or `Number($quantity)` to coerce to number if doing math.

#### `<input type="range">`

**Native behavior**: A slider. Value is a string representing a number.

**Key events**:
| Event | When It Fires |
|---|---|
| `input` | Continuously as the slider moves (every pixel) |
| `change` | When the user releases the slider |

```html
<input type="range" data-bind:volume min="0" max="100" step="1" />
<span data-text="$volume"></span>

<!-- Throttle for performance on rapid updates -->
<input type="range" data-bind:brightness
       data-on:input__throttle.50ms="@post('/brightness')" />
```

**Important attributes**: `min`, `max`, `step`, `list` (for tick marks via datalist).

#### `<input type="checkbox">`

**Native behavior**: Boolean toggle. The `checked` property is `true`/`false`. The `value` attribute is what gets submitted in forms (defaults to `"on"`).

**Key events**:
| Event | When It Fires |
|---|---|
| `change` | When checked/unchecked (click or keyboard) |
| `click` | On click (fires before `change`) |
| `input` | Same timing as `change` for checkboxes |

**Datastar binding**: `data-bind` on a checkbox binds to the **checked state** (boolean), not the value attribute.

```html
<!-- Binds $darkMode to true/false -->
<input type="checkbox" data-bind:darkMode />

<!-- React to toggle -->
<input type="checkbox" data-bind:enabled
       data-on:change="@post('/toggle')" />

<!-- Controlled checkbox (set checked from signal) -->
<input type="checkbox" data-attr:checked="$isSelected" />

<!-- Disabled based on condition -->
<input type="checkbox" data-bind:agree data-attr:disabled="$locked" />
```

**Group of checkboxes** (Datastar binds each to a separate boolean signal):

```html
<label><input type="checkbox" data-bind:optA /> Option A</label>
<label><input type="checkbox" data-bind:optB /> Option B</label>
<label><input type="checkbox" data-bind:optC /> Option C</label>
```

#### `<input type="radio">`

**Native behavior**: Mutually exclusive within a `name` group. The `value` attribute determines what the signal gets set to.

**Key events**: Same as checkbox — `change`, `click`, `input`.

**Datastar binding**: `data-bind` binds to the **signal** using the radio's `value` attribute. All radios in a group should bind to the **same signal name**.

```html
<!-- $size will be "sm", "md", or "lg" based on selection -->
<label><input type="radio" name="size" value="sm" data-bind:size /> Small</label>
<label><input type="radio" name="size" value="md" data-bind:size /> Medium</label>
<label><input type="radio" name="size" value="lg" data-bind:size /> Large</label>

<!-- Show content based on radio selection -->
<div data-show="$size === 'lg'">Large selected!</div>
```

#### `<input type="hidden">`

**Native behavior**: Not displayed. Carries a value for form submission.

**No user interaction**. Not useful with `data-bind` (no user input). Can use `data-attr:value` to set dynamically.

```html
<input type="hidden" data-attr:value="$sessionToken" name="token" />
```

#### `<input type="file">`

**Native behavior**: Opens a file picker dialog. Value is a `FileList`. **Cannot be set programmatically** (security restriction).

**Key events**:
| Event | When It Fires |
|---|---|
| `change` | When user selects file(s) |
| `cancel` | When user cancels the file picker (modern browsers) |

```html
<input type="file" data-on:change="@post('/upload')" accept=".pdf,.docx" />

<!-- Multiple files -->
<input type="file" data-on:change="@post('/upload')" multiple />

<!-- Accept only images -->
<input type="file" accept="image/*" data-on:change="@post('/upload-image')" />
```

**Important**: `data-bind` does NOT work with file inputs. Use `data-on:change` and handle file reading on the backend via the request body, or use JavaScript with `data-ref` and `data-on:change` to read files client-side.

**Important attributes**: `accept` (MIME types or extensions), `multiple`, `capture` (mobile camera: `user` front, `environment` rear).

#### `<input type="date">`, `<input type="time">`, `<input type="datetime-local">`

**Native behavior**: Browser-native date/time pickers. Value is a string (`"2026-02-13"`, `"14:30"`, `"2026-02-13T14:30"`).

**Key events**: `change` (when selection confirmed), `input` (as value changes).

```html
<input type="date" data-bind:startDate />
<input type="time" data-bind:startTime />
<input type="datetime-local" data-bind:eventTime />

<!-- With reactive min/max constraints -->
<input type="date" data-bind:endDate data-attr:min="$startDate" />
```

**Important attributes**: `min`, `max`, `step`, `required`, `readonly`.

**Related types**: `month` (value: `"2026-02"`), `week` (value: `"2026-W07"`).

#### `<input type="color">`

**Native behavior**: Color picker. Value is a hex string like `"#ff0000"`.

**Key events**: `input` (as color changes in picker), `change` (when picker closes).

```html
<input type="color" data-bind:themeColor />
<div data-style:background-color="$themeColor">Preview</div>
```

---

## Button Elements

### `<button>`

The most common interactive element. **Default type is `"submit"` inside a form, `"button"` outside a form** — always specify `type` explicitly.

**Types**:
| Type | Behavior |
|---|---|
| `button` | Does nothing natively. Fires `click`. Pure interaction trigger. |
| `submit` | Submits the nearest parent `<form>`. |
| `reset` | Resets all form controls in the nearest parent `<form>` to default. |

**Key events**:
| Event | When It Fires |
|---|---|
| `click` | Click or Enter/Space when focused |
| `mousedown` / `mouseup` | Press / release (before `click`) |
| `pointerdown` / `pointerup` | Pointer-agnostic (mouse, touch, pen) |
| `focus` / `blur` | Keyboard navigation |
| `keydown` | While focused, any key pressed |

**Datastar patterns**:

```html
<!-- Standard Datastar action button -->
<button type="button"
        data-on:click="@post('/api/items')">
    Add Item
</button>

<!-- With loading indicator -->
<button type="button"
        data-on:click="@post('/api/save')"
        data-indicator:saving
        data-attr:disabled="$saving">
    <span data-show="!$saving">Save</span>
    <span data-show="$saving">Saving...</span>
</button>

<!-- Confirm before action -->
<button type="button"
        data-on:click="if(confirm('Delete?')) @delete('/api/items/1')">
    Delete
</button>

<!-- Toggle signal -->
<button type="button" data-on:click="$showMenu = !$showMenu">
    Toggle Menu
</button>

<!-- Conditional styling -->
<button type="button"
        data-class:btn-active="$isActive"
        data-on:click="$isActive = !$isActive">
    Toggle
</button>

<!-- Disabled state from signal -->
<button data-attr:disabled="$formInvalid || $submitting">Submit</button>

<!-- With formatted path parameter -->
<button data-on:click={ datastar.DeleteSSE("/api/items/%d", itemID) }>
    Delete
</button>
```

**Accessibility**: Buttons are focusable and keyboard-activatable by default. Screen readers announce them as buttons. Use `<button>` over styled `<div>` or `<a>` for actions.

### `<input type="submit">` / `<input type="button">` / `<input type="reset">`

Legacy button variants. Prefer `<button>` — it can contain HTML children (icons, spans, etc.) while `<input>` buttons can only have plain text via the `value` attribute.

```html
<!-- Avoid this: -->
<input type="submit" value="Save" />

<!-- Prefer this: -->
<button type="submit">Save</button>
```

---

## Select & Option Elements

### `<select>`

**Native behavior**: Dropdown menu. Value is the `value` attribute of the selected `<option>`. With `multiple`, value is a list of selected values.

**Key events**:
| Event | When It Fires |
|---|---|
| `change` | When selection changes |
| `input` | Same timing as `change` for `<select>` |
| `focus` / `blur` | Focus gained / lost |

**Datastar patterns**:

```html
<!-- Single select — $color becomes "red", "green", or "blue" -->
<select data-bind:color>
    <option value="">-- Choose --</option>
    <option value="red">Red</option>
    <option value="green">Green</option>
    <option value="blue">Blue</option>
</select>

<!-- React on change -->
<select data-bind:category data-on:change="@get('/products')">
    <option value="all">All</option>
    <option value="electronics">Electronics</option>
    <option value="clothing">Clothing</option>
</select>

<!-- Conditionally disabled options (use server-rendered HTML, Datastar doesn't dynamically add options) -->
<select data-bind:country>
    <option value="us">United States</option>
    <option value="uk">United Kingdom</option>
</select>

<!-- Set disabled from signal -->
<select data-bind:priority data-attr:disabled="$locked">
    <option value="low">Low</option>
    <option value="high">High</option>
</select>

<!-- Reactively show based on selection -->
<div data-show="$category === 'other'">
    <input type="text" data-bind:otherCategory placeholder="Specify..." />
</div>
```

### `<select multiple>`

With `multiple`, the user can select multiple options (Ctrl/Cmd+click or Shift+click). **`data-bind` syncs to an array signal**.

```html
<select multiple data-bind:selectedTags>
    <option value="js">JavaScript</option>
    <option value="go">Go</option>
    <option value="py">Python</option>
</select>
```

**Important attributes**: `size` (number of visible rows), `required`, `disabled`, `autofocus`.

### `<option>`

Children of `<select>`, `<optgroup>`, or `<datalist>`. Not independently interactive.

**Key attributes**: `value` (what gets sent to the signal/server — if omitted, uses text content), `selected` (pre-selected), `disabled` (grayed out, not selectable), `label` (alternative display text).

### `<optgroup>`

Groups `<option>` elements under a label. Not selectable itself.

```html
<select data-bind:font>
    <optgroup label="Serif">
        <option value="georgia">Georgia</option>
        <option value="times">Times New Roman</option>
    </optgroup>
    <optgroup label="Sans-serif">
        <option value="arial">Arial</option>
        <option value="helvetica">Helvetica</option>
    </optgroup>
</select>
```

---

## Textarea

### `<textarea>`

**Native behavior**: Multi-line text input. Value is the text content (not the `value` attribute — the text between tags is the initial value).

**Key events**: Same as `<input type="text">` — `input`, `change`, `focus`, `blur`, `keydown`, `keyup`, `select`, `paste`.

```html
<!-- Two-way bind -->
<textarea data-bind:bio rows="5" cols="40"></textarea>

<!-- Character count -->
<textarea data-bind:message maxlength="500"></textarea>
<span data-text="$message.length + '/500'"></span>

<!-- Auto-submit on change (debounced) -->
<textarea data-bind:notes
          data-on:input__debounce.1000ms="@put('/api/notes')">
</textarea>

<!-- Disabled/readonly -->
<textarea data-bind:content data-attr:readonly="!$canEdit"></textarea>
```

**Important attributes**: `rows`, `cols`, `maxlength`, `minlength`, `placeholder`, `required`, `readonly`, `disabled`, `wrap` (`hard`/`soft`), `spellcheck`, `autocomplete`.

---

## Form Element

### `<form>`

**Native behavior**: Groups form controls. On submit, collects all named inputs and sends them as a request. **With Datastar, you almost never use native form submission** — instead you use `data-on:click` with `@post()`/`@put()` on buttons, or `data-on:submit` on the form.

**Key events**:
| Event | When It Fires |
|---|---|
| `submit` | Form submission triggered (button click or Enter in text input) |
| `reset` | Form reset triggered |
| `formdata` | When FormData is constructed (just before submission) |

**Datastar patterns**:

```html
<!-- RECOMMENDED: No <form> at all — Datastar sends signals, not form data -->
<div data-signals="{name:'', email:''}">
    <input type="text" data-bind:name />
    <input type="email" data-bind:email />
    <button type="button" data-on:click="@post('/api/users')">Submit</button>
</div>

<!-- If you need <form> (for native validation, accessibility, or Enter-to-submit): -->
<form data-on:submit__prevent="@post('/api/users')">
    <input type="text" data-bind:name required />
    <input type="email" data-bind:email required />
    <button type="submit">Submit</button>
</form>
```

**Critical**: Use `data-on:submit__prevent` (with the `__prevent` modifier) to prevent the native form submission (which would cause a full page navigation). The `__prevent` modifier calls `evt.preventDefault()`.

**When to use `<form>` vs not**:
- **Use `<form>`**: When you want browser-native validation (`required`, `pattern`, `min`/`max`), Enter-to-submit in text fields, or improved accessibility.
- **Skip `<form>`**: When Datastar signals handle all state and you use explicit `@post()` calls on buttons. This is the more common Datastar pattern.

**Important attributes**: `action` (unused with Datastar), `method` (unused), `enctype` (unused — Datastar sends JSON), `novalidate` (skip browser validation), `autocomplete`.

---

## Container Elements

### `<div>`, `<section>`, `<article>`, `<main>`, `<aside>`, `<header>`, `<footer>`, `<nav>`

These are all **generic or semantic containers**. They have no inherent interactive behavior — they are purely structural. From Datastar's perspective, they all work identically.

**Key events** (all inherited from global event handlers):
| Event | When It Fires |
|---|---|
| `click` | Click on the element or any child (bubbles up) |
| `mouseenter` / `mouseleave` | Mouse enters / exits (no bubbling) |
| `mouseover` / `mouseout` | Mouse enters / exits (bubbles) |
| `pointerenter` / `pointerleave` | Pointer agnostic version |
| `scroll` | If the element has overflow scrolling |
| `keydown` / `keyup` | Only if focused (requires `tabindex`) |
| `dragstart` / `dragover` / `drop` | If `draggable="true"` |
| `touchstart` / `touchmove` / `touchend` | Touch events (mobile) |

**Datastar patterns**:

```html
<!-- Container for signals scope -->
<div data-signals="{count: 0, name: 'World'}">
    <span data-text="$name"></span>
</div>

<!-- Initialize data from server on load -->
<div id="container" data-init="@get('/api/data')"></div>

<!-- Toggle visibility -->
<div data-show="$showPanel" class="panel">
    Panel content...
</div>

<!-- Conditional CSS class -->
<div data-class:hidden="!$visible"
     data-class:active="$isActive"
     data-class:dark="$darkMode">
    Content
</div>

<!-- Reactive inline styles -->
<div data-style:background-color="$bgColor"
     data-style:opacity="$isVisible ? '1' : '0.5'">
    Content
</div>

<!-- Click handler on container (delegation) -->
<div data-on:click="$selectedPanel = 'info'" class="card">
    Click this card
</div>

<!-- Scroll detection -->
<div style="overflow-y: auto; height: 400px;"
     data-on:scroll__throttle.100ms="$scrollPos = evt.target.scrollTop">
    Long content...
</div>

<!-- Intersection observer (lazy load, infinite scroll) -->
<div data-on-intersect="@get('/api/more-items')">
    Loading more...
</div>

<!-- Intersection observer — fire only once -->
<div data-on-intersect__once="@get('/api/analytics/view')">
    Tracked section
</div>

<!-- Periodic polling -->
<div data-on-interval.5000ms="@get('/api/status')">
    Status: <span data-text="$status"></span>
</div>

<!-- Side effects (runs whenever dependencies change) -->
<div data-effect="document.title = 'Count: ' + $count"></div>

<!-- Set an HTML attribute reactively -->
<div data-attr:id="'item-' + $itemId"
     data-attr:data-testid="'test-' + $itemId">
    Content
</div>

<!-- Remove element from backend -->
<!-- Backend calls: sse.RemoveElement("#notification") -->
<div id="notification">This can be removed by the server</div>
```

**Semantic meaning** (for accessibility and SEO — behavior is identical to `<div>` for Datastar):
| Element | Semantic Meaning |
|---|---|
| `<div>` | Generic container, no semantic meaning |
| `<section>` | Thematic grouping of content (has a heading) |
| `<article>` | Self-contained composition (blog post, comment, widget) |
| `<main>` | Primary content of the document (one per page) |
| `<aside>` | Tangentially related content (sidebar, callout) |
| `<header>` | Introductory content or navigational aids |
| `<footer>` | Footer for its nearest sectioning content |
| `<nav>` | Navigation links |

---

## Text Content Elements

### `<span>`, `<p>`, `<h1>`–`<h6>`, `<label>`, `<strong>`, `<em>`, `<small>`, `<mark>`, `<abbr>`, `<cite>`, `<q>`, `<time>`, `<data>`, `<sub>`, `<sup>`, `<del>`, `<ins>`, `<s>`, `<u>`, `<b>`, `<i>`, `<kbd>`, `<samp>`, `<var>`

These are all **text-level** or **phrasing content** elements. They are not inherently interactive (except `<label>`). Datastar treats them all the same — as targets for `data-text`, `data-show`, `data-class`, etc.

**Datastar patterns**:

```html
<!-- Set text content reactively -->
<span data-text="$username"></span>
<p data-text="$description"></p>
<h1 data-text="$pageTitle"></h1>

<!-- Formatted text -->
<span data-text="'$' + Number($price).toFixed(2)"></span>
<span data-text="$count + ' items'"></span>

<!-- Conditional visibility -->
<span data-show="$error" class="text-red">Error: <span data-text="$error"></span></span>

<!-- Conditional classes -->
<span data-class:line-through="$completed"
      data-class:text-muted="$completed">
    Task text
</span>

<!-- Clickable span (make it behave like a button) -->
<span data-on:click="@get('/api/item/1/edit')"
      style="cursor: pointer;"
      role="button" tabindex="0">
    Edit
</span>
```

### `<label>`

**Special behavior**: Clicking a `<label>` activates its associated form control (the element with matching `id` in the `for` attribute, or the first form control descendant).

```html
<!-- Explicit association via "for" -->
<label for="nameInput">Name:</label>
<input id="nameInput" type="text" data-bind:name />

<!-- Implicit association (input inside label) -->
<label>
    Name:
    <input type="text" data-bind:name />
</label>

<!-- Label with reactive text -->
<label data-text="'Quantity (' + $qty + '):'"></label>
```

**Important**: Always associate labels with form controls for accessibility. Screen readers read the label when the input is focused.

---

## List Elements

### `<ul>`, `<ol>`, `<li>`

**No inherent interactivity**. Lists are structural. Datastar interacts with them via event delegation and server-rendered content.

```html
<!-- Static list with reactive items (server replaces the whole list) -->
<ul id="items">
    <!-- Server sends <li> elements via PatchElements targeting #items -->
</ul>

<!-- List item as clickable (Northstar todo pattern) -->
<ul>
    <li class="flex items-center gap-2" id="todo-1">
        <label data-on:click="@post('/api/todos/1/toggle')">
            <!-- checkbox icon -->
        </label>
        <span data-on:click="@get('/api/todos/1/edit')">Todo text</span>
        <button data-on:click="@delete('/api/todos/1')">×</button>
    </li>
</ul>

<!-- Reactive show/hide on list items -->
<li data-show="$filter === 'all' || ($filter === 'active' && !$completed)">
    Item text
</li>
```

### `<dl>`, `<dt>`, `<dd>` (Description Lists)

Same as `<ul>`/`<li>` — no inherent interactivity, used structurally.

---

## Table Elements

### `<table>`, `<thead>`, `<tbody>`, `<tfoot>`, `<tr>`, `<th>`, `<td>`, `<caption>`, `<colgroup>`, `<col>`

**No inherent interactivity**. Tables display tabular data. Datastar interacts with tables by patching rows/cells from the server or using event handlers.

```html
<!-- Table with reactive content (server patches #table-body) -->
<table>
    <thead>
        <tr>
            <th data-on:click="@get('/api/users?sort=name')">Name</th>
            <th data-on:click="@get('/api/users?sort=email')">Email</th>
        </tr>
    </thead>
    <tbody id="table-body">
        <!-- Server patches rows here -->
    </tbody>
</table>

<!-- Sortable headers with sort indicator -->
<th data-on:click="$sortBy = 'name'; $sortDir = $sortDir === 'asc' ? 'desc' : 'asc'; @get('/api/users')"
    style="cursor: pointer;">
    Name <span data-show="$sortBy === 'name'" data-text="$sortDir === 'asc' ? '▲' : '▼'"></span>
</th>

<!-- Row click handler -->
<tr data-on:click="@get('/api/users/1')" style="cursor: pointer;"
    data-class:bg-primary="$selectedId === 1">
    <td>John</td>
    <td>john@example.com</td>
</tr>

<!-- Editable cell -->
<td>
    <span data-show="$editingCell !== 'name-1'" data-text="$name1"
          data-on:dblclick="$editingCell = 'name-1'"></span>
    <input data-show="$editingCell === 'name-1'" data-bind:name1
           data-on:blur="$editingCell = ''; @put('/api/users/1')"
           data-on:keydown="if(evt.key==='Enter'){$editingCell=''; @put('/api/users/1')}" />
</td>
```

---

## Link & Navigation

### `<a>` (Anchor)

**Native behavior**: Navigates to `href` URL on click. This is the **primary navigation element** on the web.

**Key events**:
| Event | When It Fires |
|---|---|
| `click` | Click or Enter when focused |
| `mouseenter` / `mouseleave` | Hover |
| `focus` / `blur` | Tab navigation |

**Datastar patterns**:

```html
<!-- Standard navigation (let browser handle it) -->
<a href="/about">About</a>

<!-- Prevent navigation, use Datastar action instead -->
<a href="#" data-on:click__prevent="@get('/api/panel/info')">
    Load Info
</a>

<!-- Dynamic href -->
<a data-attr:href="'/users/' + $userId">View Profile</a>

<!-- Open in new tab -->
<a href="/report.pdf" target="_blank" rel="noopener">Download Report</a>

<!-- Styled as button -->
<a href="#" role="button"
   data-on:click__prevent="@post('/api/action')"
   class="btn btn-primary">
    Perform Action
</a>

<!-- Conditional link styling -->
<a href="/dashboard"
   data-class:active="$currentPage === 'dashboard'"
   data-class:disabled="!$isLoggedIn">
    Dashboard
</a>
```

**Important attributes**: `href` (URL — required for it to be a link), `target` (`_blank`, `_self`, `_parent`, `_top`), `rel` (`noopener`, `noreferrer`, `nofollow`), `download` (triggers download instead of navigation), `hreflang`, `type`.

**When to use `<a>` vs `<button>`**:
- `<a>` — **Navigation** (going to a new page/URL). Has an `href`.
- `<button>` — **Action** (doing something on the current page). No `href`.
- Datastar blurs this line since server actions respond with SSE, but the semantic distinction matters for accessibility.

---

## Image & Media

### `<img>`

**Native behavior**: Displays an image. Not interactive by default.

**Key events**:
| Event | When It Fires |
|---|---|
| `load` | Image finished loading successfully |
| `error` | Image failed to load |
| `click` | If clickable (e.g., inside an `<a>`, or with a handler) |

```html
<!-- Reactive image source -->
<img data-attr:src="$avatarUrl" data-attr:alt="$userName + ' avatar'" />

<!-- Lazy loading (native) -->
<img src="/placeholder.jpg" loading="lazy" data-attr:src="$imageUrl" />

<!-- Error fallback -->
<img data-attr:src="$photoUrl"
     data-on:error="evt.target.src = '/fallback.png'" />

<!-- Clickable image -->
<img src="/thumb.jpg" data-on:click="$lightboxUrl = '/full.jpg'; $showLightbox = true"
     style="cursor: pointer;" />

<!-- Toggle visibility -->
<img data-show="$hasAvatar" data-attr:src="$avatarUrl" />
```

**Important attributes**: `src`, `alt` (required for accessibility), `width`, `height` (prevents layout shift), `loading` (`lazy`/`eager`), `decoding` (`async`/`sync`/`auto`), `srcset` (responsive images), `sizes`, `crossorigin`, `referrerpolicy`, `fetchpriority`.

### `<video>`

**Key events**: `play`, `pause`, `ended`, `timeupdate`, `volumechange`, `loadeddata`, `error`, `seeking`, `seeked`, `canplay`, `waiting`, `progress`.

```html
<video id="player" data-attr:src="$videoUrl" controls
       data-on:play="$isPlaying = true"
       data-on:pause="$isPlaying = false"
       data-on:ended="@post('/api/watch-complete')"
       data-on:timeupdate__throttle.1000ms="$currentTime = evt.target.currentTime">
</video>
<span data-text="Math.floor($currentTime) + 's'"></span>
```

**Important attributes**: `src`, `controls`, `autoplay`, `muted`, `loop`, `poster` (thumbnail), `preload` (`none`/`metadata`/`auto`), `width`, `height`, `playsinline`.

### `<audio>`

Same events as `<video>`. Same Datastar patterns.

```html
<audio data-attr:src="$trackUrl" controls
       data-on:ended="@post('/api/next-track')">
</audio>
```

### `<picture>` and `<source>`

Container for responsive image sources. `<picture>` wraps multiple `<source>` elements and one `<img>` fallback.

```html
<picture>
    <source srcset="/hero-large.webp" media="(min-width: 1200px)" type="image/webp" />
    <source srcset="/hero-medium.webp" media="(min-width: 768px)" type="image/webp" />
    <img src="/hero-small.jpg" alt="Hero image" />
</picture>
```

No Datastar-specific patterns — this is purely a rendering concern.

### `<canvas>`

A drawing surface for JavaScript. **No inherent content** — everything is drawn programmatically.

**Key events**: All pointer/mouse/touch/keyboard events work on canvas.

```html
<canvas id="chart" width="800" height="400"
        data-ref:chartCanvas
        data-on:click="handleCanvasClick(evt, $chartCanvas)">
</canvas>
```

Datastar can set attributes and listen to events, but drawing requires JavaScript (via `data-effect` or `data-init` with `data-ref`).

### `<svg>`

Inline SVG elements support all Datastar attributes. SVG child elements (`<rect>`, `<circle>`, `<path>`, etc.) also support events.

```html
<svg width="100" height="100">
    <circle cx="50" cy="50"
            data-attr:r="$radius"
            data-attr:fill="$circleColor"
            data-on:click="$selectedShape = 'circle'" />
</svg>
```

---

## Dialog & Disclosure

### `<dialog>`

**Native behavior**: A dialog box / modal. Has built-in `open` attribute, `showModal()` / `show()` / `close()` methods, backdrop, focus trapping, and Escape-to-close.

**Key events**:
| Event | When It Fires |
|---|---|
| `close` | Dialog is closed (via `.close()`, Escape, or form submission) |
| `cancel` | User pressed Escape (fires before `close`) |
| `click` | Click on dialog or backdrop |

```html
<!-- Datastar-controlled dialog using show/close methods -->
<dialog data-ref:myDialog
        data-effect="$showDialog ? $myDialog.showModal() : $myDialog.close()"
        data-on:close="$showDialog = false">
    <h2>Confirm Action</h2>
    <p>Are you sure?</p>
    <button data-on:click="$showDialog = false">Cancel</button>
    <button data-on:click="@delete('/api/items/1'); $showDialog = false">Delete</button>
</dialog>

<button data-on:click="$showDialog = true">Open Dialog</button>

<!-- Close on backdrop click -->
<dialog data-ref:dlg
        data-on:click="if(evt.target === $dlg) $showDialog = false">
    <div>Dialog content (clicks here don't close)</div>
</dialog>
```

**Important**: `<dialog>` with `showModal()` provides native focus trapping, Escape-to-close, and a `::backdrop` pseudo-element. This is far better than custom modal implementations.

### `<details>` / `<summary>`

**Native behavior**: A disclosure widget. `<summary>` is the always-visible toggle, the rest of `<details>` content is shown/hidden.

**Key events**:
| Event | When It Fires |
|---|---|
| `toggle` | Fires on `<details>` when opened or closed |

```html
<details data-on:toggle="if(evt.target.open) @get('/api/section/details')">
    <summary>Show more info</summary>
    <div id="details-content">
        <!-- Content loaded from server when opened -->
    </div>
</details>

<!-- Controlled open state -->
<details data-attr:open="$faqOpen">
    <summary data-on:click__prevent="$faqOpen = !$faqOpen">FAQ</summary>
    <p>Answer content...</p>
</details>
```

---

## Progress & Meter

### `<progress>`

**Native behavior**: Progress bar. No user interaction.

```html
<!-- Determinate progress (0 to max) -->
<progress data-attr:value="$uploadProgress" max="100"></progress>
<span data-text="$uploadProgress + '%'"></span>

<!-- Indeterminate (no value = animated "loading") -->
<progress data-show="$isLoading"></progress>
```

**Important attributes**: `value` (current), `max` (maximum, default 1).

### `<meter>`

**Native behavior**: Scalar measurement within a known range. Browser styles it green/yellow/red based on `low`/`high`/`optimum`. No user interaction.

```html
<meter data-attr:value="$diskUsage" min="0" max="100"
       low="50" high="80" optimum="30"></meter>
<span data-text="$diskUsage + '% disk used'"></span>
```

**Important attributes**: `value`, `min`, `max`, `low`, `high`, `optimum`.

---

## Embedded Content

### `<iframe>`

**Native behavior**: Embeds another HTML page. Heavily sandboxed for security.

**Key events**: `load`, `error`. **Cannot access events inside the iframe** from the parent (cross-origin security).

```html
<!-- Reactive iframe source -->
<iframe data-attr:src="$previewUrl" width="100%" height="500"></iframe>

<!-- Sandboxed iframe -->
<iframe src="/widget" sandbox="allow-scripts allow-same-origin"></iframe>
```

**Important attributes**: `src`, `srcdoc` (inline HTML), `sandbox`, `allow` (permissions policy), `loading` (`lazy`), `width`, `height`, `name`, `referrerpolicy`.

### `<object>` and `<embed>`

Legacy embedding elements. Rarely used in modern web development. Use `<iframe>`, `<img>`, `<video>`, or `<audio>` instead.

---

## Template & Slot

### `<template>`

**Native behavior**: Content inside `<template>` is **not rendered**. It is parsed but inert — scripts don't execute, images don't load, etc. Used as a stamp for JavaScript cloning, or by web components.

**Datastar does not process `<template>` content** — it's inert. However, the server can send HTML fragments that include any elements.

### `<slot>`

Used inside Shadow DOM (web components) to define insertion points. **Not relevant for Datastar's light DOM approach** unless you're using web components (see below).

---

## Web Components

### Custom Elements (`<my-component>`)

Any element with a hyphen in the name is a custom element. Datastar interacts with them by:

1. **Passing data in via attributes**: `data-attr:propname="$signal"`
2. **Listening for custom events**: `data-on:customevent="$signal = evt.detail.value"`
3. **Using local signals**: Prefix with `_` to keep signals local (not sent to server)

```html
<!-- Northstar pattern: Vanilla custom element -->
<reverse-component
    data-on:reverse="$_reversed = evt.detail.value"
    data-attr:name="$_name">
</reverse-component>

<!-- Lit/Sortable web component -->
<sortable-example
    data-signals="{title: 'Item Info', items: [{name: 'one'}, {name: 'two'}]}"
    data-attr:title="$title"
    data-attr:items="JSON.stringify($items)"
    data-on:change="console.log(event.detail)">
</sortable-example>
```

**Custom element lifecycle events** (from the JS class, not DOM events):
- `connectedCallback` — element added to DOM
- `disconnectedCallback` — element removed from DOM
- `attributeChangedCallback` — observed attribute changed (this is how `data-attr` flows data in)

---

## Head & Metadata

### `<head>`, `<meta>`, `<title>`, `<link>`, `<script>`, `<style>`, `<base>`

These elements live in `<head>` and are **not rendered in the page body**. Datastar typically does not interact with them, with these exceptions:

**`<title>`**: Can be changed reactively via `data-effect`:
```html
<div data-effect="document.title = 'Chat (' + $unreadCount + ')'"></div>
```

**`<script>`**: The Datastar JS itself is loaded via a script tag. It must be `type="module"` and `defer`:
```html
<script defer type="module" src="/static/datastar/datastar.js"></script>
```

The backend can execute JavaScript via `sse.ExecuteScript()`, which Datastar evaluates without needing a `<script>` tag.

**`<meta name="viewport">`**: Required for responsive design:
```html
<meta name="viewport" content="width=device-width, initial-scale=1" />
```

---

## Sectioning & Semantic

### `<figure>` / `<figcaption>`

Container for self-contained content (images, diagrams, code) with a caption. No interactivity.

```html
<figure>
    <img data-attr:src="$chartUrl" alt="Sales chart" />
    <figcaption data-text="'Sales data for ' + $selectedYear"></figcaption>
</figure>
```

### `<blockquote>`

Block-level quotation. No interactivity.

```html
<blockquote data-show="$showQuote">
    <p data-text="$quoteText"></p>
    <cite data-text="'— ' + $quoteAuthor"></cite>
</blockquote>
```

### `<pre>` / `<code>`

Preformatted text / inline code. No interactivity. `<pre>` preserves whitespace.

```html
<pre><code data-text="$codeSnippet"></code></pre>
```

### `<hr>`

Thematic break (horizontal rule). No interactivity. Can conditionally show:

```html
<hr data-show="$showDivider" />
```

### `<br>`

Line break. No interactivity. No Datastar usage.

---

## Output Element

### `<output>`

**Native behavior**: Represents the result of a calculation or user action. Associated with form controls via the `for` attribute. Semantically announces changes to screen readers.

```html
<output data-text="$total" for="price qty">Total: </output>
```

Good for accessibility — screen readers treat it as a live region.

---

## Fieldset & Legend

### `<fieldset>` / `<legend>`

Groups form controls with a visible label. **Disabling a `<fieldset>` disables all its descendant form controls**.

```html
<fieldset data-attr:disabled="$isSubmitting">
    <legend>Shipping Address</legend>
    <input type="text" data-bind:street placeholder="Street" />
    <input type="text" data-bind:city placeholder="City" />
    <input type="text" data-bind:zip placeholder="ZIP" />
</fieldset>
```

**This is extremely useful with Datastar** — a single `data-attr:disabled` on the fieldset disables the entire group during submission.

---

## Datalist

### `<datalist>`

Provides autocomplete suggestions for an `<input>`. The browser shows a dropdown of matching options as the user types.

```html
<input type="text" data-bind:city list="cities" />
<datalist id="cities">
    <option value="New York" />
    <option value="Los Angeles" />
    <option value="Chicago" />
    <option value="Houston" />
    <option value="Phoenix" />
</datalist>
```

**Not dynamically updatable via Datastar signals alone** — the `<datalist>` content must be server-rendered. The backend can patch the datalist element via SSE:

```go
// Backend
sse.PatchElements(`<datalist id="cities"><option value="Austin"/><option value="Dallas"/></datalist>`)
```

---

## Global Events Reference

Every HTML element supports these events. Use with `data-on:eventname`.

### Mouse Events

| Event | When It Fires | Notes |
|---|---|---|
| `click` | Left-click (or Enter/Space on focused interactive element) | Most common action trigger |
| `dblclick` | Double-click | |
| `contextmenu` | Right-click | Use `__prevent` to suppress native menu |
| `mousedown` | Mouse button pressed | `evt.button`: 0=left, 1=middle, 2=right |
| `mouseup` | Mouse button released | |
| `mousemove` | Mouse moves over element | **High frequency** — throttle! |
| `mouseenter` | Mouse enters element | Does NOT bubble |
| `mouseleave` | Mouse exits element | Does NOT bubble |
| `mouseover` | Mouse enters element or child | Bubbles |
| `mouseout` | Mouse exits element or child | Bubbles |

### Pointer Events (Recommended over Mouse Events)

Work with mouse, touch, and pen input.

| Event | When It Fires |
|---|---|
| `pointerdown` | Any pointer pressed |
| `pointerup` | Any pointer released |
| `pointermove` | Pointer moves |
| `pointerenter` | Pointer enters element |
| `pointerleave` | Pointer exits element |
| `pointerover` | Pointer enters element (bubbles) |
| `pointerout` | Pointer exits element (bubbles) |
| `pointercancel` | Pointer interaction cancelled |
| `gotpointercapture` | Element captured pointer |
| `lostpointercapture` | Element lost pointer capture |

### Touch Events (Mobile)

| Event | When It Fires |
|---|---|
| `touchstart` | Finger touches screen |
| `touchmove` | Finger moves on screen |
| `touchend` | Finger lifts off screen |
| `touchcancel` | Touch interrupted |

### Keyboard Events

| Event | When It Fires | Notes |
|---|---|---|
| `keydown` | Key pressed (repeats if held) | `evt.key` = the character, `evt.code` = physical key |
| `keyup` | Key released | Does not repeat |

Common `evt.key` values: `'Enter'`, `'Escape'`, `'Tab'`, `'Backspace'`, `'Delete'`, `'ArrowUp'`, `'ArrowDown'`, `'ArrowLeft'`, `'ArrowRight'`, `' '` (space), `'a'`–`'z'`, `'0'`–`'9'`.

Modifier checks: `evt.ctrlKey`, `evt.shiftKey`, `evt.altKey`, `evt.metaKey` (Cmd on Mac).

```html
<!-- Common keyboard patterns -->
<input data-on:keydown="if(evt.key==='Enter') @post('/submit')" />
<input data-on:keydown="if(evt.key==='Escape') $editing = false" />
<div data-on:keydown="if(evt.ctrlKey && evt.key==='s'){evt.preventDefault(); @post('/save')}" tabindex="0" />
```

### Focus Events

| Event | When It Fires | Notes |
|---|---|---|
| `focus` | Element receives focus | Does NOT bubble |
| `blur` | Element loses focus | Does NOT bubble |
| `focusin` | Element receives focus | Bubbles (use for delegation) |
| `focusout` | Element loses focus | Bubbles |

### Drag & Drop Events

| Event | When It Fires | Target |
|---|---|---|
| `dragstart` | Drag begins | Dragged element |
| `drag` | During drag (continuous) | Dragged element |
| `dragend` | Drag ends | Dragged element |
| `dragenter` | Dragged item enters element | Drop target |
| `dragover` | Dragged item over element (continuous) | Drop target — **must `__prevent` default to allow drop** |
| `dragleave` | Dragged item leaves element | Drop target |
| `drop` | Item dropped | Drop target — **must `__prevent` default** |

```html
<div draggable="true"
     data-on:dragstart="evt.dataTransfer.setData('text/plain', '1')">
    Drag me
</div>
<div data-on:dragover__prevent=""
     data-on:drop__prevent="@post('/api/reorder')">
    Drop here
</div>
```

### Clipboard Events

| Event | When It Fires |
|---|---|
| `copy` | User copies (Ctrl+C) |
| `cut` | User cuts (Ctrl+X) |
| `paste` | User pastes (Ctrl+V) |

### Scroll Events

| Event | When It Fires | Notes |
|---|---|---|
| `scroll` | Element scrolls | **High frequency** — always throttle |
| `scrollend` | Scrolling ends | Newer — check browser support |

### Animation & Transition Events

| Event | When It Fires |
|---|---|
| `animationstart` | CSS animation starts |
| `animationend` | CSS animation ends |
| `animationiteration` | CSS animation repeats |
| `transitionstart` | CSS transition starts |
| `transitionend` | CSS transition ends |
| `transitionrun` | CSS transition begins running |
| `transitioncancel` | CSS transition cancelled |

### Media Events (on `<video>`, `<audio>`)

| Event | When It Fires |
|---|---|
| `play` | Playback starts |
| `pause` | Playback paused |
| `ended` | Playback finished |
| `timeupdate` | Current time changed (continuous during playback) |
| `volumechange` | Volume or muted state changed |
| `loadeddata` | First frame of media loaded |
| `canplay` | Enough data to start playing |
| `waiting` | Playback stopped due to buffering |
| `seeking` | Seek operation started |
| `seeked` | Seek operation finished |
| `durationchange` | Duration changed |
| `ratechange` | Playback rate changed |
| `error` | Loading failed |

### Window/Document Events (use on body/root elements)

| Event | When It Fires | Notes |
|---|---|---|
| `load` | Page fully loaded | Use on `<body>` or `<img>` |
| `DOMContentLoaded` | HTML parsed, before images/CSS | |
| `resize` | Window resized | Throttle! |
| `beforeunload` | User about to leave page | Can prompt confirmation |
| `hashchange` | URL hash changed | |
| `popstate` | Browser history navigation | |
| `online` / `offline` | Network connectivity changed | |
| `visibilitychange` | Tab becomes visible/hidden | `document.visibilityState` |

---

## Datastar Attribute Quick Reference

### Event Modifiers

Appended to `data-on:event` with double underscores:

| Modifier | Effect | Example |
|---|---|---|
| `__prevent` | Calls `evt.preventDefault()` | `data-on:submit__prevent` |
| `__stop` | Calls `evt.stopPropagation()` | `data-on:click__stop` |
| `__capture` | Listens in capture phase | `data-on:click__capture` |
| `__once` | Handler fires only once then removes | `data-on:click__once` |
| `__passive` | Sets `{ passive: true }` (scroll perf) | `data-on:scroll__passive` |
| `__self` | Only fires if `evt.target === this element` | `data-on:click__self` |
| `__outside` | Fires when click is OUTSIDE this element | `data-on:click__outside` |
| `__debounce.Xms` | Debounce by X ms | `data-on:input__debounce.300ms` |
| `__throttle.Xms` | Throttle to at most once per X ms | `data-on:mousemove__throttle.100ms` |

Multiple modifiers can be combined: `data-on:click__prevent__stop`.

### Merge Modes (for PatchElements from backend)

When the Go backend sends HTML fragments, it can specify how they merge into the DOM:

| Mode | Go Option | Behavior |
|---|---|---|
| Morph (default) | `datastar.WithModeMorph()` | Intelligent DOM diffing via Idiomorph |
| Inner | `datastar.WithModeInner()` | Replace inner HTML |
| Outer | `datastar.WithModeOuter()` | Replace entire element |
| Prepend | `datastar.WithModePrepend()` | Add before first child |
| Append | `datastar.WithModeAppend()` | Add after last child |
| Before | `datastar.WithModeBefore()` | Add before the element |
| After | `datastar.WithModeAfter()` | Add after the element |
| Upsert attributes | `datastar.WithModeUpsertAttributes()` | Only update attributes |
| Delete | `sse.RemoveElement("#id")` | Remove element from DOM |

### Elements That Support `data-bind`

**Only these elements work with `data-bind`**:

| Element | What `data-bind` Binds To | Signal Type |
|---|---|---|
| `<input type="text">` (and text-like) | `.value` | `string` |
| `<input type="number">` | `.value` | `string` (coerce with `+$signal`) |
| `<input type="range">` | `.value` | `string` |
| `<input type="checkbox">` | `.checked` | `boolean` |
| `<input type="radio">` | `.value` of selected radio | `string` |
| `<input type="date/time/datetime-local">` | `.value` | `string` |
| `<input type="color">` | `.value` | `string` (`"#rrggbb"`) |
| `<select>` | `.value` of selected option | `string` |
| `<select multiple>` | Array of selected `.value`s | `array` |
| `<textarea>` | `.value` | `string` |

**Everything else** (`<div>`, `<span>`, `<button>`, `<img>`, `<a>`, `<table>`, etc.) — use `data-text`, `data-attr`, `data-class`, `data-style`, or `data-show` instead.

### Elements That Are Focusable by Default

These elements can receive keyboard focus without `tabindex`:

- `<a href="...">` (must have href)
- `<button>` (not disabled)
- `<input>` (not disabled, not `type="hidden"`)
- `<select>` (not disabled)
- `<textarea>` (not disabled)
- `<details>` / `<summary>`
- Elements with `contenteditable`
- `<iframe>`

To make any other element focusable: add `tabindex="0"` (in tab order) or `tabindex="-1"` (focusable via JS but not tab order).

```html
<!-- Make a div keyboard-accessible -->
<div tabindex="0" role="button"
     data-on:click="@post('/action')"
     data-on:keydown="if(evt.key==='Enter' || evt.key===' ') @post('/action')">
    Custom Button
</div>
```

### Void Elements (Self-Closing, No Children)

These elements **cannot have children** or closing tags:

`<area>`, `<base>`, `<br>`, `<col>`, `<embed>`, `<hr>`, `<img>`, `<input>`, `<link>`, `<meta>`, `<param>`, `<source>`, `<track>`, `<wbr>`

They can still use Datastar attributes:
```html
<img data-attr:src="$url" data-show="$hasImage" />
<input data-bind:name />
<hr data-show="$showDivider" />
```

### Boolean Attributes

Some HTML attributes are boolean — their presence means `true`, absence means `false`:

`disabled`, `checked`, `selected`, `required`, `readonly`, `multiple`, `autofocus`, `autoplay`, `controls`, `loop`, `muted`, `hidden`, `open` (on `<details>`/`<dialog>`), `novalidate`, `formnovalidate`, `defer`, `async`, `draggable`, `spellcheck`, `contenteditable`, `inert`

In Datastar, control them with `data-attr`:

```html
<button data-attr:disabled="$isSubmitting">Submit</button>
<input data-attr:required="$fieldRequired" />
<details data-attr:open="$expanded">...</details>
<div data-attr:inert="$modalOpen">Background content</div>
```

**`inert` attribute** (modern and very useful): Makes an element and all its descendants non-interactive (unfocusable, unclickable, invisible to assistive technology). Great for disabling background content when a modal is open.

### The `id` Attribute and Datastar Targeting

**Every element that the Go backend needs to patch MUST have an `id`**. The Datastar SDK targets elements by ID (or CSS selector):

```go
// Backend targets by ID (default — uses id from fragment's root element)
sse.PatchElements(`<div id="user-info">Updated content</div>`)

// Or target a specific selector
sse.PatchElements(`<span>New text</span>`, datastar.WithSelector("#user-name"))

// Remove by ID
sse.RemoveElementByID("notification-banner")

// Remove by selector
sse.RemoveElement(".old-items")
```

**Rule**: If the server ever needs to update, replace, or remove an element, give it a stable, unique `id`.

---

## Complete Element Decision Tree for Datastar

```
What do I need?
│
├── User types text → <input type="text"> + data-bind:signal
├── User types long text → <textarea> + data-bind:signal
├── User picks one from list → <select> + data-bind:signal
├── User toggles on/off → <input type="checkbox"> + data-bind:signal
├── User picks one from few → <input type="radio"> + data-bind:signal
├── User picks a number → <input type="number"> or <input type="range"> + data-bind:signal
├── User picks date/time → <input type="date/time/datetime-local"> + data-bind:signal
├── User picks color → <input type="color"> + data-bind:signal
├── User uploads file → <input type="file"> + data-on:change
├── User triggers action → <button type="button"> + data-on:click="@post(...)"
├── User navigates → <a href="/path">
├── User submits form → <button type="submit"> inside <form data-on:submit__prevent>
│
├── Display reactive text → <span data-text="$signal">
├── Show/hide content → <div data-show="$condition">
├── Toggle CSS class → <element data-class:name="$condition">
├── Set HTML attribute → <element data-attr:name="$expression">
├── Set inline style → <element data-style:prop="$expression">
│
├── Load data on mount → <div data-init="@get('/api/data')">
├── Stream real-time → <div data-init="@get('/api/stream', {requestCancellation:'disabled'})">
├── Poll periodically → <div data-on-interval.5000ms="@get('/api/status')">
├── Load on scroll into view → <div data-on-intersect="@get('/api/more')">
│
├── Modal dialog → <dialog> + data-ref + data-effect for show/close
├── Expandable section → <details> + <summary>
├── Disable form group → <fieldset data-attr:disabled="$submitting">
│
└── Server updates DOM → Element has id → Backend calls sse.PatchElements(...)
```
