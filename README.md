# Blow

**Better glow** — a fork of [charmbracelet/glow](https://github.com/charmbracelet/glow) optimized for tmux workflows.

## Why Blow?

tmux + TUI markdown reader = pain:

- **Block selection bleeds across panes** — tmux's mouse selection crosses pane boundaries
- **Scroll + copy jumps to bottom** — copying in tmux's scroll mode snaps back to the bottom
- **No text selection in TUI** — glow's TUI mode has no mouse drag selection
- **No search in pager** — can't search within a rendered document

Blow fixes all of these by capturing mouse events inside the TUI and handling selection, scrolling, and clipboard internally.

## Features

- **TUI by default** — `blow README.md` opens TUI mode (no `-t` flag needed)
- **Mouse drag selection** — click and drag to select text, auto-copied to clipboard on release
- **In-TUI search** — press `/` to search, `n`/`N` to navigate matches, `Esc` to clear
- **tmux-safe scrolling** — scroll stays within the TUI pane, never bleeds across tmux panes

## Build

Requires Go 1.21+.

```bash
git clone https://github.com/eomjaeeun/better-glow.git
cd better-glow
go build -o blow
```

Binary is output as `./blow`.

### Install to PATH

```bash
go build -o blow && mv blow /usr/local/bin/
```

## Usage

```bash
# Open a file in TUI mode (default)
blow README.md

# Browse markdown files in current directory
blow

# Browse markdown files in a specific directory
blow ~/docs
```

### Keybindings (Pager)

| Key | Action |
|-----|--------|
| `/` | Start search |
| `n` / `N` | Next / previous match |
| `Esc` | Clear search results |
| `q` | Quit |
| `?` | Show help |

### Mouse

- **Click + drag** to select text
- **Release** to copy selection to clipboard
- **Scroll wheel** to scroll within the document

## Config

Config file location: `~/.config/blow/blow.yml`

```yaml
style: "dark"
mouse: true
width: 80
```

## Credits

Forked from [charmbracelet/glow](https://github.com/charmbracelet/glow). Part of [Project Better](https://github.com/eomjaeeun).

## License

[MIT](LICENSE)
