# Windows and Linux Support Understanding

**Date**: 2026-07-23
**Scope**: Runtime lifecycle, native metadata, build tasks, packaging, and user documentation
**Revision**: `844bffc9006357180e818b9a00dd8c092436add5`

## Current State

The Go runtime already uses portable Wails APIs and supplies Windows-specific
quit behavior. macOS is the only release target because its compiler flags,
bundle metadata, icon generation, signing, documentation, and tests are the
only complete platform path.

## Findings

1. Generic `build` and `test` tasks export macOS deployment variables on every
   operating system.
2. The hidden-window lifecycle requires disabling automatic quit on both
   Windows and Linux; only Windows is configured today.
3. Linux needs a stable `ProgramName` matching the executable and `.desktop`
   entry to preserve correct window grouping and icon association.
4. Windows requires an `.ico`, executable manifest, version-info JSON, and a
   generated architecture-specific `.syso` resource before building.
5. Linux can use the Wails `.desktop` and AppImage generators with the existing
   512×512 tiled app icon.
6. The existing `package` task name should remain the public entry point and
   dispatch to an OS-specific implementation.
7. A host can fully execute only its native Wails build because Linux requires
   CGO system libraries and Windows embeds target-specific resources. Static
   tests should validate metadata on any host; actual packages must be built and
   smoke-tested on each target.

## Planned Integration

- Add Linux lifecycle and program-name options to `main.go`.
- Split build/test environment settings so macOS flags are Darwin-only.
- Add Windows resource generation and production `.exe` packaging.
- Add Linux `.desktop` generation and AppImage packaging.
- Keep the existing ad-hoc macOS bundle task unchanged behind platform dispatch.
- Extend packaging tests and README instructions for all three desktop targets.

## Risks

- Wails 3 remains alpha; remote URL bridge behavior must be smoke-tested against
  WebView2 and WebKitGTK, not inferred from macOS WKWebView behavior.
- Linux tray visibility depends on the desktop environment's status-notifier
  support.
- Windows users without the evergreen WebView2 runtime need installation via
  the Microsoft bootstrapper or an installer that carries it.
- AppImage creation downloads continuous-release tooling and therefore requires
  network access during packaging.
