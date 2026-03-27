package laniakea

func (bot *Bot[T]) GetSession(key string) (SceneSession, error) {
	return bot.sessionStore.Get(key)
}

func (bot *Bot[T]) SetSession(key string, session SceneSession) error {
	return bot.sessionStore.Set(key, session)
}

func (bot *Bot[T]) DeleteSession(key string) error {
	return bot.sessionStore.Delete(key)
}
func (bot *Bot[T]) FindScene(name string) (*sceneMeta, bool) {
	for _, plugin := range bot.plugins {
		scene, ok := plugin.scenes[name]
		if !ok {
			continue
		}

		steps := make(map[string]struct{}, len(scene.steps))
		for step := range scene.steps {
			steps[step] = struct{}{}
		}

		return &sceneMeta{
			Name:  scene.Name,
			Scope: scene.Scope,
			Entry: scene.Entry,
			Steps: steps,
		}, true
	}
	return nil, false
}

func (bot *Bot[T]) FindSceneSession(ctx *MsgContext) (string, SceneSession, error) {
	var zero SceneSession
	if ctx.Msg == nil {
		return "", zero, ErrMessageNil
	}

	for _, scope := range bot.sceneScopePriority {
		key, ok := buildSceneKey(scope, ctx)
		if !ok {
			continue
		}

		session, err := bot.sessionStore.Get(key)
		if err != nil {
			return "", zero, err
		}
		if session.Scene != "" {
			return key, session, nil
		}
	}

	return "", zero, ErrCantFindSession
}
func (bot *Bot[T]) BuildSceneKey(scope SceneScope, ctx *MsgContext) (string, bool) {
	return buildSceneKey(scope, ctx)
}
