# Tingly DS — Security Report

**Date**: 2026-07-21
**Status**: PASS
**Critical**: 0 · **High**: 0 · **Medium**: 0

## Secrets

Local source and configuration scanning found no API keys, tokens, private
keys, passwords, or world-writable source files. The application stores no
DeepSeek credential itself; authentication state remains in the OS WebView.

## Dependencies

The first Go vulnerability scan found `GO-2026-5024` in the required
`golang.org/x/sys v0.43.0` module. It was Windows-only and not reachable from
this program, but the dependency was nevertheless upgraded to the fixed
`v0.44.0`. A second `govulncheck ./...` reported `No vulnerabilities found`.

Wails is pinned to `v3.0.0-alpha2.114`. Its prerelease maturity is a maintenance
risk rather than a known vulnerability and is called out in the README.

## Configuration and bridge review

- No Go service or business method is exposed to remote page JavaScript.
- The only native message accepts one strict JSON shape, has a 4096-byte limit,
  requires the main frame and an HTTPS DeepSeek origin, and accepts only
  credential-free HTTP(S) targets outside DeepSeek.
- Release DevTools and native file-drop bridging are disabled.
- Camera, microphone, location, notification, and clipboard permissions remain
  under operating-system default/prompt handling.
- External navigation is delegated to the system browser instead of loading
  unrelated origins into the privileged WebView.
- Icon assets have source URLs, SHA-256 digests, and the Lobe Icons MIT notice
  in `THIRD_PARTY_NOTICES.md`.

## Residual risks

- DeepSeek controls the remote page and may change login, navigation, or upload
  behaviour; retest after material site changes.
- The local package is ad-hoc signed. Do not redistribute it as a trusted public
  release without Developer ID signing and notarisation.
- DeepSeek marks remain subject to the owner's trademark rights; this project
  is a personal shell and does not claim affiliation.
- Windows/Linux remote-URL bridge behaviour is not security-approved by this
  macOS review and must be re-audited before either target is released.
