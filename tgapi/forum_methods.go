package tgapi

type BaseForumTopicP struct {
	ChatID          int64 `json:"chat_id"`
	MessageThreadID int   `json:"message_thread_id"`
}

func (api *API) GetForumTopicIconStickers() ([]Sticker, error) {
	req := NewRequest[[]Sticker]("getForumTopicIconStickers", NoParams)
	return req.Do(api)
}

type CreateForumTopicP struct {
	ChatID            int64               `json:"chat_id"`
	Name              string              `json:"name"`
	IconColor         ForumTopicIconColor `json:"icon_color"`
	IconCustomEmojiID string              `json:"icon_custom_emoji_id"`
}

func (api *API) CreateForumTopic(params CreateForumTopicP) (ForumTopic, error) {
	req := NewRequestWithChatID[ForumTopic]("createForumTopic", params, params.ChatID)
	return req.Do(api)
}

type EditForumTopicP struct {
	BaseForumTopicP
	Name              string `json:"name"`
	IconCustomEmojiID string `json:"icon_custom_emoji_id"`
}

func (api *API) EditForumTopic(params EditForumTopicP) (bool, error) {
	req := NewRequestWithChatID[bool]("editForumTopic", params, params.ChatID)
	return req.Do(api)
}

func (api *API) CloseForumTopic(params BaseForumTopicP) (bool, error) {
	req := NewRequestWithChatID[bool]("closeForumTopic", params, params.ChatID)
	return req.Do(api)
}
func (api *API) ReopenForumTopic(params BaseForumTopicP) (bool, error) {
	req := NewRequestWithChatID[bool]("reopenForumTopic", params, params.ChatID)
	return req.Do(api)
}
func (api *API) DeleteForumTopic(params BaseForumTopicP) (bool, error) {
	req := NewRequestWithChatID[bool]("deleteForumTopic", params, params.ChatID)
	return req.Do(api)
}
func (api *API) UnpinAllForumTopicMessages(params BaseForumTopicP) (bool, error) {
	req := NewRequestWithChatID[bool]("unpinAllForumTopicMessages", params, params.ChatID)
	return req.Do(api)
}

type BaseGeneralForumTopicP struct {
	ChatID int64 `json:"chat_id"`
}

type EditGeneralForumTopicP struct {
	ChatID int64  `json:"chat_id"`
	Name   string `json:"name"`
}

func (api *API) EditGeneralForumTopic(params EditGeneralForumTopicP) (bool, error) {
	req := NewRequestWithChatID[bool]("editGeneralForumTopic", params, params.ChatID)
	return req.Do(api)
}

func (api *API) CloseGeneralForumTopic(params BaseGeneralForumTopicP) (bool, error) {
	req := NewRequestWithChatID[bool]("closeGeneralForumTopic", params, params.ChatID)
	return req.Do(api)
}
func (api *API) ReopenGeneralForumTopic(params BaseGeneralForumTopicP) (bool, error) {
	req := NewRequestWithChatID[bool]("reopenGeneralForumTopic", params, params.ChatID)
	return req.Do(api)
}
func (api *API) HideGeneralForumTopic(params BaseGeneralForumTopicP) (bool, error) {
	req := NewRequestWithChatID[bool]("hideGeneralForumTopic", params, params.ChatID)
	return req.Do(api)
}
func (api *API) UnhideGeneralForumTopic(params BaseGeneralForumTopicP) (bool, error) {
	req := NewRequestWithChatID[bool]("unhideGeneralForumTopic", params, params.ChatID)
	return req.Do(api)
}
func (api *API) UnpinAllGeneralForumTopicMessages(params BaseGeneralForumTopicP) (bool, error) {
	req := NewRequestWithChatID[bool]("unpinAllGeneralForumTopicMessages", params, params.ChatID)
	return req.Do(api)
}
