# Dock and Persistence Test Report

**Date**: 2026-07-21
**Result**: PASS

## Automated checks

- `wails3 task test`: PASS (`gofmt`, `go vet ./...`, `go test ./...`).
- `wails3 build`: PASS with production tags.
- `packaging_test.go` guards the stable macOS bundle identifier, executable
  name, and absence of `LSUIElement`.

## Packaging checks

- `wails3 package`: PASS.
- `plutil -lint bin/TinglyDS.app/Contents/Info.plist`: PASS.
- `codesign --verify --deep --strict bin/TinglyDS.app`: PASS.
- Generated bundle size: 8.1 MB on macOS arm64.
- Final metadata contains `dev.tingly.tingly-ds` / `tingly-ds` and omits
  `LSUIElement`, allowing the regular activation policy to expose the Dock.

## Manual follow-up

The generated bundle is ready for interactive confirmation of Dock click,
tray toggle, and login retention after replacing it with another build.
