package tgapi

// Chat represents a chat (private, group, supergroup, channel).
// Since: Bot API 1.0
// See https://core.telegram.org/bots/api#chat
type Chat struct {
	ID               int64    `json:"id"`
	Type             ChatType `json:"type"`
	Title            *string  `json:"title,omitempty"`
	Username         *string  `json:"username,omitempty"`
	FirstName        *string  `json:"first_name,omitempty"`
	LastName         *string  `json:"last_name,omitempty"`
	IsForum          *bool    `json:"is_forum,omitempty"`           // Since: Bot API 6.3
	IsDirectMessages *bool    `json:"is_direct_messages,omitempty"` // Since: Bot API 9.2
}

// ChatType represents the type of a chat.
type ChatType string

const (
	// ChatTypePrivate identifies a private chat.
	ChatTypePrivate ChatType = "private"
	// ChatTypeGroup identifies a basic group chat.
	ChatTypeGroup ChatType = "group"
	// ChatTypeSupergroup identifies a supergroup chat.
	ChatTypeSupergroup ChatType = "supergroup"
	// ChatTypeChannel identifies a channel chat.
	ChatTypeChannel ChatType = "channel"
)

// ChatFullInfo contains full information about a chat.
// Since: Bot API 7.5
// See https://core.telegram.org/bots/api#chatfullinfo
type ChatFullInfo struct {
	ID               int64      `json:"id"`
	Type             ChatType   `json:"type"`
	Title            string     `json:"title"`
	Username         string     `json:"username"`
	FirstName        string     `json:"first_name"`
	LastName         string     `json:"last_name"`
	IsForum          bool       `json:"is_forum"`
	IsDirectMessages bool       `json:"is_direct_messages"`
	AccentColorID    int        `json:"accent_color_id"`
	MaxReactionCount int        `json:"max_reaction_count"`
	Photo            *ChatPhoto `json:"photo,omitempty"`
	ActiveUsernames  []string   `json:"active_usernames,omitempty"`
	Birthdate        *Birthdate `json:"birthdate,omitempty"`

	BusinessIntro        *BusinessIntro        `json:"business_intro,omitempty"`
	BusinessLocation     *BusinessLocation     `json:"business_location,omitempty"`
	BusinessOpeningHours *BusinessOpeningHours `json:"business_opening_hours,omitempty"`

	PersonalChat *Chat `json:"personal_chat,omitempty"`
	ParentChat   *Chat `json:"parent_chat,omitempty"` // Since: Bot API 9.2

	AvailableReaction []ReactionType `json:"available_reaction,omitempty"`

	BackgroundCustomEmojiID        *string `json:"background_custom_emoji_id,omitempty"`
	ProfileAccentColorID           *int    `json:"profile_accent_color_id,omitempty"`
	ProfileBackgroundCustomEmojiID *string `json:"profile_background_custom_emoji_id,omitempty"`
	EmojiStatusCustomEmojiID       *string `json:"emoji_status_custom_emoji_id,omitempty"`
	EmojiStatusExpirationDate      *int    `json:"emoji_status_expiration_date,omitempty"`

	Bio                                *string `json:"bio,omitempty"`
	HasPrivateForwards                 *bool   `json:"has_private_forwards,omitempty"`
	HasRestrictedVoiceAndVideoMessages *bool   `json:"has_restricted_voice_and_video_messages,omitempty"`
	JoinToSendMessages                 *bool   `json:"join_to_send_messages,omitempty"`
	JoinByRequest                      *bool   `json:"join_by_request,omitempty"`

	Description       *string            `json:"description,omitempty"`
	InviteLink        *string            `json:"invite_link,omitempty"`
	PinnedMessage     *Message           `json:"pinned_message,omitempty"`
	Permissions       *ChatPermissions   `json:"permissions,omitempty"`
	AcceptedGiftTypes *AcceptedGiftTypes `json:"accepted_gift_types,omitempty"`

	CanSendPaidMedia             *bool   `json:"can_send_paid_media,omitempty"`
	SlowModeDelay                *int    `json:"slow_mode_delay,omitempty"`
	UnrestrictedBoostCount       *int    `json:"unrestricted_boost_count,omitempty"`
	MessageAutoDeleteTime        *int    `json:"message_auto_delete_time,omitempty"`
	HasAggressiveAntiSpamEnabled *bool   `json:"has_aggressive_anti_spam_enabled,omitempty"`
	HasHiddenMembers             *bool   `json:"has_hidden_members,omitempty"`
	HasProtectedContent          *bool   `json:"has_protected_content,omitempty"`
	HasVisibleHistory            *bool   `json:"has_visible_history,omitempty"`
	StickerSetName               *string `json:"sticker_set_name,omitempty"`
	CanSetStickerSet             *bool   `json:"can_set_sticker_set,omitempty"`
	CustomEmojiStickerSetName    *string `json:"custom_emoji_sticker_set_name,omitempty"`
	LinkedChatID                 *int64  `json:"linked_chat_id,omitempty"`

	Location             *ChatLocation     `json:"location,omitempty"`
	Rating               *UserRating       `json:"rating,omitempty"`
	FirstProfileAudio    *Audio            `json:"first_profile_audio,omitempty"`     // Since: Bot API 9.4
	UniqueGiftColors     *UniqueGiftColors `json:"unique_gift_colors,omitempty"`      // Since: Bot API 9.3
	PaidMessageStarCount *int              `json:"paid_message_star_count,omitempty"` // Since: Bot API 9.3
}

