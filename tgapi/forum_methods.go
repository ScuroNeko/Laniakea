package tgapi

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

// CloseForumTopic closes an open forum topic.
// Returns True on success.
// See https://core.telegram.org/bots/api#closeforumtopic
func (api *API) CloseForumTopic(params BaseForumTopicP) (bool, error) {
	req := NewRequestWithChatID[bool]("closeForumTopic", params, params.ChatID)
	return req.Do(api)
}

// ReopenForumTopic reopens a closed forum topic.
// Returns True on success.
// See https://core.telegram.org/bots/api#reopenforumtopic
func (api *API) ReopenForumTopic(params BaseForumTopicP) (bool, error) {
	req := NewRequestWithChatID[bool]("reopenForumTopic", params, params.ChatID)
	return req.Do(api)
}

// DeleteForumTopic deletes a forum topic.
// Returns True on success.
// See https://core.telegram.org/bots/api#deleteforumtopic
func (api *API) DeleteForumTopic(params BaseForumTopicP) (bool, error) {
	req := NewRequestWithChatID[bool]("deleteForumTopic", params, params.ChatID)
	return req.Do(api)
}

// UnpinAllForumTopicMessages clears the list of pinned messages in a forum topic.
// Returns True on success.
// See https://core.telegram.org/bots/api#unpinallforumtopicmessages
func (api *API) UnpinAllForumTopicMessages(params BaseForumTopicP) (bool, error) {
	req := NewRequestWithChatID[bool]("unpinAllForumTopicMessages", params, params.ChatID)
	return req.Do(api)
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

// CloseGeneralForumTopic closes the 'General' topic in a forum supergroup.
// Returns True on success.
// See https://core.telegram.org/bots/api#closegeneralforumtopic
func (api *API) CloseGeneralForumTopic(params BaseGeneralForumTopicP) (bool, error) {
	req := NewRequestWithChatID[bool]("closeGeneralForumTopic", params, params.ChatID)
	return req.Do(api)
}

// ReopenGeneralForumTopic reopens the 'General' topic in a forum supergroup.
// Returns True on success.
// See https://core.telegram.org/bots/api#reopengeneralforumtopic
func (api *API) ReopenGeneralForumTopic(params BaseGeneralForumTopicP) (bool, error) {
	req := NewRequestWithChatID[bool]("reopenGeneralForumTopic", params, params.ChatID)
	return req.Do(api)
}

// HideGeneralForumTopic hides the 'General' topic in a forum supergroup.
// Returns True on success.
// See https://core.telegram.org/bots/api#hidegeneralforumtopic
func (api *API) HideGeneralForumTopic(params BaseGeneralForumTopicP) (bool, error) {
	req := NewRequestWithChatID[bool]("hideGeneralForumTopic", params, params.ChatID)
	return req.Do(api)
}

// UnhideGeneralForumTopic unhides the 'General' topic in a forum supergroup.
// Returns True on success.
// See https://core.telegram.org/bots/api#unhidegeneralforumtopic
func (api *API) UnhideGeneralForumTopic(params BaseGeneralForumTopicP) (bool, error) {
	req := NewRequestWithChatID[bool]("unhideGeneralForumTopic", params, params.ChatID)
	return req.Do(api)
}

// UnpinAllGeneralForumTopicMessages clears the list of pinned messages in the 'General' topic.
// Returns True on success.
// See https://core.telegram.org/bots/api#unpinallgeneralforumtopicmessages
func (api *API) UnpinAllGeneralForumTopicMessages(params BaseGeneralForumTopicP) (bool, error) {
	req := NewRequestWithChatID[bool]("unpinAllGeneralForumTopicMessages", params, params.ChatID)
	return req.Do(api)
}
