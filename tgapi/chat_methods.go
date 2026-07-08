package tgapi

import "context"

// BanChatMember holds parameters for the banChatMember method.
// Since: Bot API 5.3
// See https://core.telegram.org/bots/api#banchatmember
type BanChatMember struct {
	ChatID         int64 `json:"chat_id"`
	UserID         int64 `json:"user_id"`
	UntilDate      int   `json:"until_date,omitempty"`
	RevokeMessages bool  `json:"revoke_messages,omitempty"`
}

// BanChatMember bans a user in a chat.
// Since: Bot API 5.3
// Returns True on success.
// See https://core.telegram.org/bots/api#banchatmember
func (api *API) BanChatMember(params BanChatMember) (bool, error) {
	req := NewRequestWithChatID[bool]("banChatMember", params, params.ChatID)
	return req.Do(api)
}

// BanChatMemberWithContext is the context-aware variant of BanChatMember.
// Since: Bot API 5.3
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#banchatmember
func (api *API) BanChatMemberWithContext(ctx context.Context, params BanChatMember) (bool, error) {
	req := NewRequestWithChatID[bool]("banChatMember", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// UnbanChatMember holds parameters for the unbanChatMember method.
// Since: Bot API 2.0
// See https://core.telegram.org/bots/api#unbanchatmember
type UnbanChatMember struct {
	ChatID       int64 `json:"chat_id"`
	UserID       int64 `json:"user_id"`
	OnlyIfBanned bool  `json:"only_if_banned"`
}

// UnbanChatMember unbans a previously banned user in a chat.
// Since: Bot API 2.0
// Returns True on success.
// See https://core.telegram.org/bots/api#unbanchatmember
func (api *API) UnbanChatMember(params UnbanChatMember) (bool, error) {
	req := NewRequestWithChatID[bool]("unbanChatMember", params, params.ChatID)
	return req.Do(api)
}

// UnbanChatMemberWithContext is the context-aware variant of UnbanChatMember.
// Since: Bot API 2.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#unbanchatmember
func (api *API) UnbanChatMemberWithContext(ctx context.Context, params UnbanChatMember) (bool, error) {
	req := NewRequestWithChatID[bool]("unbanChatMember", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// RestrictChatMember holds parameters for the restrictChatMember method.
// Since: Bot API 3.1
// See https://core.telegram.org/bots/api#restrictchatmember
type RestrictChatMember struct {
	ChatID                        int64           `json:"chat_id"`
	UserID                        int64           `json:"user_id"`
	Permissions                   ChatPermissions `json:"permissions"`
	UseIndependentChatPermissions bool            `json:"use_independent_chat_permissions,omitempty"`
	UntilDate                     int             `json:"until_date,omitempty"`
}

// RestrictChatMember restricts a user in a chat.
// Since: Bot API 3.1
// Returns True on success.
// See https://core.telegram.org/bots/api#restrictchatmember
func (api *API) RestrictChatMember(params RestrictChatMember) (bool, error) {
	req := NewRequestWithChatID[bool]("restrictChatMember", params, params.ChatID)
	return req.Do(api)
}

// RestrictChatMemberWithContext is the context-aware variant of RestrictChatMember.
// Since: Bot API 3.1
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#restrictchatmember
func (api *API) RestrictChatMemberWithContext(ctx context.Context, params RestrictChatMember) (bool, error) {
	req := NewRequestWithChatID[bool]("restrictChatMember", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// PromoteChatMember holds parameters for the promoteChatMember method.
// Since: Bot API 3.1
// See https://core.telegram.org/bots/api#promotechatmember
type PromoteChatMember struct {
	ChatID      int64 `json:"chat_id"`
	UserID      int64 `json:"user_id"`
	IsAnonymous bool  `json:"is_anonymous,omitempty"` // Since: Bot API 5.1

	CanManageChat           bool `json:"can_manage_chat,omitempty"` // Since: Bot API 5.3
	CanDeleteMessages       bool `json:"can_delete_messages,omitempty"`
	CanManageVideoChats     bool `json:"can_manage_video_chats,omitempty"` // Since: Bot API 6.0
	CanRestrictMembers      bool `json:"can_restrict_members,omitempty"`
	CanPromoteMembers       bool `json:"can_promote_members,omitempty"`
	CanChangeInfo           bool `json:"can_change_info,omitempty"`
	CanInviteUsers          bool `json:"can_invite_users,omitempty"`
	CanPostStories          bool `json:"can_post_stories,omitempty"`   // Since: Bot API 6.9
	CanEditStories          bool `json:"can_edit_stories,omitempty"`   // Since: Bot API 6.9
	CanDeleteStories        bool `json:"can_delete_stories,omitempty"` // Since: Bot API 6.9
	CanPostMessages         bool `json:"can_post_messages,omitempty"`
	CanEditMessages         bool `json:"can_edit_messages,omitempty"`
	CanPinMessages          bool `json:"can_pin_messages,omitempty"`
	CanManageTopics         bool `json:"can_manage_topics,omitempty"`          // Since: Bot API 6.3
	CanManageDirectMessages bool `json:"can_manage_direct_messages,omitempty"` // Since: Bot API 9.1
	CanManageTags           bool `json:"can_manage_tags,omitempty"`            // Since: Bot API 9.5
}

// PromoteChatMember promotes or demotes a user in a chat.
// Since: Bot API 3.1
// Returns True on success.
// See https://core.telegram.org/bots/api#promotechatmember
func (api *API) PromoteChatMember(params PromoteChatMember) (bool, error) {
	req := NewRequestWithChatID[bool]("promoteChatMember", params, params.ChatID)
	return req.Do(api)
}

// PromoteChatMemberWithContext is the context-aware variant of PromoteChatMember.
// Since: Bot API 3.1
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#promotechatmember
func (api *API) PromoteChatMemberWithContext(ctx context.Context, params PromoteChatMember) (bool, error) {
	req := NewRequestWithChatID[bool]("promoteChatMember", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SetChatAdministratorCustomTitle holds parameters for the setChatAdministratorCustomTitle method.
// Since: Bot API 5.0
// See https://core.telegram.org/bots/api#setchatadministratorcustomtitle
type SetChatAdministratorCustomTitle struct {
	ChatID      int64  `json:"chat_id"`
	UserID      int64  `json:"user_id"`
	CustomTitle string `json:"custom_title"`
}

// SetChatAdministratorCustomTitle sets a custom title for an administrator.
// Since: Bot API 5.0
// Returns True on success.
// See https://core.telegram.org/bots/api#setchatadministratorcustomtitle
func (api *API) SetChatAdministratorCustomTitle(params SetChatAdministratorCustomTitle) (bool, error) {
	req := NewRequestWithChatID[bool]("setChatAdministratorCustomTitle", params, params.ChatID)
	return req.Do(api)
}

// SetChatAdministratorCustomTitleWithContext is the context-aware variant of SetChatAdministratorCustomTitle.
// Since: Bot API 5.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#setchatadministratorcustomtitle
func (api *API) SetChatAdministratorCustomTitleWithContext(ctx context.Context, params SetChatAdministratorCustomTitle) (bool, error) {
	req := NewRequestWithChatID[bool]("setChatAdministratorCustomTitle", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SetChatMemberTag holds parameters for the setChatMemberTag method.
// Since: Bot API 9.5
// See https://core.telegram.org/bots/api#setchatmembertag
type SetChatMemberTag struct {
	ChatID int64  `json:"chat_id"`
	UserID int64  `json:"user_id"`
	Tag    string `json:"tag,omitempty"`
}

// SetChatMemberTag sets a tag for a chat member.
// Since: Bot API 9.5
// Returns True on success.
// See https://core.telegram.org/bots/api#setchatmembertag
func (api *API) SetChatMemberTag(params SetChatMemberTag) (bool, error) {
	req := NewRequestWithChatID[bool]("setChatMemberTag", params, params.ChatID)
	return req.Do(api)
}

// SetChatMemberTagWithContext is the context-aware variant of SetChatMemberTag.
// Since: Bot API 9.5
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#setchatmembertag
func (api *API) SetChatMemberTagWithContext(ctx context.Context, params SetChatMemberTag) (bool, error) {
	req := NewRequestWithChatID[bool]("setChatMemberTag", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// BanChatSenderChat holds parameters for the banChatSenderChat method.
// Since: Bot API 5.6
// See https://core.telegram.org/bots/api#banchatsenderchat
type BanChatSenderChat struct {
	ChatID       int64 `json:"chat_id"`
	SenderChatID int64 `json:"sender_chat_id"`
}

// BanChatSenderChat bans a channel chat in a supergroup or channel.
// Since: Bot API 5.6
// Returns True on success.
// See https://core.telegram.org/bots/api#banchatsenderchat
func (api *API) BanChatSenderChat(params BanChatSenderChat) (bool, error) {
	req := NewRequestWithChatID[bool]("banChatSenderChat", params, params.ChatID)
	return req.Do(api)
}

// BanChatSenderChatWithContext is the context-aware variant of BanChatSenderChat.
// Since: Bot API 5.6
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#banchatsenderchat
func (api *API) BanChatSenderChatWithContext(ctx context.Context, params BanChatSenderChat) (bool, error) {
	req := NewRequestWithChatID[bool]("banChatSenderChat", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// UnbanChatSenderChat holds parameters for the unbanChatSenderChat method.
// Since: Bot API 5.6
// See https://core.telegram.org/bots/api#unbanchatsenderchat
type UnbanChatSenderChat struct {
	ChatID       int64 `json:"chat_id"`
	SenderChatID int64 `json:"sender_chat_id"`
}

// UnbanChatSenderChat unbans a previously banned channel chat.
// Since: Bot API 5.6
// Returns True on success.
// See https://core.telegram.org/bots/api#unbanchatsenderchat
func (api *API) UnbanChatSenderChat(params UnbanChatSenderChat) (bool, error) {
	req := NewRequestWithChatID[bool]("unbanChatSenderChat", params, params.ChatID)
	return req.Do(api)
}

// UnbanChatSenderChatWithContext is the context-aware variant of UnbanChatSenderChat.
// Since: Bot API 5.6
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#unbanchatsenderchat
func (api *API) UnbanChatSenderChatWithContext(ctx context.Context, params UnbanChatSenderChat) (bool, error) {
	req := NewRequestWithChatID[bool]("unbanChatSenderChat", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SetChatPermissions holds parameters for the setChatPermissions method.
// Since: Bot API 4.4
// See https://core.telegram.org/bots/api#setchatpermissions
type SetChatPermissions struct {
	ChatID                        int64           `json:"chat_id"`
	Permissions                   ChatPermissions `json:"permissions"`
	UseIndependentChatPermissions bool            `json:"use_independent_chat_permissions,omitempty"`
}

// SetChatPermissions sets default chat permissions for all members.
// Since: Bot API 4.4
// Returns True on success.
// See https://core.telegram.org/bots/api#setchatpermissions
func (api *API) SetChatPermissions(params SetChatPermissions) (bool, error) {
	req := NewRequestWithChatID[bool]("setChatPermissions", params, params.ChatID)
	return req.Do(api)
}

// SetChatPermissionsWithContext is the context-aware variant of SetChatPermissions.
// Since: Bot API 4.4
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#setchatpermissions
func (api *API) SetChatPermissionsWithContext(ctx context.Context, params SetChatPermissions) (bool, error) {
	req := NewRequestWithChatID[bool]("setChatPermissions", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// ExportChatInviteLink holds parameters for the exportChatInviteLink method.
// Since: Bot API 3.1
// See https://core.telegram.org/bots/api#exportchatinvitelink
type ExportChatInviteLink struct {
	ChatID int64 `json:"chat_id"`
}

// ExportChatInviteLink generates a new primary invite link for a chat.
// Since: Bot API 3.1
// Returns the new invite link as string.
// See https://core.telegram.org/bots/api#exportchatinvitelink
func (api *API) ExportChatInviteLink(params ExportChatInviteLink) (string, error) {
	req := NewRequestWithChatID[string]("exportChatInviteLink", params, params.ChatID)
	return req.Do(api)
}

// ExportChatInviteLinkWithContext is the context-aware variant of ExportChatInviteLink.
// Since: Bot API 3.1
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#exportchatinvitelink
func (api *API) ExportChatInviteLinkWithContext(ctx context.Context, params ExportChatInviteLink) (string, error) {
	req := NewRequestWithChatID[string]("exportChatInviteLink", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// CreateChatInviteLink holds parameters for the createChatInviteLink method.
// Since: Bot API 5.1
// See https://core.telegram.org/bots/api#createchatinvitelink
type CreateChatInviteLink struct {
	ChatID             int64   `json:"chat_id"`
	Name               *string `json:"name,omitempty"`
	ExpireDate         int     `json:"expire_date,omitempty"`
	MemberLimit        int     `json:"member_limit,omitempty"`
	CreatesJoinRequest bool    `json:"creates_join_request,omitempty"`
}

// CreateChatInviteLink creates an additional invite link for a chat.
// Since: Bot API 5.1
// Returns the created invite link.
// See https://core.telegram.org/bots/api#createchatinvitelink
func (api *API) CreateChatInviteLink(params CreateChatInviteLink) (ChatInviteLink, error) {
	req := NewRequestWithChatID[ChatInviteLink]("createChatInviteLink", params, params.ChatID)
	return req.Do(api)
}

// CreateChatInviteLinkWithContext is the context-aware variant of CreateChatInviteLink.
// Since: Bot API 5.1
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#createchatinvitelink
func (api *API) CreateChatInviteLinkWithContext(ctx context.Context, params CreateChatInviteLink) (ChatInviteLink, error) {
	req := NewRequestWithChatID[ChatInviteLink]("createChatInviteLink", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// EditChatInviteLink holds parameters for the editChatInviteLink method.
// Since: Bot API 5.1
// See https://core.telegram.org/bots/api#editchatinvitelink
type EditChatInviteLink struct {
	ChatID     int64  `json:"chat_id"`
	InviteLink string `json:"invite_link"`

	Name               string `json:"name,omitempty"`
	ExpireDate         int    `json:"expire_date,omitempty"`
	MemberLimit        int    `json:"member_limit,omitempty"`
	CreatesJoinRequest bool   `json:"creates_join_request,omitempty"`
}

// EditChatInviteLink edits a non‑primary invite link.
// Since: Bot API 5.1
// Returns the edited invite link.
// See https://core.telegram.org/bots/api#editchatinvitelink
func (api *API) EditChatInviteLink(params EditChatInviteLink) (ChatInviteLink, error) {
	req := NewRequestWithChatID[ChatInviteLink]("editChatInviteLink", params, params.ChatID)
	return req.Do(api)
}

// EditChatInviteLinkWithContext is the context-aware variant of EditChatInviteLink.
// Since: Bot API 5.1
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#editchatinvitelink
func (api *API) EditChatInviteLinkWithContext(ctx context.Context, params EditChatInviteLink) (ChatInviteLink, error) {
	req := NewRequestWithChatID[ChatInviteLink]("editChatInviteLink", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// CreateChatSubscriptionInviteLink holds parameters for the createChatSubscriptionInviteLink method.
// Since: Bot API 8.0
// See https://core.telegram.org/bots/api#createchatsubscriptioninvitelink
type CreateChatSubscriptionInviteLink struct {
	ChatID             int64  `json:"chat_id"`
	Name               string `json:"name,omitempty"`
	SubscriptionPeriod int    `json:"subscription_period,omitempty"`
	SubscriptionPrice  int    `json:"subscription_price,omitempty"`
}

// CreateChatSubscriptionInviteLink creates a subscription invite link for a channel chat.
// Since: Bot API 8.0
// Returns the created invite link.
// See https://core.telegram.org/bots/api#createchatsubscriptioninvitelink
func (api *API) CreateChatSubscriptionInviteLink(params CreateChatSubscriptionInviteLink) (ChatInviteLink, error) {
	req := NewRequestWithChatID[ChatInviteLink]("createChatSubscriptionInviteLink", params, params.ChatID)
	return req.Do(api)
}

// CreateChatSubscriptionInviteLinkWithContext is the context-aware variant of CreateChatSubscriptionInviteLink.
// Since: Bot API 8.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#createchatsubscriptioninvitelink
func (api *API) CreateChatSubscriptionInviteLinkWithContext(ctx context.Context, params CreateChatSubscriptionInviteLink) (ChatInviteLink, error) {
	req := NewRequestWithChatID[ChatInviteLink]("createChatSubscriptionInviteLink", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// EditChatSubscriptionInviteLink holds parameters for the editChatSubscriptionInviteLink method.
// Since: Bot API 8.0
// See https://core.telegram.org/bots/api#editchatsubscriptioninvitelink
type EditChatSubscriptionInviteLink struct {
	ChatID     int64  `json:"chat_id"`
	InviteLink string `json:"invite_link"`
	Name       string `json:"name,omitempty"`
}

// EditChatSubscriptionInviteLink edits a subscription invite link.
// Since: Bot API 8.0
// Returns the edited invite link.
// See https://core.telegram.org/bots/api#editchatsubscriptioninvitelink
func (api *API) EditChatSubscriptionInviteLink(params EditChatSubscriptionInviteLink) (ChatInviteLink, error) {
	req := NewRequestWithChatID[ChatInviteLink]("editChatSubscriptionInviteLink", params, params.ChatID)
	return req.Do(api)
}

// EditChatSubscriptionInviteLinkWithContext is the context-aware variant of EditChatSubscriptionInviteLink.
// Since: Bot API 8.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#editchatsubscriptioninvitelink
func (api *API) EditChatSubscriptionInviteLinkWithContext(ctx context.Context, params EditChatSubscriptionInviteLink) (ChatInviteLink, error) {
	req := NewRequestWithChatID[ChatInviteLink]("editChatSubscriptionInviteLink", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// RevokeChatInviteLink holds parameters for the revokeChatInviteLink method.
// Since: Bot API 5.1
// See https://core.telegram.org/bots/api#revokechatinvitelink
type RevokeChatInviteLink struct {
	ChatID     int64  `json:"chat_id"`
	InviteLink string `json:"invite_link"`
}

// RevokeChatInviteLink revokes an invite link.
// Since: Bot API 5.1
// Returns the revoked invite link object.
// See https://core.telegram.org/bots/api#revokechatinvitelink
func (api *API) RevokeChatInviteLink(params RevokeChatInviteLink) (ChatInviteLink, error) {
	req := NewRequestWithChatID[ChatInviteLink]("revokeChatInviteLink", params, params.ChatID)
	return req.Do(api)
}

// RevokeChatInviteLinkWithContext is the context-aware variant of RevokeChatInviteLink.
// Since: Bot API 5.1
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#revokechatinvitelink
func (api *API) RevokeChatInviteLinkWithContext(ctx context.Context, params RevokeChatInviteLink) (ChatInviteLink, error) {
	req := NewRequestWithChatID[ChatInviteLink]("revokeChatInviteLink", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// ApproveChatJoinRequest holds parameters for the approveChatJoinRequest method.
// Since: Bot API 5.4
// See https://core.telegram.org/bots/api#approvechatjoinrequest
type ApproveChatJoinRequest struct {
	ChatID int64 `json:"chat_id"`
	UserID int64 `json:"user_id"`
}

// ApproveChatJoinRequest approves a chat join request.
// Since: Bot API 5.4
// Returns True on success.
// See https://core.telegram.org/bots/api#approvechatjoinrequest
func (api *API) ApproveChatJoinRequest(params ApproveChatJoinRequest) (bool, error) {
	req := NewRequestWithChatID[bool]("approveChatJoinRequest", params, params.ChatID)
	return req.Do(api)
}

// ApproveChatJoinRequestWithContext is the context-aware variant of ApproveChatJoinRequest.
// Since: Bot API 5.4
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#approvechatjoinrequest
func (api *API) ApproveChatJoinRequestWithContext(ctx context.Context, params ApproveChatJoinRequest) (bool, error) {
	req := NewRequestWithChatID[bool]("approveChatJoinRequest", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// DeclineChatJoinRequest holds parameters for the declineChatJoinRequest method.
// Since: Bot API 5.4
// See https://core.telegram.org/bots/api#declinechatjoinrequest
type DeclineChatJoinRequest struct {
	ChatID int64 `json:"chat_id"`
	UserID int64 `json:"user_id"`
}

// DeclineChatJoinRequest declines a chat join request.
// Since: Bot API 5.4
// Returns True on success.
// See https://core.telegram.org/bots/api#declinechatjoinrequest
func (api *API) DeclineChatJoinRequest(params DeclineChatJoinRequest) (bool, error) {
	req := NewRequestWithChatID[bool]("declineChatJoinRequest", params, params.ChatID)
	return req.Do(api)
}

// DeclineChatJoinRequestWithContext is the context-aware variant of DeclineChatJoinRequest.
// Since: Bot API 5.4
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#declinechatjoinrequest
func (api *API) DeclineChatJoinRequestWithContext(ctx context.Context, params DeclineChatJoinRequest) (bool, error) {
	req := NewRequestWithChatID[bool]("declineChatJoinRequest", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// ChatJoinRequestQueryResult is the verdict passed to answerChatJoinRequestQuery.
// Since: Bot API 10.1
type ChatJoinRequestQueryResult string

const (
	// JoinRequestApprove allows the user to join the chat.
	JoinRequestApprove ChatJoinRequestQueryResult = "approve"
	// JoinRequestDecline disallows the user to join the chat.
	JoinRequestDecline ChatJoinRequestQueryResult = "decline"
	// JoinRequestQueue leaves the decision to other administrators.
	JoinRequestQueue ChatJoinRequestQueryResult = "queue"
)

// AnswerChatJoinRequestQuery holds parameters for the answerChatJoinRequestQuery method.
// Since: Bot API 10.1
// See https://core.telegram.org/bots/api#answerchatjoinrequestquery
type AnswerChatJoinRequestQuery struct {
	ChatJoinRequestQueryID string                     `json:"chat_join_request_query_id"`
	Result                 ChatJoinRequestQueryResult `json:"result"`
}

// AnswerChatJoinRequestQuery processes a received chat join request query.
// Since: Bot API 10.1
// Returns True on success.
// See https://core.telegram.org/bots/api#answerchatjoinrequestquery
func (api *API) AnswerChatJoinRequestQuery(params AnswerChatJoinRequestQuery) (bool, error) {
	req := NewRequest[bool]("answerChatJoinRequestQuery", params)
	return req.Do(api)
}

// AnswerChatJoinRequestQueryWithContext is the context-aware variant of AnswerChatJoinRequestQuery.
// Since: Bot API 10.1
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#answerchatjoinrequestquery
func (api *API) AnswerChatJoinRequestQueryWithContext(ctx context.Context, params AnswerChatJoinRequestQuery) (bool, error) {
	req := NewRequest[bool]("answerChatJoinRequestQuery", params)
	return req.DoWithContext(ctx, api)
}

// SendChatJoinRequestWebApp holds parameters for the sendChatJoinRequestWebApp method.
// Since: Bot API 10.1
// See https://core.telegram.org/bots/api#sendchatjoinrequestwebapp
type SendChatJoinRequestWebApp struct {
	ChatJoinRequestQueryID string `json:"chat_join_request_query_id"`
	WebAppURL              string `json:"web_app_url"`
}

// SendChatJoinRequestWebApp shows a Mini App to the user before deciding a
// join request query; resolve the query with AnswerChatJoinRequestQuery based
// on the Mini App interaction.
// Since: Bot API 10.1
// Returns True on success.
// See https://core.telegram.org/bots/api#sendchatjoinrequestwebapp
func (api *API) SendChatJoinRequestWebApp(params SendChatJoinRequestWebApp) (bool, error) {
	req := NewRequest[bool]("sendChatJoinRequestWebApp", params)
	return req.Do(api)
}

// SendChatJoinRequestWebAppWithContext is the context-aware variant of SendChatJoinRequestWebApp.
// Since: Bot API 10.1
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendchatjoinrequestwebapp
func (api *API) SendChatJoinRequestWebAppWithContext(ctx context.Context, params SendChatJoinRequestWebApp) (bool, error) {
	req := NewRequest[bool]("sendChatJoinRequestWebApp", params)
	return req.DoWithContext(ctx, api)
}

// SetChatPhoto holds parameters for the setChatPhoto method.
// Since: Bot API 3.1
// See https://core.telegram.org/bots/api#setchatphoto
type SetChatPhoto struct {
	ChatID int64 `json:"chat_id"`
}

// SetChatPhoto changes the chat photo.
// Since: Bot API 3.1
// photo is the file to upload as the new photo.
// Returns True on success.
// See https://core.telegram.org/bots/api#setchatphoto
func (api *API) SetChatPhoto(params SetChatPhoto, photo UploaderFile) (bool, error) {
	uploader := NewUploader(api)
	defer func() {
		_ = uploader.Close()
	}()
	req := NewUploaderRequestWithChatID[bool]("setChatPhoto", params, params.ChatID, photo.SetType(UploaderPhotoType))
	return req.Do(uploader)
}

// DeleteChatPhoto holds parameters for the deleteChatPhoto method.
// Since: Bot API 3.1
// See https://core.telegram.org/bots/api#deletechatphoto
type DeleteChatPhoto struct {
	ChatID int64 `json:"chat_id"`
}

// DeleteChatPhoto deletes a chat photo.
// Since: Bot API 3.1
// Returns True on success.
// See https://core.telegram.org/bots/api#deletechatphoto
func (api *API) DeleteChatPhoto(params DeleteChatPhoto) (bool, error) {
	req := NewRequestWithChatID[bool]("deleteChatPhoto", params, params.ChatID)
	return req.Do(api)
}

// DeleteChatPhotoWithContext is the context-aware variant of DeleteChatPhoto.
// Since: Bot API 3.1
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#deletechatphoto
func (api *API) DeleteChatPhotoWithContext(ctx context.Context, params DeleteChatPhoto) (bool, error) {
	req := NewRequestWithChatID[bool]("deleteChatPhoto", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SetChatTitle holds parameters for the setChatTitle method.
// Since: Bot API 3.1
// See https://core.telegram.org/bots/api#setchattitle
type SetChatTitle struct {
	ChatID int64  `json:"chat_id"`
	Title  string `json:"title"`
}

// SetChatTitle changes the chat title.
// Since: Bot API 3.1
// Returns True on success.
// See https://core.telegram.org/bots/api#setchattitle
func (api *API) SetChatTitle(params SetChatTitle) (bool, error) {
	req := NewRequestWithChatID[bool]("setChatTitle", params, params.ChatID)
	return req.Do(api)
}

// SetChatTitleWithContext is the context-aware variant of SetChatTitle.
// Since: Bot API 3.1
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#setchattitle
func (api *API) SetChatTitleWithContext(ctx context.Context, params SetChatTitle) (bool, error) {
	req := NewRequestWithChatID[bool]("setChatTitle", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SetChatDescription holds parameters for the setChatDescription method.
// Since: Bot API 3.1
// See https://core.telegram.org/bots/api#setchatdescription
type SetChatDescription struct {
	ChatID      int64  `json:"chat_id"`
	Description string `json:"description"`
}

// SetChatDescription changes the chat description.
// Since: Bot API 3.1
// Returns True on success.
// See https://core.telegram.org/bots/api#setchatdescription
func (api *API) SetChatDescription(params SetChatDescription) (bool, error) {
	req := NewRequestWithChatID[bool]("setChatDescription", params, params.ChatID)
	return req.Do(api)
}

// SetChatDescriptionWithContext is the context-aware variant of SetChatDescription.
// Since: Bot API 3.1
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#setchatdescription
func (api *API) SetChatDescriptionWithContext(ctx context.Context, params SetChatDescription) (bool, error) {
	req := NewRequestWithChatID[bool]("setChatDescription", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// PinChatMessage holds parameters for the pinChatMessage method.
// Since: Bot API 3.1
// See https://core.telegram.org/bots/api#pinchatmessage
type PinChatMessage struct {
	BusinessConnectionID *string `json:"business_connection_id,omitempty"`
	ChatID               int64   `json:"chat_id"`
	MessageID            int     `json:"message_id"`
	DisableNotification  bool    `json:"disable_notification,omitempty"`
}

// PinChatMessage pins a message in a chat.
// Since: Bot API 3.1
// Returns True on success.
// See https://core.telegram.org/bots/api#pinchatmessage
func (api *API) PinChatMessage(params PinChatMessage) (bool, error) {
	req := NewRequestWithChatID[bool]("pinChatMessage", params, params.ChatID)
	return req.Do(api)
}

// PinChatMessageWithContext is the context-aware variant of PinChatMessage.
// Since: Bot API 3.1
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#pinchatmessage
func (api *API) PinChatMessageWithContext(ctx context.Context, params PinChatMessage) (bool, error) {
	req := NewRequestWithChatID[bool]("pinChatMessage", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// UnpinChatMessage holds parameters for the unpinChatMessage method.
// Since: Bot API 3.1
// See https://core.telegram.org/bots/api#unpinchatmessage
type UnpinChatMessage struct {
	BusinessConnectionID *string `json:"business_connection_id,omitempty"`
	ChatID               int64   `json:"chat_id"`
	MessageID            int     `json:"message_id"`
}

// UnpinChatMessage unpins a message in a chat.
// Since: Bot API 3.1
// Returns True on success.
// See https://core.telegram.org/bots/api#unpinchatmessage
func (api *API) UnpinChatMessage(params UnpinChatMessage) (bool, error) {
	req := NewRequestWithChatID[bool]("unpinChatMessage", params, params.ChatID)
	return req.Do(api)
}

// UnpinChatMessageWithContext is the context-aware variant of UnpinChatMessage.
// Since: Bot API 3.1
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#unpinchatmessage
func (api *API) UnpinChatMessageWithContext(ctx context.Context, params UnpinChatMessage) (bool, error) {
	req := NewRequestWithChatID[bool]("unpinChatMessage", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// UnpinAllChatMessages holds parameters for the unpinAllChatMessages method.
// Since: Bot API 5.0
// See https://core.telegram.org/bots/api#unpinallchatmessages
type UnpinAllChatMessages struct {
	ChatID int64 `json:"chat_id"`
}

// UnpinAllChatMessages unpins all pinned messages in a chat.
// Since: Bot API 5.0
// Returns True on success.
// See https://core.telegram.org/bots/api#unpinallchatmessages
func (api *API) UnpinAllChatMessages(params UnpinAllChatMessages) (bool, error) {
	req := NewRequestWithChatID[bool]("unpinAllChatMessages", params, params.ChatID)
	return req.Do(api)
}

// UnpinAllChatMessagesWithContext is the context-aware variant of UnpinAllChatMessages.
// Since: Bot API 5.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#unpinallchatmessages
func (api *API) UnpinAllChatMessagesWithContext(ctx context.Context, params UnpinAllChatMessages) (bool, error) {
	req := NewRequestWithChatID[bool]("unpinAllChatMessages", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// LeaveChat holds parameters for the leaveChat method.
// Since: Bot API 2.1
// See https://core.telegram.org/bots/api#leavechat
type LeaveChat struct {
	ChatID int64 `json:"chat_id"`
}

// LeaveChat makes the bot leave a chat.
// Since: Bot API 2.1
// Returns True on success.
// See https://core.telegram.org/bots/api#leavechat
func (api *API) LeaveChat(params LeaveChat) (bool, error) {
	req := NewRequestWithChatID[bool]("leaveChat", params, params.ChatID)
	return req.Do(api)
}

// LeaveChatWithContext is the context-aware variant of LeaveChat.
// Since: Bot API 2.1
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#leavechat
func (api *API) LeaveChatWithContext(ctx context.Context, params LeaveChat) (bool, error) {
	req := NewRequestWithChatID[bool]("leaveChat", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// GetChat holds parameters for the getChat method.
// Since: Bot API 2.1
// See https://core.telegram.org/bots/api#getchat
type GetChat struct {
	ChatID int64 `json:"chat_id"`
}

// GetChat gets up‑to‑date information about a chat.
// Since: Bot API 2.1
// See https://core.telegram.org/bots/api#getchat
func (api *API) GetChat(params GetChat) (ChatFullInfo, error) {
	req := NewRequestWithChatID[ChatFullInfo]("getChat", params, params.ChatID)
	return req.Do(api)
}

// GetChatWithContext is the context-aware variant of GetChat.
// Since: Bot API 2.1
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#getchat
func (api *API) GetChatWithContext(ctx context.Context, params GetChat) (ChatFullInfo, error) {
	req := NewRequestWithChatID[ChatFullInfo]("getChat", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// GetChatAdministrators holds parameters for the getChatAdministrators method.
// Since: Bot API 2.1
// See https://core.telegram.org/bots/api#getchatadministrators
type GetChatAdministrators struct {
	ChatID     int64 `json:"chat_id"`
	ReturnBots bool  `json:"return_bots,omitempty"` // Since: Bot API 10.0
}

// GetChatAdministrators returns a list of administrators in a chat.
// Since: Bot API 2.1
// See https://core.telegram.org/bots/api#getchatadministrators
func (api *API) GetChatAdministrators(params GetChatAdministrators) ([]ChatMember, error) {
	req := NewRequestWithChatID[[]ChatMember]("getChatAdministrators", params, params.ChatID)
	return req.Do(api)
}

// GetChatAdministratorsWithContext is the context-aware variant of GetChatAdministrators.
// Since: Bot API 2.1
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#getchatadministrators
func (api *API) GetChatAdministratorsWithContext(ctx context.Context, params GetChatAdministrators) ([]ChatMember, error) {
	req := NewRequestWithChatID[[]ChatMember]("getChatAdministrators", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// GetChatMemberCount holds parameters for the getChatMemberCount method.
// Since: Bot API 2.1
// See https://core.telegram.org/bots/api#getchatmembercount
type GetChatMemberCount struct {
	ChatID int64 `json:"chat_id"`
}

// GetChatMemberCount returns the number of members in a chat.
// Since: Bot API 2.1
// See https://core.telegram.org/bots/api#getchatmembercount
func (api *API) GetChatMemberCount(params GetChatMemberCount) (int, error) {
	req := NewRequestWithChatID[int]("getChatMemberCount", params, params.ChatID)
	return req.Do(api)
}

// GetChatMemberCountWithContext is the context-aware variant of GetChatMemberCount.
// Since: Bot API 2.1
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#getchatmembercount
func (api *API) GetChatMemberCountWithContext(ctx context.Context, params GetChatMemberCount) (int, error) {
	req := NewRequestWithChatID[int]("getChatMemberCount", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// GetChatMember holds parameters for the getChatMember method.
// Since: Bot API 2.1
// See https://core.telegram.org/bots/api#getchatmember
type GetChatMember struct {
	ChatID int64 `json:"chat_id"`
	UserID int64 `json:"user_id"`
}

// GetChatMember returns information about a member of a chat.
// Since: Bot API 2.1
// See https://core.telegram.org/bots/api#getchatmember
func (api *API) GetChatMember(params GetChatMember) (ChatMember, error) {
	req := NewRequestWithChatID[ChatMember]("getChatMember", params, params.ChatID)
	return req.Do(api)
}

// GetChatMemberWithContext is the context-aware variant of GetChatMember.
// Since: Bot API 2.1
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#getchatmember
func (api *API) GetChatMemberWithContext(ctx context.Context, params GetChatMember) (ChatMember, error) {
	req := NewRequestWithChatID[ChatMember]("getChatMember", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SetChatStickerSet holds parameters for the setChatStickerSet method.
// Since: Bot API 3.2
// See https://core.telegram.org/bots/api#setchatstickerset
type SetChatStickerSet struct {
	ChatID         int64  `json:"chat_id"`
	StickerSetName string `json:"sticker_set_name"`
}

// SetChatStickerSet associates a sticker set with a supergroup.
// Since: Bot API 3.2
// Returns True on success.
// See https://core.telegram.org/bots/api#setchatstickerset
func (api *API) SetChatStickerSet(params SetChatStickerSet) (bool, error) {
	req := NewRequestWithChatID[bool]("setChatStickerSet", params, params.ChatID)
	return req.Do(api)
}

// SetChatStickerSetWithContext is the context-aware variant of SetChatStickerSet.
// Since: Bot API 3.2
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#setchatstickerset
func (api *API) SetChatStickerSetWithContext(ctx context.Context, params SetChatStickerSet) (bool, error) {
	req := NewRequestWithChatID[bool]("setChatStickerSet", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// DeleteChatStickerSet holds parameters for the deleteChatStickerSet method.
// Since: Bot API 3.2
// See https://core.telegram.org/bots/api#deletechatstickerset
type DeleteChatStickerSet struct {
	ChatID int64 `json:"chat_id"`
}

// DeleteChatStickerSet deletes a sticker set from a supergroup.
// Since: Bot API 3.2
// Returns True on success.
// See https://core.telegram.org/bots/api#deletechatstickerset
func (api *API) DeleteChatStickerSet(params DeleteChatStickerSet) (bool, error) {
	req := NewRequestWithChatID[bool]("deleteChatStickerSet", params, params.ChatID)
	return req.Do(api)
}

// DeleteChatStickerSetWithContext is the context-aware variant of DeleteChatStickerSet.
// Since: Bot API 3.2
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#deletechatstickerset
func (api *API) DeleteChatStickerSetWithContext(ctx context.Context, params DeleteChatStickerSet) (bool, error) {
	req := NewRequestWithChatID[bool]("deleteChatStickerSet", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// GetUserChatBoosts holds parameters for the getUserChatBoosts method.
// Since: Bot API 7.0
// See https://core.telegram.org/bots/api#getuserchatboosts
type GetUserChatBoosts struct {
	ChatID int64 `json:"chat_id"`
	UserID int64 `json:"user_id"`
}

// GetUserChatBoosts returns the list of boosts a user has given to a chat.
// Since: Bot API 7.0
// See https://core.telegram.org/bots/api#getuserchatboosts
func (api *API) GetUserChatBoosts(params GetUserChatBoosts) (UserChatBoosts, error) {
	req := NewRequestWithChatID[UserChatBoosts]("getUserChatBoosts", params, params.ChatID)
	return req.Do(api)
}

// GetUserChatBoostsWithContext is the context-aware variant of GetUserChatBoosts.
// Since: Bot API 7.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#getuserchatboosts
func (api *API) GetUserChatBoostsWithContext(ctx context.Context, params GetUserChatBoosts) (UserChatBoosts, error) {
	req := NewRequestWithChatID[UserChatBoosts]("getUserChatBoosts", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// GetChatGifts holds parameters for the getChatGifts method.
// Since: Bot API 9.3
// See https://core.telegram.org/bots/api#getchatgifts
type GetChatGifts struct {
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
// Since: Bot API 9.3
// See https://core.telegram.org/bots/api#getchatgifts
func (api *API) GetChatGifts(params GetChatGifts) (OwnedGifts, error) {
	req := NewRequestWithChatID[OwnedGifts]("getChatGifts", params, params.ChatID)
	return req.Do(api)
}

// GetChatGiftsWithContext is the context-aware variant of GetChatGifts.
// Since: Bot API 9.3
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#getchatgifts
func (api *API) GetChatGiftsWithContext(ctx context.Context, params GetChatGifts) (OwnedGifts, error) {
	req := NewRequestWithChatID[OwnedGifts]("getChatGifts", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}
