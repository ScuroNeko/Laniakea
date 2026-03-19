package tgapi

import "context"

// SendMessageP holds parameters for the sendMessage method.
// See https://core.telegram.org/bots/api#sendmessage
type SendMessageP struct {
	BusinessConnectionID  string `json:"business_connection_id,omitempty"`
	ChatID                int64  `json:"chat_id"`
	MessageThreadID       int    `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID int64  `json:"direct_messages_topic_id,omitempty"`

	Text                 string              `json:"text"`
	ParseMode            ParseMode           `json:"parse_mode,omitempty"`
	Entities             []MessageEntity     `json:"entities,omitempty"`
	LinkPreviewOptions   *LinkPreviewOptions `json:"link_preview_options,omitempty"`
	DisableNotifications bool                `json:"disable_notification,omitempty"`
	ProtectContent       bool                `json:"protect_content,omitempty"`
	AllowPaidBroadcast   bool                `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID      string              `json:"message_effect_id,omitempty"`

	SuggestedPostParameters *SuggestedPostParameters `json:"suggested_post_parameters,omitempty"`
	ReplyParameters         *ReplyParameters         `json:"reply_parameters,omitempty"`
	ReplyMarkup             *ReplyMarkup             `json:"reply_markup,omitempty"`
}

// SendMessage sends a text message.
// See https://core.telegram.org/bots/api#sendmessage
func (api *API) SendMessage(params SendMessageP) (Message, error) {
	req := NewRequestWithChatID[Message, SendMessageP]("sendMessage", params, params.ChatID)
	return req.Do(api)
}

