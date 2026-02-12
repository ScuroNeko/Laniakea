package tgapi

type SendMessageP struct {
	BusinessConnectionID  string `json:"business_connection_id,omitempty"`
	ChatID                int    `json:"chat_id"`
	MessageThreadID       int    `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID int    `json:"direct_messages_topic_id,omitempty"`

	Text                 string              `json:"text"`
	ParseMode            ParseMode           `json:"parse_mode,omitempty"`
	Entities             []MessageEntity     `json:"entities,omitempty"`
	LinkPreviewOptions   *LinkPreviewOptions `json:"link_preview_options,omitempty"`
	DisableNotifications bool                `json:"disable_notifications,omitempty"`
	ProtectContent       bool                `json:"protect_content,omitempty"`
	AllowPaidBroadcast   bool                `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID      string              `json:"message_effect_id,omitempty"`

	SuggestedPostParameters *SuggestedPostParameters `json:"suggested_post_parameters,omitempty"`
	ReplyParameters         *ReplyParameters         `json:"reply_parameters,omitempty"`
	ReplyMarkup             *ReplyMarkup             `json:"reply_markup,omitempty"`
}

func (api *Api) SendMessage(params SendMessageP) (Message, error) {
	req := NewRequest[Message, SendMessageP]("sendMessage", params)
	return req.Do(api)
}

type ForwardMessageP struct {
	ChatID                int `json:"chat_id"`
	MessageThreadID       int `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID int `json:"direct_messages_topic_id,omitempty"`

	MessageID           int  `json:"message_id,omitempty"`
	FromChatID          int  `json:"from_chat_id,omitempty"`
	VideoStartTimestamp int  `json:"video_start_timestamp,omitempty"`
	DisableNotification bool `json:"disable_notification,omitempty"`
	ProtectContent      bool `json:"protect_content,omitempty"`

	MessageEffectID         string                   `json:"message_effect_id,omitempty"`
	SuggestedPostParameters *SuggestedPostParameters `json:"suggested_post_parameters,omitempty"`
}

func (api *Api) ForwardMessage(params ForwardMessageP) (Message, error) {
	req := NewRequest[Message]("forwardMessage", params)
	return req.Do(api)
}

type ForwardMessagesP struct {
	ChatID                int `json:"chat_id"`
	MessageThreadID       int `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID int `json:"direct_messages_topic_id,omitempty"`

	FromChatID          int   `json:"from_chat_id,omitempty"`
	MessageIDs          []int `json:"message_ids,omitempty"`
	DisableNotification bool  `json:"disable_notification,omitempty"`
	ProtectContent      bool  `json:"protect_content,omitempty"`
}

func (api *Api) ForwardMessages(params ForwardMessagesP) ([]int, error) {
	req := NewRequest[[]int]("forwardMessages", params)
	return req.Do(api)
}

type CopyMessageP struct {
	ChatID                int `json:"chat_id"`
	MessageThreadID       int `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID int `json:"direct_messages_topic_id,omitempty"`

	FromChatID          int       `json:"from_chat_id"`
	MessageID           int       `json:"message_id"`
	VideoStartTimestamp int       `json:"video_start_timestamp,omitempty"`
	Caption             string    `json:"caption,omitempty"`
	ParseMode           ParseMode `json:"parse_mode,omitempty"`

	CaptionEntities       []MessageEntity `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia bool            `json:"show_caption_above_media,omitempty"`
	DisableNotification   bool            `json:"disable_notification,omitempty"`
	ProtectContent        bool            `json:"protect_content,omitempty"`
	AllowPaidBroadcast    bool            `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID       string          `json:"message_effect_id,omitempty"`

	SuggestedPostParameters *SuggestedPostParameters `json:"suggested_post_parameters,omitempty"`
	ReplyParameters         *ReplyParameters         `json:"reply_parameters,omitempty"`
	ReplyMarkup             *ReplyMarkup             `json:"reply_markup,omitempty"`
}

func (api *Api) CopyMessage(params CopyMessageP) (int, error) {
	req := NewRequest[int]("copyMessage", params)
	return req.Do(api)
}

type CopyMessagesP struct {
	ChatID                int `json:"chat_id"`
	MessageThreadID       int `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID int `json:"direct_messages_topic_id,omitempty"`

	FromChatID          int   `json:"from_chat_id,omitempty"`
	MessageIDs          []int `json:"message_ids,omitempty"`
	DisableNotification bool  `json:"disable_notification,omitempty"`
	ProtectContent      bool  `json:"protect_content,omitempty"`
	RemoveCaption       bool  `json:"remove_caption,omitempty"`
}

