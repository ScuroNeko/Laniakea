package laniakea

type SceneContext struct {
	*MsgContext
	sess SceneSession
	key  string
}

func (ctx *SceneContext) Next(step string) SceneResult {
	return SceneResult{
		Action: SceneActionNext,
		Next:   step,
	}
}
func (ctx *SceneContext) Stay() SceneResult {
	return SceneResult{Action: SceneActionStay}
}
func (ctx *SceneContext) Exit() SceneResult {
	return SceneResult{Action: SceneActionExit}
}
func (ctx *SceneContext) Pass() SceneResult {
	return SceneResult{Action: SceneActionPass}
}
