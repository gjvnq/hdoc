# hdoc
A simple tool for writing nice standalone HTML documents.

# File Format

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

## Templates

TODO: something like "static" web components.

Use ```<template>``` with the attribute `name` to specify new elements that will be resolved on compile time.

Example:

```html
<template name="note">
  <div class="note">
    <strong>Note <h-val select="@data-counter"/></strong>
    <h-val select="."/>
  </div>
</template>
```

## Counters

Use ```<h-counter>``` with attributes `name`, `set` and `style` to configure a counter. The `name` attribute is the element number to count. The `style` can be `0` (for indo-arabic numbers), `i` (for lowercase roman numerals) and `I` (for uppercase roman numerals). Example:

```html
<h-counter name="h1" set="1" style="I">
```

To disable the counter on any specific element, use `no-counter` attribute to prevent counting. Example: ```<h1 no-counter>Preface</h1>```

## Table of Contents

Use ```<toc>``` to auto generate the table of contents. Use the `elems` attribute to specify which elements to count. The inner text will be used as the title. Example:

```html
<toc elems="h1 h2 h3">The Table of Contents</toc>
```

```html
<toc elems="figure">List of Figures</toc>
```

The attribute `no-toc` can be used to exclude an element from the TOC.