// ChatPhoto represents a chat photo.
// Since: Bot API 3.1
// See https://core.telegram.org/bots/api#chatphoto
type ChatPhoto struct {
	SmallFileID       string `json:"small_file_id"`
	SmallFileUniqueID string `json:"small_file_unique_id"`
	BigFileID         string `json:"big_file_id"`
	BigFileUniqueID   string `json:"big_file_unique_id"`
}

// ChatPermissions describes actions that a non‑administrator user is allowed to take in a chat.
// Since: Bot API 4.4
// See https://core.telegram.org/bots/api#chatpermissions
type ChatPermissions struct {
	CanSendMessages      bool `json:"can_send_messages"`
	CanSendAudios        bool `json:"can_send_audios"`      // Since: Bot API 6.5
	CanSendDocuments     bool `json:"can_send_documents"`   // Since: Bot API 6.5
	CanSendPhotos        bool `json:"can_send_photos"`      // Since: Bot API 6.5
	CanSendVideos        bool `json:"can_send_videos"`      // Since: Bot API 6.5
	CanSendVideoNotes    bool `json:"can_send_video_notes"` // Since: Bot API 6.5
	CanSendVoiceNotes    bool `json:"can_send_voice_notes"` // Since: Bot API 6.5
	CanSendPolls         bool `json:"can_send_polls"`
	CanSendOtherMessages bool `json:"can_send_other_messages"`
	CanAddWebPagePreview bool `json:"can_add_web_page_previews"`
	CanReactToMessages   bool `json:"can_react_to_messages"` // Since: Bot API 10.0
	CanEditTag           bool `json:"can_edit_tag"`          // Since: Bot API 9.5
	CanChangeInfo        bool `json:"can_change_info"`
	CanInviteUsers       bool `json:"can_invite_users"`
	CanPinMessages       bool `json:"can_pin_messages"`
	CanManageTopics      bool `json:"can_manage_topics"` // Since: Bot API 6.3
}

// ChatLocation represents a location to which a chat is connected.
// Since: Bot API 5.0
// See https://core.telegram.org/bots/api#chatlocation
type ChatLocation struct {
	Location Location `json:"location"`
	Address  string   `json:"address"`
}

// ChatInviteLink represents an invite link for a chat.
// Since: Bot API 5.1
// See https://core.telegram.org/bots/api#chatinvitelink
type ChatInviteLink struct {
	InviteLink        string `json:"invite_link"`
	Creator           User   `json:"creator"`
	CreateJoinRequest bool   `json:"creates_join_request"`
	IsPrimary         bool   `json:"is_primary"`
	IsRevoked         bool   `json:"is_revoked"`

	Name                    *string `json:"name,omitempty"`
	ExpireDate              *int    `json:"expire_date,omitempty"`
	MemberLimit             *int    `json:"member_limit,omitempty"`
	PendingJoinRequestCount *int    `json:"pending_join_request_count,omitempty"`
	SubscriptionPeriod      *int    `json:"subscription_period,omitempty"`
	SubscriptionPrice       *int    `json:"subscription_price,omitempty"`
}

