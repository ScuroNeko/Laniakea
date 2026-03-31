package laniakea

import (
	"context"
	"fmt"
	"time"

	"git.scuroneko.dev/scuroneko/laniakea/tgapi"
)

// HandlerEventKind identifies the kind of handler observed by runtime events.
type HandlerEventKind string

const (
	// HandlerCommandKind identifies a command handler.
	HandlerCommandKind HandlerEventKind = "command"
	// HandlerPayloadKind identifies a callback payload handler.
	HandlerPayloadKind HandlerEventKind = "payload"
	// HandlerUpdateKind identifies a generic update handler.
	HandlerUpdateKind HandlerEventKind = "update"
	// HandlerRunnerKind identifies a background runner execution.
	HandlerRunnerKind HandlerEventKind = "runner"
	// HandlerPollingKind identifies polling and getUpdates runtime work.
	HandlerPollingKind HandlerEventKind = "polling"
	// HandlerSceneKind identifies a scene runtime handler wrapper.
	HandlerSceneKind HandlerEventKind = "scene"
	// HandlerSceneStepKind identifies a scene step handler.
	HandlerSceneStepKind HandlerEventKind = "scene_step"
	// HandlerSceneCommandKind identifies a scene-local command handler.
	HandlerSceneCommandKind HandlerEventKind = "scene_command"
	// HandlerSceneMessageKind identifies a scene message fallback handler.
	HandlerSceneMessageKind HandlerEventKind = "scene_message"
)

type Event interface {
	isEvent()
}

// UpdateReceivedEvent describes an update entering the bot runtime.
type UpdateReceivedEvent struct {
	UpdateID   int
	UpdateType tgapi.UpdateType
	FromID     int64
	ChatID     int64
}

// UpdateHandledEvent describes a completed update execution path.
type UpdateHandledEvent struct {
	UpdateID   int
	UpdateType tgapi.UpdateType
	FromID     int64
	ChatID     int64
	Duration   time.Duration
	Handled    bool
}

// HandlerStartedEvent describes a handler about to execute.
type HandlerStartedEvent struct {
	UpdateID    int
	UpdateType  tgapi.UpdateType
	Plugin      string
	HandlerKind HandlerEventKind
	HandlerName string
	FromID      int64
	ChatID      int64
}

// HandlerFinishedEvent describes a handler that has completed.
type HandlerFinishedEvent struct {
	UpdateID    int
	UpdateType  tgapi.UpdateType
	Plugin      string
	HandlerKind HandlerEventKind
	HandlerName string
	FromID      int64
	ChatID      int64
	Duration    time.Duration
	Err         error
	UserFacing  bool
}

// SceneTransitionEvent describes a scene state transition.
type SceneTransitionEvent struct {
	Plugin string
	Scene  string
	From   string
	To     string
	Action SceneAction
	FromID int64
	ChatID int64
}

// PolicyCheckedEvent describes the result of a policy evaluation.
type PolicyCheckedEvent struct {
	Name     string
	Plugin   string
	FromID   int64
	ChatID   int64
	Passed   bool
	Err      error
	Internal bool
}

// RunnerFinishedEvent describes a completed background runner execution.
type RunnerFinishedEvent struct {
	Name     string
	Duration time.Duration
	Err      error
}

// PollingRetryEvent describes a polling retry after a failed getUpdates call.
type PollingRetryEvent struct {
	Attempt int
	Delay   time.Duration
	Err     error
}

// ErrorEvent describes an error routed through framework error handling.
type ErrorEvent struct {
	UpdateID    int
	UpdateType  tgapi.UpdateType
	Plugin      string
	HandlerKind HandlerEventKind
	HandlerName string
	FromID      int64
	ChatID      int64
	Err         error
	UserFacing  bool
}

func (UpdateReceivedEvent) isEvent()  {}
func (UpdateHandledEvent) isEvent()   {}
func (HandlerStartedEvent) isEvent()  {}
func (HandlerFinishedEvent) isEvent() {}
func (SceneTransitionEvent) isEvent() {}
func (PolicyCheckedEvent) isEvent()   {}
func (RunnerFinishedEvent) isEvent()  {}
func (PollingRetryEvent) isEvent()    {}
func (ErrorEvent) isEvent()           {}

// Observer receives best-effort runtime instrumentation events.
type Observer interface {
	OnReceiveUpdate(ctx context.Context, event UpdateReceivedEvent)
	OnHandledUpdate(ctx context.Context, event UpdateHandledEvent)
	OnHandlerStarted(ctx context.Context, event HandlerStartedEvent)
	OnHandlerFinished(ctx context.Context, event HandlerFinishedEvent)
	OnSceneTransition(ctx context.Context, event SceneTransitionEvent)
	OnPolicyChecked(ctx context.Context, event PolicyCheckedEvent)
	OnRunnerFinished(ctx context.Context, event RunnerFinishedEvent)
	OnPollingRetry(ctx context.Context, event PollingRetryEvent)
	OnError(ctx context.Context, event ErrorEvent)
}

func (bot *Bot[T]) safeEmitEvent(ctx context.Context, event Event) {
	if bot.observer == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			bot.logger.Errorln(fmt.Sprintf("panic in observer: %v", r))
		}
	}()
	switch e := event.(type) {
	case UpdateReceivedEvent:
		bot.observer.OnReceiveUpdate(ctx, e)
	case UpdateHandledEvent:
		bot.observer.OnHandledUpdate(ctx, e)
	case HandlerStartedEvent:
		bot.observer.OnHandlerStarted(ctx, e)
	case HandlerFinishedEvent:
		bot.observer.OnHandlerFinished(ctx, e)
	case SceneTransitionEvent:
		bot.observer.OnSceneTransition(ctx, e)
	case PolicyCheckedEvent:
		bot.observer.OnPolicyChecked(ctx, e)
	case RunnerFinishedEvent:
		bot.observer.OnRunnerFinished(ctx, e)
	case PollingRetryEvent:
		bot.observer.OnPollingRetry(ctx, e)
	case ErrorEvent:
		bot.observer.OnError(ctx, e)
	}
}
