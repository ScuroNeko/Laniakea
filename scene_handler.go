package laniakea

import (
	"errors"
	"fmt"
	"strings"
)

func (bot *Bot[T]) tryHandleScene(ctx *MsgContext) (bool, error) {
	key, session, err := bot.findSceneSession(ctx)
	if err != nil {
		if errors.Is(err, ErrCantFindSession) || errors.Is(err, ErrMessageNil) {
			return false, nil
		}
		return false, err
	}
	if session.Scene == "" {
		return false, nil
	}

	for _, plugin := range bot.plugins {
		scene, ok := plugin.scenes[session.Scene]
		if !ok {
			continue
		}
		if scene.PluginName != "" && scene.PluginName != plugin.name {
			continue
		}
		if !plugin.executeMiddlewares(ctx, bot.appData) {
			return false, nil
		}
		sceneCtx := &SceneContext{
			MsgContext: ctx,
			sess:       session,
			key:        key,
		}
		return bot.executeScene(scene, sceneCtx)
	}
	return false, ErrSceneNotFound
}

func (bot *Bot[T]) executeScene(scene *Scene[T], ctx *SceneContext) (bool, error) {
	if ctx.MsgContext == nil || ctx.sess.Scene == "" {
		return false, nil
	}

	var text string
	if ctx.Msg != nil {
		text = ctx.Msg.Text
		if text == "" {
			text = ctx.Msg.Caption
		}
	}

	text = strings.TrimSpace(text)
	prefix, cmd, args := bot.parseCommand(text)
	if cmd != "" {
		ctx.Prefix = prefix
		ctx.Text = args
		ctx.Args = strings.Fields(args)

		res, matched, err := scene.executeCommand(cmd, ctx, bot.appData)
		if err != nil {
			return false, err
		}
		if matched {
			return bot.applySceneResult(scene, ctx, res)
		}
	}
	ctx.Text = text
	ctx.Args = nil
	ctx.Prefix = ""
	if ctx.sess.Step != "" {
		res, matched, err := scene.executeStep(ctx.sess.Step, ctx, bot.appData)
		if err != nil {
			return false, err
		}
		if matched {
			return bot.applySceneResult(scene, ctx, res)
		}
	}

	res, matched, err := scene.executeMessage(ctx, bot.appData)
	if err != nil {
		return false, err
	}
	if matched {
		return bot.applySceneResult(scene, ctx, res)
	}

	return false, nil
}
func (bot *Bot[T]) applySceneResult(scene *Scene[T], ctx *SceneContext, result SceneResult) (bool, error) {
	switch result.Action {
	case SceneActionStay:
		if err := bot.sessionStore.Set(ctx.key, ctx.sess); err != nil {
			return false, err
		}
		return true, nil
	case SceneActionNext:
		if result.Next == "" {
			return false, ErrSceneStepNotFound
		}
		if _, ok := scene.steps[result.Next]; !ok {
			return false, ErrSceneStepNotFound
		}
		ctx.sess.Step = result.Next
		if err := bot.sessionStore.Set(ctx.key, ctx.sess); err != nil {
			return false, err
		}
		return true, nil
	case SceneActionExit:
		if err := bot.sessionStore.Delete(ctx.key); err != nil {
			return false, err
		}
		return true, nil
	case SceneActionPass:
		return false, nil
	default:
		return false, nil
	}
}
func buildSceneKey(scope SceneScope, ctx *MsgContext) (string, bool) {
	if ctx == nil {
		return "", false
	}

	switch scope {
	case SceneScopeUserChat:
		if ctx.Msg == nil || ctx.Msg.Chat == nil || ctx.FromID == 0 {
			return "", false
		}
		return fmt.Sprintf("user_id:%d:chat_id:%d", ctx.FromID, ctx.Msg.Chat.ID), true
	case SceneScopeChat:
		if ctx.Msg == nil || ctx.Msg.Chat == nil {
			return "", false
		}
		return fmt.Sprintf("chat_id:%d", ctx.Msg.Chat.ID), true
	case SceneScopeUser:
		if ctx.FromID == 0 {
			return "", false
		}
		return fmt.Sprintf("user_id:%d", ctx.FromID), true
	default:
		return "", false
	}
}
