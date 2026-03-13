# Laniakea

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go&style=flat-square)](https://go.dev/)
[![License: GPL-3.0](https://img.shields.io/badge/License-GPL%203.0-blue.svg?style=flat-square)](LICENSE)
![Gitea Release](https://img.shields.io/gitea/v/release/ScuroNeko/Laniakea?gitea_url=https%3A%2F%2Fgit.nix13.pw&sort=semver&display_name=release&style=flat-square&color=purple&link=https%3A%2F%2Fgit.nix13.pw%2FScuroNeko%2FLaniakea%2Freleases)

A lightweight, easy-to-use, and performant Telegram Bot API wrapper for Go. It simplifies bot development with a clean plugin system, middleware support, automatic command generation, and built-in rate limiting.

[На русском](README_RU.md)

---

## ✨ Features
*   **Simple & Intuitive API:** Designed for ease of use, based on practical examples.
*   **Plugin System:** Organize your bot's functionality into independent, reusable plugins.
*   **Command Handling:** Easily register commands and extract arguments.
*   **Middleware Support:** Run code before or after commands (e.g., logging, access control).
*   **Automatic Command Generation:** Generate help and command lists automatically.
*   **Built-in Rate Limiting:** Protect your bot from hitting Telegram API limits (supports `retry_after` handling).
*   **Context-Aware:** Pass custom database or state contexts to your handlers.
*   **Fluent Interface:** Chain methods for clean configuration (e.g., `bot.ErrorTemplate(...).AddPlugins(...)`).

---

## 📦 Installation

```bash
go get git.nix13.pw/scuroneko/laniakea
```

## 🚀 Quick Start (with step-by-step explanation)

Here is a minimal echo/ping bot example with detailed comments.
```go
package main

import (
	"log"

	"git.nix13.pw/scuroneko/laniakea" // Import the Laniakea library
)

// echo is a command handler function.
// It receives two parameters:
//   - ctx: the message context (contains info about the message, sender, chat, etc.)
//   - db: your custom database context (here we use NoDB, a placeholder for no database)
func echo(ctx *laniakea.MsgContext, db *laniakea.NoDB) {
	// Answer the user with the text they sent, without any command prefix.
	// ctx.Text contains the user's message with the command part stripped off.
	ctx.Answer(ctx.Text) // User input WITHOUT command
}

func main() {
	// 1. Create bot options. Replace "TOKEN" with your actual bot token from @BotFather.
	opts := &laniakea.BotOpts{Token: "TOKEN"}

	// 2. Initialize a new bot instance.
	//    We use laniakea.NoDB as the database context type (no database needed for this example).
	bot := laniakea.NewBot[laniakea.NoDB](opts)
	// Ensure bot resources are cleaned up on exit.
	defer bot.Close()

	// 3. Create a new plugin named "ping".
	//    Plugins help group related commands and middlewares.
	p := laniakea.NewPlugin[laniakea.NoDB]("ping")

	// 4. Add a command to the plugin.
	//    p.NewCommand(echo, "echo") creates a command that triggers the 'echo' function on the "/echo" command.
	p.AddCommand(p.NewCommand(echo, "echo"))

	// 5. Add another command using an anonymous function (closure).
	//    This command simply replies "Pong" when the user sends "/ping".
	p.AddCommand(p.NewCommand(func(ctx *laniakea.MsgContext, db *laniakea.NoDB) {
		ctx.Answer("Pong")
	}, "ping"))

	// 6. Configure the bot with a custom error template and add the plugin.
	//    ErrorTemplate sets a format string for errors (where %s will be replaced by the actual error).
	//    AddPlugins(p) registers our "ping" plugin with the bot.
	bot = bot.ErrorTemplate("Error\n\n%s").AddPlugins(p)

	// 7. Automatically generate commands like /start, /help, and a list of all registered commands.
	//    This is optional but very useful for most bots.
	if err := bot.AutoGenerateCommands(); err != nil {
		log.Println(err)
	}

	// 8. Start the bot, listening for updates (long polling).
	bot.Run()
}
```

### How It Works
1. `BotOpts`: Holds configuration like the API token.
2. `NewBot[T]`: Creates a bot instance. The type parameter T allows you to pass a custom database context (e.g., *sql.DB) that will be available in all handlers. Use laniakea.NoDB if you don't need it.
3. `NewPlugin`: Creates a logical group for commands and middlewares.
4. `AddCommand`: Registers a command. The first argument is the handler function (func(*MsgContext, T)), the second is the command name (without the slash).
5. **Handler Functions**: Receive *MsgContext (message details, methods like Answer) and your custom database context T.
6. `ErrorTemplate`: Sets a template for error messages. The %s placeholder is replaced by the actual error.
7. `AutoGenerateCommands`: Adds built-in commands (/start, /help) and a command that lists all available commands.
8. `Run()`: Starts the bot's update polling loop.

## 📖 Core Concepts
### Plugins

Plugins are the main way to organize code. A plugin can have multiple commands and middlewares.
```go
plugin := laniakea.NewPlugin[MyDB]("admin")
plugin.AddCommand(plugin.NewCommand(banUser, "ban"))
bot.AddPlugins(plugin)
```

### Commands

A command is a function that handles a specific bot command (e.g., /start).
```go
func myHandler(ctx *laniakea.MsgContext, db *MyDB) {
    // Access command arguments via ctx.Args ([]string)
    // Reply to the user: ctx.Answer("some text")
}
```

### MsgContext

Provides access to the incoming message and useful reply methods:

- `Answer(text string) *AnswerMessage`: Sends a message with parse_mode none.
- `AnswerMarkdown(text string) *AnswerMessage`: Sends a message formatted with MarkdownV2 (you handle escaping).
- `Keyboard(text string, keyboard *InlineKeyboard) *AnswerMessage`: Sends a message with parse_mode none and inline keyboard.
- `KeyboardMarkdown(text string, keyboard *InlineKeyboard) *AnswerMessage`: Sends a message formatted with MarkdownV2 (you handle escaping) and inline keyboard.
- `AnswerPhoto(photoId, text string) *AnswerMessage`: Sends a message with photo with parse_mode none.
- `AnswerPhotoMarkdown(photoId, text string) *AnswerMessage`: Sends a message formatted with MarkdownV2 (you handle escaping) with.
- `EditCallback(text string)`: Edits message with parse_mode none after clicking inline button.
- `EditCallbackMarkdown(text string)`: Edits a message formatted with MarkdownV2 (you handle escaping) after clicking inline button.
- `SendChatAction(action string)`: Sends a “typing”, “uploading photo”, etc., action.
- Fields: `Text`, `Args`, `From`, `Chat`, `Msg`, etc.
- And more methods and fields!

### Database Context

The `T` in `NewBot[T]` is a powerful feature. You can pass any type (like a database connection pool), and it will be available in every command and middleware handler.

```go
type MyDB struct { /* ... */ }
db := &MyDB{...}
bot := laniakea.NewBot[*MyDB](opts, db) // Pass db instance
```

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
Use the Use method of a plugin to add one or more middleware functions. They are executed in the order they are added.

```go
plugin := laniakea.NewPlugin[MyDB]("admin")
plugin.Use(loggingMiddleware, adminOnlyMiddleware)
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
- **Inline Keyboards**: Build keyboards using laniakea.NewKeyboard().
- **Rate Limiting**: Pass a configured utils.RateLimiter via BotOpts to handle Telegram's rate limits gracefully.
- **Custom HTTP Client**: Provide your own http.Client in BotOpts for fine-tuned control.

## 📝 License

This project is licensed under the GNU General Public License v3.0 — see the [LICENSE](LICENSE) file for details.

## 📚 Learn More
[GoDoc](https://pkg.go.dev/git.nix13.pw/scuroneko/laniakea)

[Telegram Bot API](https://core.telegram.org/bots/api)

    ✅ Built with ❤️ by scuroneko
