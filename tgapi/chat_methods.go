package tgapi

// BanChatMemberP holds parameters for the banChatMember method.
// See https://core.telegram.org/bots/api#banchatmember
type BanChatMemberP struct {
	ChatID         int64 `json:"chat_id"`
	UserID         int64 `json:"user_id"`
	UntilDate      int   `json:"until_date,omitempty"`
	RevokeMessages bool  `json:"revoke_messages,omitempty"`
}

// BanChatMember bans a user in a chat.
// Returns True on success.
// See https://core.telegram.org/bots/api#banchatmember
func (api *API) BanChatMember(params BanChatMemberP) (bool, error) {
	req := NewRequestWithChatID[bool]("banChatMember", params, params.ChatID)
	return req.Do(api)
}

// UnbanChatMemberP holds parameters for the unbanChatMember method.
// See https://core.telegram.org/bots/api#unbanchatmember
type UnbanChatMemberP struct {
	ChatID       int64 `json:"chat_id"`
	UserID       int64 `json:"user_id"`
	OnlyIfBanned bool  `json:"only_if_banned"`
}

// UnbanChatMember unbans a previously banned user in a chat.
// Returns True on success.
// See https://core.telegram.org/bots/api#unbanchatmember
func (api *API) UnbanChatMember(params UnbanChatMemberP) (bool, error) {
	req := NewRequestWithChatID[bool]("unbanChatMember", params, params.ChatID)
	return req.Do(api)
}

// RestrictChatMemberP holds parameters for the restrictChatMember method.
// See https://core.telegram.org/bots/api#restrictchatmember
type RestrictChatMemberP struct {
	ChatID                        int64           `json:"chat_id"`
	UserID                        int64           `json:"user_id"`
	Permissions                   ChatPermissions `json:"permissions"`
	UseIndependentChatPermissions bool            `json:"use_independent_chat_permissions,omitempty"`
	UntilDate                     int             `json:"until_date,omitempty"`
}

// RestrictChatMember restricts a user in a chat.
// Returns True on success.
// See https://core.telegram.org/bots/api#restrictchatmember
func (api *API) RestrictChatMember(params RestrictChatMemberP) (bool, error) {
	req := NewRequestWithChatID[bool]("restrictChatMember", params, params.ChatID)
	return req.Do(api)
}

// PromoteChatMember holds parameters for the promoteChatMember method.
// See https://core.telegram.org/bots/api#promotechatmember
type PromoteChatMember struct {
	ChatID      int64 `json:"chat_id"`
	UserID      int64 `json:"user_id"`
	IsAnonymous bool  `json:"is_anonymous,omitempty"`

	CanManageChat           bool `json:"can_manage_chat,omitempty"`
	CanDeleteMessages       bool `json:"can_delete_messages,omitempty"`
	CanManageVideoChats     bool `json:"can_manage_video_chats,omitempty"`
	CanRestrictMembers      bool `json:"can_restrict_members,omitempty"`
	CanPromoteMembers       bool `json:"can_promote_members,omitempty"`
	CanChangeInfo           bool `json:"can_change_info,omitempty"`
	CanInviteUsers          bool `json:"can_invite_users,omitempty"`
	CanPostStories          bool `json:"can_post_stories,omitempty"`
	CanEditStories          bool `json:"can_edit_stories,omitempty"`
	CanDeleteStories        bool `json:"can_delete_stories,omitempty"`
	CanPostMessages         bool `json:"can_post_messages,omitempty"`
	CanEditMessages         bool `json:"can_edit_messages,omitempty"`
	CanPinMessages          bool `json:"can_pin_messages,omitempty"`
	CanManageTopics         bool `json:"can_manage_topics,omitempty"`
	CanManageDirectMessages bool `json:"can_manage_direct_messages,omitempty"`
	CanManageTags           bool `json:"can_manage_tags,omitempty"`
}

// PromoteChatMember promotes or demotes a user in a chat.
// Returns True on success.
// See https://core.telegram.org/bots/api#promotechatmember
func (api *API) PromoteChatMember(params PromoteChatMember) (bool, error) {
	req := NewRequestWithChatID[bool]("promoteChatMember", params, params.ChatID)
	return req.Do(api)
}

