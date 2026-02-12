package tgapi

type Chat struct {
	ID               int     `json:"id"`
	Type             string  `json:"type"`
	Title            *string `json:"title,omitempty"`
	Username         *string `json:"username,omitempty"`
	FirstName        *string `json:"first_name,omitempty"`
	LastName         *string `json:"last_name,omitempty"`
	IsForum          *bool   `json:"is_forum,omitempty"`
	IsDirectMessages *bool   `json:"is_direct_messages,omitempty"`
}

type ChatType string

const (
	ChatTypePrivate    ChatType = "private"
	ChatTypeGroup      ChatType = "group"
	ChatTypeSupergroup ChatType = "supergroup"
	ChatTypeChannel    ChatType = "channel"
)

type ChatFullInfo struct {
	ID               int        `json:"id"`
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
	ParentChat   *Chat `json:"parent_chat,omitempty"`

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
	LinkedChatID                 *int    `json:"linked_chat_id,omitempty"`

	Location             *ChatLocation     `json:"location,omitempty"`
	Rating               *UserRating       `json:"rating,omitempty"`
	FirstProfileAudio    *Audio            `json:"first_profile_audio,omitempty"`
	UniqueGiftColors     *UniqueGiftColors `json:"unique_gift_colors,omitempty"`
	PaidMessageStarCount *int              `json:"paid_message_star_count,omitempty"`
}

type ChatPhoto struct {
	SmallFileID       string `json:"small_file_id"`
	SmallFileUniqueID string `json:"small_file_unique_id"`
	BigFileID         string `json:"big_file_id"`
	BigFileUniqueID   string `json:"big_file_unique_id"`
}

type ChatPermissions struct {
	CanSendMessages      bool `json:"can_send_messages"`
	CanSendAudios        bool `json:"can_send_audios"`
	CanSendDocuments     bool `json:"can_send_documents"`
	CanSendPhotos        bool `json:"can_send_photos"`
	CanSendVideoNotes    bool `json:"can_send_video_notes"`
	CanSendVoiceNotes    bool `json:"can_send_voice_notes"`
	CanSendPolls         bool `json:"can_send_polls"`
	CanSendOtherMessages bool `json:"can_send_other_messages"`
	CanAddWebPagePreview bool `json:"can_add_web_page_preview"`
	CanChangeInfo        bool `json:"can_change_info"`
	CanInviteUsers       bool `json:"can_invite_users"`
	CanPinMessages       bool `json:"can_pin_messages"`
	CanManageTopics      bool `json:"can_manage_topics"`
}
type ChatLocation struct {
	Location Location `json:"location"`
	Address  string   `json:"address"`
}
type ChatInviteLink struct {
	InviteLink        string `json:"invite_link"`
	Creator           User   `json:"creator"`
	CreateJoinRequest bool   `json:"create_join_request"`
	IsPrimary         bool   `json:"is_primary"`
	IsRevoked         bool   `json:"is_revoked"`

	Name                    *string `json:"name,omitempty"`
	ExpireDate              *int    `json:"expire_date,omitempty"`
	MemberLimit             *int    `json:"member_limit,omitempty"`
	PendingJoinRequestCount *int    `json:"pending_join_request_count,omitempty"`
	SubscriptionPeriod      *int    `json:"subscription_period,omitempty"`
	SubscriptionPrice       *int    `json:"subscription_price,omitempty"`
}

type ChatMemberStatusType string

const (
	ChatMemberStatusOwner         ChatMemberStatusType = "owner"
	ChatMemberStatusAdministrator ChatMemberStatusType = "administrator"
	ChatMemberStatusMember        ChatMemberStatusType = "member"
	ChatMemberStatusRestricted    ChatMemberStatusType = "restricted"
	ChatMemberStatusLeft          ChatMemberStatusType = "left"
	ChatMemberStatusBanned        ChatMemberStatusType = "kicked"
)

type ChatMember struct {
	Status ChatMemberStatusType `json:"status"`
	User   User                 `json:"user"`

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
	CanPostStories      *bool `json:"can_post_stories,omitempty"`
	CanEditStories      *bool `json:"can_edit_stories,omitempty"`
	CanDeleteStories    *bool `json:"can_delete_stories,omitempty"`

	CanPostMessages         *bool `json:"can_post_messages,omitempty"`
	CanEditMessages         *bool `json:"can_edit_messages,omitempty"`
	CanPinMessages          *bool `json:"can_pin_messages,omitempty"`
	CanManageTopics         *bool `json:"can_manage_topics,omitempty"`
	CanManageDirectMessages *bool `json:"can_manage_direct_messages,omitempty"`

	// Member
	UntilDate *int `json:"until_date,omitempty"`

	// Restricted
	IsMember             *bool `json:"is_member,omitempty"`
	CanSendMessages      *bool `json:"can_send_messages,omitempty"`
	CanSendAudios        *bool `json:"can_send_audios,omitempty"`
	CanSendDocuments     *bool `json:"can_send_documents,omitempty"`
	CanSendVideos        *bool `json:"can_send_videos,omitempty"`
	CanSendVideoNotes    *bool `json:"can_send_video_notes,omitempty"`
	CanSendVoiceNotes    *bool `json:"can_send_voice_notes,omitempty"`
	CanSendPolls         *bool `json:"can_send_polls,omitempty"`
	CanSendOtherMessages *bool `json:"can_send_other_messages,omitempty"`
	CanAddWebPagePreview *bool `json:"can_add_web_page_preview,omitempty"`
}

type ChatBoostSource struct {
	Source string `json:"source"`
	User   User   `json:"user"`

	// Giveaway
	GiveawayMessageID *int  `json:"giveaway_message_id,omitempty"`
	PrizeStarCount    *int  `json:"prize_star_count,omitempty"`
	IsUnclaimed       *bool `json:"is_unclaimed,omitempty"`
}

type ChatBoost struct {
	BoostID        int             `json:"boost_id"`
	AddDate        int             `json:"add_date"`
	ExpirationDate int             `json:"expiration_date"`
	Source         ChatBoostSource `json:"source"`
}
type UserChatBoosts struct {
	Boosts []ChatBoost `json:"boosts"`
}

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
}

type ChatBoostUpdated struct {
	Chat  Chat      `json:"chat"`
	Boost ChatBoost `json:"boost"`
}

type ChatBoostRemoved struct {
	Chat       Chat            `json:"chat"`
	BoostID    string          `json:"boost_id"`
	RemoveDate int             `json:"remove_date"`
	Source     ChatBoostSource `json:"source"`
}
