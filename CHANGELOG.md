# Changelog

## v1.0.0-rc.12

### Added
- `AnswerLong(...)`, `AnswerLongf(...)`, `KeyboardLong(...)`, and `SplitMessageText(...)` for explicit plain-text splitting of long replies without changing the semantics of existing single-message helpers.
- Centralized library-level validation errors in `errors.go`, including `ErrEmptyMessage`, `ErrMessageTooLong`, `ErrCaptionTooLong`, and context/target validation sentinels.
- `Bot.GetPayloadType()`, `InlineKeyboard.GetPayloadType()`, and optional strict payload decoding via `BotOpts.StrictPayloadType` / `Bot.SetStrictPayloadType(...)`.
- `MsgContext.BindArgs(...)` for binding positional command arguments into exported struct fields.
- Binding sentinels `ErrBindArgsTargetNotPointer`, `ErrBindArgsTargetNotStruct`, `ErrBindArgsUnsupportedFieldType`, and `ErrBindArgsConversion`.

### Changed
- `CommandExecutor` now returns `error`, and command, payload, and non-command update handlers now use centralized bot error handling for returned errors.
- README and README_RU examples now use the new handler signature and document the long-message helpers.
- README and README_RU now link to the project wiki, and the wiki now includes a page-priority tracker while content is being filled in.
- `AGENTS.md` now requires every change to be recorded in `CHANGELOG.md`, enforces version alignment with `utils/version.go`, and blocks breaking changes without a major-version bump.
- `AGENTS.md` now also defines a short commit-message format: one summary line plus up to three high-signal detail lines.
- `AGENTS.md` now explicitly requires each commit-message detail line to be placed on its own new line.
- `AGENTS.md` now also requires commit messages to be emitted as a plain multiline block instead of collapsed prose or list formatting.
- `AGENTS.md` now requires new or expanded project documentation to be maintained in both English and Russian whenever reasonably possible.
- `AGENTS.md` now requires all agent-created commits to be GPG-signed and to fail fast instead of falling back to unsigned commits when signing cannot be completed.
- `AGENTS.md` now also links the wiki backlog flow more tightly to `TODO.md` and `CHANGELOG.md`, requiring draft-wiki confirmation for large new ideas and synchronized completion records for backlog items.
- Added `TODO.md` to track missing framework-level concepts, with detailed notes for scenes, typed handler input, and request-scoped cancellation.
- Payload-type comments and docs now distinguish between the bot's default payload type and keyboard-local overrides.
- `MsgContext.Context()` now safely falls back to `context.Background()` when no request-scoped context is attached.
- `MsgContext` reply, edit, callback, delete, action, and draft-limiter paths now use the context accessor instead of reaching into raw internal state.
- `TODO.md` is now a short pointer file, while the detailed framework backlog lives in the wiki as `Framework-Backlog` / `Framework-Backlog-RU`.
- Version constants were bumped to `v1.0.0-rc.12`.

### Fixed
- Message and caption validation now runs before Telegram API calls, rejecting empty messages, oversized message text, and oversized captions with stable sentinel errors.
- Draft flushing and draft updates now reject oversized messages before sending invalid requests.
- Callback payload decoding now optionally enforces strict type matching, while the default tolerant mode logs Base64-to-JSON decoding in debug mode and still accepts keyboard-local payload overrides.
- Positional argument binding now leaves missing trailing struct fields at zero values, joins the remaining arguments into the final string field, and returns clearer binding errors.
- Request-scoped contexts are now created per update handler execution and safely reused through `MsgContext.Context()` even for manually constructed test contexts.
- Command and payload handlers now have regression coverage for end-to-end typed argument binding through the normal routing path.

### Breaking Changes
- `CommandExecutor[T]` changed from `func(ctx *MsgContext, db T)` to `func(ctx *MsgContext, db T) error`.
- `Plugin.NewCommand(...)`, `Plugin.NewPayload(...)`, and `Plugin.AddUpdateHandler(...)` now require handlers with the new error-returning signature.

### Tests
- Added regression tests for `MsgContext.BindArgs(...)`, including scalar conversion, tail-string binding, zero-value trailing fields, invalid targets, unsupported field types, and end-to-end command/payload binding.

## v1.0.0-rc.11

### Fixed
- `chat_boost` update decoding now accepts string `boost_id` values, matching the current Telegram Bot API schema and preventing polling failures on boosted-chat updates.

## v1.0.0-rc.10

### Added
- `Plugin.AddUpdateHandler` for routing non-command Telegram updates by `tgapi.UpdateType`.
- Derived `tgapi.Update.Type` assignment during JSON decoding, plus `tgapi.UpdateTypeUnknown` for unmatched payloads.
- `tgapi.API.OpenFileByLink(...)` and `OpenFileByLinkWithContext(...)` for streaming downloads from Telegram's file server.
- Regression tests for update dispatch, keyboard builders, localization fallback, runners, rate limiting, parse mode encoding, streaming downloads, and context isolation.
- Regression tests for bot single-run enforcement, nil plugin registration, `L10n` concurrent access, `API.Close()` idle-connection cleanup, and `tgapi` worker-pool edge cases.
- `SEMVER.md` documenting versioning expectations for the project.