// SetChatAdministratorCustomTitleP holds parameters for the setChatAdministratorCustomTitle method.
// See https://core.telegram.org/bots/api#setchatadministratorcustomtitle
type SetChatAdministratorCustomTitleP struct {
	ChatID      int64  `json:"chat_id"`
	UserID      int64  `json:"user_id"`
	CustomTitle string `json:"custom_title"`
}

// SetChatAdministratorCustomTitle sets a custom title for an administrator.
// Returns True on success.
// See https://core.telegram.org/bots/api#setchatadministratorcustomtitle
func (api *API) SetChatAdministratorCustomTitle(params SetChatAdministratorCustomTitleP) (bool, error) {
	req := NewRequestWithChatID[bool]("setChatAdministratorCustomTitle", params, params.ChatID)
	return req.Do(api)
}

// SetChatMemberTagP holds parameters for the setChatMemberTag method.
// See https://core.telegram.org/bots/api#setchatmembertag
type SetChatMemberTagP struct {
	ChatID int64  `json:"chat_id"`
	UserID int64  `json:"user_id"`
	Tag    string `json:"tag,omitempty"`
}

// SetChatMemberTag sets a tag for a chat member.
// Returns True on success.
// See https://core.telegram.org/bots/api#setchatmembertag
func (api *API) SetChatMemberTag(params SetChatMemberTagP) (bool, error) {
	req := NewRequestWithChatID[bool]("setChatMemberTag", params, params.ChatID)
	return req.Do(api)
}

// BanChatSenderChatP holds parameters for the banChatSenderChat method.
// See https://core.telegram.org/bots/api#banchatsenderchat
type BanChatSenderChatP struct {
	ChatID       int64 `json:"chat_id"`
	SenderChatID int64 `json:"sender_chat_id"`
}

// BanChatSenderChat bans a channel chat in a supergroup or channel.
// Returns True on success.
// See https://core.telegram.org/bots/api#banchatsenderchat
func (api *API) BanChatSenderChat(params BanChatSenderChatP) (bool, error) {
	req := NewRequestWithChatID[bool]("banChatSenderChat", params, params.ChatID)
	return req.Do(api)
}

// UnbanChatSenderChatP holds parameters for the unbanChatSenderChat method.
// See https://core.telegram.org/bots/api#unbanchatsenderchat
type UnbanChatSenderChatP struct {
	ChatID       int64 `json:"chat_id"`
	SenderChatID int64 `json:"sender_chat_id"`
}

// UnbanChatSenderChat unbans a previously banned channel chat.
// Returns True on success.
// See https://core.telegram.org/bots/api#unbanchatsenderchat
func (api *API) UnbanChatSenderChat(params UnbanChatSenderChatP) (bool, error) {
	req := NewRequestWithChatID[bool]("unbanChatSenderChat", params, params.ChatID)
	return req.Do(api)
}

// SetChatPermissionsP holds parameters for the setChatPermissions method.
// See https://core.telegram.org/bots/api#setchatpermissions
type SetChatPermissionsP struct {
	ChatID                        int64           `json:"chat_id"`
	Permissions                   ChatPermissions `json:"permissions"`
	UseIndependentChatPermissions bool            `json:"use_independent_chat_permissions,omitempty"`
}

// SetChatPermissions sets default chat permissions for all members.
// Returns True on success.
// See https://core.telegram.org/bots/api#setchatpermissions
func (api *API) SetChatPermissions(params SetChatPermissionsP) (bool, error) {
	req := NewRequestWithChatID[bool]("setChatPermissions", params, params.ChatID)
	return req.Do(api)
}

// ExportChatInviteLinkP holds parameters for the exportChatInviteLink method.
// See https://core.telegram.org/bots/api#exportchatinvitelink
type ExportChatInviteLinkP struct {
	ChatID int64 `json:"chat_id"`
}

// ExportChatInviteLink generates a new primary invite link for a chat.
// Returns the new invite link as string.
// See https://core.telegram.org/bots/api#exportchatinvitelink
func (api *API) ExportChatInviteLink(params ExportChatInviteLinkP) (string, error) {
	req := NewRequestWithChatID[string]("exportChatInviteLink", params, params.ChatID)
	return req.Do(api)
}

