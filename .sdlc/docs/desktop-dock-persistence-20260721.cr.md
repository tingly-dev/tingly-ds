# Dock and Persistence Code Review

**Date**: 2026-07-21
**Verdict**: APPROVED

## Review notes

- Both macOS mechanisms that hid the Dock icon were addressed: runtime
  activation policy and bundle metadata.
- No redundant Dock reopen listener was added because the pinned Wails version
  already implements the expected hidden-window behavior.
- Persistence relies on the native WebView's supported default rather than
  copying cookies or creating a custom credential database.
- A focused metadata regression test protects the storage identity and Dock
  packaging contract.
- Build-from-source documentation uses the exact pinned Wails CLI version.

No blocking findings remain.
