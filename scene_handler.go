package laniakea

import (
	"errors"
	"fmt"
	"strings"
	"time"
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

		return bot.executeScene(sceneCtx, scene)
	}
	return false, ErrSceneNotFound
}

func (bot *Bot[T]) executeScene(ctx *SceneContext, scene *Scene[T]) (bool, error) {
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

		if _, ok := scene.commands[cmd]; ok {
			startTime := time.Now()
			bot.emitSceneStarted(ctx, scene, HandlerSceneCommandKind, cmd)
			res, _, err := scene.executeCommand(cmd, ctx, bot.appData)
			if err != nil {
				bot.emitSceneFinished(ctx, scene, HandlerSceneCommandKind, cmd, startTime, err)
				bot.emitSceneError(ctx, scene, HandlerSceneCommandKind, cmd, err)
				return false, err
			}
			from := ctx.sess.Step
			ok, err := bot.applySceneResult(scene, ctx, res)
			bot.emitSceneFinished(ctx, scene, HandlerSceneCommandKind, cmd, startTime, err)
			if err != nil {
				bot.emitSceneError(ctx, scene, HandlerSceneCommandKind, cmd, err)
			}
			if ok {
				bot.emitSceneTransition(ctx, scene, from, res)
			}
			return ok, err
		}

		// Unmatched slash-commands should continue through normal bot command routing
		// instead of also triggering the active scene step or fallback handler.
		return false, nil
	}
	ctx.Text = text
	ctx.Args = nil
	ctx.Prefix = ""
	if ctx.sess.Step != "" {
		step := ctx.sess.Step
		if _, ok := scene.steps[step]; ok {
			startTime := time.Now()
			bot.emitSceneStarted(ctx, scene, HandlerSceneStepKind, step)
			res, _, err := scene.executeStep(step, ctx, bot.appData)
			if err != nil {
				bot.emitSceneFinished(ctx, scene, HandlerSceneStepKind, step, startTime, err)
				bot.emitSceneError(ctx, scene, HandlerSceneStepKind, step, err)
				return false, err
			}
			from := step
			ok, err := bot.applySceneResult(scene, ctx, res)
			bot.emitSceneFinished(ctx, scene, HandlerSceneStepKind, step, startTime, err)
			if err != nil {
				bot.emitSceneError(ctx, scene, HandlerSceneStepKind, from, err)
			}
			if ok {
				bot.emitSceneTransition(ctx, scene, from, res)
			}
			return ok, err
		}
	}

	if scene.message != nil {
		startTime := time.Now()
		bot.emitSceneStarted(ctx, scene, HandlerSceneMessageKind, "message_fallback")
		res, _, err := scene.executeMessage(ctx, bot.appData)
		if err != nil {
			bot.emitSceneFinished(ctx, scene, HandlerSceneMessageKind, "message_fallback", startTime, err)
			bot.emitSceneError(ctx, scene, HandlerSceneMessageKind, "message_fallback", err)
			return false, err
		}
		from := ctx.sess.Step
		ok, err := bot.applySceneResult(scene, ctx, res)
		bot.emitSceneFinished(ctx, scene, HandlerSceneMessageKind, "message_fallback", startTime, err)
		if err != nil {
			bot.emitSceneError(ctx, scene, HandlerSceneMessageKind, "message_fallback", err)
		}
		if ok {
			bot.emitSceneTransition(ctx, scene, from, res)
		}
		return ok, err
	}

	return false, nil
}

func (bot *Bot[T]) emitSceneStarted(ctx *SceneContext, scene *Scene[T], kind HandlerEventKind, name string) {
	bot.safeEmitEvent(ctx.Context(), HandlerStartedEvent{
		UpdateID:    ctx.Update.UpdateID,
		UpdateType:  ctx.Update.Type,
		Plugin:      scene.PluginName,
		HandlerKind: kind,
		HandlerName: name,
		FromID:      ctx.FromID,
		ChatID:      ctx.ChatID,
	})
}

func (bot *Bot[T]) emitSceneFinished(ctx *SceneContext, scene *Scene[T], kind HandlerEventKind, name string, startedAt time.Time, err error) {
	bot.safeEmitEvent(ctx.Context(), HandlerFinishedEvent{
		UpdateID:    ctx.Update.UpdateID,
		UpdateType:  ctx.Update.Type,
		Plugin:      scene.PluginName,
		HandlerKind: kind,
		HandlerName: name,
		FromID:      ctx.FromID,
		ChatID:      ctx.ChatID,
		Duration:    time.Since(startedAt),
		Err:         err,
		UserFacing:  IsUserError(err),
	})
}

func (bot *Bot[T]) emitSceneError(ctx *SceneContext, scene *Scene[T], kind HandlerEventKind, name string, err error) {
	bot.safeEmitEvent(ctx.Context(), ErrorEvent{
		UpdateID:    ctx.Update.UpdateID,
		UpdateType:  ctx.Update.Type,
		Plugin:      scene.PluginName,
		HandlerKind: kind,
		HandlerName: name,
		FromID:      ctx.FromID,
		ChatID:      ctx.ChatID,
		Err:         err,
		UserFacing:  IsUserError(err),
	})
}

func (bot *Bot[T]) emitSceneTransition(ctx *SceneContext, scene *Scene[T], from string, result SceneResult) {
	if result.Action == SceneActionPass {
		return
	}

	to := from
	switch result.Action {
	case SceneActionNext:
		to = result.Next
	case SceneActionExit:
		to = ""
	}

	bot.safeEmitEvent(ctx.Context(), SceneTransitionEvent{
		Plugin: scene.PluginName,
		Scene:  scene.Name,
		From:   from,
		To:     to,
		Action: result.Action,
		FromID: ctx.FromID,
		ChatID: ctx.ChatID,
	})
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