// CreateChatInviteLinkP holds parameters for the createChatInviteLink method.
// See https://core.telegram.org/bots/api#createchatinvitelink
type CreateChatInviteLinkP struct {
	ChatID             int64   `json:"chat_id"`
	Name               *string `json:"name,omitempty"`
	ExpireDate         int     `json:"expire_date,omitempty"`
	MemberLimit        int     `json:"member_limit,omitempty"`
	CreatesJoinRequest bool    `json:"creates_join_request,omitempty"`
}

// CreateChatInviteLink creates an additional invite link for a chat.
// Returns the created invite link.
// See https://core.telegram.org/bots/api#createchatinvitelink
func (api *API) CreateChatInviteLink(params CreateChatInviteLinkP) (ChatInviteLink, error) {
	req := NewRequestWithChatID[ChatInviteLink]("createChatInviteLink", params, params.ChatID)
	return req.Do(api)
}

// EditChatInviteLinkP holds parameters for the editChatInviteLink method.
// See https://core.telegram.org/bots/api#editchatinvitelink
type EditChatInviteLinkP struct {
	ChatID     int64  `json:"chat_id"`
	InviteLink string `json:"invite_link"`

	Name               string `json:"name,omitempty"`
	ExpireDate         int    `json:"expire_date,omitempty"`
	MemberLimit        int    `json:"member_limit,omitempty"`
	CreatesJoinRequest bool   `json:"creates_join_request,omitempty"`
}

// EditChatInviteLink edits a non‑primary invite link.
// Returns the edited invite link.
// See https://core.telegram.org/bots/api#editchatinvitelink
func (api *API) EditChatInviteLink(params EditChatInviteLinkP) (ChatInviteLink, error) {
	req := NewRequestWithChatID[ChatInviteLink]("editChatInviteLink", params, params.ChatID)
	return req.Do(api)
}

// CreateChatSubscriptionInviteLinkP holds parameters for the createChatSubscriptionInviteLink method.
// See https://core.telegram.org/bots/api#createchatsubscriptioninvitelink
type CreateChatSubscriptionInviteLinkP struct {
	ChatID             int64  `json:"chat_id"`
	Name               string `json:"name,omitempty"`
	SubscriptionPeriod int    `json:"subscription_period,omitempty"`
	SubscriptionPrice  int    `json:"subscription_price,omitempty"`
}

// CreateChatSubscriptionInviteLink creates a subscription invite link for a channel chat.
// Returns the created invite link.
// See https://core.telegram.org/bots/api#createchatsubscriptioninvitelink
func (api *API) CreateChatSubscriptionInviteLink(params CreateChatSubscriptionInviteLinkP) (ChatInviteLink, error) {
	req := NewRequestWithChatID[ChatInviteLink]("createChatSubscriptionInviteLink", params, params.ChatID)
	return req.Do(api)
}

// EditChatSubscriptionInviteLinkP holds parameters for the editChatSubscriptionInviteLink method.
// See https://core.telegram.org/bots/api#editchatsubscriptioninvitelink
type EditChatSubscriptionInviteLinkP struct {
	ChatID     int64  `json:"chat_id"`
	InviteLink string `json:"invite_link"`
	Name       string `json:"name,omitempty"`
}

// EditChatSubscriptionInviteLink edits a subscription invite link.
// Returns the edited invite link.
// See https://core.telegram.org/bots/api#editchatsubscriptioninvitelink
func (api *API) EditChatSubscriptionInviteLink(params EditChatSubscriptionInviteLinkP) (ChatInviteLink, error) {
	req := NewRequestWithChatID[ChatInviteLink]("editChatSubscriptionInviteLink", params, params.ChatID)
	return req.Do(api)
}

// RevokeChatInviteLinkP holds parameters for the revokeChatInviteLink method.
// See https://core.telegram.org/bots/api#revokechatinvitelink
type RevokeChatInviteLinkP struct {
	ChatID     int64  `json:"chat_id"`
	InviteLink string `json:"invite_link"`
}

