package tgapi

type BaseForumTopicP struct {
	ChatID          int `json:"chat_id"`
	MessageThreadID int `json:"message_thread_id"`
}

func (api *API) GetForumTopicIconSet() ([]Sticker, error) {
	req := NewRequest[[]Sticker]("getForumTopicIconSet", NoParams)
	return req.Do(api)
}

type CreateForumTopicP struct {
	ChatID            int                 `json:"chat_id"`
	Name              string              `json:"name"`
	IconColor         ForumTopicIconColor `json:"icon_color"`
	IconCustomEmojiID string              `json:"icon_custom_emoji_id"`
}

func (api *API) CreateForumTopic(params CreateForumTopicP) (ForumTopic, error) {
	req := NewRequest[ForumTopic]("createForumTopic", params)
	return req.Do(api)
}

type EditForumTopicP struct {
	BaseForumTopicP
	Name              string `json:"name"`
	IconCustomEmojiID string `json:"icon_custom_emoji_id"`
}

func (api *API) EditForumTopic(params EditForumTopicP) (bool, error) {
	req := NewRequest[bool]("editForumTopic", params)
	return req.Do(api)
}

func (api *API) CloseForumTopic(params BaseForumTopicP) (bool, error) {
	req := NewRequest[bool]("closeForumTopic", params)
	return req.Do(api)
}
func (api *API) ReopenForumTopic(params BaseForumTopicP) (bool, error) {
	req := NewRequest[bool]("reopenForumTopic", params)
	return req.Do(api)
}
func (api *API) DeleteForumTopic(params BaseForumTopicP) (bool, error) {
	req := NewRequest[bool]("deleteForumTopic", params)
	return req.Do(api)
}
func (api *API) UnpinAllForumTopicMessages(params BaseForumTopicP) (bool, error) {
	req := NewRequest[bool]("unpinAllForumTopicMessages", params)
	return req.Do(api)
}

type BaseGeneralForumTopicP struct {
	ChatID int `json:"chat_id"`
}

type EditGeneralForumTopicP struct {
	ChatID int    `json:"chat_id"`
	Name   string `json:"name"`
}

func (api *API) EditGeneralForumTopic(params EditGeneralForumTopicP) (bool, error) {
	req := NewRequest[bool]("editGeneralForumTopic", params)
	return req.Do(api)
}

func (api *API) CloseGeneralForumTopic(params BaseGeneralForumTopicP) (bool, error) {
	req := NewRequest[bool]("closeGeneralForumTopic", params)
	return req.Do(api)
}
func (api *API) ReopenGeneralForumTopic(params BaseGeneralForumTopicP) (bool, error) {
	req := NewRequest[bool]("reopenGeneralForumTopic", params)
	return req.Do(api)
}
func (api *API) HideGeneralForumTopic(params BaseGeneralForumTopicP) (bool, error) {
	req := NewRequest[bool]("hideGeneralForumTopic", params)
	return req.Do(api)
}
func (api *API) UnhideGeneralForumTopic(params BaseGeneralForumTopicP) (bool, error) {
	req := NewRequest[bool]("unhideGeneralForumTopic", params)
	return req.Do(api)
}
func (api *API) UnpinAllGeneralForumTopicMessages(params BaseGeneralForumTopicP) (bool, error) {
	req := NewRequest[bool]("unpinAllGeneralForumTopicMessages", params)
	return req.Do(api)
}