// ChatMemberStatusType indicates the status of a chat member.
type ChatMemberStatusType string

const (
	// ChatMemberStatusOwner identifies a chat owner.
	ChatMemberStatusOwner ChatMemberStatusType = "owner"
	// ChatMemberStatusAdministrator identifies a chat administrator.
	ChatMemberStatusAdministrator ChatMemberStatusType = "administrator"
	// ChatMemberStatusMember identifies a regular member.
	ChatMemberStatusMember ChatMemberStatusType = "member"
	// ChatMemberStatusRestricted identifies a restricted member.
	ChatMemberStatusRestricted ChatMemberStatusType = "restricted"
	// ChatMemberStatusLeft identifies a user who left the chat.
	ChatMemberStatusLeft ChatMemberStatusType = "left"
	// ChatMemberStatusBanned identifies a banned user.
	ChatMemberStatusBanned ChatMemberStatusType = "kicked"
)

// ChatMember contains information about one member of a chat.
// Since: Bot API 3.1
// See https://core.telegram.org/bots/api#chatmember
type ChatMember struct {
	Status ChatMemberStatusType `json:"status"`
	User   User                 `json:"user"`
	Tag    string               `json:"tag,omitempty"` // Since: Bot API 9.5

	// Owner
	IsAnonymous *bool   `json:"is_anonymous"`
	CustomTitle *string `json:"custom_title,omitempty"`

	// Administrator
	CanBeEdited         *bool `json:"can_be_edited,omitempty"`
	CanManageChat       *bool `json:"can_manage_chat,omitempty"`
	CanDeleteMessages   *bool `json:"can_delete_messages,omitempty"`
	CanManageVideoChats *bool `json:"can_manage_video_chats,omitempty"`
	CanRestrictMembers  *bool `json:"can_restrict_members,omitempty"`
	CanPromoteMembers   *bool `json:"can_promote_members,omitempty"`
	CanChangeInfo       *bool `json:"can_change_info,omitempty"`
	CanInviteUsers      *bool `json:"can_invite_users,omitempty"`
	CanPostStories      *bool `json:"can_post_stories,omitempty"`   // Since: Bot API 6.9
	CanEditStories      *bool `json:"can_edit_stories,omitempty"`   // Since: Bot API 6.9
	CanDeleteStories    *bool `json:"can_delete_stories,omitempty"` // Since: Bot API 6.9

	CanPostMessages         *bool `json:"can_post_messages,omitempty"`
	CanEditMessages         *bool `json:"can_edit_messages,omitempty"`
	CanPinMessages          *bool `json:"can_pin_messages,omitempty"`
	CanManageTopics         *bool `json:"can_manage_topics,omitempty"`          // Since: Bot API 6.3
	CanManageDirectMessages *bool `json:"can_manage_direct_messages,omitempty"` // Since: Bot API 9.1
	CanManageTags           *bool `json:"can_manage_tags,omitempty"`            // Since: Bot API 9.5

	// Member
	UntilDate *int `json:"until_date,omitempty"`

	// Restricted
	IsMember             *bool `json:"is_member,omitempty"`
	CanSendMessages      *bool `json:"can_send_messages,omitempty"`
	CanSendAudios        *bool `json:"can_send_audios,omitempty"`      // Since: Bot API 6.5
	CanSendDocuments     *bool `json:"can_send_documents,omitempty"`   // Since: Bot API 6.5
	CanSendPhotos        *bool `json:"can_send_photos,omitempty"`      // Since: Bot API 6.5
	CanSendVideos        *bool `json:"can_send_videos,omitempty"`      // Since: Bot API 6.5
	CanSendVideoNotes    *bool `json:"can_send_video_notes,omitempty"` // Since: Bot API 6.5
	CanSendVoiceNotes    *bool `json:"can_send_voice_notes,omitempty"` // Since: Bot API 6.5
	CanSendPolls         *bool `json:"can_send_polls,omitempty"`
	CanSendOtherMessages *bool `json:"can_send_other_messages,omitempty"`
	CanAddWebPagePreview *bool `json:"can_add_web_page_previews,omitempty"`
	CanReactToMessages   *bool `json:"can_react_to_messages,omitempty"` // Since: Bot API 10.0
	CanEditTag           *bool `json:"can_edit_tag,omitempty"`          // Since: Bot API 9.5
}

