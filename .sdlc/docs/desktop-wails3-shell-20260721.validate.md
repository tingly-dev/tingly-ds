# DeepSeek Desktop Shell — Goal Validation

**Date**: 2026-07-21
**Status**: PASS

## Acceptance criteria

| Criterion | Status | Evidence |
| --- | --- | --- |
| Repository layout | PASS | Go source, tests, metadata, and docs live at the repository root; generated `bin/` is ignored |
| Native lightweight shell | PASS | Wails uses the OS WebView; final `.app` is 8.1 MiB |
| Direct DeepSeek experience | PASS | Main WebView URL is `https://chat.deepseek.com/` |
| Tray-controlled lifecycle | PASS | Toggle, show/hide, reload, browser, and quit callbacks are wired; user smoke test passed |
| Preserve session while hidden | PASS | Close hook cancels destruction and hides the existing window |
| Restrained native bridge | PASS | No Wails services; only origin-validated external URL opening |
| Official DeepSeek artwork | PASS | Template tray icon embedded; Avatar compiled into `AppIcon.icns` |
| Reproducible local package | PASS | Pinned Wails version, tests, package, plist, and ad-hoc signature verified |

## Validation notes

- Explicit tray callbacks are used instead of an attached popover window, so
  authentication dialogs and file pickers do not hide the chat on focus loss.
- The app was intentionally left free of a frontend framework and bundled
  Chromium runtime.
- The package is suitable for local use. Public distribution still requires a
  Developer ID signature and Apple notarisation.
- This validation applies to macOS arm64. Windows/Linux remain future ports and
  require platform-specific remote-navigation and bridge testing.
