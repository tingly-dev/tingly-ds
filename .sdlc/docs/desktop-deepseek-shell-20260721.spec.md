# DeepSeek Desktop Shell — MVP Specification

**Status**: Superseded historical Tauri proposal
**Current implementation**: `desktop-wails3-shell-20260721.spec.md`

> This document is retained to preserve the original alternative analysis. It
> does not describe the source currently present in this repository.

## Motivation

Opening DeepSeek in a general browser adds friction. The MVP should provide an app-like, persistent, menu-bar-accessible entry point without bundling Chromium or recreating the website.

## Scope

### Functional requirements

- Load `https://chat.deepseek.com/` in a native system WebView.
- Preserve cookies, local storage and cache across launches.
- Create a system tray icon with Open/Hide, Reload and Quit actions.
- Toggle and focus the window from a tray click.
- Hide the window instead of terminating when its close button is pressed.
- Keep DeepSeek-owned HTTPS navigation inside the WebView.
- Open unrelated HTTPS links in the operating system browser.
- Show a clear load error page when initial navigation fails.

### Non-functional requirements

- Use Tauri 2 and no frontend framework.
- Do not expose Tauri commands or plugins to the remote page.
- Keep dependencies and binary size small.
- Support macOS first without preventing later Windows builds.

## Components

- `main.rs`: process entry point.
- `lib.rs`: app setup, WebView creation, tray events and close behavior.
- `navigation.rs`: allowlist policy with unit tests.
- `tauri.conf.json`: app identity, bundle targets and security defaults.
- `capabilities/default.json`: local-only window permissions; no remote URL grants.

## Security

- Only HTTPS DeepSeek hosts (`deepseek.com` and subdomains) may remain embedded.
- Localhost, custom schemes and unrelated origins are denied.
- The remote WebView is not listed in a Tauri remote capability.
- No shell, filesystem, clipboard or process APIs are exposed to page JavaScript.

## Testing and validation

- Unit-test allowed and denied navigation URLs.
- Run Rust formatting, unit tests and `cargo check`.
- Build the application bundle when the platform toolchain is available.
- Manually verify tray interaction and login only if GUI launch is authorized.

## Deferred

- Global shortcut and launch-at-login.
- Multiple AI sites/profiles.
- Native notifications, automatic updates and code signing.
- Download manager and microphone permission UI.