### Changed
- `NewBot` now returns `(*Bot[T], error)` instead of terminating the host process on configuration or startup failures.
- `Run` and `RunWithContext` now return errors; `RunWithContext` returns `ErrNoPrefixes` and `ErrNoPlugins` for invalid bot configuration.
- Polling retries now use exponential backoff instead of busy-looping on repeated `getUpdates` failures.
- `Bot` is now explicitly single-use; repeated `Run()` or `RunWithContext(...)` calls return `ErrBotAlreadyRun`.
- Database context wiring now uses `T` consistently instead of forcing `*T`; shared dependencies should typically use pointer types such as `*sql.DB`.
- `DatabaseContext`, `GetDBContext`, and `DbLogger` were updated to the new `T`-based dependency model.
- `DatabaseContext(...)` now warns once when `T` is a value type, to highlight likely unintended copying of shared dependencies.
- `AddDatabaseLoggerWriter(...)` now skips unset and nil database contexts instead of calling the writer with invalid values.
- `L10n` is now safe for concurrent use and copies added dictionary entries to avoid external mutation after registration.
- Plugin registration now snapshots commands, payloads, middlewares, and update handlers so later mutations of the original `*Plugin` do not leak into the bot.
- `AddPlugins(...)` now skips nil plugin pointers instead of panicking.
- `GetUpdateTypes()` now returns a copy instead of exposing internal slice state.
- Update handling now normalizes `MsgContext` for more Telegram update kinds and routes plugin-level update handlers with isolated context copies.
- `message`, `channel_post`, and `callback_query` remain on the command/payload flow; non-command updates can be handled through plugin update handlers.
- Command auto-generation now validates Telegram command names with the correct character set and `1..32` length limit, and emits commands in deterministic sorted order.
- Builder-style APIs were normalized to value returns for `NewCommandArg`, `NewMiddleware`, `NewRunner`, and `NewCallbackData`.
- `MenuButton` replaced `BaseMenuButton`, and `GetChatMenuButton(...)` now returns the renamed type.
- Several Telegram DTOs were tightened for optionality and serialization correctness, including `InputPaidMedia`, `MenuButton`, optional gift fields, and message entity slices.
- `tgapi.NewRequest(...)`, `NewRequestWithChatID(...)`, `NewUploaderRequest(...)`, and `NewUploaderRequestWithChatID(...)` are now documented as low-level unsafe escape hatches rather than internal helpers.
- `tgapi.API.Close()` now closes idle HTTP connections before releasing logger resources.
- Multipart form encoding now writes scalar field bytes directly instead of converting through temporary strings.
- README, README_RU, package docs, and exported godoc were updated to match the current APIs and concurrency/lifecycle model.
- Version constants were bumped to `v1.0.0-rc.10`.

### Fixed
- Required command arguments are now enforced by declared argument index, not only by total required count.
- `ParseNone` now omits `parse_mode` from JSON requests instead of serializing `"None"`.
- Upload file type detection is now case-insensitive for file extensions.
- Draft creation no longer panics when no limiter is configured, and draft flushing now rejects zero chat IDs before sending invalid requests.
- Channel posts with `SenderChat` no longer panic in the command path and now preserve the expected `MsgContext` fields.
- File logger initialization now falls back to stdout loggers instead of terminating the process on logger setup failures.
- `GetChatMenuButton` and `SetChatMenuButton` now serialize `chat_id` correctly when omitted.
- Update decoding tests now match the canonical `deleted_business_messages` model and no longer rely on the removed singular alias.

### Breaking Changes
- `NewBot[T](opts)` now returns `(*Bot[T], error)`.
- `Run()` now returns `error`.
- `RunWithContext(ctx)` now returns `error`.
- `Run()` and `RunWithContext(ctx)` are now single-use per bot instance; create a new `Bot` after they return.
- Database context handlers now receive `T` instead of `*T`. For shared dependencies, instantiate the bot with a pointer type, for example `Bot[*sql.DB]`.
- `DatabaseContext(...)` now takes `T` instead of `*T`.
- `GetDBContext()` now returns `T` instead of `*T`.
- `DbLogger[T]` now receives `T` instead of `*T`.
- `NewCommandArg(...)`, `NewMiddleware(...)`, `NewRunner(...)`, and `NewCallbackData(...)` now return values instead of pointers.
- `BaseMenuButton` was renamed to `MenuButton`, and `GetChatMenuButton(...)` now returns `MenuButton`.
- `tgapi.Update` no longer exposes the deprecated `DeletedBusinessMessage` alias; use `DeletedBusinessMessages`.

### Tests
- Added coverage for polling backoff helpers, command sorting, database logger safety checks, update handler routing, update-context isolation, channel posts with `SenderChat`, parse mode encoding, streaming downloads, and rate limiter behavior.

