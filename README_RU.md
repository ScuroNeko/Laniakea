# Laniakea

![Laniakea](assets/logo.jpg)

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go&style=flat-square)](https://go.dev/)
[![License: GPL-3.0](https://img.shields.io/badge/License-GPL%203.0-blue.svg?style=flat-square)](LICENSE)
![Gitea Release](https://img.shields.io/gitea/v/release/ScuroNeko/Laniakea?gitea_url=https%3A%2F%2Fgit.scuroneko.dev&sort=semver&display_name=release&style=flat-square&color=purple&link=https%3A%2F%2Fgit.scuroneko.dev%2FScuroNeko%2FLaniakea%2Freleases)

Легковесная, простая в использовании и производительная обёртка для Telegram Bot API на Go. Она упрощает разработку ботов благодаря чистой системе плагинов, поддержке Middleware, автоматической генерации команд и встроенному рейтлимитеру.

[English](README.md)

[Wiki](https://git.scuroneko.dev/ScuroNeko/Laniakea/wiki)

---

## ✨ Возможности

*   **Простой и интуитивный API:** Разработан для лёгкости использования, основан на практических примерах.
*   **Система плагинов:** Организуйте функциональность бота в независимые, переиспользуемые плагины.
*   **Обработка команд:** Легко регистрируйте команды и извлекайте аргументы.
*   **Поддержка промежуточных слоёв (Middleware):** Выполняйте код до или после команд (например, логирование, проверка доступа).
*   **Автоматическая генерация команд:** Генерируйте справку и списки команд автоматически.
*   **Встроенный ограничитель запросов (Rate Limiter):** Защитите бота от превышения лимитов Telegram API (с обработкой `retry_after`).
*   **Контекст данных:** Передавайте общие данные приложения или state в обработчики.
*   **Настраиваемый API:** Комбинируйте `Set...` и `Add...` helper-методы для понятной конфигурации, например `bot.SetErrorTemplate(...).AddPlugins(...)`.
*   **Polling и Webhook Runtime:** Запускайте бота через long polling с `Run()` / `RunWithContext(...)` или через webhook server, которым владеет сам бот, с `RunWebHookWithContext(...)`.

---

## 📦 Установка

```bash
go get git.scuroneko.dev/scuroneko/laniakea
```

или

```bash
go get github.com/scuroneko/laniakea
```

## 🚀 Быстрый старт (с пошаговыми комментариями)
Вот минимальный пример бота "echo/ping" с подробными комментариями.

```go
package main

import (
	"log"

	"git.scuroneko.dev/scuroneko/laniakea" // Импортируем библиотеку Laniakea
)

// echo — это функция-обработчик команды.
// Она получает два параметра:
//   - ctx: контекст сообщения (содержит информацию о сообщении, отправителе, чате и т.д.)
//   - data: ваши общие данные приложения (здесь мы используем NoData — заглушку без общих зависимостей)
func echo(ctx *laniakea.MsgContext, data laniakea.NoData) error {
	// Отвечаем пользователю текстом, который он прислал, без префикса команды.
	// ctx.Text содержит сообщение пользователя, из которого удалена часть с командой.
	ctx.Answer(ctx.Text) // Ввод пользователя БЕЗ команды
	return nil
}

func main() {
	// 1. Создаём опции бота. Замените "TOKEN" на реальный токен от @BotFather.
	opts := &laniakea.BotOpts{Token: "TOKEN"}

	// 2. Инициализируем новый экземпляр бота.
	//    Используем laniakea.NoData как тип данных приложения (общие зависимости не нужны для примера).
	bot, err := laniakea.NewBot[laniakea.NoData](opts)
	if err != nil {
		log.Fatal(err)
	}
	// Гарантируем освобождение ресурсов бота при выходе.
	defer bot.Close()

	// 3. Создаём новый плагин с именем "ping".
	//    Плагины помогают группировать связанные команды и промежуточные обработчики.
	p := laniakea.NewPlugin[laniakea.NoData]("ping")

	// 4. Добавляем команду в плагин.
	//    p.NewCommand(echo, "echo") создаёт команду, которая вызывает функцию 'echo' по команде "/echo".
	p.AddCommand(p.NewCommand(echo, "echo"))

	// 5. Добавляем ещё одну команду, используя анонимную функцию (замыкание).
	//    Эта команда просто отвечает "Pong", когда пользователь отправляет "/ping".
	p.AddCommand(p.NewCommand(func(ctx *laniakea.MsgContext, data laniakea.NoData) error {
		ctx.Answer("Pong")
		return nil
	}, "ping"))

	// 6. Настраиваем бота: задаём шаблон ошибки и добавляем плагин.
	//    SetErrorTemplate устанавливает формат для сообщений об ошибках (где %s будет заменён на текст ошибки).
	//    AddPlugins(p) регистрирует наш плагин "ping" в боте.
	bot = bot.SetErrorTemplate("Ошибка\n\n%s").AddPlugins(p)

	// 7. Автоматически генерируем команды, такие как /start, /help и список всех зарегистрированных команд.
	//    Это необязательно, но очень полезно для большинства ботов.
	if err := bot.AutoGenerateCommands(); err != nil {
		log.Println(err)
	}

	// 8. Запускаем бота, начиная прослушивание обновлений (long polling).
	if err := bot.Run(); err != nil {
		log.Fatal(err)
	}
}
```

### Как это работает
1. `BotOpts`: Содержит конфигурацию, например, токен API.
2. `NewBot[T]`: Создаёт экземпляр бота. Параметр типа T позволяет передать общие данные приложения (например, *sql.DB или контейнер сервисов), которые будут доступны во всех обработчиках. Используйте laniakea.NoData, если они не нужны.
3. `NewPlugin`: Создаёт логическую группу для команд и Middleware.
4. `AddCommand`: Регистрирует команду. Первый аргумент — функция-обработчик (`func(*MsgContext, T) error`), второй — имя команды (без слеша).
5. **Функции-обработчики**: Получают *MsgContext (детали сообщения, методы типа Answer) и ваши данные приложения типа T, а ошибку возвращают для централизованной обработки.
6. `SetErrorTemplate`: Устанавливает шаблон для сообщений об ошибках. Плейсхолдер %s заменяется на текст ошибки.
7. `AutoGenerateCommands`: Регистрирует команды из плагинов в Telegram для поддерживаемых scope.
8. `Run()`: Запускает цикл опроса обновлений бота и возвращает ошибку, если старт или polling завершился неуспешно.
9. `RunWebHookWithContext(...)`: Запускает bot-owned webhook runtime, когда Telegram должен доставлять update по HTTP вместо long polling.
10. Экземпляр `Bot` одноразовый. После завершения `Run()`, `RunWithContext()` или `RunWebHookWithContext()` для следующего запуска создавайте новый бот.

## Конфиг из файла

`BotOpts` можно не только собирать вручную или из environment, но и загружать и сохранять через file codec API.

Из коробки доступно:
- `BotOptsFileJsonCodec` для JSON-файлов.

Пример:

```go
codec := laniakea.BotOptsFileJsonCodec{}
opts, err := laniakea.LoadBotOptsFile(codec, "config.json")
if err != nil {
	log.Fatal(err)
}

bot, err := laniakea.NewBot[laniakea.NoData](opts)
if err != nil {
	log.Fatal(err)
}
```

Плейсхолдеры вида `{{ TG_TOKEN }}` внутри файла перед декодированием разворачиваются из переменных окружения.

Для других форматов можно реализовать собственный codec через интерфейс `BotOptsFileCodec`.
Из коробки сейчас поддерживается только JSON. Если нужен другой формат, например TOML, используй `BotOptsFileJsonCodec` как эталонную реализацию собственного codec.

Подробности есть в wiki: [Bot Options and Configuration RU](https://git.scuroneko.dev/ScuroNeko/Laniakea/wiki/Bot-Options-and-Configuration-RU)

## Webhook Runtime

Laniakea также поддерживает bot-owned webhook runtime через `RunWebHookWithContext(...)` и `RunWebHook(...)`.

Используй его, когда:
- Telegram должен сам отправлять update на твой HTTP endpoint вместо polling.
- Ты хочешь, чтобы webhook-update проходили через ту же внутреннюю очередь, тот же worker pool, тех же runners и тот же single-use lifecycle, что и polling.
- Ты хочешь, чтобы Laniakea сама регистрировала webhook и владела локальным HTTP server.

Практические замечания:
- Задавай `BotWebHookOpts.SecretToken` для аутентификации запросов.
- Непустой `BotWebHookOpts.SecretToken` обязателен, если включён `BotWebHookOpts.UseStatusPath`.
- Используй явный `BotWebHookOpts.Path`, а не `/`.
- Если ты переводишь уже существующий deployment с webhook-режима на long polling, сначала удали webhook через `CloseWebHook()` или `tgapi.DeleteWebhook(...)`. Пока webhook не удалён, Telegram продолжает доставку через него.
- Запускай `RunWebHookWithContext(...)` с cancelable context и после остановки runtime всё равно вызывай `Close()`.

Полное руководство есть в wiki: [Webhook Runtime](https://git.scuroneko.dev/ScuroNeko/Laniakea/wiki/Webhook-Runtime-RU)

## 📖 Основные концепции
### Плагины (Plugins)
Плагины — основной способ организации кода. Плагин может содержать несколько команд и Middleware.

```go
plugin := laniakea.NewPlugin[*MyDB]("admin")
plugin.AddCommand(plugin.NewCommand(banUser, "ban"))
bot.AddPlugins(plugin)
```

### Команды (Commands)
Команда — это функция, которая обрабатывает конкретную команду бота (например, /start).

```go
func myHandler(ctx *laniakea.MsgContext, db *MyDB) error {
    // Доступ к аргументам команды через ctx.Args ([]string)
    // Ответ пользователю: ctx.Answer("какой-то текст")
    return nil
}
```

### Контекст сообщения (MsgContext)
Предоставляет доступ к входящему сообщению и полезные методы для ответа:

- `Answer(text string)`: Отправляет сообщение с parse_mode none.
- `AnswerLong(text string) []*AnswerMessage`: Разбивает длинный plain text на несколько сообщений.
- `AnswerMarkdown(text string)`: Отправляет сообщение, отформатированное MarkdownV2 (экранирование на вашей стороне).
- `Keyboard(text string, keyboard *InlineKeyboard) *AnswerMessage`: Отправляет сообщение с parse_mode none и Inline клавиатурой.
- `KeyboardLong(text string, keyboard *InlineKeyboard) []*AnswerMessage`: Разбивает длинный plain text на несколько сообщений и вешает клавиатуру на последний chunk.
- `KeyboardMarkdown(text string, keyboard *InlineKeyboard) *AnswerMessage`: Отправляет сообщение, отформатированное MarkdownV2 (экранирование на вашей стороне), и Inline клавиатурой.
- `AnswerPhoto(photoId, text string) *AnswerMessage`: Отправляет фотографию с подписью и parse_mode none.
- `AnswerPhotoMarkdown(photoId, text string) *AnswerMessage`: Отправляет фотографию с подписью, отформатированной MarkdownV2 (экранирование на вашей стороне).
- `EditCallback(text string)`: Редактирует сообщение с `parse_mode` none после нажатия inline-кнопки.
- `EditCallbackMarkdown(text string)`: Редактирует сообщение в формате MarkdownV2 (экранирование на вашей стороне) после нажатия inline-кнопки.
- `SendAction(action tgapi.ChatActionType)`: Отправляет действие "печатает", "загружает фото" и т.д.
- Поля: `Text`, `Args`, `From`, `FromID`, `Msg`, `InlineMsgId`, `CallbackQueryId` и другие.
- И много других методов и полей!

### App Data
Параметр типа `T` в `NewBot[T]` — мощная возможность. Вы можете передать любой тип, но для разделяемых зависимостей вроде пула соединений с БД, контейнера сервисов или API-клиента обычно стоит использовать pointer type.

```go
type MyDB struct { /* ... */ }
db := &MyDB{...}
bot, err := laniakea.NewBot[*MyDB](opts)
if err != nil {
    log.Fatal(err)
}
bot.SetAppData(db)
```

### Сцены и сессии (Scenes and Sessions)

Сцены описывают многошаговые диалоги внутри плагина. Активная сцена хранится в session state, ключ которого зависит от scope, поэтому поток можно изолировать на пользователя, на чат или на пару пользователь-чат.

```go
plugin := laniakea.NewPlugin[MyDB]("signup")

plugin.NewScene("signup").
    SetScope(laniakea.SceneScopeUserChat).
    SetEntry("ask_name").
    OnStep("ask_name", func(ctx *laniakea.SceneContext, db MyDB) (laniakea.SceneResult, error) {
        if ctx.Text == "" {
            ctx.Answer("Как тебя зовут?")
            return ctx.Stay(), nil
        }

        if err := ctx.SaveData(struct {
            Name string `json:"name"`
        }{Name: ctx.Text}); err != nil {
            return laniakea.SceneResult{}, err
        }

        ctx.Answer("Приятно познакомиться.")
        return ctx.Next("done"), nil
    }).
    OnStep("done", func(ctx *laniakea.SceneContext, db MyDB) (laniakea.SceneResult, error) {
        return ctx.Exit(), nil
    })
```

- Используйте `ctx.EnterScene("signup")`, чтобы войти в entry step, настроенный у сцены.
- Используйте `ctx.EnterSceneStep("signup", "done")`, если нужен явный стартовый step.
- Из scene handler возвращайте `ctx.Stay()`, `ctx.Next(step)`, `ctx.Exit()` или `ctx.Pass()` для управления потоком.
- `SceneActionPass` не меняет текущую session state и продолжает обычный routing бота.
- Для JSON-состояния сцены используйте `SceneContext.SaveData(...)` и `SceneContext.BindData(...)`.
- Выбирайте `SceneScopeUser`, `SceneScopeChat` или `SceneScopeUserChat` в зависимости от того, насколько широко должен разделяться диалог.

### tgapi: API и Uploader

В `tgapi` есть два клиента:

- `API` для JSON-запросов (`SendMessage`, `EditMessageText`, методы с `file_id`/URL).
- `Uploader` для multipart-загрузок (`SendPhoto`, `SendDocument`, `SendVideo` с бинарными файлами).

Для продвинутых сценариев `tgapi.NewRequest(...)` и `tgapi.NewUploaderRequest(...)` остаются публичными low-level escape hatch API. Они менее безопасны, чем типизированные helper-методы: вызывающая сторона сама отвечает за корректное имя Telegram-метода и совместимые типы параметров/ответа.

## 🧩 Промежуточные слои (Middleware)
Middleware — это функции, которые выполняются перед обработчиком команды. Они идеально подходят для сквозных задач, таких как логирование, контроль доступа, ограничение скорости запросов или модификация контекста.

### Сигнатура
Функция middleware имеет ту же сигнатуру, что и обработчик команды, но должна возвращать bool:

```go
func(ctx *MsgContext, db T) bool
```

- Если возвращается true, выполняется следующий middleware (или сама команда).
- Если возвращается false, цепочка выполнения немедленно прерывается (команда не запускается).

### Добавление middleware
Используйте метод `AddMiddleware` плагина для добавления одной или нескольких функций middleware. Они выполняются в порядке добавления.

```go
plugin := laniakea.NewPlugin[*MyDB]("admin")
plugin.AddMiddleware(laniakea.NewMiddleware("logging", loggingMiddleware))
plugin.AddMiddleware(laniakea.NewMiddleware("admin-only", adminOnlyMiddleware))
plugin.AddCommand(plugin.NewCommand(banUser, "ban"))
```

### Примеры middleware

1. Логирующий middleware – логирует каждое выполнение команды.
```go
func loggingMiddleware(ctx *laniakea.MsgContext, db *MyDB) bool {
    log.Printf("Пользователь %d выполнил команду: %s", ctx.FromID, ctx.Msg.Text)
    return true // продолжаем к следующему middleware/команде
}
```

2. Middleware только для администраторов – ограничивает доступ пользователям с определённой ролью.
```go
func adminOnlyMiddleware(ctx *laniakea.MsgContext, db *MyDB) bool {
    if !db.IsAdmin(ctx.FromID) { // предполагается, что db имеет метод IsAdmin
        ctx.Answer("⛔ Доступ запрещён. Только для администраторов.")
        return false // останавливаем выполнение
    }
    return true
}
```

### Важные замечания
- Middleware может изменять MsgContext (например, добавлять пользовательские поля) перед запуском команды.

## ⚙️ Расширенная настройка
- **Инлайн-клавиатуры**: Создавайте клавиатуры с помощью `laniakea.NewInlineKeyboardJson`, `laniakea.NewInlineKeyboardBase64` или `laniakea.NewInlineKeyboard`. `Bot.SetPayloadType(...)` задаёт payload format по умолчанию, а `InlineKeyboard.SetPayloadType(...)` переопределяет его для конкретной клавиатуры.
- **Ограничение запросов**: Передайте настроенный `utils.RateLimiter` через `BotOpts` для корректной обработки лимитов Telegram.
- **Локализация**: `L10n` безопасен для конкурентного использования после подключения к боту.
- **Пользовательские update handlers**: Используйте `plugin.AddUpdateHandler(...)` для Telegram update types вне command/payload flow.
- **Жизненный цикл**: `RunWithContext(...)` и `RunWebHookWithContext(...)` не вызывают `Close()` автоматически. Завершайте бот явно и создавайте новый `Bot` для следующего запуска.

## Обработка Telegram Updates
- Команды и payload-ы обрабатываются через плагины.
- Для некомандных update-ов можно зарегистрировать обработчик через `plugin.AddUpdateHandler(updateType, handler)`.
- `message`, `channel_post` и `callback_query` остаются в command/payload flow.
- После JSON-декодирования `tgapi.Update` заполняет поле `Type`, чтобы обработчики могли явно видеть итоговый вид update.

## 📝 Лицензия
Этот проект лицензирован под GNU General Public License v3.0 - подробности см. в файле [LICENSE](LICENSE).

## 📚 Дополнительная информация
[GoDoc Laniakea](https://pkg.go.dev/git.scuroneko.dev/scuroneko/laniakea)

[Wiki](https://git.scuroneko.dev/ScuroNeko/Laniakea/wiki)

[Telegram Bot API](https://core.telegram.org/bots/api)

    ✅ Создано с ❤️ scuroneko
