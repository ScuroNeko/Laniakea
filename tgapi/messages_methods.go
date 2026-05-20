package tgapi

import "context"

// SendMessage holds parameters for the sendMessage method.
// Since: Bot API 1.0
// See https://core.telegram.org/bots/api#sendmessage
type SendMessage struct {
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
// Since: Bot API 1.0
// See https://core.telegram.org/bots/api#sendmessage
func (api *API) SendMessage(params SendMessage) (Message, error) {
	req := NewRequestWithChatID[Message, SendMessage]("sendMessage", params, params.ChatID)
	return req.Do(api)
}

// SendMessageWithContext is the context-aware variant of SendMessage.
// Since: Bot API 1.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendmessage
func (api *API) SendMessageWithContext(ctx context.Context, params SendMessage) (Message, error) {
	req := NewRequestWithChatID[Message, SendMessage]("sendMessage", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// ForwardMessage holds parameters for the forwardMessage method.
// Since: Bot API 1.0
// See https://core.telegram.org/bots/api#forwardmessage
type ForwardMessage struct {
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
// Since: Bot API 1.0
// See https://core.telegram.org/bots/api#forwardmessage
func (api *API) ForwardMessage(params ForwardMessage) (Message, error) {
	req := NewRequestWithChatID[Message]("forwardMessage", params, params.ChatID)
	return req.Do(api)
}

// ForwardMessageWithContext is the context-aware variant of ForwardMessage.
// Since: Bot API 1.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#forwardmessage
func (api *API) ForwardMessageWithContext(ctx context.Context, params ForwardMessage) (Message, error) {
	req := NewRequestWithChatID[Message]("forwardMessage", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// ForwardMessages holds parameters for the forwardMessages method.
// Since: Bot API 7.0
// See https://core.telegram.org/bots/api#forwardmessages
type ForwardMessages struct {
	ChatID                int64 `json:"chat_id"`
	MessageThreadID       int   `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID int   `json:"direct_messages_topic_id,omitempty"`

	FromChatID          int64 `json:"from_chat_id,omitempty"`
	MessageIDs          []int `json:"message_ids,omitempty"`
	DisableNotification bool  `json:"disable_notification,omitempty"`
	ProtectContent      bool  `json:"protect_content,omitempty"`
}

// ForwardMessages forwards multiple messages.
// Since: Bot API 7.0
// Returns an array of message IDs of the sent messages.
// See https://core.telegram.org/bots/api#forwardmessages
func (api *API) ForwardMessages(params ForwardMessages) ([]MessageID, error) {
	req := NewRequestWithChatID[[]MessageID]("forwardMessages", params, params.ChatID)
	return req.Do(api)
}

// ForwardMessagesWithContext is the context-aware variant of ForwardMessages.
// Since: Bot API 7.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#forwardmessages
func (api *API) ForwardMessagesWithContext(ctx context.Context, params ForwardMessages) ([]MessageID, error) {
	req := NewRequestWithChatID[[]MessageID]("forwardMessages", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// CopyMessage holds parameters for the copyMessage method.
// Since: Bot API 5.0
// See https://core.telegram.org/bots/api#copymessage
type CopyMessage struct {
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
// Since: Bot API 5.0
// Returns the MessageID of the sent copy.
// See https://core.telegram.org/bots/api#copymessage
func (api *API) CopyMessage(params CopyMessage) (int, error) {
	msgID, err := NewRequestWithChatID[MessageID]("copyMessage", params, params.ChatID).Do(api)
	if err != nil {
		return 0, err
	}
	return msgID.MessageID, nil
}

// CopyMessageWithContext is the context-aware variant of CopyMessage.
// Since: Bot API 5.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#copymessage
func (api *API) CopyMessageWithContext(ctx context.Context, params CopyMessage) (int, error) {
	msgID, err := NewRequestWithChatID[MessageID]("copyMessage", params, params.ChatID).DoWithContext(ctx, api)
	if err != nil {
		return 0, err
	}
	return msgID.MessageID, nil
}

// CopyMessages holds parameters for the copyMessages method.
// Since: Bot API 7.0
// See https://core.telegram.org/bots/api#copymessages
type CopyMessages struct {
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
// Since: Bot API 7.0
// Returns an array of message IDs of the sent copies.
// See https://core.telegram.org/bots/api#copymessages
func (api *API) CopyMessages(params CopyMessages) ([]MessageID, error) {
	req := NewRequestWithChatID[[]MessageID]("copyMessages", params, params.ChatID)
	return req.Do(api)
}

// CopyMessagesWithContext is the context-aware variant of CopyMessages.
// Since: Bot API 7.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#copymessages
func (api *API) CopyMessagesWithContext(ctx context.Context, params CopyMessages) ([]MessageID, error) {
	req := NewRequestWithChatID[[]MessageID]("copyMessages", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SendLocation holds parameters for the sendLocation method.
// Since: Bot API 1.0
// See https://core.telegram.org/bots/api#sendlocation
type SendLocation struct {
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
// Since: Bot API 1.0
// See https://core.telegram.org/bots/api#sendlocation
func (api *API) SendLocation(params SendLocation) (Message, error) {
	req := NewRequestWithChatID[Message]("sendLocation", params, params.ChatID)
	return req.Do(api)
}

// SendLocationWithContext is the context-aware variant of SendLocation.
// Since: Bot API 1.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendlocation
func (api *API) SendLocationWithContext(ctx context.Context, params SendLocation) (Message, error) {
	req := NewRequestWithChatID[Message]("sendLocation", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SendVenue holds parameters for the sendVenue method.
// Since: Bot API 2.0
// See https://core.telegram.org/bots/api#sendvenue
type SendVenue struct {
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
// Since: Bot API 2.0
// See https://core.telegram.org/bots/api#sendvenue
func (api *API) SendVenue(params SendVenue) (Message, error) {
	req := NewRequestWithChatID[Message]("sendVenue", params, params.ChatID)
	return req.Do(api)
}

// SendVenueWithContext is the context-aware variant of SendVenue.
// Since: Bot API 2.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendvenue
func (api *API) SendVenueWithContext(ctx context.Context, params SendVenue) (Message, error) {
	req := NewRequestWithChatID[Message]("sendVenue", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SendContact holds parameters for the sendContact method.
// Since: Bot API 2.0
// See https://core.telegram.org/bots/api#sendcontact
type SendContact struct {
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
// Since: Bot API 2.0
// See https://core.telegram.org/bots/api#sendcontact
func (api *API) SendContact(params SendContact) (Message, error) {
	req := NewRequestWithChatID[Message]("sendContact", params, params.ChatID)
	return req.Do(api)
}

// SendContactWithContext is the context-aware variant of SendContact.
// Since: Bot API 2.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendcontact
func (api *API) SendContactWithContext(ctx context.Context, params SendContact) (Message, error) {
	req := NewRequestWithChatID[Message]("sendContact", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SendPoll holds parameters for the sendPoll method.
// Since: Bot API 4.2
// See https://core.telegram.org/bots/api#sendpoll
type SendPoll struct {
	BusinessConnectionID string `json:"business_connection_id,omitempty"`
	ChatID               int64  `json:"chat_id"`
	MessageThreadID      int    `json:"message_thread_id,omitempty"`

	Question               string            `json:"question"`
	QuestionParseMode      ParseMode         `json:"question_parse_mode,omitempty"`
	QuestionEntities       []MessageEntity   `json:"question_entities,omitempty"`
	Options                []InputPollOption `json:"options"`
	IsAnonymous            bool              `json:"is_anonymous,omitempty"`
	Type                   PollType          `json:"type"`
	AllowsMultipleAnswers  bool              `json:"allows_multiple_answers,omitempty"`
	AllowsRevoting         bool              `json:"allows_revoting,omitempty"`
	ShuffleOptions         bool              `json:"shuffle_options,omitempty"`
	AllowAddingOptions     bool              `json:"allow_adding_options,omitempty"`
	HideResultsUntilCloses bool              `json:"hide_results_until_closes,omitempty"`
	MembersOnly            bool              `json:"members_only,omitempty"`  // Since: Bot API 10.0
	CountryCodes           []string          `json:"country_codes,omitempty"` // Since: Bot API 10.0
	CorrectOptionIDs       []int             `json:"correct_option_ids,omitempty"`
	Explanation            string            `json:"explanation,omitempty"`
	ExplanationParseMode   ParseMode         `json:"explanation_parse_mode,omitempty"`
	ExplanationEntities    []MessageEntity   `json:"explanation_entities,omitempty"`
	ExplanationMedia       *InputPollMedia   `json:"explanation_media,omitempty"`
	Media                  *InputPollMedia   `json:"media,omitempty"`
	OpenPeriod             int               `json:"open_period,omitempty"`
	CloseDate              int               `json:"close_date"`
	IsClosed               bool              `json:"is_closed,omitempty"`

	Description          string          `json:"description"`
	DescriptionParseMode ParseMode       `json:"description_parse_mode,omitempty"`
	DescriptionEntities  []MessageEntity `json:"description_entities,omitempty"`

	DisableNotification bool   `json:"disable_notification,omitempty"`
	ProtectContent      bool   `json:"protect_content,omitempty"`
	AllowPaidBroadcast  bool   `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID     string `json:"message_effect_id,omitempty"`

	ReplyParameters *ReplyParameters `json:"reply_parameters,omitempty"`
	ReplyMarkup     *ReplyMarkup     `json:"reply_markup,omitempty"`
}

// SendPoll sends a native poll.
// Since: Bot API 4.2
// See https://core.telegram.org/bots/api#sendpoll
func (api *API) SendPoll(params SendPoll) (Message, error) {
	req := NewRequestWithChatID[Message]("sendPoll", params, params.ChatID)
	return req.Do(api)
}

// SendPollWithContext is the context-aware variant of SendPoll.
// Since: Bot API 4.2
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendpoll
func (api *API) SendPollWithContext(ctx context.Context, params SendPoll) (Message, error) {
	req := NewRequestWithChatID[Message]("sendPoll", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SendChecklist holds parameters for the sendChecklist method.
// Since: Bot API 9.1
// See https://core.telegram.org/bots/api#sendchecklist
type SendChecklist struct {
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
// Since: Bot API 9.1
// See https://core.telegram.org/bots/api#sendchecklist
func (api *API) SendChecklist(params SendChecklist) (Message, error) {
	req := NewRequestWithChatID[Message]("sendChecklist", params, params.ChatID)
	return req.Do(api)
}

// SendChecklistWithContext is the context-aware variant of SendChecklist.
// Since: Bot API 9.1
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendchecklist
func (api *API) SendChecklistWithContext(ctx context.Context, params SendChecklist) (Message, error) {
	req := NewRequestWithChatID[Message]("sendChecklist", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SendDice holds parameters for the sendDice method.
// Since: Bot API 4.7
// See https://core.telegram.org/bots/api#senddice
type SendDice struct {
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
// Since: Bot API 4.7
// See https://core.telegram.org/bots/api#senddice
func (api *API) SendDice(params SendDice) (Message, error) {
	req := NewRequestWithChatID[Message]("sendDice", params, params.ChatID)
	return req.Do(api)
}

// SendDiceWithContext is the context-aware variant of SendDice.
// Since: Bot API 4.7
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#senddice
func (api *API) SendDiceWithContext(ctx context.Context, params SendDice) (Message, error) {
	req := NewRequestWithChatID[Message]("sendDice", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SendMessageDraft holds parameters for the sendMessageDraft method.
// Since: Bot API 9.1
// See https://core.telegram.org/bots/api#sendmessagedraft
type SendMessageDraft struct {
	ChatID          int64           `json:"chat_id"`
	MessageThreadID int             `json:"message_thread_id,omitempty"`
	DraftID         uint64          `json:"draft_id"`
	Text            string          `json:"text"`
	ParseMode       ParseMode       `json:"parse_mode,omitempty"`
	Entities        []MessageEntity `json:"entities,omitempty"`
}

// SendMessageDraft sends or updates a draft message in the target chat.
// Since: Bot API 9.1
// Returns True on success.
// See https://core.telegram.org/bots/api#sendmessagedraft
func (api *API) SendMessageDraft(params SendMessageDraft) (bool, error) {
	req := NewRequestWithChatID[bool]("sendMessageDraft", params, params.ChatID)
	return req.Do(api)
}

// SendMessageDraftWithContext is the context-aware variant of SendMessageDraft.
// Since: Bot API 9.1
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendmessagedraft
func (api *API) SendMessageDraftWithContext(ctx context.Context, params SendMessageDraft) (bool, error) {
	req := NewRequestWithChatID[bool]("sendMessageDraft", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SendChatAction holds parameters for the sendChatAction method.
// Since: Bot API 1.0
// See https://core.telegram.org/bots/api#sendchataction
type SendChatAction struct {
	BusinessConnectionID string         `json:"business_connection_id,omitempty"`
	ChatID               int64          `json:"chat_id"`
	MessageThreadID      int            `json:"message_thread_id,omitempty"`
	Action               ChatActionType `json:"action"`
}

// SendChatAction sends a chat action (typing, uploading photo, etc.).
// Since: Bot API 1.0
// Returns True on success.
// See https://core.telegram.org/bots/api#sendchataction
func (api *API) SendChatAction(params SendChatAction) (bool, error) {
	req := NewRequestWithChatID[bool]("sendChatAction", params, params.ChatID)
	return req.Do(api)
}

// SendChatActionWithContext is the context-aware variant of SendChatAction.
// Since: Bot API 1.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendchataction
func (api *API) SendChatActionWithContext(ctx context.Context, params SendChatAction) (bool, error) {
	req := NewRequestWithChatID[bool]("sendChatAction", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SetMessageReaction holds parameters for the setMessageReaction method.
// Since: Bot API 7.0
// See https://core.telegram.org/bots/api#setmessagereaction
type SetMessageReaction struct {
	ChatID    int64          `json:"chat_id"`
	MessageID int            `json:"message_id"`
	Reaction  []ReactionType `json:"reaction"`
	IsBig     bool           `json:"is_big,omitempty"`
}

// SetMessageReaction changes the chosen reaction on a message.
// Since: Bot API 7.0
// Returns True on success.
// See https://core.telegram.org/bots/api#setmessagereaction
func (api *API) SetMessageReaction(params SetMessageReaction) (bool, error) {
	req := NewRequestWithChatID[bool]("setMessageReaction", params, params.ChatID)
	return req.Do(api)
}

// SetMessageReactionWithContext is the context-aware variant of SetMessageReaction.
// Since: Bot API 7.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#setmessagereaction
func (api *API) SetMessageReactionWithContext(ctx context.Context, params SetMessageReaction) (bool, error) {
	req := NewRequestWithChatID[bool]("setMessageReaction", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// EditMessageText holds parameters for the editMessageText method.
// Since: Bot API 2.0
// See https://core.telegram.org/bots/api#editmessagetext
type EditMessageText struct {
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
// Since: Bot API 2.0
// If inline_message_id is provided, returns a boolean success flag;
// otherwise returns the edited Message.
// See https://core.telegram.org/bots/api#editmessagetext
func (api *API) EditMessageText(params EditMessageText) (Message, bool, error) {
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
// Since: Bot API 2.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#editmessagetext
func (api *API) EditMessageTextWithContext(ctx context.Context, params EditMessageText) (Message, bool, error) {
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

// EditMessageCaption holds parameters for the editMessageCaption method.
// Since: Bot API 2.0
// See https://core.telegram.org/bots/api#editmessagecaption
type EditMessageCaption struct {
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
// Since: Bot API 2.0
// If inline_message_id is provided, returns a boolean success flag;
// otherwise returns the edited Message.
// See https://core.telegram.org/bots/api#editmessagecaption
func (api *API) EditMessageCaption(params EditMessageCaption) (Message, bool, error) {
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
// Since: Bot API 2.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#editmessagecaption
func (api *API) EditMessageCaptionWithContext(ctx context.Context, params EditMessageCaption) (Message, bool, error) {
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

// EditMessageMedia holds parameters for the editMessageMedia method.
// Since: Bot API 4.0
// See https://core.telegram.org/bots/api#editmessagemedia
type EditMessageMedia struct {
	BusinessConnectionID string                `json:"business_connection_id,omitempty"`
	ChatID               int64                 `json:"chat_id,omitempty"`
	MessageID            int                   `json:"message_id,omitempty"`
	InlineMessageID      string                `json:"inline_message_id,omitempty"`
	Media                InputMedia            `json:"media"`
	ReplyMarkup          *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// EditMessageMedia edits media messages.
// Since: Bot API 4.0
// If inline_message_id is provided, returns a boolean success flag;
// otherwise returns the edited Message.
// See https://core.telegram.org/bots/api#editmessagemedia
func (api *API) EditMessageMedia(params EditMessageMedia) (Message, bool, error) {
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
// Since: Bot API 4.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#editmessagemedia
func (api *API) EditMessageMediaWithContext(ctx context.Context, params EditMessageMedia) (Message, bool, error) {
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

// EditMessageLiveLocation holds parameters for the editMessageLiveLocation method.
// Since: Bot API 3.4
// See https://core.telegram.org/bots/api#editmessagelivelocation
type EditMessageLiveLocation struct {
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
// Since: Bot API 3.4
// If inline_message_id is provided, returns a boolean success flag;
// otherwise returns the edited Message.
// See https://core.telegram.org/bots/api#editmessagelivelocation
func (api *API) EditMessageLiveLocation(params EditMessageLiveLocation) (Message, bool, error) {
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
// Since: Bot API 3.4
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#editmessagelivelocation
func (api *API) EditMessageLiveLocationWithContext(ctx context.Context, params EditMessageLiveLocation) (Message, bool, error) {
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

// StopMessageLiveLocation holds parameters for the stopMessageLiveLocation method.
// Since: Bot API 3.4
// See https://core.telegram.org/bots/api#stopmessagelivelocation
type StopMessageLiveLocation struct {
	BusinessConnectionID string                `json:"business_connection_id,omitempty"`
	ChatID               int64                 `json:"chat_id,omitempty"`
	MessageID            int                   `json:"message_id,omitempty"`
	InlineMessageID      string                `json:"inline_message_id,omitempty"`
	ReplyMarkup          *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// StopMessageLiveLocation stops a live location message.
// Since: Bot API 3.4
// If inline_message_id is provided, returns a boolean success flag;
// otherwise returns the edited Message.
// See https://core.telegram.org/bots/api#stopmessagelivelocation
func (api *API) StopMessageLiveLocation(params StopMessageLiveLocation) (Message, bool, error) {
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
// Since: Bot API 3.4
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#stopmessagelivelocation
func (api *API) StopMessageLiveLocationWithContext(ctx context.Context, params StopMessageLiveLocation) (Message, bool, error) {
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

// EditMessageChecklist holds parameters for the editMessageChecklist method.
// Since: Bot API 9.1
// See https://core.telegram.org/bots/api#editmessagechecklist
type EditMessageChecklist struct {
	BusinessConnectionID string                `json:"business_connection_id"`
	ChatID               int64                 `json:"chat_id"`
	MessageID            int                   `json:"message_id"`
	Checklist            InputChecklist        `json:"checklist"`
	ReplyMarkup          *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// EditMessageChecklist edits a checklist message.
// Since: Bot API 9.1
// See https://core.telegram.org/bots/api#editmessagechecklist
func (api *API) EditMessageChecklist(params EditMessageChecklist) (Message, error) {
	req := NewRequestWithChatID[Message]("editMessageChecklist", params, params.ChatID)
	return req.Do(api)
}

// EditMessageChecklistWithContext is the context-aware variant of EditMessageChecklist.
// Since: Bot API 9.1
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#editmessagechecklist
func (api *API) EditMessageChecklistWithContext(ctx context.Context, params EditMessageChecklist) (Message, error) {
	req := NewRequestWithChatID[Message]("editMessageChecklist", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// EditMessageReplyMarkup holds parameters for the editMessageReplyMarkup method.
// Since: Bot API 2.0
// See https://core.telegram.org/bots/api#editmessagereplymarkup
type EditMessageReplyMarkup struct {
	BusinessConnectionID string                `json:"business_connection_id,omitempty"`
	ChatID               int64                 `json:"chat_id,omitempty"`
	MessageID            int                   `json:"message_id,omitempty"`
	InlineMessageID      string                `json:"inline_message_id,omitempty"`
	ReplyMarkup          *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// EditMessageReplyMarkup edits only the reply markup of messages.
// Since: Bot API 2.0
// If inline_message_id is provided, returns a boolean success flag;
// otherwise returns the edited Message.
// See https://core.telegram.org/bots/api#editmessagereplymarkup
func (api *API) EditMessageReplyMarkup(params EditMessageReplyMarkup) (Message, bool, error) {
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
// Since: Bot API 2.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#editmessagereplymarkup
func (api *API) EditMessageReplyMarkupWithContext(ctx context.Context, params EditMessageReplyMarkup) (Message, bool, error) {
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

// StopPoll holds parameters for the stopPoll method.
// Since: Bot API 4.2
// See https://core.telegram.org/bots/api#stoppoll
type StopPoll struct {
	BusinessConnectionID string                `json:"business_connection_id,omitempty"`
	ChatID               int64                 `json:"chat_id"`
	MessageID            int                   `json:"message_id"`
	ReplyMarkup          *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// StopPoll stops a poll that was sent by the bot.
// Since: Bot API 4.2
// Returns the stopped Poll.
// See https://core.telegram.org/bots/api#stoppoll
func (api *API) StopPoll(params StopPoll) (Poll, error) {
	req := NewRequestWithChatID[Poll]("stopPoll", params, params.ChatID)
	return req.Do(api)
}

// StopPollWithContext is the context-aware variant of StopPoll.
// Since: Bot API 4.2
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#stoppoll
func (api *API) StopPollWithContext(ctx context.Context, params StopPoll) (Poll, error) {
	req := NewRequestWithChatID[Poll]("stopPoll", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// ApproveSuggestedPost holds parameters for the approveSuggestedPost method.
// Since: Bot API 9.2
// See https://core.telegram.org/bots/api#approvesuggestedpost
type ApproveSuggestedPost struct {
	ChatID    int64 `json:"chat_id"`
	MessageID int   `json:"message_id"`
	SendDate  int   `json:"send_date,omitempty"`
}

// ApproveSuggestedPost approves a suggested channel post.
// Since: Bot API 9.2
// Returns True on success.
// See https://core.telegram.org/bots/api#approvesuggestedpost
func (api *API) ApproveSuggestedPost(params ApproveSuggestedPost) (bool, error) {
	req := NewRequestWithChatID[bool]("approveSuggestedPost", params, params.ChatID)
	return req.Do(api)
}

// ApproveSuggestedPostWithContext is the context-aware variant of ApproveSuggestedPost.
// Since: Bot API 9.2
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#approvesuggestedpost
func (api *API) ApproveSuggestedPostWithContext(ctx context.Context, params ApproveSuggestedPost) (bool, error) {
	req := NewRequestWithChatID[bool]("approveSuggestedPost", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// DeclineSuggestedPost holds parameters for the declineSuggestedPost method.
// Since: Bot API 9.2
// See https://core.telegram.org/bots/api#declinesuggestedpost
type DeclineSuggestedPost struct {
	ChatID    int64  `json:"chat_id"`
	MessageID int    `json:"message_id"`
	Comment   string `json:"comment,omitempty"`
}

// DeclineSuggestedPost declines a suggested channel post.
// Since: Bot API 9.2
// Returns True on success.
// See https://core.telegram.org/bots/api#declinesuggestedpost
func (api *API) DeclineSuggestedPost(params DeclineSuggestedPost) (bool, error) {
	req := NewRequestWithChatID[bool]("declineSuggestedPost", params, params.ChatID)
	return req.Do(api)
}

// DeclineSuggestedPostWithContext is the context-aware variant of DeclineSuggestedPost.
// Since: Bot API 9.2
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#declinesuggestedpost
func (api *API) DeclineSuggestedPostWithContext(ctx context.Context, params DeclineSuggestedPost) (bool, error) {
	req := NewRequestWithChatID[bool]("declineSuggestedPost", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// DeleteMessage holds parameters for the deleteMessage method.
// Since: Bot API 3.0
// See https://core.telegram.org/bots/api#deletemessage
type DeleteMessage struct {
	ChatID    int64 `json:"chat_id"`
	MessageID int   `json:"message_id"`
}

// DeleteMessage deletes a message.
// Since: Bot API 3.0
// Returns True on success.
// See https://core.telegram.org/bots/api#deletemessage
func (api *API) DeleteMessage(params DeleteMessage) (bool, error) {
	req := NewRequestWithChatID[bool]("deleteMessage", params, params.ChatID)
	return req.Do(api)
}

// DeleteMessageWithContext is the context-aware variant of DeleteMessage.
// Since: Bot API 3.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#deletemessage
func (api *API) DeleteMessageWithContext(ctx context.Context, params DeleteMessage) (bool, error) {
	req := NewRequestWithChatID[bool]("deleteMessage", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// DeleteMessages holds parameters for the deleteMessages method.
// Since: Bot API 7.0
// See https://core.telegram.org/bots/api#deletemessages
type DeleteMessages struct {
	ChatID     int64 `json:"chat_id"`
	MessageIDs []int `json:"message_ids"`
}

// DeleteMessages deletes multiple messages at once.
// Since: Bot API 7.0
// Returns True on success.
// See https://core.telegram.org/bots/api#deletemessages
func (api *API) DeleteMessages(params DeleteMessages) (bool, error) {
	req := NewRequestWithChatID[bool]("deleteMessages", params, params.ChatID)
	return req.Do(api)
}

// DeleteMessagesWithContext is the context-aware variant of DeleteMessages.
// Since: Bot API 7.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#deletemessages
func (api *API) DeleteMessagesWithContext(ctx context.Context, params DeleteMessages) (bool, error) {
	req := NewRequestWithChatID[bool]("deleteMessages", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// AnswerCallbackQuery holds parameters for the answerCallbackQuery method.
// Since: Bot API 2.0
// See https://core.telegram.org/bots/api#answercallbackquery
type AnswerCallbackQuery struct {
	CallbackQueryID string `json:"callback_query_id"`
	Text            string `json:"text,omitempty"`
	ShowAlert       bool   `json:"show_alert,omitempty"`
	URL             string `json:"url,omitempty"`
	CacheTime       int    `json:"cache_time,omitempty"`
}

// AnswerCallbackQuery sends answers to callback queries sent from inline keyboards.
// Since: Bot API 2.0
// Returns True on success.
// See https://core.telegram.org/bots/api#answercallbackquery
func (api *API) AnswerCallbackQuery(params AnswerCallbackQuery) (bool, error) {
	req := NewRequest[bool]("answerCallbackQuery", params)
	return req.Do(api)
}

// AnswerCallbackQueryWithContext is the context-aware variant of AnswerCallbackQuery.
// Since: Bot API 2.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#answercallbackquery
func (api *API) AnswerCallbackQueryWithContext(ctx context.Context, params AnswerCallbackQuery) (bool, error) {
	req := NewRequest[bool]("answerCallbackQuery", params)
	return req.DoWithContext(ctx, api)
}

// AnswerGuestQuery holds parameters for the answerGuestQuery method.
// Since: Bot API 10.0
// See https://core.telegram.org/bots/api#answerguestquery
type AnswerGuestQuery struct {
	GuestQueryID string            `json:"guest_query_id"`
	Result       InlineQueryResult `json:"result"`
}

// AnswerGuestQuery answers a guest query.
// Since: Bot API 10.0
// See https://core.telegram.org/bots/api#answerguestquery
func (api *API) AnswerGuestQuery(params AnswerGuestQuery) (SentGuestMessage, error) {
	req := NewRequest[SentGuestMessage]("answerGuestQuery", params)
	return req.Do(api)
}

// AnswerGuestQueryWithContext is the context-aware variant of AnswerGuestQuery.
// Since: Bot API 10.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#answerguestquery
func (api *API) AnswerGuestQueryWithContext(ctx context.Context, params AnswerGuestQuery) (SentGuestMessage, error) {
	req := NewRequest[SentGuestMessage]("answerGuestQuery", params)
	return req.DoWithContext(ctx, api)
}

// DeleteAllMessageReactions holds parameters for the deleteAllMessageReactions method.
// Since: Bot API 10.0
// See https://core.telegram.org/bots/api#deleteallmessagereactions
type DeleteAllMessageReactions struct {
	ChatID      int64 `json:"chat_id"`
	UserID      int64 `json:"user_id,omitempty"`
	ActorChatID int64 `json:"actor_chat_id,omitempty"`
}

// DeleteAllMessageReactions deletes all reactions on a message.
// Since: Bot API 10.0
// Returns True on success.
// See https://core.telegram.org/bots/api#deleteallmessagereactions
func (api *API) DeleteAllMessageReactions(params DeleteAllMessageReactions) (bool, error) {
	req := NewRequest[bool]("deleteAllMessageReactions", params)
	return req.Do(api)
}

// DeleteAllMessageReactionWithContext is the context-aware variant of DeleteAllMessageReactions.
// Since: Bot API 10.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#deleteallmessagereactions
func (api *API) DeleteAllMessageReactionWithContext(ctx context.Context, params DeleteAllMessageReactions) (bool, error) {
	req := NewRequest[bool]("deleteAllMessageReactions", params)
	return req.DoWithContext(ctx, api)
}

// DeleteMessageReaction holds parameters for the deleteMessageReaction method.
// Since: Bot API 10.0
// See https://core.telegram.org/bots/api#deletemessagereaction
type DeleteMessageReaction struct {
	ChatID      int64 `json:"chat_id"`
	MessageID   int   `json:"message_id"`
	UserID      int64 `json:"user_id,omitempty"`
	ActorChatID int64 `json:"actor_chat_id,omitempty"`
}

// DeleteMessageReaction deletes a reaction on a message.
// Since: Bot API 10.0
// Returns True on success.
// See https://core.telegram.org/bots/api#deletemessagereaction
func (api *API) DeleteMessageReaction(params DeleteMessageReaction) (bool, error) {
	req := NewRequest[bool]("deleteMessageReaction", params)
	return req.Do(api)
}

// DeleteMessageReactionWithContext is the context-aware variant of DeleteMessageReaction.
// Since: Bot API 10.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#deletemessagereaction
func (api *API) DeleteMessageReactionWithContext(ctx context.Context, params DeleteMessageReaction) (bool, error) {
	req := NewRequest[bool]("deleteMessageReaction", params)
	return req.DoWithContext(ctx, api)
}
