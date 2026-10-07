#!/usr/bin/env bash
# Generate GitHub Release body text for a given tag.
#
# Usage: release-notes.sh <tag>
#
# The body lists each downloadable asset with platform-specific usage guidance.
# The macOS section is explicit about Gatekeeper/quarantine because the bundle
# is ad-hoc signed and not Developer-ID notarised.
set -euo pipefail

tag="${1:?tag argument is required, e.g. v0.1.0}"

cat <<EOF
# Tingly DS ${tag}

A lightweight Go/Wails 3 desktop shell for https://chat.deepseek.com/. This
release ships prebuilt binaries for macOS, Windows, and Linux. They are
unsigned release baselines; see the per-platform notes below.

## Downloads

| Platform | File | Architecture |
| --- | --- | --- |
| macOS | \`TinglyDS-macOS-arm64.zip\` | Apple Silicon (arm64) |
| Windows | \`TinglyDS-Windows-x86_64.exe\` | x86_64 |
| Linux | \`TinglyDS-Linux-x86_64.AppImage\` | x86_64 |

> The macOS build is Apple Silicon (arm64) only. Intel (x86_64) Mac users can
> build from source (see below) by running \`wails3 task build\` on the Intel
> host.

## macOS

1. Download and unzip \`TinglyDS-macOS-arm64.zip\`. You get
   \`TinglyDS.app\`.
2. Move \`TinglyDS.app\` into \`/Applications\`.
3. The first launch will be **blocked by Gatekeeper** because the app is ad-hoc
   signed, not Developer-ID signed or notarised. To allow it:
   - **Option A (simplest):** Right-click (or Control-click) \`TinglyDS.app\`
     and choose **Open**. In the dialog, click **Open** again. macOS remembers
     the choice afterwards and a normal double-click works.
   - **Option B (terminal):** Remove the quarantine attribute once:
     \`\`\`sh
     xattr -dr com.apple.quarantine /Applications/TinglyDS.app
     \`\`\`
   - **Option C (Settings):** Attempt to open the app, then go to
     **System Settings → Privacy & Security** and click **Open Anyway** under
     the "TinglyDS was blocked" notice.
4. Requires **macOS 12 (Monterey) or newer**.

The tray/menu-bar icon may also need **System Settings → Privacy & Security →
Accessibility** if it does not appear.

## Windows

1. Download \`TinglyDS-Windows-x86_64.exe\`.
2. SmartScreen may warn "Windows protected your PC" because the executable is
   unsigned. Click **More info → Run anyway**.
3. Requires **Windows 10 or 11** with the **Microsoft Edge WebView2 Runtime**.
   Most current installs include the evergreen runtime; if not, install it from
   Microsoft: https://developer.microsoft.com/microsoft-edge/webview2/
4. The app runs with standard user privileges; no administrator rights are
   requested.

## Linux

1. Download \`TinglyDS-Linux-x86_64.AppImage\`.
2. Make it executable and run it:
   \`\`\`sh
   chmod +x TinglyDS-Linux-x86_64.AppImage
   ./TinglyDS-Linux-x86_64.AppImage
   \`\`\`
3. Requires **WebKitGTK** and a desktop environment with status-notifier /
   AppIndicator support. On some GNOME setups you also need an AppIndicator
   extension for the tray icon.
4. If AppImage mounting fails, install FUSE (e.g. \`libfuse2\` on Ubuntu/Debian)
   or run with \`--appimage-extract-and-run\`.

## User data and login

Each platform WebView persists its own data store (cookies, local storage,
cache) outside the application artifact, so DeepSeek login survives replacing
the app. The packaged identity (\`dev.tingly.tingly-ds\`) must stay stable for
the WebView to find the same profile.

## Build from source

See the [README](https://github.com/tingly-dev/tingly-ds#readme) for full build
instructions. Quick start:

\`\`\`sh
git clone https://github.com/tingly-dev/tingly-ds.git
cd tingly-ds
git checkout ${tag}
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0-beta.28
wails3 task build
\`\`\`
EOF