// RevokeChatInviteLink revokes an invite link.
// Returns the revoked invite link object.
// See https://core.telegram.org/bots/api#revokechatinvitelink
func (api *API) RevokeChatInviteLink(params RevokeChatInviteLinkP) (ChatInviteLink, error) {
	req := NewRequestWithChatID[ChatInviteLink]("revokeChatInviteLink", params, params.ChatID)
	return req.Do(api)
}

// ApproveChatJoinRequestP holds parameters for the approveChatJoinRequest method.
// See https://core.telegram.org/bots/api#approvechatjoinrequest
type ApproveChatJoinRequestP struct {
	ChatID int64 `json:"chat_id"`
	UserID int64 `json:"user_id"`
}

// ApproveChatJoinRequest approves a chat join request.
// Returns True on success.
// See https://core.telegram.org/bots/api#approvechatjoinrequest
func (api *API) ApproveChatJoinRequest(params ApproveChatJoinRequestP) (bool, error) {
	req := NewRequestWithChatID[bool]("approveChatJoinRequest", params, params.ChatID)
	return req.Do(api)
}

// DeclineChatJoinRequestP holds parameters for the declineChatJoinRequest method.
// See https://core.telegram.org/bots/api#declinechatjoinrequest
type DeclineChatJoinRequestP struct {
	ChatID int64 `json:"chat_id"`
	UserID int64 `json:"user_id"`
}

// DeclineChatJoinRequest declines a chat join request.
// Returns True on success.
// See https://core.telegram.org/bots/api#declinechatjoinrequest
func (api *API) DeclineChatJoinRequest(params DeclineChatJoinRequestP) (bool, error) {
	req := NewRequestWithChatID[bool]("declineChatJoinRequest", params, params.ChatID)
	return req.Do(api)
}

// SetChatPhotoP holds parameters for the setChatPhoto method.
// See https://core.telegram.org/bots/api#setchatphoto
type SetChatPhotoP struct {
	ChatID int64 `json:"chat_id"`
}

// SetChatPhoto changes the chat photo.
// photo is the file to upload as the new photo.
// Returns True on success.
// See https://core.telegram.org/bots/api#setchatphoto
func (api *API) SetChatPhoto(params SetChatPhotoP, photo UploaderFile) (bool, error) {
	uploader := NewUploader(api)
	defer func() {
		_ = uploader.Close()
	}()
	req := NewUploaderRequestWithChatID[bool]("setChatPhoto", params, params.ChatID, photo.SetType(UploaderPhotoType))
	return req.Do(uploader)
}

// DeleteChatPhotoP holds parameters for the deleteChatPhoto method.
// See https://core.telegram.org/bots/api#deletechatphoto
type DeleteChatPhotoP struct {
	ChatID int64 `json:"chat_id"`
}

// DeleteChatPhoto deletes a chat photo.
// Returns True on success.
// See https://core.telegram.org/bots/api#deletechatphoto
func (api *API) DeleteChatPhoto(params DeleteChatPhotoP) (bool, error) {
	req := NewRequestWithChatID[bool]("deleteChatPhoto", params, params.ChatID)
	return req.Do(api)
}

// SetChatTitleP holds parameters for the setChatTitle method.
// See https://core.telegram.org/bots/api#setchattitle
type SetChatTitleP struct {
	ChatID int64  `json:"chat_id"`
	Title  string `json:"title"`
}

// SetChatTitle changes the chat title.
// Returns True on success.
// See https://core.telegram.org/bots/api#setchattitle
func (api *API) SetChatTitle(params SetChatTitleP) (bool, error) {
	req := NewRequestWithChatID[bool]("setChatTitle", params, params.ChatID)
	return req.Do(api)
}

// SetChatDescriptionP holds parameters for the setChatDescription method.
// See https://core.telegram.org/bots/api#setchatdescription
type SetChatDescriptionP struct {
	ChatID      int64  `json:"chat_id"`
	Description string `json:"description"`
}

// SetChatDescription changes the chat description.
// Returns True on success.
// See https://core.telegram.org/bots/api#setchatdescription
func (api *API) SetChatDescription(params SetChatDescriptionP) (bool, error) {
	req := NewRequestWithChatID[bool]("setChatDescription", params, params.ChatID)
	return req.Do(api)
}

