package tgapi

type BaseForumTopicP struct {
	ChatID          int `json:"chat_id"`
	MessageThreadID int `json:"message_thread_id"`
}

func (api *Api) GetForumTopicIconSet() ([]Sticker, error) {
	req := NewRequest[[]Sticker]("getForumTopicIconSet", NoParams)
	return req.Do(api)
}

type CreateForumTopicP struct {
	ChatID            int                 `json:"chat_id"`
	Name              string              `json:"name"`
	IconColor         ForumTopicIconColor `json:"icon_color"`
	IconCustomEmojiID string              `json:"icon_custom_emoji_id"`
}

func (api *Api) CreateForumTopic(params CreateForumTopicP) (ForumTopic, error) {
	req := NewRequest[ForumTopic]("createForumTopic", params)
	return req.Do(api)
}

type EditForumTopicP struct {
	BaseForumTopicP
	Name              string `json:"name"`
	IconCustomEmojiID string `json:"icon_custom_emoji_id"`
}

func (api *Api) EditForumTopic(params EditForumTopicP) (bool, error) {
	req := NewRequest[bool]("editForumTopic", params)
	return req.Do(api)
}

func (api *Api) CloseForumTopic(params BaseForumTopicP) (bool, error) {
	req := NewRequest[bool]("closeForumTopic", params)
	return req.Do(api)
}
func (api *Api) ReopenForumTopic(params BaseForumTopicP) (bool, error) {
	req := NewRequest[bool]("reopenForumTopic", params)
	return req.Do(api)
}
func (api *Api) DeleteForumTopic(params BaseForumTopicP) (bool, error) {
	req := NewRequest[bool]("deleteForumTopic", params)
	return req.Do(api)
}
func (api *Api) UnpinAllForumTopicMessages(params BaseForumTopicP) (bool, error) {
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

func (api *Api) EditGeneralForumTopic(params EditGeneralForumTopicP) (bool, error) {
	req := NewRequest[bool]("editGeneralForumTopic", params)
	return req.Do(api)
}

func (api *Api) CloseGeneralForumTopic(params BaseGeneralForumTopicP) (bool, error) {
	req := NewRequest[bool]("closeGeneralForumTopic", params)
	return req.Do(api)
}
func (api *Api) ReopenGeneralForumTopic(params BaseGeneralForumTopicP) (bool, error) {
	req := NewRequest[bool]("reopenGeneralForumTopic", params)
	return req.Do(api)
}
func (api *Api) HideGeneralForumTopic(params BaseGeneralForumTopicP) (bool, error) {
	req := NewRequest[bool]("hideGeneralForumTopic", params)
	return req.Do(api)
}
func (api *Api) UnhideGeneralForumTopic(params BaseGeneralForumTopicP) (bool, error) {
	req := NewRequest[bool]("unhideGeneralForumTopic", params)
	return req.Do(api)
}
func (api *Api) UnpinAllGeneralForumTopicMessages(params BaseGeneralForumTopicP) (bool, error) {
	req := NewRequest[bool]("unpinAllGeneralForumTopicMessages", params)
	return req.Do(api)
}
