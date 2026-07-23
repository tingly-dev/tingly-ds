# Windows and Linux Support Validation Report

**Date**: 2026-07-23
**Specification**: `.sdlc/docs/desktop-windows-linux-support-20260723.spec.md`
**Overall Status**: PARTIAL

## Criteria Results

| Criterion | Status | Evidence |
| --- | --- | --- |
| Linux hide-on-close lifecycle | Passed (static) | `application.LinuxOptions.DisableQuitOnLastWindowClosed` is enabled |
| Stable Linux program identity | Passed (static) | `ProgramName` and desktop `Exec` are `tingly-ds` |
| Native platform build dispatch | Passed | Taskfile parsed and macOS dispatch executed |
| macOS `.app` remains valid | Passed | package, strict signature verification, plist lint |
| Windows resource-bearing GUI executable | Passed (build) | amd64 and arm64 PE GUI cross-builds succeeded with `.syso` |
| Linux AppImage path | Passed (static) | native-only task uses absolute inputs and a correctly named icon |
| Stable package identity tests | Passed | macOS, Windows, and Linux metadata tests passed |
| No Windows elevation | Passed | manifest and test enforce `asInvoker` and reject administrator level |
| Documentation covers all targets | Passed | dependencies, artifacts, persistence, security, and smoke tests documented |
| Windows WebView2 runtime flow | Not run | Windows host unavailable |
| Linux WebKitGTK/AppImage runtime flow | Not run | Linux host unavailable |

## Findings

No blocking static, test, or build defect remains. The task cannot be marked
fully validated because native Windows/Linux UI behavior is part of the
acceptance criteria and cannot be observed on the macOS host.

## Required Native Follow-up

On Windows and Linux, run `wails3 task test`, `wails3 package`, launch the output,
and complete the README native release checklist. In particular, verify tray
support, hide-on-close, login persistence, external-link routing, and WebView
permissions.

## Conclusion

The implementation satisfies the code, metadata, build-definition, and
host-feasible artifact criteria. Release claims for Windows and Linux should be
held until their target-native smoke tests pass.
