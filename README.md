# DeepSeek Shell — Wails 3

A small Go/Wails 3 desktop wrapper for `https://chat.deepseek.com/`. It uses the
operating system WebView rather than bundling Chromium and is controlled from a
menu-bar/system-tray entry.

## Behaviour

- The window opens on launch and keeps the normal WebView login session.
- Left-click the DeepSeek icon in the macOS menu bar to show or hide the window.
- Right-click the tray entry for Show, Hide, Reload, Open in Browser, and Quit.
- Closing the window or pressing Escape hides it without ending the process.
- DeepSeek links stay embedded; unrelated HTTP(S) anchor links open in the
  default browser.

## Requirements

- macOS 12 or newer for the provided package metadata.
- Go 1.25 or newer (the pinned Wails module raises the module directive to
  Go 1.25 during `go mod tidy`).
- Wails CLI matching the pinned module version:

  ```sh
  go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0-alpha2.114
  ```

Wails 3 is still prerelease software. Re-run the full test and package workflow
before changing its pinned version.

The currently implemented and verified release target is macOS arm64. The code
keeps later Windows/Linux ports straightforward, but those targets are not yet
claimed as supported; Wails alpha remote-URL injection and origin metadata must
be validated on each target before release.

## Develop

```sh
go mod download
wails3 task test
wails3 task run
```

For file watching:

```sh
wails3 dev -config ./build/config.yml
```

## Package for macOS

```sh
wails3 package
```

The ad-hoc signed bundle is written to `bin/DeepSeekShell.app`. Ad-hoc signing
is suitable for local use. Distribution to other Macs still requires a
Developer ID signature and Apple notarisation.

The app and tray marks are DeepSeek assets downloaded from
[LobeHub Icons](https://lobehub.com/icons/deepseek). Their exact source URLs,
checksums, and Lobe Icons licence notice are recorded in
`THIRD_PARTY_NOTICES.md`.

## Security model

- No Go service is registered for page JavaScript.
- Wails simple event emission, native file-drop bridging, and release DevTools
  are disabled.
- Web permissions use the operating-system prompt/default policy.
- The sole raw bridge message can only request that a validated HTTP(S) URL be
  opened in the system browser, and only from the main DeepSeek frame.

External-link interception is deliberately small and best-effort. DeepSeek can
change its frontend or login flow, so login, uploads, downloads, and external
links should be manually checked after major website or WebView updates.
