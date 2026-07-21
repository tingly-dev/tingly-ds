# Project Understanding

The repository began as an uncommitted Go/Wails 3 application on branch
`main` and was imported in focused source and test commits. The implementation
now lives directly at the repository root; `bin/` contains generated artifacts
and is ignored. No Rust source is present. Earlier Tauri specification and
research files are retained as historical alternatives rather than active
architecture.

The architecture is intentionally narrow: Wails owns lifecycle and system
integration, while DeepSeek's existing web application owns the chat UI,
authentication, and persisted WebView state. The project contains one native
entry point, a pure navigation-policy module, embedded artwork, focused tests,
and macOS packaging metadata.

The principal maintenance boundaries are the pinned Wails alpha API, changes
to DeepSeek's remote frontend, native permission prompts, and platform-specific
remote-URL bridge behaviour. macOS arm64 is the only verified release target.
