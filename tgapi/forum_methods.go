package tgapi

import "context"

// BaseForumTopicP contains common fields for forum topic operations that require a chat ID and a message thread ID.
type BaseForumTopicP struct {
	ChatID          int64 `json:"chat_id"`
	MessageThreadID int   `json:"message_thread_id"`
}

// GetForumTopicIconStickers returns the list of custom emoji that can be used as a forum topic icon.
// See https://core.telegram.org/bots/api#getforumtopiciconstickers
func (api *API) GetForumTopicIconStickers() ([]Sticker, error) {
	req := NewRequest[[]Sticker]("getForumTopicIconStickers", NoParams)
	return req.Do(api)
}

// GetForumTopicIconStickersWithContext is the context-aware variant of GetForumTopicIconStickers.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#getforumtopiciconstickers
func (api *API) GetForumTopicIconStickersWithContext(ctx context.Context) ([]Sticker, error) {
	req := NewRequest[[]Sticker]("getForumTopicIconStickers", NoParams)
	return req.DoWithContext(ctx, api)
}

// CreateForumTopicP holds parameters for the createForumTopic method.
// See https://core.telegram.org/bots/api#createforumtopic
type CreateForumTopicP struct {
	ChatID            int64               `json:"chat_id"`
	Name              string              `json:"name"`
	IconColor         ForumTopicIconColor `json:"icon_color"`
	IconCustomEmojiID string              `json:"icon_custom_emoji_id"`
}

// CreateForumTopic creates a topic in a forum supergroup.
// Returns the created ForumTopic on success.
// See https://core.telegram.org/bots/api#createforumtopic
func (api *API) CreateForumTopic(params CreateForumTopicP) (ForumTopic, error) {
	req := NewRequestWithChatID[ForumTopic]("createForumTopic", params, params.ChatID)
	return req.Do(api)
}

