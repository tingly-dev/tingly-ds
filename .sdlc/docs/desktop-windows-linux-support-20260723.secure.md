# Windows and Linux Support Security Report

**Date**: 2026-07-23
**Status**: PASSED with dependency-scan limitation

## Secrets

A tracked-text credential keyword scan found only security documentation and the
intentional `user:secret@example.com` negative-test fixture. No application
secret, token, key, password, or signing credential was introduced.

## Permissions and Runtime Surface

- The Windows manifest requests `asInvoker` with `uiAccess=false`; it neither
  requests nor requires administrator privileges.
- Linux packaging installs no privileged service and introduces no install or
  removal scripts.
- No page-facing Go service, new raw bridge command, or additional WebView
  permission was introduced.
- Existing external URL validation and HTTPS-only DeepSeek entry URL are
  unchanged.

## Packaging Risks

- The Windows `.exe` and Linux AppImage baseline artifacts are unsigned. The
  README explicitly requires platform signing before public distribution.
- Wails AppImage generation downloads continuous-release tools from GitHub at
  build time. This is a supply-chain reproducibility risk; controlled release CI
  should pin/verify downloaded tool hashes before treating Linux artifacts as
  production releases.
- Generated Windows resources contain only public metadata and icon assets.
- Windows build intermediates and Linux AppImage work directories are ignored,
  reducing accidental artifact commits.

## Dependency Scan

`govulncheck` and `gitleaks` are not installed in the current environment, so no
CVE database-backed or dedicated entropy-based scan ran. `go vet`, unit tests,
and manual source/config review passed. The pinned Wails 3 prerelease dependency
remains an acknowledged risk and target-native testing is required.

## Conclusion

No new critical or high security issue was found. Public Linux release hardening
should make AppImage tooling reproducible, and CI should add `govulncheck` plus a
dedicated secret scanner.
