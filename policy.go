package laniakea

import (
	"errors"
	"fmt"

	"git.scuroneko.dev/scuroneko/laniakea/tgapi"
)

// Policy defines a reusable authorization rule for the current update context.
type Policy[T AppData] func(ctx *MsgContext, data T) error

// RequirePolicy adapts a Policy into a blocking middleware.
func RequirePolicy[T AppData](name string, p Policy[T]) Middleware[T] {
	return NewMiddleware(name, func(ctx *MsgContext, data T) bool {
		if err := p(ctx, data); err != nil {
			ctx.emitPolicyChecked(PolicyCheckedEvent{
				Name:     name,
				FromID:   ctx.FromID,
				ChatID:   ctx.ChatID,
				Passed:   false,
				Err:      err,
				Internal: IsInternalError(err),
			})
			ctx.error(err)
			return false
		}
		ctx.emitPolicyChecked(PolicyCheckedEvent{
			Name:   name,
			FromID: ctx.FromID,
			ChatID: ctx.ChatID,
			Passed: true,
		})
		return true
	})
}

// AllPolicies composes policies that all must succeed.
func AllPolicies[T AppData](policies ...Policy[T]) Policy[T] {
	return func(ctx *MsgContext, data T) error {
		for _, p := range policies {
			if err := p(ctx, data); err != nil {
				return err
			}
		}
		return nil
	}
}

// AnyPolicy composes policies where at least one must succeed.
func AnyPolicy[T AppData](policies ...Policy[T]) Policy[T] {
	return func(ctx *MsgContext, data T) error {
		var firstDeny error
		var internalErr error
		for _, p := range policies {
			err := p(ctx, data)
			if err == nil {
				return nil
			}
			if IsInternalError(err) {
				if internalErr == nil {
					internalErr = err
				}
				continue
			}
			if firstDeny == nil {
				firstDeny = err
			}
		}
		if internalErr != nil {
			return internalErr
		}
		if firstDeny != nil {
			return firstDeny
		}
		return AsUserError(errors.New("no policy matched"))
	}
}

// NotPolicy inverts a policy deny result while preserving internal failures.
func NotPolicy[T AppData](policy Policy[T]) Policy[T] {
	return func(ctx *MsgContext, data T) error {
		var err error
		if err = policy(ctx, data); err == nil {
			return AsUserError(errors.New("the action is not allowed due to policy violation"))
		}
		if IsInternalError(err) {
			return err
		}
		return nil
	}
}

// RequirePrivateChat allows execution only in private chats.
func RequirePrivateChat[T AppData]() Policy[T] {
	return func(ctx *MsgContext, data T) error {
		if ctx.Msg == nil || ctx.Msg.Chat == nil {
			return AsInternalError(errors.New("private-chat policy requires message chat context"))
		}

		if ctx.Msg.Chat.Type != tgapi.ChatTypePrivate {
			return AsUserError(errors.New("this action is only available in private chat"))
		}

		return nil
	}
}

// RequireGroupChat allows execution only in group or supergroup chats.
func RequireGroupChat[T AppData]() Policy[T] {
	return func(ctx *MsgContext, data T) error {
		if ctx.Msg == nil || ctx.Msg.Chat == nil {
			return AsInternalError(errors.New("group-chat policy requires message chat context"))
		}

		if ctx.Msg.Chat.Type != tgapi.ChatTypeGroup && ctx.Msg.Chat.Type != tgapi.ChatTypeSupergroup {
			return AsUserError(errors.New("this action is only available in group chats"))
		}

		return nil
	}
}

// RequireSupergroupChat allows execution only in supergroup chats.
func RequireSupergroupChat[T AppData]() Policy[T] {
	return func(ctx *MsgContext, data T) error {
		if ctx.Msg == nil || ctx.Msg.Chat == nil {
			return AsInternalError(errors.New("supergroup-chat policy requires message chat context"))
		}

		if ctx.Msg.Chat.Type != tgapi.ChatTypeSupergroup {
			return AsUserError(errors.New("this action is only available in supergroup chats"))
		}

		return nil
	}
}

// RequireChatAdmin allows execution only for chat administrators or owners.
func RequireChatAdmin[T AppData]() Policy[T] {
	return func(ctx *MsgContext, data T) error {
		if ctx.FromID == 0 || ctx.ChatID == 0 {
			return AsInternalError(errors.New("chat-admin policy requires message chat context"))
		}

		member, err := ctx.Api.GetChatMember(tgapi.GetChatMember{
			ChatID: ctx.ChatID,
			UserID: ctx.FromID,
		})
		if err != nil {
			return AsInternalError(fmt.Errorf("failed to fetch chat member status: %w", err))
		}

		if member.Status != tgapi.ChatMemberStatusAdministrator && member.Status != tgapi.ChatMemberStatusOwner {
			return AsUserError(errors.New("this action is only available to chat admins"))
		}

		return nil
	}
}

// RequireChatCreator allows execution only for the chat owner.
func RequireChatCreator[T AppData]() Policy[T] {
	return func(ctx *MsgContext, data T) error {
		if ctx.FromID == 0 || ctx.ChatID == 0 {
			return AsInternalError(errors.New("chat-creator policy requires message chat context"))
		}

		member, err := ctx.Api.GetChatMember(tgapi.GetChatMember{
			ChatID: ctx.ChatID,
			UserID: ctx.FromID,
		})
		if err != nil {
			return AsInternalError(fmt.Errorf("failed to fetch chat creator: %w", err))
		}

		if member.Status != tgapi.ChatMemberStatusOwner {
			return AsUserError(errors.New("this action is only available to the chat creator"))
		}

		return nil
	}
}

// RequireBotAdmin allows execution only when the bot is an admin in the chat.
func RequireBotAdmin[T AppData]() Policy[T] {
	return func(ctx *MsgContext, data T) error {
		if ctx.ChatID == 0 {
			return AsInternalError(errors.New("bot-admin policy requires message chat context"))
		}

		bot, err := ctx.Api.GetMe()
		if err != nil {
			return AsInternalError(fmt.Errorf("failed to fetch bot info: %w", err))
		}

		member, err := ctx.Api.GetChatMember(tgapi.GetChatMember{
			ChatID: ctx.ChatID,
			UserID: bot.ID,
		})
		if err != nil {
			return AsInternalError(fmt.Errorf("failed to fetch bot member status: %w", err))
		}

		if member.Status != tgapi.ChatMemberStatusAdministrator && member.Status != tgapi.ChatMemberStatusOwner {
			return AsUserError(errors.New("this action requires the bot to be an admin in the chat"))
		}

		return nil
	}
}

// RequireCallbackFromUser allows execution only for callback queries sent by non-bot users.
func RequireCallbackFromUser[T AppData]() Policy[T] {
	return func(ctx *MsgContext, data T) error {
		if ctx.Update.CallbackQuery == nil {
			return AsInternalError(errors.New("callback-user policy requires callback query context"))
		}
		if ctx.Update.CallbackQuery.From.IsBot {
			return AsUserError(errors.New("this action is only available to human users"))
		}
		return nil
	}
}
