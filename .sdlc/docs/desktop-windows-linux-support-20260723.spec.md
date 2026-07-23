# Windows and Linux Release Support Specification

**Date**: 2026-07-23
**Status**: Approved
**Architecture**: `.sdlc/arch/cross-platform-20260723.arch.md`

## Motivation

Tingly DS runtime code is mostly platform-neutral, but its build flags,
packaging metadata, validation, and documentation currently make macOS arm64
the only usable release path. Windows and Linux users need native artifacts
that preserve the same lightweight system-WebView design and tray lifecycle.

## Proposal

Provide a release-ready baseline on all three desktop operating systems:
retain the macOS `.app`, produce a resource-bearing Windows `.exe`, and produce
a Linux AppImage with freedesktop desktop metadata. Keep `wails3 package` as the
single public command and dispatch it to the host-specific implementation.
Installer matrices (NSIS, deb, rpm, and Arch packages) are explicitly deferred.

## Scope

### Included

- Windows amd64/arm64 production `.exe` build with application icon, UTF-8
  manifest, stable version metadata, and GUI subsystem flag.
- Linux native production build and AppImage with application icon and
  `.desktop` metadata.
- Linux hide-on-close lifecycle and stable desktop program identity.
- Host-neutral development/test tasks and platform-specific package dispatch.
- Static tests protecting package identity and metadata.
- Build, runtime dependency, persistence, and platform caveat documentation.

### Excluded

- Code signing, notarization, Windows Authenticode, and Linux repository signing.
- Windows NSIS/MSIX installers.
- Linux deb/rpm/Arch packages.
- Cross-compiling Linux CGO binaries or claiming runtime verification on a host
  where the app was not launched.
- Changes to DeepSeek navigation or bridge security policy.

## Functional Requirements

1. Closing the sole window or pressing Escape hides it without ending the
   process on macOS, Windows, and Linux; explicit Quit ends the process.
2. The colour tray icon is used on Windows/Linux and the template icon remains
   on macOS.
3. Linux sets both `DisableQuitOnLastWindowClosed` and `ProgramName`.
4. `wails3 task build` produces the correct native executable name on each host.
5. `wails3 package` dispatches by host OS:
   - macOS: `bin/TinglyDS.app`
   - Windows: `bin/tingly-ds.exe`
   - Linux: architecture-named AppImage in `bin/`
6. Windows release builds generate and consume an architecture-specific `.syso`
   resource, then remove that generated intermediate.
7. Linux packaging generates `build/linux/tingly-ds.desktop` and packages the
   binary and `build/icons/deepseek-dock.png` with Wails' AppImage generator.
8. Package identity remains `dev.tingly.tingly-ds` and executable/program name
   remains `tingly-ds`.

## Build Interfaces

| Task | Host | Output |
| --- | --- | --- |
| `build` | all | Native debug/release executable under `bin/` |
| `package` | all | Dispatch to the matching package task |
| `package:darwin` | macOS | Ad-hoc signed `.app` |
| `package:windows` | Windows | Production GUI `.exe` with Win32 resources |
| `package:linux` | Linux | AppImage |
| `generate:windows:resources` | Windows | Temporary `wails_windows_<arch>.syso` |
| `generate:linux:desktop` | Linux | `build/linux/tingly-ds.desktop` |

`ARCH` defaults to the host architecture and may be provided explicitly where
Wails resource generation supports it. Native package tasks remain restricted
to their matching host because Wails Linux builds require CGO/WebKitGTK and
Windows release verification requires WebView2.

## Metadata

### Windows

- `build/windows/wails.exe.manifest`: Windows 10/11 compatibility, per-monitor
  DPI awareness, long path awareness, UTF-8 active code page, and `asInvoker`.
- `build/windows/info.json`: company, product, description, copyright, and
  semantic version fields used by `wails3 generate syso`.
- `build/windows/icon.ico`: generated from the existing tiled app artwork and
  checked into the repository for deterministic native builds.

### Linux

- Desktop entry name: `Tingly DS`
- Exec/program name: `tingly-ds`
- Icon name: `tingly-ds`
- Category: `Network;InstantMessaging;`
- Startup notification enabled; terminal disabled.

## Error Handling

- Task commands fail immediately when icon/resource generation, compilation, or
  packaging fails.
- Platform-restricted package tasks cannot run accidentally on another host.
- Generated Windows `.syso` is removed after successful build and ignored by
  Git in case an interrupted build leaves it behind.
- Runtime URL/browser errors continue to be logged without widening bridge
  privileges.

## Security Considerations

- Windows manifest requests standard user privileges (`asInvoker`), never
  elevation.
- No additional Go service or page-facing bridge is introduced.
- AppImage packaging downloads Wails' configured Linux tooling; release builds
  should run in a controlled CI environment and record artifact hashes.
- Unsigned outputs are documented as local/testing artifacts, not trusted public
  releases.

## Testing and Validation

### Host-independent

- `gofmt`, `go vet ./...`, and `go test ./...` pass.
- Tests verify stable identities and required/non-privileged Windows manifest
  settings plus Linux desktop entry fields.
- Task listing exposes all platform package tasks.

### Native target smoke tests

For each target, launch the produced artifact and verify:

1. window opens and DeepSeek login works;
2. close/Escape hides without quitting;
3. tray click and menu actions show/hide/reload/open browser/quit;
4. login persists after quit, relaunch, and artifact replacement;
5. external links open in the default browser while DeepSeek links stay in-app;
6. upload, download, microphone, camera, clipboard, and notification behavior is
   checked against the native WebView policy.

A macOS-only development session may report static and macOS package results,
but must not claim Windows/Linux runtime validation.

## Dependencies

- Windows: Go 1.25+, matching Wails CLI, WebView2 runtime. Resource generation
  uses only the Wails CLI; no installer tooling is required.
- Linux: Go 1.25+, matching Wails CLI, C compiler, GTK3, WebKitGTK, and the
  system-tray/AppIndicator development libraries required by Wails. AppImage
  generation also needs network access and common archive/FUSE-compatible tools.
- macOS: existing Xcode Command Line Tools workflow.

## Implementation Phases

1. Add runtime Linux options.
2. Add Windows/Linux metadata and generated Windows icon.
3. refactor Taskfile into platform-neutral build/test and platform package tasks.
4. Extend metadata tests and `.gitignore`.
5. Rewrite README platform requirements, build/package instructions, persistence
   paths, artifact descriptions, and verification limitations.
6. Run static checks and native macOS packaging; document unexecuted target
   smoke tests.

## Acceptance Criteria

- All host-independent checks pass on the development host.
- macOS behavior and packaging remain regression-free.
- Windows/Linux task definitions and metadata are test-covered.
- Native CI or maintainers can generate the specified artifacts using only the
  documented dependencies and commands.
- Documentation distinguishes implemented packaging paths from target-native
  runtime verification and unsigned distribution status.

## Alternatives Considered

- **Full installer matrix now**: deferred to avoid adding NSIS and nfpm scripts
  before baseline native artifacts are exercised.
- **Runtime support only**: rejected because users need identifiable launchable
  artifacts rather than raw binaries with missing desktop integration.
- **Bundle Chromium**: rejected; it conflicts with the project's lightweight
  system-WebView architecture.
