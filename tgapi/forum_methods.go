package tgapi

import "context"

// BaseForumTopic contains common fields for forum topic operations that require a chat ID and a message thread ID.
// Since: Bot API 6.3
type BaseForumTopic struct {
	ChatID          int64 `json:"chat_id"`
	MessageThreadID int   `json:"message_thread_id"`
}

// GetForumTopicIconStickers returns the list of custom emoji that can be used as a forum topic icon.
// Since: Bot API 6.3
// See https://core.telegram.org/bots/api#getforumtopiciconstickers
func (api *API) GetForumTopicIconStickers() ([]Sticker, error) {
	req := NewRequest[[]Sticker]("getForumTopicIconStickers", NoParams)
	return req.Do(api)
}

// GetForumTopicIconStickersWithContext is the context-aware variant of GetForumTopicIconStickers.
// Since: Bot API 6.3
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#getforumtopiciconstickers
func (api *API) GetForumTopicIconStickersWithContext(ctx context.Context) ([]Sticker, error) {
	req := NewRequest[[]Sticker]("getForumTopicIconStickers", NoParams)
	return req.DoWithContext(ctx, api)
}

// CreateForumTopic holds parameters for the createForumTopic method.
// Since: Bot API 6.3
// See https://core.telegram.org/bots/api#createforumtopic
type CreateForumTopic struct {
	ChatID            int64               `json:"chat_id"`
	Name              string              `json:"name"`
	IconColor         ForumTopicIconColor `json:"icon_color"`
	IconCustomEmojiID string              `json:"icon_custom_emoji_id"`
}

// CreateForumTopic creates a topic in a forum supergroup.
// Since: Bot API 6.3
// Returns the created ForumTopic on success.
// See https://core.telegram.org/bots/api#createforumtopic
func (api *API) CreateForumTopic(params CreateForumTopic) (ForumTopic, error) {
	req := NewRequestWithChatID[ForumTopic]("createForumTopic", params, params.ChatID)
	return req.Do(api)
}

