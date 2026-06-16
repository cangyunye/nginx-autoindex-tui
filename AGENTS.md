# nginx-autoindex-tui

A Go + Bubble Tea terminal file browser and downloader for nginx autoindex servers.  
Parses HTML (goquery) and JSON autoindex listings, renders an interactive table, supports
index navigation, mark-mode batch download via wget, and directory history stack.

## Project

| | |
|---|---|
| **Language** | Go 1.21+ |
| **TUI** | `github.com/charmbracelet/bubbletea` (Elm architecture) |
| **Styling** | `github.com/charmbracelet/lipgloss` + `lipgloss/table` |
| **HTML parsing** | `github.com/PuerkitoBio/goquery` |
| **CLI** | `github.com/spf13/cobra` |
| **Download** | `os/exec` → system `wget` |
| **Entry point** | `main.go` (CLI parsing → `tea.NewProgram`) |

## Commands

```bash
make build          # 优化编译 (-ldflags "-w -s"), 输出 nginx-autoindex-tui
make test           # go test ./tests/... -v
make tidy           # go mod tidy
make clean          # 删除二进制 + 构建缓存
```

## Architecture

Four internal packages:

- **`main.go`** — CLI entry: parses flags (`--force`, `--concurrent`, `--output-dir`,
  `--insecure`, `--user-agent`, `--border-style`, `--theme`) and positional `[URL]` via cobra,
  creates the TUI model, launches `tea.NewProgram`.

- **`internal/fetcher/`** — HTTP layer. `FetchHTML(rawURL)` performs GET with 30 s timeout,
  auto-prepends `http://` when scheme missing. `ResolveURL(base, href)` uses
  `url.ResolveReference`. Detects Content-Type to decide HTML vs JSON parsing.

- **`internal/parser/`** — Parses autoindex listings. `Parse(r io.Reader)` handles HTML
  (goquery → `<h1>` title + `<pre>` / `<a>` entries) and JSON (`name`/`type`/`mtime`/`size`).
  Output: `*Page` with `Title` + `[]Entry` (Href, Name, DateTime, Size, IsDir).
  Display names are URL-decoded; hrefs stay encoded for requests.

- **`internal/tui/`** — Bubble Tea model + styles.
  - **`model.go`**: Elm-architected `Model` (Init/Update/View). Holds current URL,
    parsed page, cursor, index-input buffer, history stack, loading/error state,
    download queue (buffered channel + goroutine pool, concurrency default 3),
    mark-mode state, overwrite-confirmation flow, directory cache (LRU, 5 entries),
    download-history panel.
  - **`styles.go`**: Lip Gloss styles for title, help text, errors, meta info.
    Supports `dark` / `light` / `mono` themes and `normal` / `rounded` / `ascii`
    border styles.

## Conventions

- **Elm architecture**: Every screen state change goes through `Update(tea.Msg)`.
  View is a pure rendering function. No mutation outside Update.
- **Parser is independently testable**: `Parse` takes `io.Reader`, not URL/network.
  Tests use table-driven cases over fixed HTML/JSON fixtures.
- **URL safety**: All paths are resolved via `url.ResolveReference` and validated
  against the root URL prefix. `../` cannot escape the root.
- **Display vs request encoding**: Show users URL-decoded names; keep hrefs encoded
  for HTTP requests.
- **Download via wget**: Delegated to system `wget` via `os/exec`. No Go-native
  file writing. Max 3 concurrent wget processes by default.
- **Mark mode**: Separate state (`v` to enter, `x` to download marked, `c` to clear,
  `v`/`Esc`/`q` to exit keeping marks).
- **Index navigation**: Three-digit number input with real-time cursor sync.
  `Backspace` edits last digit; `Esc` clears input.
- **History stack**: Simple push/pop linear stack (no forward/backward branching).
- **Overwrite confirmation**: Single-file: `[Y]es/[N]o/[A]lways`. Batch: aggregate
  count + `[Y]es all/[N]o all/[A]sk each`. `--force` flag skips all prompts.
- **Controlled dependencies**: Cobra used only for flag/positional-arg parsing,
  not for subcommands or config (viper).

## Notes

(Add project-specific observations here over time.)
