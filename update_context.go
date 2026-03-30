package laniakea

import "git.scuroneko.dev/scuroneko/laniakea/tgapi"

func (bot *Bot[T]) handleUpdate(u *tgapi.Update, ctx *MsgContext) {
	for _, plugin := range bot.plugins {
		handler, ok := plugin.handlers[u.Type]
		if !ok {
			continue
		}

		pluginCtx := cloneMsgContext(ctx)
		if plugin.logger != nil {
			pluginCtx.Logger = plugin.logger
		}
		if !plugin.executeMiddlewares(pluginCtx, bot.appData) {
			continue
		}
		if err := handler(pluginCtx, bot.appData); err != nil {
			pluginCtx.error(err)
		}
	}
}

func (bot *Bot[T]) prepareUpdateCtx(u *tgapi.Update, ctx *MsgContext) {
	var from *tgapi.User
	var chat *tgapi.Chat
	switch u.Type {
	case tgapi.UpdateTypeMessage:
		if u.Message != nil {
			ctx.Msg = u.Message
			if u.Message.Chat != nil {
				chat = u.Message.Chat
			}
			if u.Message.From != nil {
				from = u.Message.From
			}
		}
	case tgapi.UpdateTypeEditedMessage:
		if u.EditedMessage != nil {
			ctx.Msg = u.EditedMessage
			if u.EditedMessage.Chat != nil {
				chat = u.EditedMessage.Chat
			}
			if u.EditedMessage.From != nil {
				from = u.EditedMessage.From
			}
		}
	case tgapi.UpdateTypeChannelPost:
		if u.ChannelPost != nil {
			ctx.Msg = u.ChannelPost
			if u.ChannelPost.Chat != nil {
				chat = u.ChannelPost.Chat
			}
			if u.ChannelPost.From != nil {
				from = u.ChannelPost.From
			}
		}
	case tgapi.UpdateTypeEditedChannelPost:
		if u.EditedChannelPost != nil {
			ctx.Msg = u.EditedChannelPost
			if u.EditedChannelPost.Chat != nil {
				chat = u.EditedChannelPost.Chat
			}
			if u.EditedChannelPost.From != nil {
				from = u.EditedChannelPost.From
			}
		}
	case tgapi.UpdateTypeBusinessMessage:
		if u.BusinessMessage != nil {
			ctx.Msg = u.BusinessMessage
			if u.BusinessMessage.Chat != nil {
				chat = u.BusinessMessage.Chat
			}
			if u.BusinessMessage.From != nil {
				from = u.BusinessMessage.From
			}
		}
	case tgapi.UpdateTypeEditedBusinessMessage:
		if u.EditedBusinessMessage != nil {
			ctx.Msg = u.EditedBusinessMessage
			if u.EditedBusinessMessage.Chat != nil {
				chat = u.EditedBusinessMessage.Chat
			}
			if u.EditedBusinessMessage.From != nil {
				from = u.EditedBusinessMessage.From
			}
		}
	case tgapi.UpdateTypeInlineQuery:
		if u.InlineQuery != nil {
			from = &u.InlineQuery.From
		}
	case tgapi.UpdateTypeChosenInlineResult:
		if u.ChosenInlineResult != nil {
			from = &u.ChosenInlineResult.From
		}
	case tgapi.UpdateTypeCallbackQuery:
		if u.CallbackQuery != nil {
			if u.CallbackQuery.Message != nil {
				ctx.Msg = u.CallbackQuery.Message
				ctx.CallbackMsgId = u.CallbackQuery.Message.MessageID
				if u.CallbackQuery.Message.Chat != nil {
					chat = u.CallbackQuery.Message.Chat
				}
			}
			if u.CallbackQuery.InlineMessageID != nil {
				ctx.InlineMsgId = *u.CallbackQuery.InlineMessageID
			}
			ctx.CallbackQueryId = u.CallbackQuery.ID
			from = &u.CallbackQuery.From
		}
	case tgapi.UpdateTypeShippingQuery:
		if u.ShippingQuery != nil {
			from = &u.ShippingQuery.From
		}
	case tgapi.UpdateTypePreCheckoutQuery:
		if u.PreCheckoutQuery != nil {
			from = &u.PreCheckoutQuery.From
		}
	case tgapi.UpdateTypePurchasedPaidMedia:
		if u.PurchasedPaidMedia != nil {
			from = &u.PurchasedPaidMedia.From
		}
	case tgapi.UpdateTypeMyChatMember:
		if u.MyChatMember != nil {
			from = &u.MyChatMember.From
			chat = &u.MyChatMember.Chat
		}
	case tgapi.UpdateTypeChatMember:
		if u.ChatMember != nil {
			from = &u.ChatMember.From
			chat = &u.ChatMember.Chat
		}
	case tgapi.UpdateTypeChatJoinRequest:
		if u.ChatJoinRequest != nil {
			from = &u.ChatJoinRequest.From
			chat = &u.ChatJoinRequest.Chat

		}
	case tgapi.UpdateTypeBusinessConnection:
		if u.BusinessConnection != nil {
			from = &u.BusinessConnection.User
		}
	case tgapi.UpdateTypePollAnswer:
		if u.PollAnswer != nil {
			from = &u.PollAnswer.User
		}
	case tgapi.UpdateTypeMessageReaction:
		if u.MessageReaction != nil {
			from = u.MessageReaction.User
			chat = u.MessageReaction.Chat
		}
	case tgapi.UpdateTypeChatBoost:
		if u.ChatBoost != nil {
			from = &u.ChatBoost.Boost.Source.User
			chat = &u.ChatBoost.Chat
		}
	case tgapi.UpdateTypeRemovedChatBoost:
		if u.RemovedChatBoost != nil {
			from = &u.RemovedChatBoost.Source.User
			chat = &u.RemovedChatBoost.Chat
		}
	}
	if ctx.Msg != nil && from == nil {
		from = ctx.Msg.From
	}
	if from != nil {
		ctx.From = from
		ctx.FromID = from.ID
	} else {
		ctx.FromID = 0
	}
	if chat != nil {
		ctx.Chat = chat
		ctx.ChatID = chat.ID
	} else {
		ctx.ChatID = 0
	}
}
