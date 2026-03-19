package tgapi

import "context"

// AnswerInlineQueryP holds parameters for the answerInlineQuery method.
// See https://core.telegram.org/bots/api#answerinlinequery
type AnswerInlineQueryP struct {
	InlineQueryID string                    `json:"inline_query_id"`
	Results       []InlineQueryResult       `json:"results"`
	CacheTime     int                       `json:"cache_time,omitempty"`
	IsPersonal    bool                      `json:"is_personal,omitempty"`
	NextOffset    string                    `json:"next_offset,omitempty"`
	Button        *InlineQueryResultsButton `json:"button,omitempty"`
}

// AnswerInlineQuery sends answers to an inline query.
// Returns true on success.
// See https://core.telegram.org/bots/api#answerinlinequery
func (api *API) AnswerInlineQuery(params AnswerInlineQueryP) (bool, error) {
	req := NewRequest[bool]("answerInlineQuery", params)
	return req.Do(api)
}

// AnswerInlineQueryWithContext is the context-aware variant of AnswerInlineQuery.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#answerinlinequery
func (api *API) AnswerInlineQueryWithContext(ctx context.Context, params AnswerInlineQueryP) (bool, error) {
	req := NewRequest[bool]("answerInlineQuery", params)
	return req.DoWithContext(ctx, api)
}

// AnswerWebAppQueryP holds parameters for the answerWebAppQuery method.
// See https://core.telegram.org/bots/api#answerwebappquery
type AnswerWebAppQueryP struct {
	WebAppQueryID string            `json:"web_app_query_id"`
	Result        InlineQueryResult `json:"result"`
}

// AnswerWebAppQuery sets the result of a Web App interaction.
// See https://core.telegram.org/bots/api#answerwebappquery
func (api *API) AnswerWebAppQuery(params AnswerWebAppQueryP) (SentWebAppMessage, error) {
	req := NewRequest[SentWebAppMessage]("answerWebAppQuery", params)
	return req.Do(api)
}

// AnswerWebAppQueryWithContext is the context-aware variant of AnswerWebAppQuery.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#answerwebappquery
func (api *API) AnswerWebAppQueryWithContext(ctx context.Context, params AnswerWebAppQueryP) (SentWebAppMessage, error) {
	req := NewRequest[SentWebAppMessage]("answerWebAppQuery", params)
	return req.DoWithContext(ctx, api)
}

// SavePreparedInlineMessageP holds parameters for the savePreparedInlineMessage method.
// See https://core.telegram.org/bots/api#savepreparedinlinemessage
type SavePreparedInlineMessageP struct {
	UserID            int64             `json:"user_id"`
	Result            InlineQueryResult `json:"result"`
	AllowUserChats    bool              `json:"allow_user_chats,omitempty"`
	AllowBotChats     bool              `json:"allow_bot_chats,omitempty"`
	AllowGroupChats   bool              `json:"allow_group_chats,omitempty"`
	AllowChannelChats bool              `json:"allow_channel_chats,omitempty"`
}

// SavePreparedInlineMessage stores a prepared message for Mini App users.
// See https://core.telegram.org/bots/api#savepreparedinlinemessage
func (api *API) SavePreparedInlineMessage(params SavePreparedInlineMessageP) (PreparedInlineMessage, error) {
	req := NewRequest[PreparedInlineMessage]("savePreparedInlineMessage", params)
	return req.Do(api)
}

// SavePreparedInlineMessageWithContext is the context-aware variant of SavePreparedInlineMessage.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#savepreparedinlinemessage
func (api *API) SavePreparedInlineMessageWithContext(ctx context.Context, params SavePreparedInlineMessageP) (PreparedInlineMessage, error) {
	req := NewRequest[PreparedInlineMessage]("savePreparedInlineMessage", params)
	return req.DoWithContext(ctx, api)
}