// ChatBoostSource describes the source of a chat boost.
// Since: Bot API 7.0
// See https://core.telegram.org/bots/api#chatboostsource
type ChatBoostSource struct {
	Source string `json:"source"`
	User   User   `json:"user"`

	// Giveaway
	GiveawayMessageID *int  `json:"giveaway_message_id,omitempty"`
	PrizeStarCount    *int  `json:"prize_star_count,omitempty"`
	IsUnclaimed       *bool `json:"is_unclaimed,omitempty"`
}

// ChatBoost represents a boost added to a chat.
// Since: Bot API 7.0
// See https://core.telegram.org/bots/api#chatboost
type ChatBoost struct {
	BoostID        string          `json:"boost_id"`
	AddDate        int             `json:"add_date"`
	ExpirationDate int             `json:"expiration_date"`
	Source         ChatBoostSource `json:"source"`
}

// UserChatBoosts represents a list of boosts a user has given to a chat.
// Since: Bot API 7.0
// See https://core.telegram.org/bots/api#userchatboosts
type UserChatBoosts struct {
	Boosts []ChatBoost `json:"boosts"`
}

// ChatBoostAdded describes a service message about a user boosting a chat.
// Since: Bot API 7.1
type ChatBoostAdded struct {
	BoostCount int `json:"boost_count"`
}

// ChatBackground represents a chat background.
// Since: Bot API 7.5
type ChatBackground struct {
	Type BackgroundType `json:"type"`
}

// ChatOwnerLeft describes a service message about a chat owner leaving.
// Since: Bot API 9.4
// See https://core.telegram.org/bots/api#chatownerleft
type ChatOwnerLeft struct {
	NewOwner *User `json:"new_owner,omitempty"`
}

// ChatOwnerChanged describes a service message about a chat owner change.
// Since: Bot API 9.4
// See https://core.telegram.org/bots/api#chatownerchanged
type ChatOwnerChanged struct {
	NewOwner User `json:"new_owner"`
}

// ChatAdministratorRights represents the rights of an administrator in a chat.
// Since: Bot API 6.0
// See https://core.telegram.org/bots/api#chatadministratorrights
type ChatAdministratorRights struct {
	IsAnonymous         bool `json:"is_anonymous"`
	CanManageChat       bool `json:"can_manage_chat"`
	CanDeleteMessages   bool `json:"can_delete_messages"`
	CanManageVideoChats bool `json:"can_manage_video_chats"`
	CanRestrictMembers  bool `json:"can_restrict_members"`
	CanPromoteMembers   bool `json:"can_promote_members"`
	CanChangeInfo       bool `json:"can_change_info"`
	CanInviteUsers      bool `json:"can_invite_users"`
	CanPostStories      bool `json:"can_post_stories"`
	CanEditStories      bool `json:"can_edit_stories"`
	CanDeleteStories    bool `json:"can_delete_stories"`

	CanPostMessages         *bool `json:"can_post_messages,omitempty"`
	CanEditMessages         *bool `json:"can_edit_messages,omitempty"`
	CanPinMessages          *bool `json:"can_pin_messages,omitempty"`
	CanManageTopics         *bool `json:"can_manage_topics,omitempty"`
	CanManageDirectMessages *bool `json:"can_manage_direct_messages,omitempty"`
	CanManageTags           *bool `json:"can_manage_tags,omitempty"`
}

// ChatBoostUpdated represents a boost added to a chat or changed.
// Since: Bot API 7.0
// See https://core.telegram.org/bots/api#chatboostupdated
type ChatBoostUpdated struct {
	Chat  Chat      `json:"chat"`
	Boost ChatBoost `json:"boost"`
}

// ChatBoostRemoved represents a boost removed from a chat.
// Since: Bot API 7.0
// See https://core.telegram.org/bots/api#chatboostremoved
type ChatBoostRemoved struct {
	Chat       Chat            `json:"chat"`
	BoostID    string          `json:"boost_id"`
	RemoveDate int             `json:"remove_date"`
	Source     ChatBoostSource `json:"source"`
}
