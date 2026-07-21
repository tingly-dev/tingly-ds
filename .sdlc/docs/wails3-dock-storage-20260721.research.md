# Wails 3 Dock and WebKit Storage Research

**Date**: 2026-07-21
**Version inspected**: `github.com/wailsapp/wails/v3@v3.0.0-alpha2.114`

## Findings

- `ActivationPolicyRegular` is Wails' default policy for applications with a
  user interface. `ActivationPolicyAccessory` intentionally omits the Dock.
- macOS `LSUIElement=true` also marks the bundle as an agent application and
  suppresses its Dock presence. Both settings must agree on regular-app UX.
- Wails installs an `ApplicationShouldHandleReopen` common-event handler on
  macOS. When no window is visible, it shows all hidden windows. Tingly DS does
  not need an additional competing handler.
- Wails creates `WKWebViewConfiguration` without selecting a non-persistent
  website data store. WebKit therefore uses its default persistent store.
- Website data is external to the `.app` bundle. Stable packaged identity
  (`dev.tingly.tingly-ds`, executable `tingly-ds`) lets rebuilds reuse it.
- A directly executed development binary may have a different process identity
  and WebKit profile than the packaged `.app`; documentation should not present
  the development launcher as the daily-use installation.

## Decision

Use regular activation policy, remove `LSUIElement`, rely on Wails' built-in
Dock reopen behavior, and guard the stable package identity with an automated
metadata test. Document the packaged-app workflow and the dev/profile boundary.
