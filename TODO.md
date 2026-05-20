# TODO

## v1.0.0 pre-release review

Findings from the full-repo review against `AGENTS.md` priorities. Build, vet, tests, and lint are clean; items below are public-API and godoc hygiene before the stable tag.

### Major — close before 1.0.0 tag

- [X] **M1. `BotPayloadType*` are `var`, must be `const`** — `bot.go:50-59`. Public sentinels are user-mutable globals. `KeyboardButtonStyle*` in `keyboard.go:10-17` already uses `const`; match the pattern.
- [X] **M2. `Observer` method naming asymmetry** — `observer.go:147-157`. `OnReceiveUpdate` → `OnUpdateReceived`; `OnHandledUpdate` → `OnUpdateHandled` to match `UpdateReceivedEvent` / `UpdateHandledEvent` and the rest of the `OnX` pattern. Breaking after 1.0.
- [X] **M3. Uploader returns ad-hoc error string instead of `*ResponseError`** — `tgapi/uploader_api.go:183`. `tgapi/api.go:258-292` returns `*ResponseError`; uploader must do the same so `errors.As(err, &tgapi.ResponseError{})` works for upload paths too.
- [X] **M4. `BotOptsFileJSON` is missing `PollTimeout`** — `bot_opts_loader.go:35-46`, plus `FromBytes`/`ToBytes` mapping. File round-trip silently drops `PollTimeout`.
- [X] **M5. Stale `Bot.Updates` godoc** — `methods.go:11-44`. Claims "30-second timeout" and "empty slice if none"; in reality timeout is `bot.pollTimeout` and the function returns `nil` on error.
- [X] **M6. Self-contradicting `NewRandomDraftProvider` godoc** — `drafts.go:50-59`. Says "cryptographically secure random numbers" but uses `math/rand/v2` (the underlying generator type correctly notes it is not crypto-secure).
- [X] **M7. `Draft.Delete` godoc says "internal method"** — `drafts.go:190-201`. Method is exported; either rewrite the godoc with a public-intent description or unexport.
- [X] **M8. Russian comments in production code**
  - `msg_handler.go:28` — "Ищем команду по точному совпадению"
  - `tgapi/uploader_api.go:181` — "Повторяем запрос"
- [X] **M9. `MessageContext.Error` godoc references unexported helper** — `msg_context.go:540`. "Error is an alias for error()" — rewrite to describe the centralized handler error path and `IsUserError` gating.
- [X] **M10. `Scene` and `SceneSession` mix exported fields with setters**
  - `Scene` exports `Name/Scope/Entry/PluginName` and also has `SetScope/SetEntry`; `PluginName` is framework-assigned but publicly mutable.
  - `SceneSession` exports `Data []byte` and also has `Set/Get/HasData/ClearData/BindData/SaveData`.
  - Pick one model per type before 1.0.0.
- [X] **M11. Constant-time compare for webhook secret** — `bot_webhook.go:296` (update handler) and `bot_webhook.go:341` (`/status`). Use `subtle.ConstantTimeCompare`.

### Minor — can slip to 1.0.x

- [X] Strip `// Internal helper …` godoc from unexported funcs (~23 occurrences in repo); `AGENTS.md` explicitly forbids godoc-style comments on unexported declarations without a strong reason.
- [X] `Plugin.AddCommand` godoc references unexported field `.command` — `plugins.go:48-49`.
- [X] `Runner` builder naming: `runner.Once(true)`, `runner.Async(true)` read awkwardly; consider `SetOnce`/`SetAsync` to match `Set*` on other types, or zero-arg `Once()` + paired `Repeat(every)`.
- [X] Typo in webhook error string: `bot_webhook.go:143` — "MaxConnections must between 1 and 100" (missing `be`).
- [X] `RunWebhookWithContext` uses inline `errors.New(...)` instead of `Err*` sentinels (`bot_webhook.go:131-156`); rest of the package uses sentinels from `errors.go`.
- [X] `tgapi.UpdateTypeManagedBot` (`tgapi/types.go:61`) has no godoc.
- [X] `Bot.GetAPI`, `Bot.GetUploader`, `InlineKeyboard.GetMaxRow` have no godoc.
- [X] `Bot.L10n` godoc says "Returns empty string if translation not found"; actually returns the key (`l10n.go:48-59`).
- [X] `Bot.handle` panic recovery only logs — emit `ErrorEvent` so observers see panics (`handler.go:18-23`).
- [X] `handleCallback` vs `handleMessage` differ in plugin-logger assignment: callback assigns unconditionally then falls back to bot logger (`msg_handler.go:209-212`); message only assigns if non-nil (`msg_handler.go:35-37`). Align.
- [X] `SetCallbackData` godoc says "default payload type is JSON" — actually the zero `BotPayloadType` falls through to the `default` branch (which happens to be JSON). Either document the zero-value behavior explicitly or initialize the builder with the bot's default (`keyboard.go:106-122`).
- [X] `commands.go:62-66` — empty `case CommandValueAny:` next to `default: regex = nil` looks like an incomplete switch. Merge or add a one-line comment.
- [X] `Bot.SetDebug` does not call `configMutable` unlike sibling setters; if intentional, note it in godoc.

### Tests to add after the fixes

- `BotOptsFileJSON` round-trip for `PollTimeout` (after M4).
- Uploader 4xx/429 surfaces `*tgapi.ResponseError` (after M3).
- `Bot.handle` panic → observer receives `ErrorEvent` (after panic-recovery fix).
- Webhook `/status` with wrong `SecretToken` returns 403 / `403`-equivalent (after M11), incl. a constant-time-compare smoke.
- Table-driven `parseCommand` cases for `/cmd@botname` and stripping behavior.

---

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
