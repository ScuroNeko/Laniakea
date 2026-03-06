package tgapi

type BanChatMemberP struct {
	ChatID         int64 `json:"chat_id"`
	UserID         int   `json:"user_id"`
	UntilDate      int   `json:"until_date,omitempty"`
	RevokeMessages bool  `json:"revoke_messages,omitempty"`
}

func (api *API) BanChatMember(params BanChatMemberP) (bool, error) {
	req := NewRequestWithChatID[bool]("banChatMember", params, params.ChatID)
	return req.Do(api)
}

type UnbanChatMemberP struct {
	ChatID       int64 `json:"chat_id"`
	UserID       int   `json:"user_id"`
	OnlyIfBanned bool  `json:"only_if_banned"`
}

func (api *API) UnbanChatMember(params UnbanChatMemberP) (bool, error) {
	req := NewRequestWithChatID[bool]("unbanChatMember", params, params.ChatID)
	return req.Do(api)
}

type RestrictChatMemberP struct {
	ChatID                        int64           `json:"chat_id"`
	UserID                        int             `json:"user_id"`
	Permissions                   ChatPermissions `json:"permissions"`
	UseIndependentChatPermissions bool            `json:"use_independent_chat_permissions,omitempty"`
	UntilDate                     int             `json:"until_date,omitempty"`
}

func (api *API) RestrictChatMember(params RestrictChatMemberP) (bool, error) {
	req := NewRequestWithChatID[bool]("restrictChatMember", params, params.ChatID)
	return req.Do(api)
}

