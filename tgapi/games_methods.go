package tgapi

import "context"

// SendGameP holds parameters for the sendGame method.
// See https://core.telegram.org/bots/api#sendgame
type SendGameP struct {
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
// See https://core.telegram.org/bots/api#sendgame
func (api *API) SendGame(params SendGameP) (Message, error) {
	req := NewRequestWithChatID[Message]("sendGame", params, params.ChatID)
	return req.Do(api)
}

// SendGameWithContext is the context-aware variant of SendGame.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendgame
func (api *API) SendGameWithContext(ctx context.Context, params SendGameP) (Message, error) {
	req := NewRequestWithChatID[Message]("sendGame", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SetGameScoreP holds parameters for the setGameScore method.
// See https://core.telegram.org/bots/api#setgamescore
type SetGameScoreP struct {
	UserID             int64  `json:"user_id"`
	Score              int    `json:"score"`
	Force              bool   `json:"force,omitempty"`
	DisableEditMessage bool   `json:"disable_edit_message,omitempty"`
	ChatID             int64  `json:"chat_id,omitempty"`
	MessageID          int    `json:"message_id,omitempty"`
	InlineMessageID    string `json:"inline_message_id,omitempty"`
}

// SetGameScore sets a user's score in a game message.
// If inline_message_id is provided, returns a boolean success flag.
// Otherwise returns the edited Message.
// See https://core.telegram.org/bots/api#setgamescore
func (api *API) SetGameScore(params SetGameScoreP) (Message, bool, error) {
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
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#setgamescore
func (api *API) SetGameScoreWithContext(ctx context.Context, params SetGameScoreP) (Message, bool, error) {
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

// GetGameHighScoresP holds parameters for the getGameHighScores method.
// See https://core.telegram.org/bots/api#getgamehighscores
type GetGameHighScoresP struct {
	UserID          int64  `json:"user_id"`
	ChatID          int64  `json:"chat_id,omitempty"`
	MessageID       int    `json:"message_id,omitempty"`
	InlineMessageID string `json:"inline_message_id,omitempty"`
}

// GetGameHighScores returns game high score data for a user.
// See https://core.telegram.org/bots/api#getgamehighscores
func (api *API) GetGameHighScores(params GetGameHighScoresP) ([]GameHighScore, error) {
	req := NewRequestWithChatID[[]GameHighScore]("getGameHighScores", params, params.ChatID)
	return req.Do(api)
}

// GetGameHighScoresWithContext is the context-aware variant of GetGameHighScores.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#getgamehighscores
func (api *API) GetGameHighScoresWithContext(ctx context.Context, params GetGameHighScoresP) ([]GameHighScore, error) {
	req := NewRequestWithChatID[[]GameHighScore]("getGameHighScores", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}