func (api *Api) CopyMessages(params CopyMessagesP) ([]int, error) {
	req := NewRequest[[]int]("copyMessages", params)
	return req.Do(api)
}

type SendLocationP struct {
	BusinessConnectionID  int `json:"business_connection_id,omitempty"`
	ChatID                int `json:"chat_id"`
	MessageThreadID       int `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID int `json:"direct_messages_topic_id,omitempty"`

	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	HorizontalAccuracy   float64 `json:"horizontal_accuracy,omitempty"`
	LivePeriod           int     `json:"live_period,omitempty"`
	Heading              int     `json:"heading,omitempty"`
	ProximityAlertRadius int     `json:"proximity_alert_radius,omitempty"`

	DisableNotification bool   `json:"disable_notification,omitempty"`
	ProtectContent      bool   `json:"protect_content,omitempty"`
	AllowPaidBroadcast  bool   `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID     string `json:"message_effect_id,omitempty"`

	SuggestedPostParameters *SuggestedPostParameters `json:"suggested_post_parameters,omitempty"`
	ReplyParameters         *ReplyParameters         `json:"reply_parameters,omitempty"`
	ReplyMarkup             *ReplyMarkup             `json:"reply_markup,omitempty"`
}

func (api *Api) SendLocation(params SendLocationP) (Message, error) {
	req := NewRequest[Message]("sendLocation", params)
	return req.Do(api)
}

type SendVenueP struct {
	BusinessConnectionID  int `json:"business_connection_id,omitempty"`
	ChatID                int `json:"chat_id"`
	MessageThreadID       int `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID int `json:"direct_messages_topic_id,omitempty"`

	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	Title           string  `json:"title"`
	Address         string  `json:"address"`
	FoursquareID    string  `json:"foursquare_id,omitempty"`
	FoursquareType  string  `json:"foursquare_type,omitempty"`
	GooglePlaceID   string  `json:"google_place_id,omitempty"`
	GooglePlaceType string  `json:"google_place_type,omitempty"`

	DisableNotification bool   `json:"disable_notification,omitempty"`
	ProtectContent      bool   `json:"protect_content,omitempty"`
	AllowPaidBroadcast  bool   `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID     string `json:"message_effect_id,omitempty"`

	SuggestedPostParameters *SuggestedPostParameters `json:"suggested_post_parameters,omitempty"`
	ReplyParameters         *ReplyParameters         `json:"reply_parameters,omitempty"`
	ReplyMarkup             *ReplyMarkup             `json:"reply_markup,omitempty"`
}

func (api *Api) SendVenue(params SendVenueP) (Message, error) {
	req := NewRequest[Message]("sendVenue", params)
	return req.Do(api)
}

type SendContactP struct {
	BusinessConnectionID  int `json:"business_connection_id,omitempty"`
	ChatID                int `json:"chat_id"`
	MessageThreadID       int `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID int `json:"direct_messages_topic_id,omitempty"`

	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name,omitempty"`
	Vcard       string `json:"vcard"`

	DisableNotification bool   `json:"disable_notification,omitempty"`
	ProtectContent      bool   `json:"protect_content,omitempty"`
	AllowPaidBroadcast  bool   `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID     string `json:"message_effect_id,omitempty"`

	SuggestedPostParameters *SuggestedPostParameters `json:"suggested_post_parameters,omitempty"`
	ReplyParameters         *ReplyParameters         `json:"reply_parameters,omitempty"`
	ReplyMarkup             *ReplyMarkup             `json:"reply_markup,omitempty"`
}

func (api *Api) SendContact(params SendContactP) (Message, error) {
	req := NewRequest[Message]("sendContact", params)
	return req.Do(api)
}

type SendPollP struct {
	BusinessConnectionID int `json:"business_connection_id,omitempty"`
	ChatID               int `json:"chat_id"`
	MessageThreadID      int `json:"message_thread_id,omitempty"`

	Question              string            `json:"question"`
	QuestionParseMode     ParseMode         `json:"question_mode,omitempty"`
	QuestionEntities      []MessageEntity   `json:"question_entities,omitempty"`
	Options               []InputPollOption `json:"options"`
	IsAnonymous           bool              `json:"is_anonymous,omitempty"`
	Type                  PollType          `json:"type"`
	AllowsMultipleAnswers bool              `json:"allows_multiple_answers,omitempty"`
	CorrectOptionID       int               `json:"correct_option_id,omitempty"`
	Explanation           string            `json:"explanation,omitempty"`
	ExplanationParseMode  ParseMode         `json:"explanation_parse_mode,omitempty"`
	ExplanationEntities   []MessageEntity   `json:"explanation_entities,omitempty"`
	OpenPeriod            int               `json:"open_period,omitempty"`
	CloseDate             int               `json:"close_date"`
	IsClosed              bool              `json:"is_closed,omitempty"`

	DisableNotification bool   `json:"disable_notification,omitempty"`
	ProtectContent      bool   `json:"protect_content,omitempty"`
	AllowPaidBroadcast  bool   `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID     string `json:"message_effect_id,omitempty"`

	ReplyParameters *ReplyParameters `json:"reply_parameters,omitempty"`
	ReplyMarkup     *ReplyMarkup     `json:"reply_markup,omitempty"`
}

