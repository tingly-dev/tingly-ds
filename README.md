# Tingly DS

Tingly DS is a small Go/Wails 3 desktop wrapper for
`https://chat.deepseek.com/`. Its Go module is
`github.com/tingly-dev/tingly-ds`. It uses the operating system WebView rather
than bundling Chromium and behaves like a normal Dock application while also
providing a menu-bar/system-tray entry.

## Behaviour

- The window opens on launch and the app remains available from both the Dock
  and the menu bar.
- Clicking the Dock icon restores the window after it has been hidden.
- Left-click the DeepSeek icon in the macOS menu bar to show or hide the window.
- Right-click the tray entry for Show, Hide, Reload, Open in Browser, and Quit.
- Closing the window or pressing Escape hides it without ending the process.
- DeepSeek links stay embedded; unrelated HTTP(S) anchor links open in the
  default browser.

## User data and login persistence

Tingly DS uses WebKit's normal persistent website-data store. DeepSeek cookies,
local storage, and cache live in the user's macOS Library, outside the `.app`
bundle and outside this source tree. Rebuilding or replacing `TinglyDS.app`
therefore does not reset the login session.

The packaged identity must remain stable for WebKit to find the same profile:

- bundle identifier: `dev.tingly.tingly-ds`
- executable: `tingly-ds`

Use `bin/TinglyDS.app` (or copy it to `/Applications`) for normal use. A binary
started directly with `wails3 task run` is a development process and may use a
separate WebKit profile. Changing the bundle identifier/executable name, or
manually clearing WebKit website data, can also require signing in again.

## Requirements

- macOS 12 or newer for the provided package metadata.
- Xcode Command Line Tools (`xcode-select --install`) for the compiler,
  `iconutil`, and `codesign`.
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

## Build from source

Clone the repository and install the exact Wails CLI version used by the Go
module:

```sh
git clone https://github.com/tingly-dev/tingly-ds.git
cd tingly-ds
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0-alpha2.114
go mod download
```

Run the checks and create the application bundle:

```sh
wails3 task test
wails3 package
open bin/TinglyDS.app
```

For a stable daily-use location, quit Tingly DS and copy the complete
`bin/TinglyDS.app` bundle into `/Applications`. Future builds can replace that
bundle without deleting the WebKit login data stored in the user Library.

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

The ad-hoc signed bundle is written to `bin/TinglyDS.app`. The application has
both a Dock icon and a menu-bar icon. Ad-hoc signing is suitable for local use.
Distribution to other Macs still requires a
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
