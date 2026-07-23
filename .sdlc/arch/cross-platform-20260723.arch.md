# Tingly DS Cross-Platform Architecture

**Last Updated**: 2026-07-23
**Cache Level**: Module
**Expires**: 2026-08-06
**Branch**: main
**Hash**: 844bffc9006357180e818b9a00dd8c092436add5
**Parent**: overview-20260721.arch.md

## Overview

Tingly DS is a Go/Wails 3 remote-WebView shell. The runtime code is largely
portable, but the current build and release workflow is macOS-only. Completing
Windows and Linux support primarily requires platform packaging metadata,
platform-neutral build tasks, and explicit Linux lifecycle configuration.

## Components

| Component | Location | Cross-platform role |
| --- | --- | --- |
| Native shell | `main.go` | Creates the app, WebView window, tray, menu, and hide-on-close lifecycle |
| Navigation policy | `navigation.go` | Platform-neutral external-link validation and browser handoff |
| Embedded icons | `assets.go`, `build/icons/` | Dock/window and tray images embedded in the Go binary |
| Build tasks | `Taskfile.yml` | Currently injects macOS compiler flags globally and packages only `.app` |
| macOS metadata | `build/darwin/Info.plist` | Stable bundle identity and icon metadata |
| Packaging tests | `packaging_test.go` | Currently verifies only the macOS identity |

## Dependencies

- Wails 3 `v3.0.0-alpha2.114` supplies native windows, tray integration,
  WebView backends, and platform artifact generators.
- Windows uses WebView2 and needs generated `.syso` resources plus WebView2
  runtime installation/bootstrap.
- Linux uses GTK3/WebKitGTK and needs CGO, a `.desktop` entry, and AppImage
  packaging dependencies.
- macOS uses WKWebView and the existing `.app` bundle workflow.

## Data Flow

All targets create the same persistent WebView at `https://chat.deepseek.com/`.
The native shell controls visibility and the system tray. DeepSeek-domain links
stay embedded, while validated external HTTP(S) links are passed to the target
OS default browser. Website data persistence is provided by each native WebView
backend rather than by repository-managed files.

## Key Patterns

- Runtime configuration is centralized in `application.Options`.
- Window close is cancelled and converted to hide so the WebView session stays
  alive; explicit Quit bypasses that behavior.
- macOS uses a monochrome template tray icon; Windows and Linux use the colour
  tray icon.
- Packaging identity must remain `dev.tingly.tingly-ds` / `tingly-ds` across
  releases so native WebView profile locations remain stable.

## Platform Gaps

- `Taskfile.yml` applies macOS deployment flags to generic build/test tasks.
- There is no Windows icon, manifest, version metadata, `.syso` generation, or
  `.exe` packaging task.
- There is no Linux `ProgramName`, lifecycle option, `.desktop` metadata, or
  AppImage task.
- The generic `package` task is restricted to Darwin.
- Documentation claims only macOS support and describes only macOS storage and
  prerequisites.
- Packaging tests do not protect Windows/Linux identity metadata.

## Integration Points

- Windows: Wails WebView2 backend, generated Win32 resources, Windows system
  tray, default browser, and optional WebView2 bootstrapper.
- Linux: Wails GTK/WebKitGTK backend, freedesktop `.desktop` integration,
  AppImage, status notifier/system tray support, and default browser.
- macOS: existing WKWebView, menu bar, Dock, and signed `.app` bundle.