// SendMessageWithContext is the context-aware variant of SendMessage.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendmessage
func (api *API) SendMessageWithContext(ctx context.Context, params SendMessageP) (Message, error) {
	req := NewRequestWithChatID[Message, SendMessageP]("sendMessage", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// ForwardMessageP holds parameters for the forwardMessage method.
// See https://core.telegram.org/bots/api#forwardmessage
type ForwardMessageP struct {
	ChatID                int64 `json:"chat_id"`
	MessageThreadID       int   `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID int   `json:"direct_messages_topic_id,omitempty"`

	MessageID           int   `json:"message_id,omitempty"`
	FromChatID          int64 `json:"from_chat_id,omitempty"`
	VideoStartTimestamp int   `json:"video_start_timestamp,omitempty"`
	DisableNotification bool  `json:"disable_notification,omitempty"`
	ProtectContent      bool  `json:"protect_content,omitempty"`

	MessageEffectID         string                   `json:"message_effect_id,omitempty"`
	SuggestedPostParameters *SuggestedPostParameters `json:"suggested_post_parameters,omitempty"`
}

// ForwardMessage forwards a message.
// See https://core.telegram.org/bots/api#forwardmessage
func (api *API) ForwardMessage(params ForwardMessageP) (Message, error) {
	req := NewRequestWithChatID[Message]("forwardMessage", params, params.ChatID)
	return req.Do(api)
}

// ForwardMessageWithContext is the context-aware variant of ForwardMessage.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#forwardmessage
func (api *API) ForwardMessageWithContext(ctx context.Context, params ForwardMessageP) (Message, error) {
	req := NewRequestWithChatID[Message]("forwardMessage", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// ForwardMessagesP holds parameters for the forwardMessages method.
// See https://core.telegram.org/bots/api#forwardmessages
type ForwardMessagesP struct {
	ChatID                int64 `json:"chat_id"`
	MessageThreadID       int   `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID int   `json:"direct_messages_topic_id,omitempty"`

	FromChatID          int64 `json:"from_chat_id,omitempty"`
	MessageIDs          []int `json:"message_ids,omitempty"`
	DisableNotification bool  `json:"disable_notification,omitempty"`
	ProtectContent      bool  `json:"protect_content,omitempty"`
}

// ForwardMessages forwards multiple messages.
// Returns an array of message IDs of the sent messages.
// See https://core.telegram.org/bots/api#forwardmessages
func (api *API) ForwardMessages(params ForwardMessagesP) ([]MessageID, error) {
	req := NewRequestWithChatID[[]MessageID]("forwardMessages", params, params.ChatID)
	return req.Do(api)
}

// ForwardMessagesWithContext is the context-aware variant of ForwardMessages.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#forwardmessages
func (api *API) ForwardMessagesWithContext(ctx context.Context, params ForwardMessagesP) ([]MessageID, error) {
	req := NewRequestWithChatID[[]MessageID]("forwardMessages", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// CopyMessageP holds parameters for the copyMessage method.
// See https://core.telegram.org/bots/api#copymessage
type CopyMessageP struct {
	ChatID                int64 `json:"chat_id"`
	MessageThreadID       int   `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID int   `json:"direct_messages_topic_id,omitempty"`

	FromChatID          int64     `json:"from_chat_id"`
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

// CopyMessage copies a message.
// Returns the MessageID of the sent copy.
// See https://core.telegram.org/bots/api#copymessage
func (api *API) CopyMessage(params CopyMessageP) (int, error) {
	msgID, err := NewRequestWithChatID[MessageID]("copyMessage", params, params.ChatID).Do(api)
	if err != nil {
		return 0, err
	}
	return msgID.MessageID, nil
}

// CopyMessageWithContext is the context-aware variant of CopyMessage.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#copymessage
func (api *API) CopyMessageWithContext(ctx context.Context, params CopyMessageP) (int, error) {
	msgID, err := NewRequestWithChatID[MessageID]("copyMessage", params, params.ChatID).DoWithContext(ctx, api)
	if err != nil {
		return 0, err
	}
	return msgID.MessageID, nil
}

// CopyMessagesP holds parameters for the copyMessages method.
// See https://core.telegram.org/bots/api#copymessages
type CopyMessagesP struct {
	ChatID                int64 `json:"chat_id"`
	MessageThreadID       int   `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID int   `json:"direct_messages_topic_id,omitempty"`

	FromChatID          int64 `json:"from_chat_id,omitempty"`
	MessageIDs          []int `json:"message_ids,omitempty"`
	DisableNotification bool  `json:"disable_notification,omitempty"`
	ProtectContent      bool  `json:"protect_content,omitempty"`
	RemoveCaption       bool  `json:"remove_caption,omitempty"`
}

// CopyMessages copies multiple messages.
// Returns an array of message IDs of the sent copies.
// See https://core.telegram.org/bots/api#copymessages
func (api *API) CopyMessages(params CopyMessagesP) ([]MessageID, error) {
	req := NewRequestWithChatID[[]MessageID]("copyMessages", params, params.ChatID)
	return req.Do(api)
}

// CopyMessagesWithContext is the context-aware variant of CopyMessages.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#copymessages
func (api *API) CopyMessagesWithContext(ctx context.Context, params CopyMessagesP) ([]MessageID, error) {
	req := NewRequestWithChatID[[]MessageID]("copyMessages", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SendLocationP holds parameters for the sendLocation method.
// See https://core.telegram.org/bots/api#sendlocation
type SendLocationP struct {
	BusinessConnectionID  string `json:"business_connection_id,omitempty"`
	ChatID                int64  `json:"chat_id"`
	MessageThreadID       int    `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID int    `json:"direct_messages_topic_id,omitempty"`

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

// SendLocation sends a point on the map.
// See https://core.telegram.org/bots/api#sendlocation
func (api *API) SendLocation(params SendLocationP) (Message, error) {
	req := NewRequestWithChatID[Message]("sendLocation", params, params.ChatID)
	return req.Do(api)
}

// SendLocationWithContext is the context-aware variant of SendLocation.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendlocation
func (api *API) SendLocationWithContext(ctx context.Context, params SendLocationP) (Message, error) {
	req := NewRequestWithChatID[Message]("sendLocation", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SendVenueP holds parameters for the sendVenue method.
// See https://core.telegram.org/bots/api#sendvenue
type SendVenueP struct {
	BusinessConnectionID  string `json:"business_connection_id,omitempty"`
	ChatID                int64  `json:"chat_id"`
	MessageThreadID       int    `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID int    `json:"direct_messages_topic_id,omitempty"`

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

// SendVenue sends information about a venue.
// See https://core.telegram.org/bots/api#sendvenue
func (api *API) SendVenue(params SendVenueP) (Message, error) {
	req := NewRequestWithChatID[Message]("sendVenue", params, params.ChatID)
	return req.Do(api)
}

// SendVenueWithContext is the context-aware variant of SendVenue.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendvenue
func (api *API) SendVenueWithContext(ctx context.Context, params SendVenueP) (Message, error) {
	req := NewRequestWithChatID[Message]("sendVenue", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SendContactP holds parameters for the sendContact method.
// See https://core.telegram.org/bots/api#sendcontact
type SendContactP struct {
	BusinessConnectionID  string `json:"business_connection_id,omitempty"`
	ChatID                int64  `json:"chat_id"`
	MessageThreadID       int    `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID int    `json:"direct_messages_topic_id,omitempty"`

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

// SendContact sends a phone contact.
// See https://core.telegram.org/bots/api#sendcontact
func (api *API) SendContact(params SendContactP) (Message, error) {
	req := NewRequestWithChatID[Message]("sendContact", params, params.ChatID)
	return req.Do(api)
}

// SendContactWithContext is the context-aware variant of SendContact.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendcontact
func (api *API) SendContactWithContext(ctx context.Context, params SendContactP) (Message, error) {
	req := NewRequestWithChatID[Message]("sendContact", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SendPollP holds parameters for the sendPoll method.
// See https://core.telegram.org/bots/api#sendpoll
type SendPollP struct {
	BusinessConnectionID string `json:"business_connection_id,omitempty"`
	ChatID               int64  `json:"chat_id"`
	MessageThreadID      int    `json:"message_thread_id,omitempty"`

	Question              string            `json:"question"`
	QuestionParseMode     ParseMode         `json:"question_parse_mode,omitempty"`
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

// SendPoll sends a native poll.
// See https://core.telegram.org/bots/api#sendpoll
func (api *API) SendPoll(params SendPollP) (Message, error) {
	req := NewRequestWithChatID[Message]("sendPoll", params, params.ChatID)
	return req.Do(api)
}

// SendPollWithContext is the context-aware variant of SendPoll.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendpoll
func (api *API) SendPollWithContext(ctx context.Context, params SendPollP) (Message, error) {
	req := NewRequestWithChatID[Message]("sendPoll", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SendChecklistP holds parameters for the sendChecklist method.
// See https://core.telegram.org/bots/api#sendchecklist
type SendChecklistP struct {
	BusinessConnectionID string         `json:"business_connection_id"`
	ChatID               int64          `json:"chat_id"`
	Checklist            InputChecklist `json:"checklist"`

	DisableNotification bool   `json:"disable_notification,omitempty"`
	ProtectContent      bool   `json:"protect_content,omitempty"`
	MessageEffectID     string `json:"message_effect_id,omitempty"`

	ReplyParameters *ReplyParameters `json:"reply_parameters,omitempty"`
	ReplyMarkup     *ReplyMarkup     `json:"reply_markup,omitempty"`
}

// SendChecklist sends a checklist.
// See https://core.telegram.org/bots/api#sendchecklist
func (api *API) SendChecklist(params SendChecklistP) (Message, error) {
	req := NewRequestWithChatID[Message]("sendChecklist", params, params.ChatID)
	return req.Do(api)
}

// SendChecklistWithContext is the context-aware variant of SendChecklist.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendchecklist
func (api *API) SendChecklistWithContext(ctx context.Context, params SendChecklistP) (Message, error) {
	req := NewRequestWithChatID[Message]("sendChecklist", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SendDiceP holds parameters for the sendDice method.
// See https://core.telegram.org/bots/api#senddice
type SendDiceP struct {
	BusinessConnectionID  string `json:"business_connection_id,omitempty"`
	ChatID                int64  `json:"chat_id"`
	MessageThreadID       int    `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID int    `json:"direct_messages_topic_id,omitempty"`

	Emoji string `json:"emoji,omitempty"`

	DisableNotification bool   `json:"disable_notification,omitempty"`
	ProtectContent      bool   `json:"protect_content,omitempty"`
	AllowPaidBroadcast  bool   `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID     string `json:"message_effect_id,omitempty"`

	SuggestedPostParameters *SuggestedPostParameters `json:"suggested_post_parameters,omitempty"`
	ReplyParameters         *ReplyParameters         `json:"reply_parameters,omitempty"`
	ReplyMarkup             *ReplyMarkup             `json:"reply_markup,omitempty"`
}

// SendDice sends a dice, which will have a random value.
// See https://core.telegram.org/bots/api#senddice
func (api *API) SendDice(params SendDiceP) (Message, error) {
	req := NewRequestWithChatID[Message]("sendDice", params, params.ChatID)
	return req.Do(api)
}

// SendDiceWithContext is the context-aware variant of SendDice.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#senddice
func (api *API) SendDiceWithContext(ctx context.Context, params SendDiceP) (Message, error) {
	req := NewRequestWithChatID[Message]("sendDice", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SendMessageDraftP holds parameters for the sendMessageDraft method.
// See https://core.telegram.org/bots/api#sendmessagedraft
type SendMessageDraftP struct {
	ChatID          int64           `json:"chat_id"`
	MessageThreadID int             `json:"message_thread_id,omitempty"`
	DraftID         uint64          `json:"draft_id"`
	Text            string          `json:"text"`
	ParseMode       ParseMode       `json:"parse_mode,omitempty"`
	Entities        []MessageEntity `json:"entities,omitempty"`
}

// SendMessageDraft sends or updates a draft message in the target chat.
// Returns True on success.
// See https://core.telegram.org/bots/api#sendmessagedraft
func (api *API) SendMessageDraft(params SendMessageDraftP) (bool, error) {
	req := NewRequestWithChatID[bool]("sendMessageDraft", params, params.ChatID)
	return req.Do(api)
}

// SendMessageDraftWithContext is the context-aware variant of SendMessageDraft.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendmessagedraft
func (api *API) SendMessageDraftWithContext(ctx context.Context, params SendMessageDraftP) (bool, error) {
	req := NewRequestWithChatID[bool]("sendMessageDraft", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SendChatActionP holds parameters for the sendChatAction method.
// See https://core.telegram.org/bots/api#sendchataction
type SendChatActionP struct {
	BusinessConnectionID string         `json:"business_connection_id,omitempty"`
	ChatID               int64          `json:"chat_id"`
	MessageThreadID      int            `json:"message_thread_id,omitempty"`
	Action               ChatActionType `json:"action"`
}

// SendChatAction sends a chat action (typing, uploading photo, etc.).
// Returns True on success.
// See https://core.telegram.org/bots/api#sendchataction
func (api *API) SendChatAction(params SendChatActionP) (bool, error) {
	req := NewRequestWithChatID[bool]("sendChatAction", params, params.ChatID)
	return req.Do(api)
}

// SendChatActionWithContext is the context-aware variant of SendChatAction.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendchataction
func (api *API) SendChatActionWithContext(ctx context.Context, params SendChatActionP) (bool, error) {
	req := NewRequestWithChatID[bool]("sendChatAction", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SetMessageReactionP holds parameters for the setMessageReaction method.
// See https://core.telegram.org/bots/api#setmessagereaction
type SetMessageReactionP struct {
	ChatID    int64          `json:"chat_id"`
	MessageId int            `json:"message_id"`
	Reaction  []ReactionType `json:"reaction"`
	IsBig     bool           `json:"is_big,omitempty"`
}

// SetMessageReaction changes the chosen reaction on a message.
// Returns True on success.
// See https://core.telegram.org/bots/api#setmessagereaction
func (api *API) SetMessageReaction(params SetMessageReactionP) (bool, error) {
	req := NewRequestWithChatID[bool]("setMessageReaction", params, params.ChatID)
	return req.Do(api)
}

// SetMessageReactionWithContext is the context-aware variant of SetMessageReaction.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#setmessagereaction
func (api *API) SetMessageReactionWithContext(ctx context.Context, params SetMessageReactionP) (bool, error) {
	req := NewRequestWithChatID[bool]("setMessageReaction", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// EditMessageTextP holds parameters for the editMessageText method.
// See https://core.telegram.org/bots/api#editmessagetext
type EditMessageTextP struct {
	BusinessConnectionID string              `json:"business_connection_id,omitempty"`
	ChatID               int64               `json:"chat_id,omitempty"`
	MessageID            int                 `json:"message_id,omitempty"`
	InlineMessageID      string              `json:"inline_message_id,omitempty"`
	Text                 string              `json:"text"`
	ParseMode            ParseMode           `json:"parse_mode,omitempty"`
	Entities             []MessageEntity     `json:"entities,omitempty"`
	LinkPreviewOptions   *LinkPreviewOptions `json:"link_preview_options,omitempty"`
	ReplyMarkup          *ReplyMarkup        `json:"reply_markup,omitempty"`
}

// EditMessageText edits text messages.
// If inline_message_id is provided, returns a boolean success flag;
// otherwise returns the edited Message.
// See https://core.telegram.org/bots/api#editmessagetext
func (api *API) EditMessageText(params EditMessageTextP) (Message, bool, error) {
	var zero Message
	if params.InlineMessageID != "" {
		req := NewRequestWithChatID[bool]("editMessageText", params, params.ChatID)
		res, err := req.Do(api)
		return zero, res, err
	}
	req := NewRequestWithChatID[Message]("editMessageText", params, params.ChatID)
	res, err := req.Do(api)
	return res, false, err
}

// EditMessageTextWithContext is the context-aware variant of EditMessageText.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#editmessagetext
func (api *API) EditMessageTextWithContext(ctx context.Context, params EditMessageTextP) (Message, bool, error) {
	var zero Message
	if params.InlineMessageID != "" {
		req := NewRequestWithChatID[bool]("editMessageText", params, params.ChatID)
		res, err := req.DoWithContext(ctx, api)
		return zero, res, err
	}
	req := NewRequestWithChatID[Message]("editMessageText", params, params.ChatID)
	res, err := req.DoWithContext(ctx, api)
	return res, false, err
}

// EditMessageCaptionP holds parameters for the editMessageCaption method.
// See https://core.telegram.org/bots/api#editmessagecaption
type EditMessageCaptionP struct {
	BusinessConnectionID  string          `json:"business_connection_id,omitempty"`
	ChatID                int64           `json:"chat_id,omitempty"`
	MessageID             int             `json:"message_id,omitempty"`
	InlineMessageID       string          `json:"inline_message_id,omitempty"`
	Caption               string          `json:"caption"`
	ParseMode             ParseMode       `json:"parse_mode,omitempty"`
	CaptionEntities       []MessageEntity `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia bool            `json:"show_caption_above_media,omitempty"`
	ReplyMarkup           *ReplyMarkup    `json:"reply_markup,omitempty"`
}

// EditMessageCaption edits captions of messages.
// If inline_message_id is provided, returns a boolean success flag;
// otherwise returns the edited Message.
// See https://core.telegram.org/bots/api#editmessagecaption
func (api *API) EditMessageCaption(params EditMessageCaptionP) (Message, bool, error) {
	var zero Message
	if params.InlineMessageID != "" {
		req := NewRequestWithChatID[bool]("editMessageCaption", params, params.ChatID)
		res, err := req.Do(api)
		return zero, res, err
	}
	req := NewRequestWithChatID[Message]("editMessageCaption", params, params.ChatID)
	res, err := req.Do(api)
	return res, false, err
}

// EditMessageCaptionWithContext is the context-aware variant of EditMessageCaption.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#editmessagecaption
func (api *API) EditMessageCaptionWithContext(ctx context.Context, params EditMessageCaptionP) (Message, bool, error) {
	var zero Message
	if params.InlineMessageID != "" {
		req := NewRequestWithChatID[bool]("editMessageCaption", params, params.ChatID)
		res, err := req.DoWithContext(ctx, api)
		return zero, res, err
	}
	req := NewRequestWithChatID[Message]("editMessageCaption", params, params.ChatID)
	res, err := req.DoWithContext(ctx, api)
	return res, false, err
}

// EditMessageMediaP holds parameters for the editMessageMedia method.
// See https://core.telegram.org/bots/api#editmessagemedia
type EditMessageMediaP struct {
	BusinessConnectionID string                `json:"business_connection_id,omitempty"`
	ChatID               int64                 `json:"chat_id,omitempty"`
	MessageID            int                   `json:"message_id,omitempty"`
	InlineMessageID      string                `json:"inline_message_id,omitempty"`
	Media                InputMedia            `json:"media"`
	ReplyMarkup          *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// EditMessageMedia edits media messages.
// If inline_message_id is provided, returns a boolean success flag;
// otherwise returns the edited Message.
// See https://core.telegram.org/bots/api#editmessagemedia
func (api *API) EditMessageMedia(params EditMessageMediaP) (Message, bool, error) {
	var zero Message
	if params.InlineMessageID != "" {
		req := NewRequestWithChatID[bool]("editMessageMedia", params, params.ChatID)
		res, err := req.Do(api)
		return zero, res, err
	}
	req := NewRequestWithChatID[Message]("editMessageMedia", params, params.ChatID)
	res, err := req.Do(api)
	return res, false, err
}

// EditMessageMediaWithContext is the context-aware variant of EditMessageMedia.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#editmessagemedia
func (api *API) EditMessageMediaWithContext(ctx context.Context, params EditMessageMediaP) (Message, bool, error) {
	var zero Message
	if params.InlineMessageID != "" {
		req := NewRequestWithChatID[bool]("editMessageMedia", params, params.ChatID)
		res, err := req.DoWithContext(ctx, api)
		return zero, res, err
	}
	req := NewRequestWithChatID[Message]("editMessageMedia", params, params.ChatID)
	res, err := req.DoWithContext(ctx, api)
	return res, false, err
}

// EditMessageLiveLocationP holds parameters for the editMessageLiveLocation method.
// See https://core.telegram.org/bots/api#editmessagelivelocation
type EditMessageLiveLocationP struct {
	BusinessConnectionID string `json:"business_connection_id,omitempty"`
	ChatID               int64  `json:"chat_id,omitempty"`
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

// EditMessageLiveLocation edits live location messages.
// If inline_message_id is provided, returns a boolean success flag;
// otherwise returns the edited Message.
// See https://core.telegram.org/bots/api#editmessagelivelocation
func (api *API) EditMessageLiveLocation(params EditMessageLiveLocationP) (Message, bool, error) {
	var zero Message
	if params.InlineMessageID != "" {
		req := NewRequestWithChatID[bool]("editMessageLiveLocation", params, params.ChatID)
		res, err := req.Do(api)
		return zero, res, err
	}
	req := NewRequestWithChatID[Message]("editMessageLiveLocation", params, params.ChatID)
	res, err := req.Do(api)
	return res, false, err
}

// EditMessageLiveLocationWithContext is the context-aware variant of EditMessageLiveLocation.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#editmessagelivelocation
func (api *API) EditMessageLiveLocationWithContext(ctx context.Context, params EditMessageLiveLocationP) (Message, bool, error) {
	var zero Message
	if params.InlineMessageID != "" {
		req := NewRequestWithChatID[bool]("editMessageLiveLocation", params, params.ChatID)
		res, err := req.DoWithContext(ctx, api)
		return zero, res, err
	}
	req := NewRequestWithChatID[Message]("editMessageLiveLocation", params, params.ChatID)
	res, err := req.DoWithContext(ctx, api)
	return res, false, err
}

// StopMessageLiveLocationP holds parameters for the stopMessageLiveLocation method.
// See https://core.telegram.org/bots/api#stopmessagelivelocation
type StopMessageLiveLocationP struct {
	BusinessConnectionID string                `json:"business_connection_id,omitempty"`
	ChatID               int64                 `json:"chat_id,omitempty"`
	MessageID            int                   `json:"message_id,omitempty"`
	InlineMessageID      string                `json:"inline_message_id,omitempty"`
	ReplyMarkup          *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// StopMessageLiveLocation stops a live location message.
// If inline_message_id is provided, returns a boolean success flag;
// otherwise returns the edited Message.
// See https://core.telegram.org/bots/api#stopmessagelivelocation
func (api *API) StopMessageLiveLocation(params StopMessageLiveLocationP) (Message, bool, error) {
	var zero Message
	if params.InlineMessageID != "" {
		req := NewRequestWithChatID[bool]("stopMessageLiveLocation", params, params.ChatID)
		res, err := req.Do(api)
		return zero, res, err
	}
	req := NewRequestWithChatID[Message]("stopMessageLiveLocation", params, params.ChatID)
	res, err := req.Do(api)
	return res, false, err
}

// StopMessageLiveLocationWithContext is the context-aware variant of StopMessageLiveLocation.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#stopmessagelivelocation
func (api *API) StopMessageLiveLocationWithContext(ctx context.Context, params StopMessageLiveLocationP) (Message, bool, error) {
	var zero Message
	if params.InlineMessageID != "" {
		req := NewRequestWithChatID[bool]("stopMessageLiveLocation", params, params.ChatID)
		res, err := req.DoWithContext(ctx, api)
		return zero, res, err
	}
	req := NewRequestWithChatID[Message]("stopMessageLiveLocation", params, params.ChatID)
	res, err := req.DoWithContext(ctx, api)
	return res, false, err
}

// EditMessageChecklistP holds parameters for the editMessageChecklist method.
type EditMessageChecklistP struct {
	BusinessConnectionID string                `json:"business_connection_id"`
	ChatID               int64                 `json:"chat_id"`
	MessageID            int                   `json:"message_id"`
	Checklist            InputChecklist        `json:"checklist"`
	ReplyMarkup          *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// EditMessageChecklist edits a checklist message.
// See https://core.telegram.org/bots/api#editmessagechecklist
func (api *API) EditMessageChecklist(params EditMessageChecklistP) (Message, error) {
	req := NewRequestWithChatID[Message]("editMessageChecklist", params, params.ChatID)
	return req.Do(api)
}

// EditMessageChecklistWithContext is the context-aware variant of EditMessageChecklist.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#editmessagechecklist
func (api *API) EditMessageChecklistWithContext(ctx context.Context, params EditMessageChecklistP) (Message, error) {
	req := NewRequestWithChatID[Message]("editMessageChecklist", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// EditMessageReplyMarkupP holds parameters for the editMessageReplyMarkup method.
// See https://core.telegram.org/bots/api#editmessagereplymarkup
type EditMessageReplyMarkupP struct {
	BusinessConnectionID string                `json:"business_connection_id,omitempty"`
	ChatID               int64                 `json:"chat_id,omitempty"`
	MessageID            int                   `json:"message_id,omitempty"`
	InlineMessageID      string                `json:"inline_message_id,omitempty"`
	ReplyMarkup          *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// EditMessageReplyMarkup edits only the reply markup of messages.
// If inline_message_id is provided, returns a boolean success flag;
// otherwise returns the edited Message.
// See https://core.telegram.org/bots/api#editmessagereplymarkup
func (api *API) EditMessageReplyMarkup(params EditMessageReplyMarkupP) (Message, bool, error) {
	var zero Message
	if params.InlineMessageID != "" {
		req := NewRequestWithChatID[bool]("editMessageReplyMarkup", params, params.ChatID)
		res, err := req.Do(api)
		return zero, res, err
	}
	req := NewRequestWithChatID[Message]("editMessageReplyMarkup", params, params.ChatID)
	res, err := req.Do(api)
	return res, false, err
}

// EditMessageReplyMarkupWithContext is the context-aware variant of EditMessageReplyMarkup.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#editmessagereplymarkup
func (api *API) EditMessageReplyMarkupWithContext(ctx context.Context, params EditMessageReplyMarkupP) (Message, bool, error) {
	var zero Message
	if params.InlineMessageID != "" {
		req := NewRequestWithChatID[bool]("editMessageReplyMarkup", params, params.ChatID)
		res, err := req.DoWithContext(ctx, api)
		return zero, res, err
	}
	req := NewRequestWithChatID[Message]("editMessageReplyMarkup", params, params.ChatID)
	res, err := req.DoWithContext(ctx, api)
	return res, false, err
}

// StopPollP holds parameters for the stopPoll method.
// See https://core.telegram.org/bots/api#stoppoll
type StopPollP struct {
	BusinessConnectionID string                `json:"business_connection_id,omitempty"`
	ChatID               int64                 `json:"chat_id"`
	MessageID            int                   `json:"message_id"`
	ReplyMarkup          *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// StopPoll stops a poll that was sent by the bot.
// Returns the stopped Poll.
// See https://core.telegram.org/bots/api#stoppoll
func (api *API) StopPoll(params StopPollP) (Poll, error) {
	req := NewRequestWithChatID[Poll]("stopPoll", params, params.ChatID)
	return req.Do(api)
}

// StopPollWithContext is the context-aware variant of StopPoll.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#stoppoll
func (api *API) StopPollWithContext(ctx context.Context, params StopPollP) (Poll, error) {
	req := NewRequestWithChatID[Poll]("stopPoll", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// ApproveSuggestedPostP holds parameters for the approveSuggestedPost method.
// See https://core.telegram.org/bots/api#approvesuggestedpost
type ApproveSuggestedPostP struct {
	ChatID    int64 `json:"chat_id"`
	MessageID int   `json:"message_id"`
	SendDate  int   `json:"send_date,omitempty"`
}

// ApproveSuggestedPost approves a suggested channel post.
// Returns True on success.
// See https://core.telegram.org/bots/api#approvesuggestedpost
func (api *API) ApproveSuggestedPost(params ApproveSuggestedPostP) (bool, error) {
	req := NewRequestWithChatID[bool]("approveSuggestedPost", params, params.ChatID)
	return req.Do(api)
}

// ApproveSuggestedPostWithContext is the context-aware variant of ApproveSuggestedPost.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#approvesuggestedpost
func (api *API) ApproveSuggestedPostWithContext(ctx context.Context, params ApproveSuggestedPostP) (bool, error) {
	req := NewRequestWithChatID[bool]("approveSuggestedPost", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// DeclineSuggestedPostP holds parameters for the declineSuggestedPost method.
// See https://core.telegram.org/bots/api#declinesuggestedpost
type DeclineSuggestedPostP struct {
	ChatID    int64  `json:"chat_id"`
	MessageID int    `json:"message_id"`
	Comment   string `json:"comment,omitempty"`
}

// DeclineSuggestedPost declines a suggested channel post.
// Returns True on success.
// See https://core.telegram.org/bots/api#declinesuggestedpost
func (api *API) DeclineSuggestedPost(params DeclineSuggestedPostP) (bool, error) {
	req := NewRequestWithChatID[bool]("declineSuggestedPost", params, params.ChatID)
	return req.Do(api)
}

// DeclineSuggestedPostWithContext is the context-aware variant of DeclineSuggestedPost.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#declinesuggestedpost
func (api *API) DeclineSuggestedPostWithContext(ctx context.Context, params DeclineSuggestedPostP) (bool, error) {
	req := NewRequestWithChatID[bool]("declineSuggestedPost", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// DeleteMessageP holds parameters for the deleteMessage method.
// See https://core.telegram.org/bots/api#deletemessage
type DeleteMessageP struct {
	ChatID    int64 `json:"chat_id"`
	MessageID int   `json:"message_id"`
}

// DeleteMessage deletes a message.
// Returns True on success.
// See https://core.telegram.org/bots/api#deletemessage
func (api *API) DeleteMessage(params DeleteMessageP) (bool, error) {
	req := NewRequestWithChatID[bool]("deleteMessage", params, params.ChatID)
	return req.Do(api)
}

// DeleteMessageWithContext is the context-aware variant of DeleteMessage.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#deletemessage
func (api *API) DeleteMessageWithContext(ctx context.Context, params DeleteMessageP) (bool, error) {
	req := NewRequestWithChatID[bool]("deleteMessage", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// DeleteMessagesP holds parameters for the deleteMessages method.
// See https://core.telegram.org/bots/api#deletemessages
type DeleteMessagesP struct {
	ChatID     int64 `json:"chat_id"`
	MessageIDs []int `json:"message_ids"`
}

// DeleteMessages deletes multiple messages at once.
// Returns True on success.
// See https://core.telegram.org/bots/api#deletemessages
func (api *API) DeleteMessages(params DeleteMessagesP) (bool, error) {
	req := NewRequestWithChatID[bool]("deleteMessages", params, params.ChatID)
	return req.Do(api)
}

// DeleteMessagesWithContext is the context-aware variant of DeleteMessages.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#deletemessages
func (api *API) DeleteMessagesWithContext(ctx context.Context, params DeleteMessagesP) (bool, error) {
	req := NewRequestWithChatID[bool]("deleteMessages", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// AnswerCallbackQueryP holds parameters for the answerCallbackQuery method.
// See https://core.telegram.org/bots/api#answercallbackquery
type AnswerCallbackQueryP struct {
	CallbackQueryID string `json:"callback_query_id"`
	Text            string `json:"text,omitempty"`
	ShowAlert       bool   `json:"show_alert,omitempty"`
	URL             string `json:"url,omitempty"`
	CacheTime       int    `json:"cache_time,omitempty"`
}

// AnswerCallbackQuery sends answers to callback queries sent from inline keyboards.
// Returns True on success.
// See https://core.telegram.org/bots/api#answercallbackquery
func (api *API) AnswerCallbackQuery(params AnswerCallbackQueryP) (bool, error) {
	req := NewRequest[bool]("answerCallbackQuery", params)
	return req.Do(api)
}

// AnswerCallbackQueryWithContext is the context-aware variant of AnswerCallbackQuery.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#answercallbackquery
func (api *API) AnswerCallbackQueryWithContext(ctx context.Context, params AnswerCallbackQueryP) (bool, error) {
	req := NewRequest[bool]("answerCallbackQuery", params)
	return req.DoWithContext(ctx, api)
}