// PinChatMessageP holds parameters for the pinChatMessage method.
// See https://core.telegram.org/bots/api#pinchatmessage
type PinChatMessageP struct {
	BusinessConnectionID *string `json:"business_connection_id,omitempty"`
	ChatID               int64   `json:"chat_id"`
	MessageID            int     `json:"message_id"`
	DisableNotification  bool    `json:"disable_notification,omitempty"`
}

// PinChatMessage pins a message in a chat.
// Returns True on success.
// See https://core.telegram.org/bots/api#pinchatmessage
func (api *API) PinChatMessage(params PinChatMessageP) (bool, error) {
	req := NewRequestWithChatID[bool]("pinChatMessage", params, params.ChatID)
	return req.Do(api)
}

// UnpinChatMessageP holds parameters for the unpinChatMessage method.
// See https://core.telegram.org/bots/api#unpinchatmessage
type UnpinChatMessageP struct {
	BusinessConnectionID *string `json:"business_connection_id,omitempty"`
	ChatID               int64   `json:"chat_id"`
	MessageID            int     `json:"message_id"`
}

// UnpinChatMessage unpins a message in a chat.
// Returns True on success.
// See https://core.telegram.org/bots/api#unpinchatmessage
func (api *API) UnpinChatMessage(params UnpinChatMessageP) (bool, error) {
	req := NewRequestWithChatID[bool]("unpinChatMessage", params, params.ChatID)
	return req.Do(api)
}

// UnpinAllChatMessagesP holds parameters for the unpinAllChatMessages method.
// See https://core.telegram.org/bots/api#unpinallchatmessages
type UnpinAllChatMessagesP struct {
	ChatID int64 `json:"chat_id"`
}

// UnpinAllChatMessages unpins all pinned messages in a chat.
// Returns True on success.
// See https://core.telegram.org/bots/api#unpinallchatmessages
func (api *API) UnpinAllChatMessages(params UnpinAllChatMessagesP) (bool, error) {
	req := NewRequestWithChatID[bool]("unpinAllChatMessages", params, params.ChatID)
	return req.Do(api)
}

// LeaveChatP holds parameters for the leaveChat method.
// See https://core.telegram.org/bots/api#leavechat
type LeaveChatP struct {
	ChatID int64 `json:"chat_id"`
}

// LeaveChat makes the bot leave a chat.
// Returns True on success.
// See https://core.telegram.org/bots/api#leavechat
func (api *API) LeaveChat(params LeaveChatP) (bool, error) {
	req := NewRequestWithChatID[bool]("leaveChat", params, params.ChatID) // fixed method name
	return req.Do(api)
}

// GetChatP holds parameters for the getChat method.
// See https://core.telegram.org/bots/api#getchat
type GetChatP struct {
	ChatID int64 `json:"chat_id"`
}

// GetChat gets up‑to‑date information about a chat.
// See https://core.telegram.org/bots/api#getchat
func (api *API) GetChat(params GetChatP) (ChatFullInfo, error) {
	req := NewRequestWithChatID[ChatFullInfo]("getChat", params, params.ChatID) // fixed method name
	return req.Do(api)
}

// GetChatAdministratorsP holds parameters for the getChatAdministrators method.
// See https://core.telegram.org/bots/api#getchatadministrators
type GetChatAdministratorsP struct {
	ChatID int64 `json:"chat_id"`
}

// GetChatAdministrators returns a list of administrators in a chat.
// See https://core.telegram.org/bots/api#getchatadministrators
func (api *API) GetChatAdministrators(params GetChatAdministratorsP) ([]ChatMember, error) {
	req := NewRequestWithChatID[[]ChatMember]("getChatAdministrators", params, params.ChatID)
	return req.Do(api)
}

// GetChatMembersCountP holds parameters for the getChatMemberCount method.
// See https://core.telegram.org/bots/api#getchatmembercount
type GetChatMembersCountP struct {
	ChatID int64 `json:"chat_id"`
}

