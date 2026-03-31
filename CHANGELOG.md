# Changelog

## v1.0.0-rc.14

### Changed
- Added missing godoc for the exported observer `Event` marker interface.

## v1.0.0-rc.13

### Added
- `AsUserError(...)`, `AsInternalError(...)`, `IsUserError(...)`, and `IsInternalError(...)` for explicitly marking centralized handler errors as user-visible or internal-only without breaking the existing default error flow.
- `Policy[T]`, `RequirePolicy(...)`, and built-in chat and callback policy helpers for expressing reusable authorization rules through the existing middleware pipeline.
- `Bot.UsePolicy(...)` and `Plugin.UsePolicy(...)` as shorthand for registering policies as middleware.
- `AllPolicies(...)`, `AnyPolicy(...)`, and `NotPolicy(...)` for composing reusable authorization rules without introducing a second execution pipeline.

### Changed
- Bot configuration mutators now treat the bot as configuration-frozen after the first run begins and ignore late mutation attempts for bot-level config such as prefixes, payload defaults, plugins, middleware, runners, localization, scene session wiring, and database context injection.
- `MsgContext` godoc and field comments now describe the normalized update contract more explicitly, including when `Msg`, `From`, callback target fields, `Text`, and `Args` are expected to be populated.
- `MsgContext` normalization now also carries `Chat` and `ChatID` for more Telegram update kinds, allowing policy and update handlers to rely on normalized chat identity outside message-only flows.
- `MsgContext.Error(...)` and returned handler errors now suppress the automatic user reply when the error is explicitly marked with `AsInternalError(...)`, while keeping the previous user-visible default for unclassified errors.
- Godoc, README examples, and regression-test naming now consistently describe the shared generic dependency model as app data, including `NoData` and `SetAppData(...)`.
- Observer configuration now treats `SetObserver(nil)` as clearing instrumentation instead of leaving the previous observer attached.
- Observer lifecycle events now cover generic update handlers and scene command, step, and message-fallback handlers with logical handler names and durations.
- `RequirePolicy(...)` now emits `PolicyCheckedEvent` for both passed and denied policy decisions.
- Scene command, step, and message-fallback flows now emit observer `ErrorEvent`s with scene-specific handler kinds and logical handler names.
- Scene transition observer events now use the same transition payload for scene command, step, and message-fallback flows.
- Observer error emission now also covers generic update handlers, callback payload decode failures, runner failures, and polling retries, including dedicated runner and polling handler kinds in `ErrorEvent`.
- `TODO.md` and the framework backlog pages now mark the observability model as completed for `v1.0.0-rc.13`.
- `tgapi.Chat.Type` now uses the typed `tgapi.ChatType` enum in public DTOs and tests instead of raw string casts.

### Tests
- Added regression coverage for the bot configuration freeze model, including ignored post-run mutations for core bot configuration methods and late registration paths.
- Added table-driven update-contract coverage for `prepareUpdateCtx(...)`, including message-backed, callback-backed, user-backed, and no-user update kinds.
- Added regression tests for policy middleware blocking, built-in private-chat policy decisions, normalized chat identity, and admin checks that use normalized `ChatID` and `FromID`.
- Added regression tests for policy composition semantics, including all-of, any-of, and deny inversion with preserved internal failures.
- Added regression tests for `SetObserver(...)`, `GetObserver()`, and clearing the observer with `SetObserver(nil)`.
- Added observer regression tests for generic update-handler errors, callback payload decode failures, runner failure events, and polling retry emission.
- Added observer regression tests for update and scene handler lifecycle events and `PolicyCheckedEvent` emission.
- Added regression tests proving that `edited_message` and `edited_channel_post` stay out of command routing and continue through generic update handlers.
- Added callback-routing regression tests for both chat-message and inline-message callback targets, including `CallbackQueryId`, `CallbackMsgId`, `InlineMsgId`, and payload-argument guarantees.
- Added regression tests for the new error-visibility model in both message and callback flows, including silent internal-only errors and explicit user-visible callback replies.

## v1.0.0-rc.12

### Added
- `AnswerLong(...)`, `AnswerLongf(...)`, `KeyboardLong(...)`, and `SplitMessageText(...)` for explicit plain-text splitting of long replies without changing the semantics of existing single-message helpers.
- Centralized library-level validation errors in `errors.go`, including `ErrEmptyMessage`, `ErrMessageTooLong`, `ErrCaptionTooLong`, and context/target validation sentinels.
- `Bot.GetPayloadType()`, `InlineKeyboard.GetPayloadType()`, and optional strict payload decoding via `BotOpts.StrictPayloadType` / `Bot.SetStrictPayloadType(...)`.
- `MsgContext.BindArgs(...)` for binding positional command arguments into exported struct fields.
- Binding sentinels `ErrBindArgsTargetNotPointer`, `ErrBindArgsTargetNotStruct`, `ErrBindArgsUnsupportedFieldType`, and `ErrBindArgsConversion`.
- Work-in-progress scene/session support, including plugin scene registration, scoped scene sessions, scene entry/exit APIs on `MsgContext`, default in-memory session storage, scene-local routing before normal command handling, and state helpers on `SceneContext`.

### Changed
- `CommandExecutor` now returns `error`, and command, payload, and non-command update handlers now use centralized bot error handling for returned errors.
- README and README_RU examples now use the new handler signature and document the long-message helpers.
- README and README_RU now link to the project wiki, and the wiki now includes a page-priority tracker while content is being filled in.
- README and README_RU now document scenes, session scopes, scene state helpers, and `SceneActionPass` semantics.
- `TODO.md` and the framework backlog pages now group the remaining framework work into explicit priority 1, 2, and 3 buckets.
- Payload-type comments and docs now distinguish between the bot's default payload type and keyboard-local overrides.
- Scene runtime sentinel errors now have explicit godoc comments.
- Public scene structs now document their exported fields more explicitly.
- `MsgContext.Context()` now safely falls back to `context.Background()` when no request-scoped context is attached.
- `MsgContext` reply, edit, callback, delete, action, and draft-limiter paths now use the context accessor instead of reaching into raw internal state.
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
- Added scene regression tests for runtime guards, scene-local command handling, and `SceneActionPass` preserving session state.
- Added scene regression tests for message fallback handling, user-scoped session lookup without `Msg`, and custom `SessionStore` error propagation.

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
