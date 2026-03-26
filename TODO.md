# TODO

This file tracks framework-level backlog items that are about missing concepts in the library itself, not just missing documentation.

## High-Priority Core Concepts

### 1. Conversation / Scene Model

Current state:
- The framework is strong at handling a single update through commands, payloads, middleware, and update handlers.
- It already has useful lower-level building blocks such as `MsgContext`, drafts, payload routing, plugins, and update handlers.
- It does not yet provide a first-class concept for long-lived user interaction flows.

Why this matters:
- Many Telegram bots quickly move beyond isolated commands and need stateful multi-step flows.
- Real bots often need concepts like "wait for the user's next message", "user is currently on step 3 of 5", or "button press moves the user to the next scene state".
- Without a scene model, library users end up building their own mini-framework on top of Laniakea.

What is missing:
- A way to route updates to an active scene before normal command routing.
- A way to persist conversation state per user or per chat.
- A way to describe steps and transitions without hand-rolling state machines around middleware and storage.
- A way to enter, continue, cancel, and complete a conversation flow explicitly.
- A way to support modal chat flows where the user is "inside" a scene and ordinary text is treated as scene input until an explicit escape command exits the mode.

Possible API direction:
- `Scene`, `Step`, and `SessionStore` concepts.
- `bot.AddScene(...)` or a dedicated scene registry.
- `ctx.Scene()`, `ctx.NextStep(...)`, `ctx.ExitScene()`, or similar state-transition helpers.
- Routing rule: active scene first, then normal command/payload flow if no scene claims the update.
- Storage-backed per-user or per-chat state with a clean interface for custom persistence.
- Scene-local escape and passthrough commands, so flows like `/startrp` can put a user into a dedicated chat mode where most messages go straight to the scene, while commands like `/exit` or a small whitelist still retain special meaning.

Important design constraints:
- This should be additive and optional.
- It should not replace plugins, commands, or handlers as the normal framework entry points.
- It should work with existing middleware and `MsgContext` instead of introducing a second incompatible execution model.

Practical target:
- Make stateful bot flows a first-class, framework-supported pattern instead of a userland convention.
- Cover both step-based forms and mode-based chat flows without forcing users to build custom routing layers around active sessions.

### 2. Typed Handler Input Model

Current state:
- Commands and payloads currently expose parsed text through `ctx.Text` and `ctx.Args`.
- `CommandArg` provides basic argument validation and shape checks.
- Handlers still do most non-trivial parsing manually.

Why this matters:
- As bots grow, handlers often start with repetitive `ctx.Args` parsing boilerplate.
- Validation logic tends to spread across handlers instead of living in one predictable binding layer.
- The current model is simple and honest, but it does not help enough once commands become more structured.

What is missing:
- A first-class way to bind command or payload arguments into a typed Go value.
- A framework-level pattern for conversion errors and validation errors beyond raw string handling.
- A low-friction way to move from positional arguments to a structured input object.

Possible API direction:
- A lightweight binding API such as `ctx.BindArgs(&input)`.
- Or explicit typed command registration such as `NewCommandTyped(...)`.
- Positional mapping into structs, optional fields, basic conversion support, and integration with current validation flow.
- Unified binding and validation failures routed through the current centralized error path.

Example of the kind of user code this should enable:

```go
type BanInput struct {
	UserID int
	Reason string
}

func ban(ctx *laniakea.MsgContext, db *App) error {
	var input BanInput
	if err := ctx.BindArgs(&input); err != nil {
		return err
	}

	return db.Ban(input.UserID, input.Reason)
}
```

Important design constraints:
- Avoid a reflection-heavy, magical subsystem.
- Keep the current `ctx.Args` model as the minimal baseline.
- Treat typed binding as an ergonomic layer on top of the current command model, not a replacement for it.

Practical target:
- Remove repetitive parsing boilerplate while preserving the framework's explicit, Go-like feel.

### 3. Request Context / Cancellation Model

Current state:
- `RunWithContext(...)` controls bot runtime lifecycle and graceful shutdown.
- `tgapi` already supports context-aware methods.
- Regular handlers do not receive a first-class request-scoped `context.Context`.

Why this matters:
- Handler business logic often needs cancellation-aware database calls, HTTP calls, or downstream service calls.
- The framework already has a good runtime cancellation story, but it does not flow naturally into user code inside handlers.
- In modern Go APIs, `context.Context` is a standard part of operational correctness.

What is missing:
- A clean request-scoped context that follows each update through handler execution.
- A standard way for application code to stop work when the bot is shutting down or the update processing context is canceled.
- A direct bridge between bot lifecycle control and service-layer cancellation.

Possible API direction:
- Prefer a non-breaking approach by exposing context through `MsgContext`, for example `ctx.Context()`.
- Build the context from the update-processing lifecycle so it is meaningful during graceful shutdown.
- Make it natural to pass that context into database methods, HTTP clients, and `tgapi.WithContext(...)` calls.

Why this should probably not be a signature change:
- Changing handler signatures to accept `context.Context` directly would be a public breaking change.
- A `MsgContext` accessor would preserve compatibility while still giving handlers an idiomatic Go cancellation path.

Practical target:
- Let handler code participate naturally in cancellation and graceful shutdown without forcing users to invent their own context plumbing.

## Secondary Backlog

- Webhook runtime model: the library has a solid polling model, but no first-class webhook execution model at the framework level.
- Service layer and dependency graph model: `DatabaseContext(T)` is intentionally minimal, but there is no stronger framework concept for application services or scoped dependencies.
- User-facing vs internal error model: the framework has a unified error flow, but it does not yet distinguish well between user-visible, internal-only, retryable, or silent errors.
- Authorization and policy model: middleware can implement auth and permissions, but there is no explicit framework concept for access policies, roles, or capability checks.
- Observability model: logging is strong, but metrics, tracing, and structured framework hooks are still missing as first-class concepts.
- Plugin composition contract: plugins are a good grouping unit, but there is no explicit model for plugin dependencies, shared capabilities, or composition contracts.
- Update schema contract: update handling exists, but there is no formal framework-level concept describing which `MsgContext` fields are guaranteed in which update kinds.
- Configuration freeze model: the framework already has real commit points like `AddPlugins(...)`, but this is still more of an implementation truth than an explicit top-level concept.

## Suggested Priority

1. Request context / cancellation model
2. Conversation / scene model
3. Typed handler input model
