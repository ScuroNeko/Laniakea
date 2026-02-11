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

func (api *Api) CreateForumTopic(p CreateForumTopicP) (ForumTopic, error) {
	req := NewRequest[ForumTopic]("createForumTopic", p)
	return req.Do(api)
}

type EditForumTopicP struct {
	BaseForumTopicP
	Name              string `json:"name"`
	IconCustomEmojiID string `json:"icon_custom_emoji_id"`
}

func (api *Api) EditForumTopic(p EditForumTopicP) (bool, error) {
	req := NewRequest[bool]("editForumTopic", p)
	return req.Do(api)
}

func (api *Api) CloseForumTopic(p BaseForumTopicP) (bool, error) {
	req := NewRequest[bool]("closeForumTopic", p)
	return req.Do(api)
}
func (api *Api) ReopenForumTopic(p BaseForumTopicP) (bool, error) {
	req := NewRequest[bool]("reopenForumTopic", p)
	return req.Do(api)
}
func (api *Api) DeleteForumTopic(p BaseForumTopicP) (bool, error) {
	req := NewRequest[bool]("deleteForumTopic", p)
	return req.Do(api)
}
func (api *Api) UnpinAllForumTopicMessages(p BaseForumTopicP) (bool, error) {
	req := NewRequest[bool]("unpinAllForumTopicMessages", p)
	return req.Do(api)
}

type BaseGeneralForumTopicP struct {
	ChatID int `json:"chat_id"`
}

type EditGeneralForumTopicP struct {
	ChatID int    `json:"chat_id"`
	Name   string `json:"name"`
}

func (api *Api) EditGeneralForumTopic(p EditGeneralForumTopicP) (bool, error) {
	req := NewRequest[bool]("editGeneralForumTopic", p)
	return req.Do(api)
}

func (api *Api) CloseGeneralForumTopic(p BaseGeneralForumTopicP) (bool, error) {
	req := NewRequest[bool]("closeGeneralForumTopic", p)
	return req.Do(api)
}
func (api *Api) ReopenGeneralForumTopic(p BaseGeneralForumTopicP) (bool, error) {
	req := NewRequest[bool]("reopenGeneralForumTopic", p)
	return req.Do(api)
}
func (api *Api) HideGeneralForumTopic(p BaseGeneralForumTopicP) (bool, error) {
	req := NewRequest[bool]("hideGeneralForumTopic", p)
	return req.Do(api)
}
func (api *Api) UnhideGeneralForumTopic(p BaseGeneralForumTopicP) (bool, error) {
	req := NewRequest[bool]("unhideGeneralForumTopic", p)
	return req.Do(api)
}
func (api *Api) UnpinAllGeneralForumTopicMessages(p BaseGeneralForumTopicP) (bool, error) {
	req := NewRequest[bool]("unpinAllGeneralForumTopicMessages", p)
	return req.Do(api)
}
