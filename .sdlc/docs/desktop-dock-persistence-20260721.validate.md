# Dock and Persistence Validation

**Date**: 2026-07-21
**Verdict**: PASS

- The app uses Wails' regular macOS activation policy.
- Package metadata no longer opts into agent-only (`LSUIElement`) behavior.
- Wails' built-in macOS reopen handler restores hidden windows from the Dock.
- WebKit's default persistent website-data store remains enabled.
- Rebuild-sensitive identity fields are stable and regression-tested.
- README contains clone, dependency, test, package, launch, and installation
  instructions plus the development/profile distinction.