// CreateForumTopicWithContext is the context-aware variant of CreateForumTopic.
// Since: Bot API 6.3
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#createforumtopic
func (api *API) CreateForumTopicWithContext(ctx context.Context, params CreateForumTopic) (ForumTopic, error) {
	req := NewRequestWithChatID[ForumTopic]("createForumTopic", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// EditForumTopic holds parameters for the editForumTopic method.
// Since: Bot API 6.3
// See https://core.telegram.org/bots/api#editforumtopic
type EditForumTopic struct {
	BaseForumTopic
	Name              string `json:"name"`
	IconCustomEmojiID string `json:"icon_custom_emoji_id"`
}

// EditForumTopic edits name and icon of a forum topic.
// Since: Bot API 6.3
// Returns True on success.
// See https://core.telegram.org/bots/api#editforumtopic
func (api *API) EditForumTopic(params EditForumTopic) (bool, error) {
	req := NewRequestWithChatID[bool]("editForumTopic", params, params.ChatID)
	return req.Do(api)
}

// EditForumTopicWithContext is the context-aware variant of EditForumTopic.
// Since: Bot API 6.3
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#editforumtopic
func (api *API) EditForumTopicWithContext(ctx context.Context, params EditForumTopic) (bool, error) {
	req := NewRequestWithChatID[bool]("editForumTopic", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// CloseForumTopic closes an open forum topic.
// Since: Bot API 6.3
// Returns True on success.
// See https://core.telegram.org/bots/api#closeforumtopic
func (api *API) CloseForumTopic(params BaseForumTopic) (bool, error) {
	req := NewRequestWithChatID[bool]("closeForumTopic", params, params.ChatID)
	return req.Do(api)
}

// CloseForumTopicWithContext is the context-aware variant of CloseForumTopic.
// Since: Bot API 6.3
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#closeforumtopic
func (api *API) CloseForumTopicWithContext(ctx context.Context, params BaseForumTopic) (bool, error) {
	req := NewRequestWithChatID[bool]("closeForumTopic", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// ReopenForumTopic reopens a closed forum topic.
// Since: Bot API 6.3
// Returns True on success.
// See https://core.telegram.org/bots/api#reopenforumtopic
func (api *API) ReopenForumTopic(params BaseForumTopic) (bool, error) {
	req := NewRequestWithChatID[bool]("reopenForumTopic", params, params.ChatID)
	return req.Do(api)
}

// ReopenForumTopicWithContext is the context-aware variant of ReopenForumTopic.
// Since: Bot API 6.3
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#reopenforumtopic
func (api *API) ReopenForumTopicWithContext(ctx context.Context, params BaseForumTopic) (bool, error) {
	req := NewRequestWithChatID[bool]("reopenForumTopic", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// DeleteForumTopic deletes a forum topic.
// Since: Bot API 6.3
// Returns True on success.
// See https://core.telegram.org/bots/api#deleteforumtopic
func (api *API) DeleteForumTopic(params BaseForumTopic) (bool, error) {
	req := NewRequestWithChatID[bool]("deleteForumTopic", params, params.ChatID)
	return req.Do(api)
}

// DeleteForumTopicWithContext is the context-aware variant of DeleteForumTopic.
// Since: Bot API 6.3
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#deleteforumtopic
func (api *API) DeleteForumTopicWithContext(ctx context.Context, params BaseForumTopic) (bool, error) {
	req := NewRequestWithChatID[bool]("deleteForumTopic", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// UnpinAllForumTopicMessages clears the list of pinned messages in a forum topic.
// Since: Bot API 6.3
// Returns True on success.
// See https://core.telegram.org/bots/api#unpinallforumtopicmessages
func (api *API) UnpinAllForumTopicMessages(params BaseForumTopic) (bool, error) {
	req := NewRequestWithChatID[bool]("unpinAllForumTopicMessages", params, params.ChatID)
	return req.Do(api)
}

// UnpinAllForumTopicMessagesWithContext is the context-aware variant of UnpinAllForumTopicMessages.
// Since: Bot API 6.3
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#unpinallforumtopicmessages
func (api *API) UnpinAllForumTopicMessagesWithContext(ctx context.Context, params BaseForumTopic) (bool, error) {
	req := NewRequestWithChatID[bool]("unpinAllForumTopicMessages", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// BaseGeneralForumTopic contains common fields for general forum topic operations that require a chat ID.
// Since: Bot API 6.4
type BaseGeneralForumTopic struct {
	ChatID int64 `json:"chat_id"`
}

// EditGeneralForumTopic holds parameters for the editGeneralForumTopic method.
// Since: Bot API 6.4
// See https://core.telegram.org/bots/api#editgeneralforumtopic
type EditGeneralForumTopic struct {
	ChatID int64  `json:"chat_id"`
	Name   string `json:"name"`
}

// EditGeneralForumTopic edits the name of the 'General' topic in a forum supergroup.
// Since: Bot API 6.4
// Returns True on success.
// See https://core.telegram.org/bots/api#editgeneralforumtopic
func (api *API) EditGeneralForumTopic(params EditGeneralForumTopic) (bool, error) {
	req := NewRequestWithChatID[bool]("editGeneralForumTopic", params, params.ChatID)
	return req.Do(api)
}

// EditGeneralForumTopicWithContext is the context-aware variant of EditGeneralForumTopic.
// Since: Bot API 6.4
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#editgeneralforumtopic
func (api *API) EditGeneralForumTopicWithContext(ctx context.Context, params EditGeneralForumTopic) (bool, error) {
	req := NewRequestWithChatID[bool]("editGeneralForumTopic", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// CloseGeneralForumTopic closes the 'General' topic in a forum supergroup.
// Since: Bot API 6.4
// Returns True on success.
// See https://core.telegram.org/bots/api#closegeneralforumtopic
func (api *API) CloseGeneralForumTopic(params BaseGeneralForumTopic) (bool, error) {
	req := NewRequestWithChatID[bool]("closeGeneralForumTopic", params, params.ChatID)
	return req.Do(api)
}

// CloseGeneralForumTopicWithContext is the context-aware variant of CloseGeneralForumTopic.
// Since: Bot API 6.4
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#closegeneralforumtopic
func (api *API) CloseGeneralForumTopicWithContext(ctx context.Context, params BaseGeneralForumTopic) (bool, error) {
	req := NewRequestWithChatID[bool]("closeGeneralForumTopic", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// ReopenGeneralForumTopic reopens the 'General' topic in a forum supergroup.
// Since: Bot API 6.4
// Returns True on success.
// See https://core.telegram.org/bots/api#reopengeneralforumtopic
func (api *API) ReopenGeneralForumTopic(params BaseGeneralForumTopic) (bool, error) {
	req := NewRequestWithChatID[bool]("reopenGeneralForumTopic", params, params.ChatID)
	return req.Do(api)
}

// ReopenGeneralForumTopicWithContext is the context-aware variant of ReopenGeneralForumTopic.
// Since: Bot API 6.4
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#reopengeneralforumtopic
func (api *API) ReopenGeneralForumTopicWithContext(ctx context.Context, params BaseGeneralForumTopic) (bool, error) {
	req := NewRequestWithChatID[bool]("reopenGeneralForumTopic", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// HideGeneralForumTopic hides the 'General' topic in a forum supergroup.
// Since: Bot API 6.4
// Returns True on success.
// See https://core.telegram.org/bots/api#hidegeneralforumtopic
func (api *API) HideGeneralForumTopic(params BaseGeneralForumTopic) (bool, error) {
	req := NewRequestWithChatID[bool]("hideGeneralForumTopic", params, params.ChatID)
	return req.Do(api)
}

// HideGeneralForumTopicWithContext is the context-aware variant of HideGeneralForumTopic.
// Since: Bot API 6.4
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#hidegeneralforumtopic
func (api *API) HideGeneralForumTopicWithContext(ctx context.Context, params BaseGeneralForumTopic) (bool, error) {
	req := NewRequestWithChatID[bool]("hideGeneralForumTopic", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// UnhideGeneralForumTopic unhides the 'General' topic in a forum supergroup.
// Since: Bot API 6.4
// Returns True on success.
// See https://core.telegram.org/bots/api#unhidegeneralforumtopic
func (api *API) UnhideGeneralForumTopic(params BaseGeneralForumTopic) (bool, error) {
	req := NewRequestWithChatID[bool]("unhideGeneralForumTopic", params, params.ChatID)
	return req.Do(api)
}

// UnhideGeneralForumTopicWithContext is the context-aware variant of UnhideGeneralForumTopic.
// Since: Bot API 6.4
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#unhidegeneralforumtopic
func (api *API) UnhideGeneralForumTopicWithContext(ctx context.Context, params BaseGeneralForumTopic) (bool, error) {
	req := NewRequestWithChatID[bool]("unhideGeneralForumTopic", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// UnpinAllGeneralForumTopicMessages clears the list of pinned messages in the 'General' topic.
// Since: Bot API 6.4
// Returns True on success.
// See https://core.telegram.org/bots/api#unpinallgeneralforumtopicmessages
func (api *API) UnpinAllGeneralForumTopicMessages(params BaseGeneralForumTopic) (bool, error) {
	req := NewRequestWithChatID[bool]("unpinAllGeneralForumTopicMessages", params, params.ChatID)
	return req.Do(api)
}

// UnpinAllGeneralForumTopicMessagesWithContext is the context-aware variant of UnpinAllGeneralForumTopicMessages.
// Since: Bot API 6.4
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#unpinallgeneralforumtopicmessages
func (api *API) UnpinAllGeneralForumTopicMessagesWithContext(ctx context.Context, params BaseGeneralForumTopic) (bool, error) {
	req := NewRequestWithChatID[bool]("unpinAllGeneralForumTopicMessages", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}
