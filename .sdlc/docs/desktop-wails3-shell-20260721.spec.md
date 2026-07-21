# Tingly DS (Wails 3) Specification

**Version**: 0.1.0
**Status**: Implemented and validated on macOS arm64
**Last Updated**: 2026-07-21

## Motivation

DeepSeek Chat has no lightweight desktop entry point. Opening a general browser
adds friction, while bundling Chromium would add unnecessary size and runtime
overhead. A Wails 3 shell can reuse the operating system WebView and expose the
site through a persistent tray-controlled window.

## Proposal

Build a macOS-first Go application directly at the repository root. The
application loads `https://chat.deepseek.com/` directly, owns no chat UI,
registers no Go services, and exposes only native window/tray actions. Closing
the window hides it; the process exits only from the tray menu.

Use explicit tray callbacks rather than `AttachWindow`, so file pickers and
authentication dialogs do not make the chat window disappear on focus loss.

## Scope

### Functional requirements

1. Load DeepSeek Chat in a persistent native WebView.
2. Show the main window on normal application launch.
3. Appear as a normal application in the macOS Dock while retaining the tray.
4. Restore the hidden window when its Dock icon is clicked.
5. Toggle and focus the window with a left click on the tray icon.
6. Provide tray actions: Open/Hide, Reload, Open in Browser, and Quit.
7. Hide rather than destroy the window when the close button is used.
8. Hide the window on Escape.
9. Keep DeepSeek HTTPS links in the WebView.
10. Open unrelated HTTP(S) anchor links in the default browser.
11. Preserve normal WebView cookies, local storage, and cache between launches
    and rebuilds by keeping the packaged application identity stable.
12. Use the DeepSeek mark for both the native application and tray entry.

### Non-functional requirements

- Pin `github.com/wailsapp/wails/v3` to `v3.0.0-alpha2.114` and use Go 1.25+.
- Use `github.com/tingly-dev/tingly-ds` as the Go module path.
- Use no frontend framework and no bundled browser engine.
- Keep source, tests, build metadata, and SDLC records at the repository root.
- Support macOS arm64 immediately and avoid blocking later Windows builds.
- Build without warnings from `go vet` and pass `go test ./...`.

## Components

| File | Responsibility |
| --- | --- |
| `main.go` | Application, window, tray, menu, close hook, and raw-message wiring |
| `assets.go`, `build/icons/` | Embedded tray marks and macOS application-icon source |
| `navigation.go` | URL allowlist, external-link message validation, injected click helper |
| `assets_test.go`, `navigation_test.go` | Icon integrity and navigation-policy tests |
| `go.mod` / `go.sum` | Reproducible Go dependency graph |
| `Taskfile.yml`, `build/` | Wails development and macOS package tasks |
| `README.md`, `THIRD_PARTY_NOTICES.md` | Usage, limitations, and asset provenance |

## Lifecycle

```text
launch -> create WebView -> show window
dock click while hidden -> show window (Wails macOS common-event handling)
tray left click -> show/focus OR hide
window close / Escape -> hide, keep WebView alive
tray Reload -> reload current page
tray Quit -> terminate application
```

## Native message contract

The injected page helper may send only this message shape:

```json
{"type":"open-external","url":"https://example.com/path"}
```

The Go handler accepts it only when all conditions hold:

- payload is within a small fixed size limit and valid JSON;
- message comes from the main frame;
- sender origin is HTTPS on `deepseek.com` or one of its subdomains;
- target uses HTTP or HTTPS, has no embedded credentials, and is not a
  DeepSeek URL.

Accepted targets are passed to the operating system browser. Everything else
is ignored.

## Security

- Register no Wails services or page-callable business methods.
- Leave `AllowSimpleEventEmit` disabled.
- Disable native file-drop bridging and release DevTools.
- Configure all WebView permission kinds as `PermissionDefault`; this makes
  Windows use prompts rather than its legacy blanket grant and lets macOS TCC
  handle user consent.
- Do not inject credentials, tokens, cookies, or request headers.
- Treat JavaScript link interception as a user-experience policy, not a hard
  security boundary.

## Error handling

- Return a non-zero process exit if Wails fails to start.
- Ignore malformed/untrusted raw messages without side effects.
- Log an external-browser launch failure in development builds.
- Let WebView render its normal network error page when DeepSeek is unavailable.

## Testing and validation

- Unit-test exact DeepSeek hosts, subdomains, lookalikes, schemes, credentials,
  malformed origins, malformed payloads, and oversized messages.
- Run `gofmt`, `go test ./...`, `go vet ./...`, and `go build`.
- Run the Wails production package task and inspect the generated `.app`.
- Verify the generated `.icns`, bundle icon key, and embedded tray PNGs.
- Verify `LSUIElement` is absent and the bundle identifier/executable remain
  `dev.tingly.tingly-ds` / `tingly-ds`.
- Rebuild and replace the `.app`, then verify the DeepSeek login remains.
- If GUI execution is available, verify initial load, tray toggle, close-to-hide,
  reload, external links, and file input.

## Deferred

- Global shortcuts, launch at login, updater, signing/notarisation, downloads,
  multiple AI sites, profiles, and original custom product artwork.

## Alternatives considered

- `SystemTray.AttachWindow`: rejected for the primary window because it hides on
  focus loss.
- Wails 2: stable, but lacks the first-class Wails 3 tray/window APIs used here.
- Tauri 2: researched as an alternative; its historical proposal is retained
  in `.sdlc/docs/`, but no Rust source is part of the current repository.
