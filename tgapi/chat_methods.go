package tgapi

type BanChatMemberP struct {
	ChatID         int   `json:"chat_id"`
	UserID         int   `json:"user_id"`
	UntilDate      int64 `json:"until_date,omitempty"`
	RevokeMessages bool  `json:"revoke_messages,omitempty"`
}

func (api *Api) BanChatMember(p BanChatMemberP) (bool, error) {
	req := NewRequest[bool]("banChatMember", p)
	return req.Do(api)
}

type UnbanChatMemberP struct {
	ChatID       int  `json:"chat_id"`
	UserID       int  `json:"user_id"`
	OnlyIfBanned bool `json:"only_if_banned"`
}

func (api *Api) UnbanChatMember(p UnbanChatMemberP) (bool, error) {
	req := NewRequest[bool]("unbanChatMember", p)
	return req.Do(api)
}

type RestrictChatMemberP struct {
	ChatID                        int             `json:"chat_id"`
	UserID                        int             `json:"user_id"`
	Permissions                   ChatPermissions `json:"permissions"`
	UseIndependentChatPermissions bool            `json:"use_independent_chat_permissions,omitempty"`
	UntilDate                     int             `json:"until_date,omitempty"`
}

func (api *Api) RestrictChatMember(p RestrictChatMemberP) (bool, error) {
	req := NewRequest[bool]("restrictChatMember", p)
	return req.Do(api)
}

type PromoteChatMember struct {
	ChatID      int  `json:"chat_id"`
	UserID      int  `json:"user_id"`
	IsAnonymous bool `json:"is_anonymous,omitempty"`

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
}

func (api *Api) PromoteChatMember(p PromoteChatMember) (bool, error) {
	req := NewRequest[bool]("promoteChatMember", p)
	return req.Do(api)
}

type SetChatAdministratorCustomTitleP struct {
	ChatID      int    `json:"chat_id"`
	UserID      int    `json:"user_id"`
	CustomTitle string `json:"custom_title"`
}

func (api *Api) SetChatAdministratorCustomTitle(p SetChatAdministratorCustomTitleP) (bool, error) {
	req := NewRequest[bool]("setChatAdministratorCustomTitle", p)
	return req.Do(api)
}

type BanChatSenderChatP struct {
	ChatID       int `json:"chat_id"`
	SenderChatID int `json:"sender_chat_id"`
}

func (api *Api) BanChatSenderChat(p BanChatSenderChatP) (bool, error) {
	req := NewRequest[bool]("banChatSenderChat", p)
	return req.Do(api)
}

type UnbanChatSenderChatP struct {
	ChatID       int `json:"chat_id"`
	SenderChatID int `json:"sender_chat_id"`
}

func (api *Api) UnbanChatSenderChat(p BanChatSenderChatP) (bool, error) {
	req := NewRequest[bool]("unbanChatSenderChat", p)
	return req.Do(api)
}

type SetChatPermissionsP struct {
	ChatID                        int             `json:"chat_id"`
	Permissions                   ChatPermissions `json:"permissions"`
	UseIndependentChatPermissions bool            `json:"use_independent_chat_permissions,omitempty"`
}

func (api *Api) SetChatPermissions(p SetChatPermissionsP) (bool, error) {
	req := NewRequest[bool]("setChatPermissions", p)
	return req.Do(api)
}

type ExportChatInviteLinkP struct {
	ChatID int `json:"chat_id"`
}

func (api *Api) ExportChatInviteLink(p ExportChatInviteLinkP) (string, error) {
	req := NewRequest[string]("exportChatInviteLink", p)
	return req.Do(api)
}

type CreateChatInviteLinkP struct {
	ChatID             int    `json:"chat_id"`
	Name               string `json:"name,omitempty"`
	ExpireDate         int    `json:"expire_date,omitempty"`
	MemberLimit        int    `json:"member_limit,omitempty"`
	CreatesJoinRequest int    `json:"creates_join_request,omitempty"`
}

func (api *Api) CreateChatInviteLink(p CreateChatInviteLinkP) (ChatInviteLink, error) {
	req := NewRequest[ChatInviteLink]("createChatInviteLink", p)
	return req.Do(api)
}

type EditChatInviteLinkP struct {
	ChatID     int    `json:"chat_id"`
	InviteLink string `json:"invite_link"`

	Name               string `json:"name,omitempty"`
	ExpireDate         int    `json:"expire_date,omitempty"`
	MemberLimit        int    `json:"member_limit,omitempty"`
	CreatesJoinRequest int    `json:"creates_join_request,omitempty"`
}

func (api *Api) EditChatInviteLink(p EditChatInviteLinkP) (ChatInviteLink, error) {
	req := NewRequest[ChatInviteLink]("editChatInviteLink", p)
	return req.Do(api)
}

type CreateChatSubscriptionInviteLinkP struct {
	ChatID             int    `json:"chat_id"`
	Name               string `json:"name,omitempty"`
	SubscriptionPeriod int    `json:"subscription_period,omitempty"`
	SubscriptionPrice  int    `json:"subscription_price,omitempty"`
}

func (api *Api) CreateChatSubscriptionInviteLink(p CreateChatSubscriptionInviteLinkP) (ChatInviteLink, error) {
	req := NewRequest[ChatInviteLink]("createChatSubscriptionInviteLink", p)
	return req.Do(api)
}

type EditChatSubscriptionInviteLinkP struct {
	ChatID     int    `json:"chat_id"`
	InviteLink string `json:"invite_link"`
	Name       string `json:"name,omitempty"`
}

