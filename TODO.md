# TODO

The framework backlog has moved to the wiki.

Primary page:

- https://git.scuroneko.dev/ScuroNeko/Laniakea/wiki/Framework-Backlog

Russian page:

- https://git.scuroneko.dev/ScuroNeko/Laniakea/wiki/Framework-Backlog-RU

Current priority split:

- `Priority 1`: webhook runtime model, authorization and policy model, observability model.
- `Priority 2`: service layer and dependency graph model, plugin composition contract.

Completed former high-priority items:

- `[v1.0.0-rc.13] Update schema contract`: documented and tested the normalized `MsgContext` update-routing contract, including routing categories and per-update field guarantees.
- `[v1.0.0-rc.13] User-facing vs internal error model`: added explicit user-visible vs internal-only error markers and updated centralized handler error routing accordingly.
- `[v1.0.0-rc.13] Configuration freeze model`: formalized bot configuration freeze after first run, documented lifecycle commit points, and added regression coverage for ignored late mutations.
- `1. Conversation / Scene Model`: completed in `v1.0.0-rc.12`.
- `2. Typed Handler Input Model`: completed in `v1.0.0-rc.12`.
- `3. Request Context / Cancellation Model`: completed in `v1.0.0-rc.12`.