## v1.0.0-rc.7

### Added
- Package-level logger helpers: `utils.CreateLogger(prefix, level)` and `utils.CreateFileLogger(prefix, level, filePath)`.
- `MsgContext.Logger`, populated from the matched plugin and falling back to the bot logger.
- Plugin lifecycle/configuration APIs: `SetLogger`, `RemoveLogger`, `SetOnClose`, and `Close`.
- `Bot.CloseRemote(ctx)` as the explicit wrapper for Telegram Bot API close.

### Changed
- Logger initialization is now unified across `Bot`, `tgapi.API`, and `tgapi.Uploader`.
- `Bot.Close()` now performs local resource teardown only and invokes `Plugin.Close()` for registered plugins.
- Local `tgapi.API` shutdown was renamed to `Close()`.
- Telegram Bot API close wrappers in `tgapi.API` were renamed to `CloseRemote()` and `CloseRemoteWithContext()`.
- `Bot.Debug()` now updates log levels for the bot logger, request logger, and already registered plugin loggers.
- `Bot.AddPlugins()` now creates a default plugin logger automatically when one is not provided.
- `Bot.AddDatabaseLoggerWriter()` now also attaches the writer to already registered plugin loggers.
- GoDoc was expanded for the new shutdown and logging APIs, and plugin registration is now documented as a configuration commit point.

### Breaking Changes
- `(*Bot).Close(ctx context.Context)` was replaced with `(*Bot).Close()`.
- `(*tgapi.API).CloseApi()` was renamed to `(*tgapi.API).Close()`.
- `(*tgapi.API).Close()` was renamed to `(*tgapi.API).CloseRemote()`.
- `(*tgapi.API).CloseWithContext()` was renamed to `(*tgapi.API).CloseRemoteWithContext(ctx)`.

### Migration
- Replace `bot.Close(ctx)` with `bot.Close()`.
- If you need Telegram Bot API close, use `bot.CloseRemote(ctx)`.
- Replace `api.CloseApi()` with `api.Close()`.
- Replace `api.Close()` with `api.CloseRemote()`.
- Replace `api.CloseWithContext(ctx)` with `api.CloseRemoteWithContext(ctx)`.
- Configure plugin loggers and `OnClose` hooks before calling `bot.AddPlugins(...)`.

### Tests
- Updated tests for the new shutdown and logging behavior.

### Notes
- Registering a plugin via `AddPlugins(...)` is a configuration commit point; the plugin should not be mutated through the original `*Plugin` afterward.
- If plugin loggers must receive a database writer, call `AddDatabaseLoggerWriter(...)` after registering plugins.

## v1.0.0-rc.4

### Added
- `WithContext` variants across `tgapi` API and uploader methods so callers can pass cancellation and deadline contexts consistently.
- `UploaderCertificateType`, `UploadSetWebhookP`, `Uploader.SetWebhook(...)`, and `Uploader.SetWebhookWithContext(...)` for multipart webhook certificate uploads.
- Missing media thumbnail fields where applicable.

### Changed
- GoDoc for context-aware methods was improved, and `See` references now point to method-specific Telegram Bot API anchors.
- `EditMessageTextP` now includes `entities` and `link_preview_options`.
- `EditMessageCaptionP` now includes `caption_entities` and `show_caption_above_media`.
- `StopPollP` now uses `reply_markup` and no longer carries `inline_message_id`.
- `SendStickerP` now includes reply and suggested-post related fields.
- `SendDocumentP` now includes `disable_content_type_detection`.
- `SendInvoiceP` no longer includes unsupported `business_connection_id`.
- `SetWebhookP` no longer carries `certificate`; GoDoc now points to uploader-based certificate upload.
- Existing non-context methods remain available, and the `Do(...)` call style is preserved.

### Breaking Changes
- Users sending webhook certificates through JSON `SetWebhookP.Certificate` must migrate to `Uploader.SetWebhook(...)`.

## v1.0.0-rc.3

### Fixed
- The update polling loop no longer logs or retries after `context.Canceled` during shutdown.
- Extra retry delay was removed from canceled polling requests so `RunWithContext` can exit immediately while stopping.

### Changed
- Shutdown behavior remains explicit: callers are still responsible for invoking `Close()` after `RunWithContext` returns.

## v1.0.0-rc.2

### Fixed
- Fixed a shutdown crash caused by `DatabaseWriter` calling `Close()` through an uninitialized embedded logger writer.
- Fixed bot shutdown hanging during Telegram long polling by making update polling use a cancelable context.
- Reduced the chance of container termination with exit code `137` during shutdown by allowing `getUpdates` to stop promptly on cancellation.

### Changed
- Switched the project to use the local `laniakea` replacement for the shutdown fix.
- Documentation now clarifies that `RunWithContext` does not close resources automatically and callers must invoke `Close()` explicitly.
- `Updates` documentation now describes context-driven cancellation behavior.

### Tests
- Added regression tests for database logger writer shutdown behavior.