func (api *Api) SendPoll(params SendPollP) (Message, error) {
	req := NewRequest[Message]("sendPoll", params)
	return req.Do(api)
}

type SendChecklistP struct {
	BusinessConnectionID int            `json:"business_connection_id"`
	ChatID               int            `json:"chat_id"`
	Checklist            InputChecklist `json:"checklist"`

	DisableNotification bool   `json:"disable_notification,omitempty"`
	ProtectContent      bool   `json:"protect_content,omitempty"`
	MessageEffectID     string `json:"message_effect_id,omitempty"`

	ReplyParameters *ReplyParameters `json:"reply_parameters,omitempty"`
	ReplyMarkup     *ReplyMarkup     `json:"reply_markup,omitempty"`
}

func (api *Api) SendChecklist(params SendChecklistP) (Message, error) {
	req := NewRequest[Message]("sendChecklist", params)
	return req.Do(api)
}

type SendDiceP struct {
	BusinessConnectionID  int `json:"business_connection_id,omitempty"`
	ChatID                int `json:"chat_id"`
	MessageThreadID       int `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID int `json:"direct_messages_topic_id,omitempty"`

	Emoji string `json:"emoji,omitempty"`

	DisableNotification bool   `json:"disable_notification,omitempty"`
	ProtectContent      bool   `json:"protect_content,omitempty"`
	AllowPaidBroadcast  bool   `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID     string `json:"message_effect_id,omitempty"`

	SuggestedPostParameters *SuggestedPostParameters `json:"suggested_post_parameters,omitempty"`
	ReplyParameters         *ReplyParameters         `json:"reply_parameters,omitempty"`
	ReplyMarkup             *ReplyMarkup             `json:"reply_markup,omitempty"`
}

func (api *Api) SendDice(params SendDiceP) (Message, error) {
	req := NewRequest[Message]("sendDice", params)
	return req.Do(api)
}

type SendMessageDraftP struct {
	ChatID          int             `json:"chat_id"`
	MessageThreadID int             `json:"message_thread_id,omitempty"`
	DraftID         int             `json:"draft_id"`
	Text            string          `json:"text"`
	ParseMode       ParseMode       `json:"parse_mode,omitempty"`
	Entities        []MessageEntity `json:"entities,omitempty"`
}

func (api *Api) SendMessageDraft(params SendMessageDraftP) (bool, error) {
	req := NewRequest[bool]("sendMessageDraft", params)
	return req.Do(api)
}

type SendChatActionP struct {
	BusinessConnectionID string         `json:"business_connection_id,omitempty"`
	ChatID               int            `json:"chat_id"`
	MessageThreadID      int            `json:"message_thread_id,omitempty"`
	Action               ChatActionType `json:"action"`
}

func (api *Api) SendChatAction(params SendChatActionP) (bool, error) {
	req := NewRequest[bool]("sendChatAction", params)
	return req.Do(api)
}

type SetMessageReactionP struct {
	ChatId    int            `json:"chat_id"`
	MessageId int            `json:"message_id"`
	Reaction  []ReactionType `json:"reaction"`
	IsBig     bool           `json:"is_big,omitempty"`
}

func (api *Api) SetMessageReaction(params SetMessageReactionP) (bool, error) {
	req := NewRequest[bool]("setMessageReaction", params)
	return req.Do(api)
}

