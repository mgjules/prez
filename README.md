# prez

My presentation slides using [presenterm](https://github.com/mfontanini/presenterm), a terminal-based presentation tool.

## About

This repository contains markdown-based presentations on various software development topics. All presentations are designed to be viewed with the presenterm tool.

## Installation

To view these presentations, install presenterm:

```bash
# Using Cargo
cargo install presenterm

# Using Homebrew
brew install presenterm

# Using Arch Linux (AUR)
yay -S presenterm
```

## Usage

To view a presentation:

```bash
presenterm presentation-file.md
```

For example:

```bash
presenterm developer-ergonomics.md
```

### Keyboard Controls

- <kbd>→</kbd> / <kbd>Space</kbd> / <kbd>Enter</kbd> / <kbd>n</kbd>: Next slide
- <kbd>←</kbd> / <kbd>Backspace</kbd> / <kbd>p</kbd>: Previous slide
- <kbd>q</kbd> / <kbd>Esc</kbd>: Quit
- <kbd>r</kbd>: Reload presentation
- <kbd>?</kbd> / <kbd>h</kbd>: Help

## Sample Presentations

- [Developer Ergonomics](developer-ergonomics.md) - Optimizing workflow and reducing fatigue
- [PDE (Personalized Development Environment)](2024-05-18-pde.md) - Customizing your development tools

## Creating Presentations

Create markdown files with slide separators (---). See [reference.md](reference.md) for syntax examples.

```markdown
# Slide 1

Content for first slide

---

# Slide 2

Content for second slide
```

You can also use YAML frontmatter at the beginning of your presentation to customize options:

```markdown
---
author: Your Name
title: Presentation Title
---

# First Slide
Content goes here
```

