# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & run

```bash
go build                     # build binary
go run .                     # run directly
go test ./...                # run all tests
```

## Architecture

A TUI application for inspecting listening TCP/UDP ports, built with [Bubble Tea v2](https://charm.land/bubbletea/v2).

```
main.go                  → entrypoint: init i18n, then start the TUI
src/
  connection/net.go      → data layer: calls gopsutil to enumerate listening ports + process details
  langs/langs.go         → i18n layer: wraps go-i18n with TOML locale files embedded via embed.FS
  view/
    table.go             → Bubble Tea model + program: table rendering, 1s auto-refresh tick, key handling
    help.go              → key binding definitions (j/k/↑/↓, r refresh, p pause, q quit, h help)
    format.go            → display formatters for protocol, family, CPU%, memory (humanize.Bytes)
    msg.go               → custom tea.Msg type (tickMsg) for the refresh timer
locales/
  fs.go                  → `embed` directive for active.*.toml files
  active.en.toml         → English column headers
  active.zh_CN.toml      → Chinese column headers
```

**Data flow:** `connection.ListeningPorts()` returns `[]ListenPortInfo` → `convertToRows()` maps to `[]table.Row` → table widget renders with columns defined by `getColumns()`.

**Auto-refresh:** Each `tickMsg` (1s interval via `tea.Tick`) triggers a fresh `connection.ListeningPorts()` call. Press `p` to toggle pause.

The Bubble Tea model is defined in `src/view/table.go` (not in `model.go`, which is effectively empty). The `Render()` function in `table.go` is the TUI entry point called from `main.go`.

## i18n

Locale files are TOML and embedded at compile time via `//go:embed`. The `langs.Setup()` call in main initializes the bundle and localizer. Currently hardcoded to English; `initLocalizer(language.SimplifiedChinese.String())` is commented out. Column headers are translated via `langs.Localize("KEY")` at TUI startup.