type PromoteChatMember struct {
	ChatID      int64 `json:"chat_id"`
	UserID      int   `json:"user_id"`
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

func (api *API) PromoteChatMember(params PromoteChatMember) (bool, error) {
	req := NewRequestWithChatID[bool]("promoteChatMember", params, params.ChatID)
	return req.Do(api)
}

type SetChatAdministratorCustomTitleP struct {
	ChatID      int64  `json:"chat_id"`
	UserID      int    `json:"user_id"`
	CustomTitle string `json:"custom_title"`
}

func (api *API) SetChatAdministratorCustomTitle(params SetChatAdministratorCustomTitleP) (bool, error) {
	req := NewRequestWithChatID[bool]("setChatAdministratorCustomTitle", params, params.ChatID)
	return req.Do(api)
}

type SetChatMemberTagP struct {
	ChatID int64  `json:"chat_id"`
	UserID int    `json:"user_id"`
	Tag    string `json:"tag,omitempty"`
}

func (api *API) SetChatMemberTag(params SetChatMemberTagP) (bool, error) {
	req := NewRequestWithChatID[bool]("setChatMemberTag", params, params.ChatID)
	return req.Do(api)
}

type BanChatSenderChatP struct {
	ChatID       int64 `json:"chat_id"`
	SenderChatID int64 `json:"sender_chat_id"`
}

func (api *API) BanChatSenderChat(params BanChatSenderChatP) (bool, error) {
	req := NewRequestWithChatID[bool]("banChatSenderChat", params, params.ChatID)
	return req.Do(api)
}

type UnbanChatSenderChatP struct {
	ChatID       int64 `json:"chat_id"`
	SenderChatID int64 `json:"sender_chat_id"`
}

func (api *API) UnbanChatSenderChat(params BanChatSenderChatP) (bool, error) {
	req := NewRequestWithChatID[bool]("unbanChatSenderChat", params, params.ChatID)
	return req.Do(api)
}

type SetChatPermissionsP struct {
	ChatID                        int64           `json:"chat_id"`
	Permissions                   ChatPermissions `json:"permissions"`
	UseIndependentChatPermissions bool            `json:"use_independent_chat_permissions,omitempty"`
}

func (api *API) SetChatPermissions(params SetChatPermissionsP) (bool, error) {
	req := NewRequestWithChatID[bool]("setChatPermissions", params, params.ChatID)
	return req.Do(api)
}

type ExportChatInviteLinkP struct {
	ChatID int64 `json:"chat_id"`
}

func (api *API) ExportChatInviteLink(params ExportChatInviteLinkP) (string, error) {
	req := NewRequestWithChatID[string]("exportChatInviteLink", params, params.ChatID)
	return req.Do(api)
}

type CreateChatInviteLinkP struct {
	ChatID             int64   `json:"chat_id"`
	Name               *string `json:"name,omitempty"`
	ExpireDate         int     `json:"expire_date,omitempty"`
	MemberLimit        int     `json:"member_limit,omitempty"`
	CreatesJoinRequest int     `json:"creates_join_request,omitempty"`
}

func (api *API) CreateChatInviteLink(params CreateChatInviteLinkP) (ChatInviteLink, error) {
	req := NewRequestWithChatID[ChatInviteLink]("createChatInviteLink", params, params.ChatID)
	return req.Do(api)
}

type EditChatInviteLinkP struct {
	ChatID     int64  `json:"chat_id"`
	InviteLink string `json:"invite_link"`

	Name               string `json:"name,omitempty"`
	ExpireDate         int    `json:"expire_date,omitempty"`
	MemberLimit        int    `json:"member_limit,omitempty"`
	CreatesJoinRequest int    `json:"creates_join_request,omitempty"`
}

func (api *API) EditChatInviteLink(params EditChatInviteLinkP) (ChatInviteLink, error) {
	req := NewRequestWithChatID[ChatInviteLink]("editChatInviteLink", params, params.ChatID)
	return req.Do(api)
}

type CreateChatSubscriptionInviteLinkP struct {
	ChatID             int64  `json:"chat_id"`
	Name               string `json:"name,omitempty"`
	SubscriptionPeriod int    `json:"subscription_period,omitempty"`
	SubscriptionPrice  int    `json:"subscription_price,omitempty"`
}

func (api *API) CreateChatSubscriptionInviteLink(params CreateChatSubscriptionInviteLinkP) (ChatInviteLink, error) {
	req := NewRequestWithChatID[ChatInviteLink]("createChatSubscriptionInviteLink", params, params.ChatID)
	return req.Do(api)
}

type EditChatSubscriptionInviteLinkP struct {
	ChatID     int64  `json:"chat_id"`
	InviteLink string `json:"invite_link"`
	Name       string `json:"name,omitempty"`
}

func (api *API) EditChatSubscriptionInviteLink(params EditChatSubscriptionInviteLinkP) (ChatInviteLink, error) {
	req := NewRequestWithChatID[ChatInviteLink]("editChatSubscriptionInviteLink", params, params.ChatID)
	return req.Do(api)
}

type RevokeChatInviteLinkP struct {
	ChatID     int64  `json:"chat_id"`
	InviteLink string `json:"invite_link"`
}

func (api *API) RevokeChatInviteLink(params RevokeChatInviteLinkP) (ChatInviteLink, error) {
	req := NewRequestWithChatID[ChatInviteLink]("revokeChatInviteLink", params, params.ChatID)
	return req.Do(api)
}

type ApproveChatJoinRequestP struct {
	ChatID int64 `json:"chat_id"`
	UserID int   `json:"user_id"`
}

func (api *API) ApproveChatJoinRequest(params ApproveChatJoinRequestP) (bool, error) {
	req := NewRequestWithChatID[bool]("approveChatJoinRequest", params, params.ChatID)
	return req.Do(api)
}

type DeclineChatJoinRequestP struct {
	ChatID int64 `json:"chat_id"`
	UserID int   `json:"user_id"`
}

func (api *API) DeclineChatJoinRequest(params DeclineChatJoinRequestP) (bool, error) {
	req := NewRequestWithChatID[bool]("declineChatJoinRequest", params, params.ChatID)
	return req.Do(api)
}

func (api *API) SetChatPhoto() {
	uploader := NewUploader(api)
	defer uploader.Close()
}

type DeleteChatPhotoP struct {
	ChatID int64 `json:"chat_id"`
}

func (api *API) DeleteChatPhoto(params DeleteChatPhotoP) (bool, error) {
	req := NewRequestWithChatID[bool]("deleteChatPhoto", params, params.ChatID)
	return req.Do(api)
}

type SetChatTitleP struct {
	ChatID int64  `json:"chat_id"`
	Title  string `json:"title"`
}

func (api *API) SetChatTitle(params SetChatTitleP) (bool, error) {
	req := NewRequestWithChatID[bool]("setChatTitle", params, params.ChatID)
	return req.Do(api)
}

type SetChatDescriptionP struct {
	ChatID      int64  `json:"chat_id"`
	Description string `json:"description"`
}

func (api *API) SetChatDescription(params SetChatDescriptionP) (bool, error) {
	req := NewRequestWithChatID[bool]("setChatDescription", params, params.ChatID)
	return req.Do(api)
}

type PinChatMessageP struct {
	BusinessConnectionID *string `json:"business_connection_id,omitempty"`
	ChatID               int64   `json:"chat_id"`
	MessageID            int     `json:"message_id"`
	DisableNotification  bool    `json:"disable_notification,omitempty"`
}

func (api *API) PinChatMessage(params PinChatMessageP) (bool, error) {
	req := NewRequestWithChatID[bool]("pinChatMessage", params, params.ChatID)
	return req.Do(api)
}

type UnpinChatMessageP struct {
	BusinessConnectionID *string `json:"business_connection_id,omitempty"`
	ChatID               int64   `json:"chat_id"`
	MessageID            int     `json:"message_id"`
}

func (api *API) UnpinChatMessage(params UnpinChatMessageP) (bool, error) {
	req := NewRequestWithChatID[bool]("unpinChatMessage", params, params.ChatID)
	return req.Do(api)
}

type UnpinAllChatMessagesP struct {
	ChatID int64 `json:"chat_id"`
}

func (api *API) UnpinAllChatMessages(params UnpinAllChatMessagesP) (bool, error) {
	req := NewRequestWithChatID[bool]("unpinAllChatMessages", params, params.ChatID)
	return req.Do(api)
}

type LeaveChatP struct {
	ChatID int64 `json:"chat_id"`
}

func (api *API) LeaveChat(params LeaveChatP) (bool, error) {
	req := NewRequestWithChatID[bool]("leaveChatP", params, params.ChatID)
	return req.Do(api)
}

type GetChatP struct {
	ChatID int64 `json:"chat_id"`
}

func (api *API) GetChatP(params GetChatP) (ChatFullInfo, error) {
	req := NewRequestWithChatID[ChatFullInfo]("getChatP", params, params.ChatID)
	return req.Do(api)
}

type GetChatAdministratorsP struct {
	ChatID int64 `json:"chat_id"`
}

func (api *API) GetChatAdministrators(params GetChatAdministratorsP) ([]ChatMember, error) {
	req := NewRequestWithChatID[[]ChatMember]("getChatAdministrators", params, params.ChatID)
	return req.Do(api)
}

type GetChatMembersCountP struct {
	ChatID int64 `json:"chat_id"`
}

func (api *API) GetChatMemberCount(params GetChatMembersCountP) (int, error) {
	req := NewRequestWithChatID[int]("getChatMemberCount", params, params.ChatID)
	return req.Do(api)
}

type GetChatMemberP struct {
	ChatID int64 `json:"chat_id"`
	UserID int   `json:"user_id"`
}

func (api *API) GetChatMember(params GetChatMemberP) (ChatMember, error) {
	req := NewRequestWithChatID[ChatMember]("getChatMember", params, params.ChatID)
	return req.Do(api)
}

type SetChatStickerSetP struct {
	ChatID         int64  `json:"chat_id"`
	StickerSetName string `json:"sticker_set_name"`
}

func (api *API) SetChatStickerSet(params SetChatStickerSetP) (bool, error) {
	req := NewRequestWithChatID[bool]("setChatStickerSet", params, params.ChatID)
	return req.Do(api)
}

type DeleteChatStickerSetP struct {
	ChatID int64 `json:"chat_id"`
}

func (api *API) DeleteChatStickerSet(params DeleteChatStickerSetP) (bool, error) {
	req := NewRequestWithChatID[bool]("deleteChatStickerSet", params, params.ChatID)
	return req.Do(api)
}

type GetUserChatBoostsP struct {
	ChatID int64 `json:"chat_id"`
	UserID int   `json:"user_id"`
}

func (api *API) GetUserChatBoosts(params GetUserChatBoostsP) (UserChatBoosts, error) {
	req := NewRequestWithChatID[UserChatBoosts]("getUserChatBoosts", params, params.ChatID)
	return req.Do(api)
}

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

func (api *API) GetChatGifts(params GetChatGiftsP) (OwnedGifts, error) {
	req := NewRequestWithChatID[OwnedGifts]("getChatGifts", params, params.ChatID)
	return req.Do(api)
}
