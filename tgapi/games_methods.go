package tgapi

import "context"

// SendGame holds parameters for the sendGame method.
// Since: Bot API 2.2
// See https://core.telegram.org/bots/api#sendgame
type SendGame struct {
	BusinessConnectionID string `json:"business_connection_id,omitempty"`
	ChatID               int64  `json:"chat_id"`
	MessageThreadID      int    `json:"message_thread_id,omitempty"`

	GameShortName string `json:"game_short_name"`

	DisableNotification bool                  `json:"disable_notification,omitempty"`
	ProtectContent      bool                  `json:"protect_content,omitempty"`
	AllowPaidBroadcast  bool                  `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID     string                `json:"message_effect_id,omitempty"`
	ReplyParameters     *ReplyParameters      `json:"reply_parameters,omitempty"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// SendGame sends a game message.
// Since: Bot API 2.2
// See https://core.telegram.org/bots/api#sendgame
func (api *API) SendGame(params SendGame) (Message, error) {
	req := NewRequestWithChatID[Message]("sendGame", params, params.ChatID)
	return req.Do(api)
}

// SendGameWithContext is the context-aware variant of SendGame.
// Since: Bot API 2.2
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendgame
func (api *API) SendGameWithContext(ctx context.Context, params SendGame) (Message, error) {
	req := NewRequestWithChatID[Message]("sendGame", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SetGameScore holds parameters for the setGameScore method.
// Since: Bot API 2.2
// See https://core.telegram.org/bots/api#setgamescore
type SetGameScore struct {
	UserID             int64  `json:"user_id"`
	Score              int    `json:"score"`
	Force              bool   `json:"force,omitempty"`
	DisableEditMessage bool   `json:"disable_edit_message,omitempty"`
	ChatID             int64  `json:"chat_id,omitempty"`
	MessageID          int    `json:"message_id,omitempty"`
	InlineMessageID    string `json:"inline_message_id,omitempty"`
}

// SetGameScore sets a user's score in a game message.
// Since: Bot API 2.2
// If inline_message_id is provided, returns a boolean success flag.
// Otherwise returns the edited Message.
// See https://core.telegram.org/bots/api#setgamescore
func (api *API) SetGameScore(params SetGameScore) (Message, bool, error) {
	var zero Message
	if params.InlineMessageID != "" {
		req := NewRequestWithChatID[bool]("setGameScore", params, params.ChatID)
		res, err := req.Do(api)
		return zero, res, err
	}
	req := NewRequestWithChatID[Message]("setGameScore", params, params.ChatID)
	res, err := req.Do(api)
	return res, false, err
}

// SetGameScoreWithContext is the context-aware variant of SetGameScore.
// Since: Bot API 2.2
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#setgamescore
func (api *API) SetGameScoreWithContext(ctx context.Context, params SetGameScore) (Message, bool, error) {
	var zero Message
	if params.InlineMessageID != "" {
		req := NewRequestWithChatID[bool]("setGameScore", params, params.ChatID)
		res, err := req.DoWithContext(ctx, api)
		return zero, res, err
	}
	req := NewRequestWithChatID[Message]("setGameScore", params, params.ChatID)
	res, err := req.DoWithContext(ctx, api)
	return res, false, err
}

// GetGameHighScores holds parameters for the getGameHighScores method.
// Since: Bot API 2.2
// See https://core.telegram.org/bots/api#getgamehighscores
type GetGameHighScores struct {
	UserID          int64  `json:"user_id"`
	ChatID          int64  `json:"chat_id,omitempty"`
	MessageID       int    `json:"message_id,omitempty"`
	InlineMessageID string `json:"inline_message_id,omitempty"`
}

// GetGameHighScores returns game high score data for a user.
// Since: Bot API 2.2
// See https://core.telegram.org/bots/api#getgamehighscores
func (api *API) GetGameHighScores(params GetGameHighScores) ([]GameHighScore, error) {
	req := NewRequestWithChatID[[]GameHighScore]("getGameHighScores", params, params.ChatID)
	return req.Do(api)
}

// GetGameHighScoresWithContext is the context-aware variant of GetGameHighScores.
// Since: Bot API 2.2
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#getgamehighscores
func (api *API) GetGameHighScoresWithContext(ctx context.Context, params GetGameHighScores) ([]GameHighScore, error) {
	req := NewRequestWithChatID[[]GameHighScore]("getGameHighScores", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}
