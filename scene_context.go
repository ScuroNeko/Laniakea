package laniakea

// SceneContext wraps MsgContext with scene session state for scene handlers.
type SceneContext struct {
	*MsgContext
	sess SceneSession
	key  string
}

// Next advances the current scene to step.
func (ctx *SceneContext) Next(step string) SceneResult {
	return SceneResult{
		Action: SceneActionNext,
		Next:   step,
	}
}

// Stay keeps the current scene step active.
func (ctx *SceneContext) Stay() SceneResult {
	return SceneResult{Action: SceneActionStay}
}

// Exit leaves the current scene.
func (ctx *SceneContext) Exit() SceneResult {
	return SceneResult{Action: SceneActionExit}
}

// Pass stops scene handling and lets normal routing continue.
func (ctx *SceneContext) Pass() SceneResult {
	return SceneResult{Action: SceneActionPass}
}

// BindData unmarshals the current scene session payload into v.
func (ctx *SceneContext) BindData(v any) error {
	return ctx.sess.BindData(v)
}

// SaveData marshals v and stores it in the current scene session payload.
func (ctx *SceneContext) SaveData(v any) error {
	return ctx.sess.SaveData(v)
}
