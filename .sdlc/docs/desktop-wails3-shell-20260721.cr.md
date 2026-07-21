# Tingly DS — Code Review

**Date**: 2026-07-21
**Target**: macOS arm64 Wails 3 MVP
**Assessment**: PASS
**Critical**: 0 · **Major**: 0 · **Minor**: 3

## Minor observations

[1] 💡 External-link handling is intentionally best-effort — `navigation.go`

Only ordinary anchor clicks are redirected to the system browser. Scripted
navigation, form submissions, redirects, and `window.open` can still navigate
the WebView. This is documented as a UX helper rather than a navigation
security boundary; the native bridge remains origin-validated.

[2] 💡 Package tasks retain generated directories — `Taskfile.yml`

The package task overwrites every required executable, plist, icon, and
signature, but does not delete the prior `.app` or `.iconset` first. The current
bundle contains only expected resources and passes strict signature checks.

[3] 💡 GUI lifecycle is smoke-tested rather than unit-tested — `main.go`

URL policy and embedded assets are automated, while tray clicks, close-to-hide,
reload, and native packaging rely on a real macOS run. The user completed that
smoke test successfully.

## Issue resolved during review

The injected JavaScript originally treated any HTTPS port and a trailing-dot
hostname differently from the Go allowlist. It now applies the same HTTPS,
default-port, trailing-dot, and DeepSeek-subdomain rules as `isDeepSeekURL`.

## Out-of-scope platform findings

The pinned Wails alpha exposes different remote-URL JavaScript injection and
raw-message origin behaviour on Windows/Linux. The current Taskfile also has
macOS-specific CGO flags. These platforms were never claimed as release targets
for this iteration and are explicitly marked as requiring a separate port and
security validation.

## Strengths

- Close-to-hide and explicit-quit paths preserve and terminate the WebView as
  intended.
- No page-callable Go service is registered.
- Raw bridge input has strict origin, frame, length, JSON, scheme, host, and
  credential validation.
- Official icon provenance and checksums are recorded; `.icns`, plist, and
  ad-hoc signature are verified.

## Approval

The macOS MVP is approved. There are no blocking findings.