func (api *Api) EditChatSubscriptionInviteLink(p EditChatSubscriptionInviteLinkP) (ChatInviteLink, error) {
	req := NewRequest[ChatInviteLink]("editChatSubscriptionInviteLink", p)
	return req.Do(api)
}

type RevokeChatInviteLinkP struct {
	ChatID     int    `json:"chat_id"`
	InviteLink string `json:"invite_link"`
}

func (api *Api) RevokeChatInviteLink(p RevokeChatInviteLinkP) (ChatInviteLink, error) {
	req := NewRequest[ChatInviteLink]("revokeChatInviteLink", p)
	return req.Do(api)
}

type ApproveChatJoinRequestP struct {
	ChatID int `json:"chat_id"`
	UserID int `json:"user_id"`
}

func (api *Api) ApproveChatJoinRequest(p ApproveChatJoinRequestP) (bool, error) {
	req := NewRequest[bool]("approveChatJoinRequest", p)
	return req.Do(api)
}

type DeclineChatJoinRequestP struct {
	ChatID int `json:"chat_id"`
	UserID int `json:"user_id"`
}

func (api *Api) DeclineChatJoinRequest(p DeclineChatJoinRequestP) (bool, error) {
	req := NewRequest[bool]("declineChatJoinRequest", p)
	return req.Do(api)
}

func (api *Api) SetChatPhoto() {
	uploader := NewUploader(api)
	defer uploader.Close()
}

type DeleteChatPhotoP struct {
	ChatID int `json:"chat_id"`
}

func (api *Api) DeleteChatPhoto(p DeleteChatPhotoP) (bool, error) {
	req := NewRequest[bool]("deleteChatPhoto", p)
	return req.Do(api)
}

type SetChatTitleP struct {
	ChatID int    `json:"chat_id"`
	Title  string `json:"title"`
}

func (api *Api) SetChatTitle(p SetChatTitleP) (bool, error) {
	req := NewRequest[bool]("setChatTitle", p)
	return req.Do(api)
}

type SetChatDescriptionP struct {
	ChatID      int    `json:"chat_id"`
	Description string `json:"description"`
}

func (api *Api) SetChatDescription(p SetChatDescriptionP) (bool, error) {
	req := NewRequest[bool]("setChatDescription", p)
	return req.Do(api)
}

type PinChatMessageP struct {
	BusinessConnectionID string `json:"business_connection_id,omitempty"`
	ChatID               int    `json:"chat_id"`
	MessageID            int    `json:"message_id"`
	DisableNotification  bool   `json:"disable_notification,omitempty"`
}

func (api *Api) PinChatMessage(p PinChatMessageP) (bool, error) {
	req := NewRequest[bool]("pinChatMessage", p)
	return req.Do(api)
}

type UnpinChatMessageP struct {
	BusinessConnectionID string `json:"business_connection_id,omitempty"`
	ChatID               int    `json:"chat_id"`
	MessageID            int    `json:"message_id"`
}

func (api *Api) UnpinChatMessage(p UnpinChatMessageP) (bool, error) {
	req := NewRequest[bool]("unpinChatMessage", p)
	return req.Do(api)
}

type UnpinAllChatMessagesP struct {
	ChatID int `json:"chat_id"`
}

func (api *Api) UnpinAllChatMessages(p UnpinAllChatMessagesP) (bool, error) {
	req := NewRequest[bool]("unpinAllChatMessages", p)
	return req.Do(api)
}

type LeaveChatP struct {
	ChatID int `json:"chat_id"`
}

func (api *Api) LeaveChat(p LeaveChatP) (bool, error) {
	req := NewRequest[bool]("leaveChatP", p)
	return req.Do(api)
}

type GetChatP struct {
	ChatID int `json:"chat_id"`
}

func (api *Api) GetChatP(p GetChatP) (ChatFullInfo, error) {
	req := NewRequest[ChatFullInfo]("getChatP", p)
	return req.Do(api)
}

type GetChatAdministratorsP struct {
	ChatID int `json:"chat_id"`
}

func (api *Api) GetChatAdministrators(p GetChatAdministratorsP) ([]BaseChatMember, error) {
	req := NewRequest[[]BaseChatMember]("getChatAdministrators", p)
	return req.Do(api)
}

type GetChatMembersCountP struct {
	ChatID int `json:"chat_id"`
}

func (api *Api) GetChatMemberCount(p GetChatMembersCountP) (int, error) {
	req := NewRequest[int]("getChatMemberCount", p)
	return req.Do(api)
}

type GetChatMemberP struct {
	ChatID int `json:"chat_id"`
	UserID int `json:"user_id"`
}

func (api *Api) GetChatMember(p GetChatMemberP) (BaseChatMember, error) {
	req := NewRequest[BaseChatMember]("getChatMember", p)
	return req.Do(api)
}

type SetChatStickerSetP struct {
	ChatID         int    `json:"chat_id"`
	StickerSetName string `json:"sticker_set_name"`
}

func (api *Api) SetChatStickerSet(p SetChatStickerSetP) (bool, error) {
	req := NewRequest[bool]("setChatStickerSet", p)
	return req.Do(api)
}

type DeleteChatStickerSetP struct {
	ChatID int `json:"chat_id"`
}

func (api *Api) DeleteChatStickerSet(p DeleteChatStickerSetP) (bool, error) {
	req := NewRequest[bool]("deleteChatStickerSet", p)
	return req.Do(api)
}

type GetUserChatBoostsP struct {
	ChatID int `json:"chat_id"`
	UserID int `json:"user_id"`
}

func (api *Api) GetUserChatBoosts(p GetUserChatBoostsP) (UserChatBoosts, error) {
	req := NewRequest[UserChatBoosts]("getUserChatBoosts", p)
	return req.Do(api)
}