// CreateForumTopicWithContext is the context-aware variant of CreateForumTopic.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#createforumtopic
func (api *API) CreateForumTopicWithContext(ctx context.Context, params CreateForumTopicP) (ForumTopic, error) {
	req := NewRequestWithChatID[ForumTopic]("createForumTopic", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// EditForumTopicP holds parameters for the editForumTopic method.
// See https://core.telegram.org/bots/api#editforumtopic
type EditForumTopicP struct {
	BaseForumTopicP
	Name              string `json:"name"`
	IconCustomEmojiID string `json:"icon_custom_emoji_id"`
}

// EditForumTopic edits name and icon of a forum topic.
// Returns True on success.
// See https://core.telegram.org/bots/api#editforumtopic
func (api *API) EditForumTopic(params EditForumTopicP) (bool, error) {
	req := NewRequestWithChatID[bool]("editForumTopic", params, params.ChatID)
	return req.Do(api)
}

// EditForumTopicWithContext is the context-aware variant of EditForumTopic.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#editforumtopic
func (api *API) EditForumTopicWithContext(ctx context.Context, params EditForumTopicP) (bool, error) {
	req := NewRequestWithChatID[bool]("editForumTopic", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// CloseForumTopic closes an open forum topic.
// Returns True on success.
// See https://core.telegram.org/bots/api#closeforumtopic
func (api *API) CloseForumTopic(params BaseForumTopicP) (bool, error) {
	req := NewRequestWithChatID[bool]("closeForumTopic", params, params.ChatID)
	return req.Do(api)
}

// CloseForumTopicWithContext is the context-aware variant of CloseForumTopic.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#closeforumtopic
func (api *API) CloseForumTopicWithContext(ctx context.Context, params BaseForumTopicP) (bool, error) {
	req := NewRequestWithChatID[bool]("closeForumTopic", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// ReopenForumTopic reopens a closed forum topic.
// Returns True on success.
// See https://core.telegram.org/bots/api#reopenforumtopic
func (api *API) ReopenForumTopic(params BaseForumTopicP) (bool, error) {
	req := NewRequestWithChatID[bool]("reopenForumTopic", params, params.ChatID)
	return req.Do(api)
}

// ReopenForumTopicWithContext is the context-aware variant of ReopenForumTopic.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#reopenforumtopic
func (api *API) ReopenForumTopicWithContext(ctx context.Context, params BaseForumTopicP) (bool, error) {
	req := NewRequestWithChatID[bool]("reopenForumTopic", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// DeleteForumTopic deletes a forum topic.
// Returns True on success.
// See https://core.telegram.org/bots/api#deleteforumtopic
func (api *API) DeleteForumTopic(params BaseForumTopicP) (bool, error) {
	req := NewRequestWithChatID[bool]("deleteForumTopic", params, params.ChatID)
	return req.Do(api)
}

// DeleteForumTopicWithContext is the context-aware variant of DeleteForumTopic.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#deleteforumtopic
func (api *API) DeleteForumTopicWithContext(ctx context.Context, params BaseForumTopicP) (bool, error) {
	req := NewRequestWithChatID[bool]("deleteForumTopic", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// UnpinAllForumTopicMessages clears the list of pinned messages in a forum topic.
// Returns True on success.
// See https://core.telegram.org/bots/api#unpinallforumtopicmessages
func (api *API) UnpinAllForumTopicMessages(params BaseForumTopicP) (bool, error) {
	req := NewRequestWithChatID[bool]("unpinAllForumTopicMessages", params, params.ChatID)
	return req.Do(api)
}

// UnpinAllForumTopicMessagesWithContext is the context-aware variant of UnpinAllForumTopicMessages.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#unpinallforumtopicmessages
func (api *API) UnpinAllForumTopicMessagesWithContext(ctx context.Context, params BaseForumTopicP) (bool, error) {
	req := NewRequestWithChatID[bool]("unpinAllForumTopicMessages", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// BaseGeneralForumTopicP contains common fields for general forum topic operations that require a chat ID.
type BaseGeneralForumTopicP struct {
	ChatID int64 `json:"chat_id"`
}

// EditGeneralForumTopicP holds parameters for the editGeneralForumTopic method.
// See https://core.telegram.org/bots/api#editgeneralforumtopic
type EditGeneralForumTopicP struct {
	ChatID int64  `json:"chat_id"`
	Name   string `json:"name"`
}

// EditGeneralForumTopic edits the name of the 'General' topic in a forum supergroup.
// Returns True on success.
// See https://core.telegram.org/bots/api#editgeneralforumtopic
func (api *API) EditGeneralForumTopic(params EditGeneralForumTopicP) (bool, error) {
	req := NewRequestWithChatID[bool]("editGeneralForumTopic", params, params.ChatID)
	return req.Do(api)
}

// EditGeneralForumTopicWithContext is the context-aware variant of EditGeneralForumTopic.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#editgeneralforumtopic
func (api *API) EditGeneralForumTopicWithContext(ctx context.Context, params EditGeneralForumTopicP) (bool, error) {
	req := NewRequestWithChatID[bool]("editGeneralForumTopic", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// CloseGeneralForumTopic closes the 'General' topic in a forum supergroup.
// Returns True on success.
// See https://core.telegram.org/bots/api#closegeneralforumtopic
func (api *API) CloseGeneralForumTopic(params BaseGeneralForumTopicP) (bool, error) {
	req := NewRequestWithChatID[bool]("closeGeneralForumTopic", params, params.ChatID)
	return req.Do(api)
}

// CloseGeneralForumTopicWithContext is the context-aware variant of CloseGeneralForumTopic.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#closegeneralforumtopic
func (api *API) CloseGeneralForumTopicWithContext(ctx context.Context, params BaseGeneralForumTopicP) (bool, error) {
	req := NewRequestWithChatID[bool]("closeGeneralForumTopic", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// ReopenGeneralForumTopic reopens the 'General' topic in a forum supergroup.
// Returns True on success.
// See https://core.telegram.org/bots/api#reopengeneralforumtopic
func (api *API) ReopenGeneralForumTopic(params BaseGeneralForumTopicP) (bool, error) {
	req := NewRequestWithChatID[bool]("reopenGeneralForumTopic", params, params.ChatID)
	return req.Do(api)
}

// ReopenGeneralForumTopicWithContext is the context-aware variant of ReopenGeneralForumTopic.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#reopengeneralforumtopic
func (api *API) ReopenGeneralForumTopicWithContext(ctx context.Context, params BaseGeneralForumTopicP) (bool, error) {
	req := NewRequestWithChatID[bool]("reopenGeneralForumTopic", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// HideGeneralForumTopic hides the 'General' topic in a forum supergroup.
// Returns True on success.
// See https://core.telegram.org/bots/api#hidegeneralforumtopic
func (api *API) HideGeneralForumTopic(params BaseGeneralForumTopicP) (bool, error) {
	req := NewRequestWithChatID[bool]("hideGeneralForumTopic", params, params.ChatID)
	return req.Do(api)
}

// HideGeneralForumTopicWithContext is the context-aware variant of HideGeneralForumTopic.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#hidegeneralforumtopic
func (api *API) HideGeneralForumTopicWithContext(ctx context.Context, params BaseGeneralForumTopicP) (bool, error) {
	req := NewRequestWithChatID[bool]("hideGeneralForumTopic", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// UnhideGeneralForumTopic unhides the 'General' topic in a forum supergroup.
// Returns True on success.
// See https://core.telegram.org/bots/api#unhidegeneralforumtopic
func (api *API) UnhideGeneralForumTopic(params BaseGeneralForumTopicP) (bool, error) {
	req := NewRequestWithChatID[bool]("unhideGeneralForumTopic", params, params.ChatID)
	return req.Do(api)
}

// UnhideGeneralForumTopicWithContext is the context-aware variant of UnhideGeneralForumTopic.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#unhidegeneralforumtopic
func (api *API) UnhideGeneralForumTopicWithContext(ctx context.Context, params BaseGeneralForumTopicP) (bool, error) {
	req := NewRequestWithChatID[bool]("unhideGeneralForumTopic", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// UnpinAllGeneralForumTopicMessages clears the list of pinned messages in the 'General' topic.
// Returns True on success.
// See https://core.telegram.org/bots/api#unpinallgeneralforumtopicmessages
func (api *API) UnpinAllGeneralForumTopicMessages(params BaseGeneralForumTopicP) (bool, error) {
	req := NewRequestWithChatID[bool]("unpinAllGeneralForumTopicMessages", params, params.ChatID)
	return req.Do(api)
}

// UnpinAllGeneralForumTopicMessagesWithContext is the context-aware variant of UnpinAllGeneralForumTopicMessages.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#unpinallgeneralforumtopicmessages
func (api *API) UnpinAllGeneralForumTopicMessagesWithContext(ctx context.Context, params BaseGeneralForumTopicP) (bool, error) {
	req := NewRequestWithChatID[bool]("unpinAllGeneralForumTopicMessages", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}
