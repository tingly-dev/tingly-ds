# Windows and Linux Support Test Report

**Date**: 2026-07-23
**Status**: PASSED with target-native limitations
**Host**: macOS arm64

## Automated Checks

| Check | Result | Evidence |
| --- | --- | --- |
| Go formatting | Passed | `wails3 task test` ran `gofmt` |
| Go vet | Passed | `go vet ./...` exited 0 |
| Unit/metadata tests | Passed | `go test ./...` exited 0 |
| Taskfile parse/list | Passed | all public package tasks listed |
| Diff whitespace | Passed | `git diff --check` exited 0 |

The Go test linker emitted non-fatal macOS SDK-version warnings from Wails CGO
objects. Tests still completed successfully; these warnings are unrelated to the
Windows/Linux metadata changes.

## Build Checks

- macOS arm64 production binary: passed; identified as a Mach-O arm64 executable.
- macOS `.app`: passed; ad-hoc signature, strict code-sign verification, and
  `Info.plist` lint all passed.
- Windows amd64: passed cross-compilation with generated `.syso`; identified as
  a PE32+ GUI x86-64 executable.
- Windows arm64: passed cross-compilation with generated `.syso`; identified as
  a PE32+ GUI AArch64 executable.
- Embedded Windows strings confirmed `Tingly DS`,
  `dev.tingly.tingly-ds`, and `asInvoker` in the arm64 artifact.
- Linux: not compiled or packaged on macOS because Wails Linux requires native
  CGO, GTK, and WebKitGTK libraries. Desktop metadata and Taskfile definitions
  are covered by host-independent tests.

## Scope Not Executed

There is no automated UI/E2E harness in this repository. Windows WebView2 and
Linux WebKitGTK launch, tray behavior, login persistence, bridge behavior,
permissions, uploads, and downloads require the documented native smoke tests.

## Conclusion

Static checks and all builds feasible from the development host passed. Native
Windows and Linux runtime acceptance remains required before publishing their
artifacts.
