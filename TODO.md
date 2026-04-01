# TODO

The framework backlog has moved to the wiki.

Primary page:

- https://git.scuroneko.dev/ScuroNeko/Laniakea/wiki/Framework-Backlog

Russian page:

- https://git.scuroneko.dev/ScuroNeko/Laniakea/wiki/Framework-Backlog-RU

Current priority split:

- `Partial`: none.
- `Ideas`: service layer and dependency graph model, plugin composition contract.

Completed former high-priority items:

- `[v1.0.0-rc.14] Webhook runtime model.`
- `[v1.0.0-rc.13] Observability model`: added first-class `Observer` events for update, command, payload, scene, policy, runner, polling, and centralized error flows, with safe event dispatch and regression coverage for the new runtime hooks.
- `[v1.0.0-rc.13] Authorization and policy model`: added first-class `Policy[T]`, middleware integration through `RequirePolicy(...)`, plugin and bot policy registration helpers, built-in Telegram-aware policies, and composable `AllPolicies(...)`, `AnyPolicy(...)`, and `NotPolicy(...)` helpers with regression coverage.
- `[v1.0.0-rc.13] Update schema contract`: documented and tested the normalized `MsgContext` update-routing contract, including routing categories and per-update field guarantees.
- `[v1.0.0-rc.13] User-facing vs internal error model`: added explicit user-visible vs internal-only error markers and updated centralized handler error routing accordingly.
- `[v1.0.0-rc.13] Configuration freeze model`: formalized bot configuration freeze after first run, documented lifecycle commit points, and added regression coverage for ignored late mutations.
- `[v1.0.0-rc.12] Conversation / Scene Model`.
- `[v1.0.0-rc.12] Typed Handler Input Model`.
- `[v1.0.0-rc.12] Request Context / Cancellation Model`.
