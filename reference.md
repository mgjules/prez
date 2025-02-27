# Welcome to Presenterm

A terminal-based presentation tool

```go
package main

import "fmt"

func main() {
  fmt.Println("Written in Go!")
}
```

---

## Everything is markdown

In fact this entire presentation is a markdown file

---

# h1

## h2

### h3

#### h4

##### h5

###### h6

---

# Markdown components

You can use everything in markdown!

- Like bulleted list
- You know the deal

1. Numbered lists too

---

# Tables

| Tables | Too    |
| ------ | ------ |
| Even   | Tables |

---

# Graphs

```
digraph {
    rankdir = LR;
    a -> b;
    b -> c;
}
```

```
┌───┐     ┌───┐     ┌───┐
│ a │ ──▶ │ b │ ──▶ │ c │
└───┘     └───┘     └───┘
```

---

All you need to do is separate slides with triple dashes
`---` on a separate line, like so:

```markdown
# Slide 1

Some stuff

---

# Slide 2

Some other stuff
```

---

# Presenterm Features

You can use YAML frontmatter at the beginning of your presentation:

```yaml
---
author: Your Name
title: Presentation Title
date: 2025-02-27
---
```

---

# Customization

Presenterm supports themes and custom colors that can be defined in a config file.

---

# Slide Transitions

Presenterm supports different slide transition effects that can be configured.

---

# Terminal Fidelity

Enjoy true terminal colors and Unicode support
