# prez

My presentation slides using [slides](https://github.com/maaslalani/slides), a terminal-based presentation tool.

## About

This repository contains markdown-based presentations on various software development topics. All presentations are designed to be viewed with the slides tool.

## Installation

To view these presentations, install slides:

```bash
# Using Go
go install github.com/maaslalani/slides@latest

# Using Homebrew
brew install slides

# Using Arch Linux
yay -S slides
```

## Usage

To view a presentation:

```bash
slides presentation-file.md
```

For example:

```bash
slides developer-ergonomics.md
```

### Keyboard Controls

- <kbd>→</kbd> / <kbd>Space</kbd>: Next slide
- <kbd>←</kbd> / <kbd>Backspace</kbd>: Previous slide
- <kbd>q</kbd>: Quit
- <kbd>r</kbd>: Reload presentation
- <kbd>?</kbd>: Help

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

