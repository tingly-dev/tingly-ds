# Wails 3 DeepSeek Shell — Research Summary

**Topic:** A lightweight, menu-bar-accessible wrapper for DeepSeek Chat.

**Research question:** Which Wails 3 APIs and lifecycle pattern provide a
small, usable shell without exposing unnecessary native capabilities to a
remote page?

## Key findings

- The locally installed CLI and module are both
  `github.com/wailsapp/wails/v3@v3.0.0-alpha2.114`; the implementation pins that
  version to avoid Alpha API drift.
- `WebviewWindowOptions.URL` loads an external HTTPS application directly.
- `SystemTray.AttachWindow` automatically hides its window when focus is lost.
  That is ideal for a transient popover, but can be disruptive when a chat app
  opens a native file chooser or authentication window.
- Wails injects only its minimal core bridge into a directly loaded remote
  page. Registering no `Services`, leaving `AllowSimpleEventEmit` disabled, and
  validating `RawMessageHandler` input keeps the native attack surface narrow.
- macOS `ActivationPolicyAccessory` provides a menu-bar-only application, while
  the Windows option can prevent quitting when the last window is hidden.

## Options analysed

### Option 1: `SystemTray.AttachWindow`

- Pros: very little code; automatic positioning and toggle behaviour.
- Cons: also hides on focus loss, which is undesirable for uploads, system
  dialogs, and longer chat sessions.
- Best for: a small command palette or status popover.

### Option 2: Explicit tray callbacks

- Pros: predictable app-like window lifecycle; close-to-hide; file dialogs do
  not implicitly dismiss the app; tray menu actions remain explicit.
- Cons: slightly more lifecycle code and no automatic tray-relative position.
- Best for: this persistent chat shell.

## Recommendation

Use an ordinary resizable window with explicit tray toggle callbacks. Start it
visible on first launch so login and loading failures are discoverable, hide it
instead of destroying it on close, and keep the process alive until the tray
Quit action is selected. Inject a small navigation helper that captures
unrelated HTTP(S) links and sends them to a Go raw-message handler; validate
both the sender origin and target URL before opening the system browser.

## Sources

- https://v3.wails.io/features/windows/options/
- https://v3.wails.io/whats-new/
- https://v3.wails.io/status/
- Local Wails module source at `v3.0.0-alpha2.114`.

## Further investigation

- Manually verify DeepSeek login, file upload, clipboard, and download flows in
  a signed macOS build.
- Recheck APIs before upgrading to a newer Wails 3 Alpha/Beta release.
