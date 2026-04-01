# Laniakea

![Laniakea](assets/logo.jpg)

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go&style=flat-square)](https://go.dev/)
[![License: GPL-3.0](https://img.shields.io/badge/License-GPL%203.0-blue.svg?style=flat-square)](LICENSE)
![Gitea Release](https://img.shields.io/gitea/v/release/ScuroNeko/Laniakea?gitea_url=https%3A%2F%2Fgit.scuroneko.dev&sort=semver&display_name=release&style=flat-square&color=purple&link=https%3A%2F%2Fgit.scuroneko.dev%2FScuroNeko%2FLaniakea%2Freleases)

A lightweight, easy-to-use, and performant Telegram Bot API wrapper for Go. It simplifies bot development with a clean plugin system, middleware support, automatic command generation, and built-in rate limiting.

[На русском](README_RU.md)

[Wiki](https://git.scuroneko.dev/ScuroNeko/Laniakea/wiki)

---

## ✨ Features
*   **Simple & Intuitive API:** Designed for ease of use, based on practical examples.
*   **Plugin System:** Organize your bot's functionality into independent, reusable plugins.
*   **Command Handling:** Easily register commands and extract arguments.
*   **Middleware Support:** Run code before or after commands (e.g., logging, access control).
*   **Automatic Command Generation:** Generate help and command lists automatically.
*   **Built-in Rate Limiting:** Protect your bot from hitting Telegram API limits (supports `retry_after` handling).
*   **Context-Aware:** Pass custom application data or state contexts to your handlers.
*   **Configurable API:** Mix `Set...` and `Add...` helpers to configure bots clearly (for example, `bot.SetErrorTemplate(...).AddPlugins(...)`).
*   **Polling and Webhook Runtime:** Run bots through long polling with `Run()` / `RunWithContext(...)` or through a bot-owned webhook server with `RunWebHookWithContext(...)`.

---

## 📦 Installation

```bash
go get git.scuroneko.dev/scuroneko/laniakea
```

or

```bash
go get github.com/scuroneko/laniakea
```

## 🚀 Quick Start (with step-by-step explanation)

Here is a minimal echo/ping bot example with detailed comments.
```go
package main

import (
	"log"

	"git.scuroneko.dev/scuroneko/laniakea" // Import the Laniakea library
)

// echo is a command handler function.
// It receives two parameters:
//   - ctx: the message context (contains info about the message, sender, chat, etc.)
//   - data: your shared application data (here we use NoData, a placeholder for no shared data)
func echo(ctx *laniakea.MsgContext, data laniakea.NoData) error {
	// Answer the user with the text they sent, without any command prefix.
	// ctx.Text contains the user's message with the command part stripped off.
	ctx.Answer(ctx.Text) // User input WITHOUT command
	return nil
}

func main() {
	// 1. Create bot options. Replace "TOKEN" with your actual bot token from @BotFather.
	opts := &laniakea.BotOpts{Token: "TOKEN"}

	// 2. Initialize a new bot instance.
	//    We use laniakea.NoData as the application data type (no shared data needed for this example).
	bot, err := laniakea.NewBot[laniakea.NoData](opts)
	if err != nil {
		log.Fatal(err)
	}
	// Ensure bot resources are cleaned up on exit.
	defer bot.Close()

	// 3. Create a new plugin named "ping".
	//    Plugins help group related commands and middlewares.
	p := laniakea.NewPlugin[laniakea.NoData]("ping")

	// 4. Add a command to the plugin.
	//    p.NewCommand(echo, "echo") creates a command that triggers the 'echo' function on the "/echo" command.
	p.AddCommand(p.NewCommand(echo, "echo"))

	// 5. Add another command using an anonymous function (closure).
	//    This command simply replies "Pong" when the user sends "/ping".
	p.AddCommand(p.NewCommand(func(ctx *laniakea.MsgContext, data laniakea.NoData) error {
		ctx.Answer("Pong")
		return nil
	}, "ping"))

	// 6. Configure the bot with a custom error template and add the plugin.
	//    SetErrorTemplate sets a format string for errors (where %s will be replaced by the actual error).
	//    AddPlugins(p) registers our "ping" plugin with the bot.
	bot = bot.SetErrorTemplate("Error\n\n%s").AddPlugins(p)

	// 7. Automatically generate commands like /start, /help, and a list of all registered commands.
	//    This is optional but very useful for most bots.
	if err := bot.AutoGenerateCommands(); err != nil {
		log.Println(err)
	}

	// 8. Start the bot, listening for updates (long polling).
	if err := bot.Run(); err != nil {
		log.Fatal(err)
	}
}
```

### How It Works
1. `BotOpts`: Holds configuration like the API token.
2. `NewBot[T]`: Creates a bot instance. The type parameter T allows you to pass custom shared application data (for example, *sql.DB or a service container) that will be available in all handlers. Use laniakea.NoData if you don't need it.
3. `NewPlugin`: Creates a logical group for commands and middlewares.
4. `AddCommand`: Registers a command. The first argument is the handler function (`func(*MsgContext, T) error`), the second is the command name (without the slash).
5. **Handler Functions**: Receive *MsgContext (message details, methods like Answer) and your custom application data T, and return an error for centralized error handling.
6. `SetErrorTemplate`: Sets a template for error messages. The %s placeholder is replaced by the actual error.
7. `AutoGenerateCommands`: Registers plugin-defined commands with Telegram across the supported scopes.
8. `Run()`: Starts the bot's update polling loop and returns an error if startup or polling fails.
9. `RunWebHookWithContext(...)`: Starts the bot-owned webhook runtime when Telegram should deliver updates over HTTP instead of long polling.
10. A `Bot` instance is single-use. After `Run()`, `RunWithContext()`, or `RunWebHookWithContext()` returns, create a new bot instance for the next session.

## Webhook Runtime

Laniakea also supports a bot-owned webhook runtime through `RunWebHookWithContext(...)` and `RunWebHook(...)`.

Use it when:
- Telegram should push updates to your HTTP endpoint instead of your bot polling for them.
- You want webhook-delivered updates to reuse the same internal queue, worker pool, runners, and single-use lifecycle as polling.
- You want Laniakea to register the webhook and own the local HTTP server.

Production notes:
- Set `BotWebHookOpts.SecretToken` for request authentication.
- `BotWebHookOpts.SecretToken` is required when `BotWebHookOpts.UseStatusPath` is enabled.
- Keep `BotWebHookOpts.Path` specific instead of serving webhook traffic on `/`.
- If you switch an existing deployment from webhook mode to long polling, delete the webhook first with `CloseWebHook()` or `tgapi.DeleteWebhook(...)`. Telegram keeps webhook delivery active until it is removed.
- Use `RunWebHookWithContext(...)` with a cancelable context, then call `Close()` after runtime shutdown.

See the full guide in the wiki: [Webhook Runtime](https://git.scuroneko.dev/ScuroNeko/Laniakea/wiki/Webhook-Runtime)

## 📖 Core Concepts
### Plugins

Plugins are the main way to organize code. A plugin can have multiple commands and middlewares.
```go
plugin := laniakea.NewPlugin[*MyDB]("admin")
plugin.AddCommand(plugin.NewCommand(banUser, "ban"))
bot.AddPlugins(plugin)
```

### Commands

A command is a function that handles a specific bot command (e.g., /start).
```go
func myHandler(ctx *laniakea.MsgContext, db *MyDB) error {
    // Access command arguments via ctx.Args ([]string)
    // Reply to the user: ctx.Answer("some text")
    return nil
}
```

### MsgContext

Provides access to the incoming message and useful reply methods:

- `Answer(text string) *AnswerMessage`: Sends a message with parse_mode none.
- `AnswerLong(text string) []*AnswerMessage`: Splits long plain text into multiple messages.
- `AnswerMarkdown(text string) *AnswerMessage`: Sends a message formatted with MarkdownV2 (you handle escaping).
- `Keyboard(text string, keyboard *InlineKeyboard) *AnswerMessage`: Sends a message with parse_mode none and inline keyboard.
- `KeyboardLong(text string, keyboard *InlineKeyboard) []*AnswerMessage`: Splits long plain text into multiple messages and attaches the keyboard to the final chunk.
- `KeyboardMarkdown(text string, keyboard *InlineKeyboard) *AnswerMessage`: Sends a message formatted with MarkdownV2 (you handle escaping) and inline keyboard.
- `AnswerPhoto(photoId, text string) *AnswerMessage`: Sends a message with photo with parse_mode none.
- `AnswerPhotoMarkdown(photoId, text string) *AnswerMessage`: Sends a photo with MarkdownV2 caption (you handle escaping).
- `EditCallback(text string)`: Edits message with parse_mode none after clicking inline button.
- `EditCallbackMarkdown(text string)`: Edits a message formatted with MarkdownV2 (you handle escaping) after clicking inline button.
- `SendAction(action tgapi.ChatActionType)`: Sends a “typing”, “uploading photo”, etc., action.
- Fields: `Text`, `Args`, `From`, `FromID`, `Msg`, `InlineMsgId`, `CallbackQueryId`, etc.
- And more methods and fields!

### tgapi: API and Uploader

`tgapi` provides two clients:

- `API` for JSON requests (e.g., `SendMessage`, `EditMessageText`, methods using file_id/URL).
- `Uploader` for multipart uploads (e.g., `SendPhoto`, `SendDocument`, `SendVideo` with binary files).

This split keeps method intent explicit: JSON-only calls go through `API`, file uploads go through `Uploader`.

For advanced cases, `tgapi.NewRequest(...)` and `tgapi.NewUploaderRequest(...)` remain public as low-level escape hatches. They are intentionally less safe than method-specific helpers: callers must supply the correct Telegram method name and compatible request/response types themselves.

### App Data

The `T` in `NewBot[T]` is a powerful feature. You can pass any type, but shared dependencies such as database pools, service containers, or API clients should usually use a pointer type.

```go
type MyDB struct { /* ... */ }
db := &MyDB{...}
bot, err := laniakea.NewBot[*MyDB](opts)
if err != nil {
    log.Fatal(err)
}
bot.SetAppData(db)
```

### Scenes and Sessions

Scenes model multi-step conversations inside a plugin. Each active scene is stored in a session keyed by scope, so you can isolate flows per user, per chat, or per user-chat pair.

```go
plugin := laniakea.NewPlugin[MyDB]("signup")

plugin.NewScene("signup").
    SetScope(laniakea.SceneScopeUserChat).
    SetEntry("ask_name").
    OnStep("ask_name", func(ctx *laniakea.SceneContext, db MyDB) (laniakea.SceneResult, error) {
        if ctx.Text == "" {
            ctx.Answer("What is your name?")
            return ctx.Stay(), nil
        }

        if err := ctx.SaveData(struct {
            Name string `json:"name"`
        }{Name: ctx.Text}); err != nil {
            return laniakea.SceneResult{}, err
        }

        ctx.Answer("Nice to meet you.")
        return ctx.Next("done"), nil
    }).
    OnStep("done", func(ctx *laniakea.SceneContext, db MyDB) (laniakea.SceneResult, error) {
        return ctx.Exit(), nil
    })
```

- Use `ctx.EnterScene("signup")` to enter the configured entry step.
- Use `ctx.EnterSceneStep("signup", "done")` when you need an explicit starting step.
- Return `ctx.Stay()`, `ctx.Next(step)`, `ctx.Exit()`, or `ctx.Pass()` from scene handlers to control flow.
- `SceneActionPass` keeps the current session unchanged and continues normal bot routing.
- Use `SceneContext.SaveData(...)` and `SceneContext.BindData(...)` for JSON session state.
- Use `SceneScopeUser`, `SceneScopeChat`, or `SceneScopeUserChat` depending on how widely a conversation should be shared.

## 🧩 Middleware
Middleware are functions that run before a command handler. They are perfect for cross-cutting concerns like logging, access control, rate limiting, or modifying the context.

### Signature
A middleware function has the same signature as a command handler, but it must return a bool:

```go
func(ctx *MsgContext, db T) bool
```

- If it returns true, the next middleware (or the command) will be executed.
- If it returns false, the execution chain stops immediately (the command will not run).

### Adding Middleware
Use `AddMiddleware` on a plugin to add one or more shared middleware functions. They are executed in the order they are added.

```go
plugin := laniakea.NewPlugin[*MyDB]("admin")
plugin.AddMiddleware(laniakea.NewMiddleware("logging", loggingMiddleware))
plugin.AddMiddleware(laniakea.NewMiddleware("admin-only", adminOnlyMiddleware))
plugin.AddCommand(plugin.NewCommand(banUser, "ban"))
```

### Example Middlewares

1. Logging Middleware – logs every command execution.
```go
func loggingMiddleware(ctx *laniakea.MsgContext, db *MyDB) bool {
    log.Printf("User %d executed command: %s", ctx.FromID, ctx.Msg.Text)
    return true // continue to next middleware/command
}
```

2. Admin-Only Middleware – restricts access to users with a specific role.
```go
func adminOnlyMiddleware(ctx *laniakea.MsgContext, db *MyDB) bool {
    if !db.IsAdmin(ctx.FromID) { // assume db has IsAdmin method
        ctx.Answer("⛔ Access denied. Admins only.")
        return false // stop execution
    }
    return true
}
```

### Important Notes
- Middleware can modify the MsgContext (e.g., add custom fields) before the command runs.

## ⚙️ Advanced Configuration
- **Inline Keyboards**: Build keyboards using `laniakea.NewInlineKeyboardJson`, `laniakea.NewInlineKeyboardBase64`, or `laniakea.NewInlineKeyboard`. `Bot.SetPayloadType(...)` defines the default payload format, and `InlineKeyboard.SetPayloadType(...)` overrides it for one keyboard.
- **Rate Limiting**: Pass a configured utils.RateLimiter via BotOpts to handle Telegram's rate limits gracefully.
- **Localization**: `L10n` is safe for concurrent use once attached to the bot.
- **Custom Update Handlers**: Use `plugin.AddUpdateHandler(...)` for Telegram update types that are not part of the command/payload flow.
- **Lifecycle**: `RunWithContext(...)` and `RunWebHookWithContext(...)` do not call `Close()` for you. Shut the bot down explicitly, and create a fresh `Bot` for the next run.

## Telegram Update Handling
- Commands and payloads are handled through plugins.
- Non-command updates can be routed with `plugin.AddUpdateHandler(updateType, handler)`.
- `message`, `channel_post`, and `callback_query` stay on the command/payload flow.
- `tgapi.Update` exposes a derived `Type` field after JSON unmarshalling so handlers can inspect the effective update kind directly.

## 📝 License

This project is licensed under the GNU General Public License v3.0 — see the [LICENSE](LICENSE) file for details.

## 📚 Learn More
[GoDoc](https://pkg.go.dev/git.scuroneko.dev/scuroneko/laniakea)

[Wiki](https://git.scuroneko.dev/ScuroNeko/Laniakea/wiki)

[Telegram Bot API](https://core.telegram.org/bots/api)

    ✅ Built with ❤️ by scuroneko