// Message update methods

type EditMessageTextP struct {
	BusinessConnectionID string       `json:"business_connection_id,omitempty"`
	ChatID               int          `json:"chat_id,omitempty"`
	MessageID            int          `json:"message_id,omitempty"`
	InlineMessageID      string       `json:"inline_message_id,omitempty"`
	Text                 string       `json:"text"`
	ParseMode            ParseMode    `json:"parse_mode,omitempty"`
	ReplyMarkup          *ReplyMarkup `json:"reply_markup,omitempty"`
}

// EditMessageText If inline message, first return will be zero-valued, and second will boolean
// Otherwise, first return will be Message, and second false
func (api *Api) EditMessageText(params EditMessageTextP) (Message, bool, error) {
	var zero Message
	if params.InlineMessageID != "" {
		req := NewRequest[bool]("editMessageText", params)
		res, err := req.Do(api)
		return zero, res, err
	}
	req := NewRequest[Message]("editMessageText", params)
	res, err := req.Do(api)
	return res, false, err
}

type EditMessageCaptionP struct {
	BusinessConnectionID string       `json:"business_connection_id,omitempty"`
	ChatID               int          `json:"chat_id,omitempty"`
	MessageID            int          `json:"message_id,omitempty"`
	InlineMessageID      string       `json:"inline_message_id,omitempty"`
	Caption              string       `json:"caption"`
	ParseMode            ParseMode    `json:"parse_mode,omitempty"`
	ReplyMarkup          *ReplyMarkup `json:"reply_markup,omitempty"`
}

// EditMessageCaption If inline message, first return will be zero-valued, and second will boolean
// Otherwise, first return will be Message, and second false
func (api *Api) EditMessageCaption(params EditMessageCaptionP) (Message, bool, error) {
	var zero Message
	if params.InlineMessageID != "" {
		req := NewRequest[bool]("editMessageCaption", params)
		res, err := req.Do(api)
		return zero, res, err
	}
	req := NewRequest[Message]("editMessageCaption", params)
	res, err := req.Do(api)
	return res, false, err
}

