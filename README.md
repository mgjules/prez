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
presenterm --config-file config.yaml presentation-file.md
```

For example:

```bash
presenterm --config-file config.yaml 2025-01-03-developer-ergonomics.md
```

### Keyboard Controls

- <kbd>→</kbd> / <kbd>Space</kbd> / <kbd>Enter</kbd> / <kbd>n</kbd>: Next slide
- <kbd>←</kbd> / <kbd>Backspace</kbd> / <kbd>p</kbd>: Previous slide
- <kbd>q</kbd> / <kbd>Esc</kbd>: Quit
- <kbd>r</kbd>: Reload presentation
- <kbd>?</kbd> / <kbd>h</kbd>: Help

## Sample Presentations

- [Developer Ergonomics](2025-01-03-developer-ergonomics.md) - Optimizing workflow and reducing fatigue
- [PDE (Personalized Development Environment)](2024-05-18-pde.md) - Customizing your development tools
- [Stand Out from the Crowd](2023-29-04-stand-out-from-the-crowd.md) - Tips to improve your CV

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
date: 2025-01-03
theme:
  name: tokyonight-storm
---

# First Slide
Content goes here
```

### Custom Configuration

This repository includes a `config.yaml` file with custom settings for presenterm:

- Tokyo Night Storm theme by default
- Keyboard shortcuts
- Image protocol settings
- Code snippet execution options

Always use the `--config-file` flag to ensure consistent presentation appearance.

