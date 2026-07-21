# Dock and Persistence Security Review

**Date**: 2026-07-21
**Verdict**: PASS

- No application-owned credential or token store was introduced.
- Authentication data remains in the operating-system WebKit website-data
  store and is scoped by the packaged application identity.
- No cookie values, storage paths, or account data are logged or exposed to Go.
- Dock lifecycle changes add no page-callable native API and do not broaden the
  existing URL/message bridge.
- Documentation explains that clearing WebKit data removes the login session
  without recommending unsafe copying of browser-profile files.
