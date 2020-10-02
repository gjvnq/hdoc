# hdoc
A simple tool for writing nice standalone HTML documents.

## Features

  * Custom elements (like Web Components, but static)
  * Auto include local files via data URIs.
  * Auto numbering of sections.
  * Auto section link on heading.
  * Auto table of contents.
  * Password protect files (requires JS for decryption when viewing).

## File Format

The file format is basically an HTML shorthand. 

```html
<hdoc>
  <opts>
    <!--- See Options subsection -->
  </opts>
  <head>
    <!-- same as HTML -->
  </head>
  <!-- same as HTML <body> -->
</hdoc>
```

## Processing Order

1. Includes
2. `<abbr>` and `<dfn>`
3. Generate `id` attributes
4. Templates
5. Counters

### Templates

Templates are a very small subsection of XSLT. Use the `priority` attribute to go around order issues. The default templates have priority 0 and the user ones have priority 1 unless specified otherwise.

Example:

```html
<template match="note">
  <p>
    <copy-of select="@*"/>
    <strong>Note: </strong>
    <copy-of select="node()|*"/>
  </p>
</template>
```

### Counters

Use ```<h-counter>``` with attributes `name`, `set` and `style` to configure a counter. The `name` attribute is the element number to count. Example:

```html
<h-counter name="h1" set="1" style="I">
```

To disable the counter on any specific element, use `no-counter` attribute to prevent counting. Example: ```<h1 no-counter>Preface</h1>```

Possible styles:

  * `0`: indo-arabic numerals.
  * `i`: lowercase roman numerals.
  * `I`: uppercase roman numerals.
  * `a`: lowercase latin letters.
  * `A`: uppercase latin letters.
  * `α`: lowercase greek letters.

To display the counter, use ```<h-counter>``` using without the `set` attribute.

To auto reset, use: `<auto-reset>` inside `<opts>`. Use the attributes `name` to specify which counter and `elems` to specify which elements will reset the counter, and `start` to specify the starting value (usually `0` or `1`). Example:

```html
<auto-reset name="h3" elems="h1 h2" start="1">
```

### Table of Contents

Use ```<toc>``` to auto generate the table of contents. Use the `elems` attribute to specify which elements to count. The inner text will be used as the title. Example:

```html
<toc elems="h1 h2 h3">The Table of Contents</toc>
```

```html
<toc elems="figure">List of Figures</toc>
```

The attribute `no-toc` can be used to exclude an element from the TOC.

### Abbreviations

Use `<abbr w="WORD"/>` (must have no inner content) will be automatically be replaced by the propper `<abbr>` tag. If the `w` (word) attribute is absent on `<abbr>` with the `title` attribute, it will be infered via the text conents. Example:

```html
<abbr title="Uniform Resource Name">URN</abbr> lorem ipsum dolor est ... <abbr w="URN"/> ...
<abbr title="Uniform Resource Location" w="ABC">URL</abbr> lorem ipsum dolor est ... <abbr w="ABC"/> ...
```

Becomes:

```html
<abbr title="Uniform Resource Name">URN</abbr> lorem ipsum dolor est ... <abbr title="Uniform Resource Name">URN</abbr> ...
<abbr title="Uniform Resource Location">URL</abbr> lorem ipsum dolor est ... <abbr title="Uniform Resource Location">URL</abbr> ...
```

All `<abbr>` will be mobile-friendly ([see approach](https://bitsofco.de/making-abbr-work-for-touchscreen-keyboard-mouse/)).

### Definitions

Use `<dfn w="WORD"/>` (must have no inner content) to include and link to the definition of a word. If the `w` (word) attribute is absent on the contentful `<dfn>`, it will be infered via the text conents. Example:

```html
<dfn>HTML</dfn> is a standard ....
In an <dfn w="HTML"/> document, ...
```

Becomes:

```html
<dfn id="dfn-HTML">HTML</dfn> is a standard ....
In an <a href="dfn-HTML" class="dfn-a">HTML</a> document, ...
```

### Syntax Highlight

### Include

Use `<include>` with the `src` attribute to include another file into this one. Use the `escape` attribute to avoid HTML injections.