// GetChatMemberCount returns the number of members in a chat.
// See https://core.telegram.org/bots/api#getchatmembercount
func (api *API) GetChatMemberCount(params GetChatMembersCountP) (int, error) {
	req := NewRequestWithChatID[int]("getChatMemberCount", params, params.ChatID)
	return req.Do(api)
}

// GetChatMemberP holds parameters for the getChatMember method.
// See https://core.telegram.org/bots/api#getchatmember
type GetChatMemberP struct {
	ChatID int64 `json:"chat_id"`
	UserID int64 `json:"user_id"`
}

// GetChatMember returns information about a member of a chat.
// See https://core.telegram.org/bots/api#getchatmember
func (api *API) GetChatMember(params GetChatMemberP) (ChatMember, error) {
	req := NewRequestWithChatID[ChatMember]("getChatMember", params, params.ChatID)
	return req.Do(api)
}

// SetChatStickerSetP holds parameters for the setChatStickerSet method.
// See https://core.telegram.org/bots/api#setchatstickerset
type SetChatStickerSetP struct {
	ChatID         int64  `json:"chat_id"`
	StickerSetName string `json:"sticker_set_name"`
}

// SetChatStickerSet associates a sticker set with a supergroup.
// Returns True on success.
// See https://core.telegram.org/bots/api#setchatstickerset
func (api *API) SetChatStickerSet(params SetChatStickerSetP) (bool, error) {
	req := NewRequestWithChatID[bool]("setChatStickerSet", params, params.ChatID)
	return req.Do(api)
}

// DeleteChatStickerSetP holds parameters for the deleteChatStickerSet method.
// See https://core.telegram.org/bots/api#deletechatstickerset
type DeleteChatStickerSetP struct {
	ChatID int64 `json:"chat_id"`
}

// DeleteChatStickerSet deletes a sticker set from a supergroup.
// Returns True on success.
// See https://core.telegram.org/bots/api#deletechatstickerset
func (api *API) DeleteChatStickerSet(params DeleteChatStickerSetP) (bool, error) {
	req := NewRequestWithChatID[bool]("deleteChatStickerSet", params, params.ChatID)
	return req.Do(api)
}

// GetUserChatBoostsP holds parameters for the getUserChatBoosts method.
// See https://core.telegram.org/bots/api#getuserchatboosts
type GetUserChatBoostsP struct {
	ChatID int64 `json:"chat_id"`
	UserID int64 `json:"user_id"`
}

// GetUserChatBoosts returns the list of boosts a user has given to a chat.
// See https://core.telegram.org/bots/api#getuserchatboosts
func (api *API) GetUserChatBoosts(params GetUserChatBoostsP) (UserChatBoosts, error) {
	req := NewRequestWithChatID[UserChatBoosts]("getUserChatBoosts", params, params.ChatID)
	return req.Do(api)
}

// GetChatGiftsP holds parameters for the getChatGifts method.
// See https://core.telegram.org/bots/api#getchatgifts
type GetChatGiftsP struct {
	ChatID                      int64  `json:"chat_id"`
	ExcludeUnsaved              bool   `json:"exclude_unsaved,omitempty"`
	ExcludeSaved                bool   `json:"exclude_saved,omitempty"`
	ExcludeUnlimited            bool   `json:"exclude_unlimited,omitempty"`
	ExcludeLimitedUpgradable    bool   `json:"exclude_limited_upgradable,omitempty"`
	ExcludeLimitedNonUpgradable bool   `json:"exclude_limited_non_upgradable,omitempty"`
	ExcludeUnique               bool   `json:"exclude_unique,omitempty"`
	ExcludeFromBlockchain       bool   `json:"exclude_from_blockchain,omitempty"`
	SortByPrice                 bool   `json:"sort_by_price,omitempty"`
	Offset                      string `json:"offset,omitempty"`
	Limit                       int    `json:"limit,omitempty"`
}

// GetChatGifts returns gifts owned by a chat.
// See https://core.telegram.org/bots/api#getchatgifts
func (api *API) GetChatGifts(params GetChatGiftsP) (OwnedGifts, error) {
	req := NewRequestWithChatID[OwnedGifts]("getChatGifts", params, params.ChatID)
	return req.Do(api)
}