type EditMessageMediaP struct {
	BusinessConnectionID string                `json:"business_connection_id,omitempty"`
	ChatID               int                   `json:"chat_id,omitempty"`
	MessageID            int                   `json:"message_id,omitempty"`
	InlineMessageID      string                `json:"inline_message_id,omitempty"`
	Message              InputMedia            `json:"message"`
	ReplyMarkup          *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// EditMessageMedia If inline message, first return will be zero-valued, and second will boolean
// Otherwise, first return will be Message, and second false
func (api *Api) EditMessageMedia(params EditMessageMediaP) (Message, bool, error) {
	var zero Message
	if params.InlineMessageID != "" {
		req := NewRequest[bool]("editMessageMedia", params)
		res, err := req.Do(api)
		return zero, res, err
	}
	req := NewRequest[Message]("editMessageMedia", params)
	res, err := req.Do(api)
	return res, false, err
}

type EditMessageLiveLocationP struct {
	BusinessConnectionID string `json:"business_connection_id,omitempty"`
	ChatID               int    `json:"chat_id,omitempty"`
	MessageID            int    `json:"message_id,omitempty"`
	InlineMessageID      string `json:"inline_message_id,omitempty"`

	Latitude             float64               `json:"latitude"`
	Longitude            float64               `json:"longitude"`
	LivePeriod           int                   `json:"live_period,omitempty"`
	HorizontalAccuracy   float64               `json:"horizontal_accuracy,omitempty"`
	Heading              int                   `json:"heading,omitempty"`
	ProximityAlertRadius int                   `json:"proximity_alert_radius,omitempty"`
	ReplyMarkup          *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// EditMessageLiveLocation If inline message, first return will be zero-valued, and second will boolean
// Otherwise, first return will be Message, and second false
func (api *Api) EditMessageLiveLocation(params EditMessageLiveLocationP) (Message, bool, error) {
	var zero Message
	if params.InlineMessageID != "" {
		req := NewRequest[bool]("editMessageLiveLocation", params)
		res, err := req.Do(api)
		return zero, res, err
	}
	req := NewRequest[Message]("editMessageLiveLocation", params)
	res, err := req.Do(api)
	return res, false, err
}

type StopMessageLiveLocationP struct {
	BusinessConnectionID string                `json:"business_connection_id,omitempty"`
	ChatID               int                   `json:"chat_id,omitempty"`
	MessageID            int                   `json:"message_id,omitempty"`
	InlineMessageID      string                `json:"inline_message_id,omitempty"`
	ReplyMarkup          *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// StopMessageLiveLocation If inline message, first return will be zero-valued, and second will boolean
// Otherwise, first return will be Message, and second false
func (api *Api) StopMessageLiveLocation(params StopMessageLiveLocationP) (Message, bool, error) {
	var zero Message
	if params.InlineMessageID != "" {
		req := NewRequest[bool]("stopMessageLiveLocation", params)
		res, err := req.Do(api)
		return zero, res, err
	}
	req := NewRequest[Message]("stopMessageLiveLocation", params)
	res, err := req.Do(api)
	return res, false, err
}

type EditMessageChecklistP struct {
	BusinessConnectionID string                `json:"business_connection_id"`
	ChatID               int                   `json:"chat_id"`
	MessageID            int                   `json:"message_id"`
	Checklist            InputChecklist        `json:"checklist"`
	ReplyMarkup          *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

func (api *Api) EditMessageChecklist(params EditMessageChecklistP) (Message, error) {
	req := NewRequest[Message]("editMessageChecklist", params)
	return req.Do(api)
}

type EditMessageReplyMarkupP struct {
	BusinessConnectionID string                `json:"business_connection_id,omitempty"`
	ChatID               int                   `json:"chat_id,omitempty"`
	MessageID            int                   `json:"message_id,omitempty"`
	InlineMessageID      string                `json:"inline_message_id,omitempty"`
	ReplyMarkup          *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

func (api *Api) EditMessageReplyMarkup(params EditMessageReplyMarkupP) (Message, bool, error) {
	var zero Message
	if params.InlineMessageID != "" {
		req := NewRequest[bool]("editMessageReplyMarkup", params)
		res, err := req.Do(api)
		return zero, res, err
	}
	req := NewRequest[Message]("editMessageReplyMarkup", params)
	res, err := req.Do(api)
	return res, false, err
}

type StopPollP struct {
	BusinessConnectionID string `json:"business_connection_id,omitempty"`
	ChatID               int    `json:"chat_id"`
	MessageID            int    `json:"message_id"`
	InlineMessageID      string `json:"inline_message_id,omitempty"`
}

func (api *Api) StopPoll(params StopPollP) (Poll, error) {
	req := NewRequest[Poll]("stopPoll", params)
	return req.Do(api)
}

type ApproveSuggestedPostP struct {
	ChatID    int `json:"chat_id"`
	MessageID int `json:"message_id"`
	SendDate  int `json:"send_date,omitempty"`
}

func (api *Api) ApproveSuggestedPost(params ApproveSuggestedPostP) (bool, error) {
	req := NewRequest[bool]("approveSuggestedPost", params)
	return req.Do(api)
}

type DeclineSuggestedPostP struct {
	ChatID    int    `json:"chat_id"`
	MessageID int    `json:"message_id"`
	Comment   string `json:"comment,omitempty"`
}

func (api *Api) DeclineSuggestedPost(params DeclineSuggestedPostP) (bool, error) {
	req := NewRequest[bool]("declineSuggestedPost", params)
	return req.Do(api)
}

type DeleteMessageP struct {
	ChatID    int `json:"chat_id"`
	MessageID int `json:"message_id"`
}

func (api *Api) DeleteMessage(params DeleteMessageP) (bool, error) {
	req := NewRequest[bool]("deleteMessage", params)
	return req.Do(api)
}

type DeleteMessagesP struct {
	ChatID     int   `json:"chat_id"`
	MessageIDs []int `json:"message_ids"`
}

func (api *Api) DeleteMessages(params DeleteMessagesP) (bool, error) {
	req := NewRequest[bool]("deleteMessages", params)
	return req.Do(api)
}

type AnswerCallbackQueryP struct {
	CallbackQueryID string `json:"callback_query_id"`
	Text            string `json:"text,omitempty"`
	ShowAlert       bool   `json:"show_alert,omitempty"`
	URL             string `json:"url,omitempty"`
	CacheTime       int    `json:"cache_time,omitempty"`
}

func (api *Api) AnswerCallbackQuery(params AnswerCallbackQueryP) (bool, error) {
	req := NewRequest[bool]("answerCallbackQuery", params)
	return req.Do(api)
}
