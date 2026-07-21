# Tingly DS — Test Report

**Date**: 2026-07-21
**Status**: PASS
**Target**: repository-root Go/Wails 3 module

## Automated results

| Check | Result | Evidence |
| --- | --- | --- |
| Formatting | PASS | `gofmt` completed through `wails3 task test` |
| Static analysis | PASS | `go vet ./...` |
| Unit tests | PASS | `go test ./...` |
| Race detector | PASS | `go test -race ./...` |
| Coverage | PASS | 40.5% overall; GUI entry point is not unit-driven |
| Navigation policy | PASS | `externalURLFromMessage` 96.2%; `isDeepSeekURL` 100% |
| Embedded icon integrity | PASS | Both 640×640 tray PNGs decode in `assets_test.go` |
| Production build | PASS | arm64 Mach-O built with the Wails `production` tag |
| macOS package | PASS | `wails3 package` generated `bin/TinglyDS.app` |
| Bundle metadata | PASS | `plutil -lint`; `CFBundleIconFile=AppIcon.icns` |
| Signature | PASS | `codesign --verify --deep --strict` |

All commands above were rerun successfully from the repository root after the
project was moved out of its temporary subdirectory.

The final bundle is 8.1 MiB. Its executable is 7.8 MiB and the generated
multi-resolution application icon is 288 KiB.

## Policy cases covered

- Exact DeepSeek host, root host, nested subdomain, TLS port, and trailing dot.
- Lookalike prefixes/suffixes, HTTP sender, non-TLS port, credentials, relative
  URL, and malformed URL.
- Wrong window, nil origin, subframe, lookalike origin, malformed JSON, multiple
  JSON values, unknown field/action, unsupported scheme, credentials, and an
  oversized bridge message.

## Manual result

The user launched and used the application on 2026-07-21 and reported no
problems. The final icon-bearing bundle was also launched successfully after
packaging. Website-specific flows should be smoke-tested again after material
DeepSeek frontend or WebKit changes.

Windows and Linux were not execution-tested and are not release targets for
this iteration. In particular, the pinned Wails alpha has platform differences
around JavaScript injection for remote URLs and raw-message origin metadata.
