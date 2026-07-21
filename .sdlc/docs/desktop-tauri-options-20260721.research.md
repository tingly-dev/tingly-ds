# Tauri Desktop Shell Research

**Status**: Historical alternative; superseded by the implemented Wails 3 shell

## Research question

What is the smallest practical implementation for a tray-accessible DeepSeek desktop app with persistent login?

## Findings

- Safari web apps are simplest but do not provide the desired custom tray workflow.
- A native Swift/AppKit shell is smallest on macOS but would need a separate Windows implementation.
- Tauri 2 uses the operating system WebView, supports tray events, and can load a remote HTTPS URL directly.
- Electron includes Chromium and is unnecessarily large for a one-site shell.

## Recommendation

At the time of research, Tauri 2 was the recommended Rust option. The user
subsequently selected and implemented the Go/Wails 3 variant. This document is
retained for comparison, not as the current build plan.

## Risks

- DeepSeek can change its domains or block embedded WebViews.
- OAuth popups, downloads, microphone permissions and CAPTCHA need real-site validation.
- Runtime memory will be dominated by the DeepSeek page rather than the small native wrapper.
