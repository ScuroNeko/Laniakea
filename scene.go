package laniakea

import (
	"encoding/json"
	"maps"
	"sync"
)

// SceneHandler handles a scene step, scene command, or fallback message.
type SceneHandler[T any] func(ctx *SceneContext, db T) (SceneResult, error)

// Scene defines a multi-step conversational flow.
type Scene[T any] struct {
	name       string
	scope      SceneScope
	entry      string
	pluginName string

	steps    map[string]SceneHandler[T]
	commands map[string]SceneHandler[T]
	payloads map[string]SceneHandler[T]
	message  SceneHandler[T]
}

// NewScene creates a new scene with user-chat scope by default.
func NewScene[T any](name string) *Scene[T] {
	return &Scene[T]{
		name:     name,
		scope:    SceneScopeUserChat,
		entry:    "",
		steps:    make(map[string]SceneHandler[T]),
		commands: make(map[string]SceneHandler[T]),
		payloads: make(map[string]SceneHandler[T]),
		message:  nil,
	}
}

// SetScope changes how scene sessions are keyed and shared.
func (s *Scene[T]) SetScope(scope SceneScope) *Scene[T] {
	s.scope = scope
	return s
}

// SetEntry sets the initial step entered by MessageContext.EnterScene.
func (s *Scene[T]) SetEntry(step string) *Scene[T] {
	s.entry = step
	return s
}

func (s *Scene[T]) setPluginName(name string) *Scene[T] {
	s.pluginName = name
	return s
}

// OnStep registers a handler for a named scene step.
func (s *Scene[T]) OnStep(step string, handler SceneHandler[T]) *Scene[T] {
	s.steps[step] = handler
	return s
}

// OnCommand registers a command handler active while the scene is running.
func (s *Scene[T]) OnCommand(cmd string, handler SceneHandler[T]) *Scene[T] {
	s.commands[cmd] = handler
	return s
}

// OnPayload registers a callback payload handler active while the scene is running.
func (s *Scene[T]) OnPayload(cmd string, handler SceneHandler[T]) *Scene[T] {
	s.payloads[cmd] = handler
	return s
}

// OnMessage registers a fallback handler used when no scene command or step matches.
func (s *Scene[T]) OnMessage(handler SceneHandler[T]) *Scene[T] {
	s.message = handler
	return s
}

func (s *Scene[T]) executeCommand(cmd string, ctx *SceneContext, db T) (SceneResult, bool, error) {
	handler, ok := s.commands[cmd]
	if !ok {
		return SceneResult{}, false, nil
	}
	result, err := handler(ctx, db)
	return result, true, err
}
func (s *Scene[T]) executePayload(cmd string, ctx *SceneContext, db T) (SceneResult, bool, error) {
	handler, ok := s.payloads[cmd]
	if !ok {
		return SceneResult{}, false, nil
	}
	result, err := handler(ctx, db)
	return result, true, err
}
func (s *Scene[T]) executeStep(step string, ctx *SceneContext, db T) (SceneResult, bool, error) {
	handler, ok := s.steps[step]
	if !ok {
		return SceneResult{}, false, nil
	}
	result, err := handler(ctx, db)
	return result, true, err
}
func (s *Scene[T]) executeMessage(ctx *SceneContext, db T) (SceneResult, bool, error) {
	if s.message == nil {
		return SceneResult{}, false, nil
	}
	result, err := s.message(ctx, db)
	return result, true, err
}

func (s *Scene[T]) clone() *Scene[T] {
	if s == nil {
		return nil
	}

	cloned := *s
	cloned.steps = make(map[string]SceneHandler[T], len(s.steps))
	cloned.commands = make(map[string]SceneHandler[T], len(s.commands))
	cloned.payloads = make(map[string]SceneHandler[T], len(s.payloads))

	maps.Copy(cloned.steps, s.steps)
	maps.Copy(cloned.commands, s.commands)
	maps.Copy(cloned.payloads, s.payloads)

	return &cloned
}

// SceneSession stores the active scene state for one session key.
type SceneSession struct {
	// Scene is the registered scene name for the active session.
	Scene string
	// Step is the current step name inside the active scene.
	Step string
	// data stores opaque session payload bytes, typically JSON.
	data []byte
}

// SetData stores arbitrary opaque session data.
func (s *SceneSession) SetData(data []byte) {
	s.data = data
}

// GetData returns the raw session data payload.
func (s *SceneSession) GetData() []byte {
	return s.data
}

// HasData reports whether the session has a non-empty data payload.
func (s *SceneSession) HasData() bool {
	return len(s.data) > 0
}

// ClearData removes any stored session data.
func (s *SceneSession) ClearData() {
	s.data = nil
}

// BindData unmarshals the stored JSON payload into v.
func (s *SceneSession) BindData(v any) error {
	if len(s.data) == 0 {
		return nil
	}
	return json.Unmarshal(s.data, v)
}

// SaveData marshals v as JSON and stores it in the session.
func (s *SceneSession) SaveData(v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	s.data = data
	return nil
}

// SessionStore persists scene sessions by key.
type SessionStore interface {
	Get(key string) (SceneSession, error)
	Set(key string, session SceneSession) error
	Delete(key string) error
}

// MemorySessionStore stores scene sessions in memory.
type MemorySessionStore struct {
	store map[string]SceneSession
	mu    sync.RWMutex
}

// NewMemorySessionStore creates an empty in-memory session store.
func NewMemorySessionStore() *MemorySessionStore {
	return &MemorySessionStore{
		store: make(map[string]SceneSession),
	}
}

// Get returns the session stored under key, or the zero session when absent.
func (s *MemorySessionStore) Get(key string) (SceneSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if session, ok := s.store[key]; ok {
		return session, nil
	}
	return SceneSession{}, nil
}

// Set stores session under key.
func (s *MemorySessionStore) Set(key string, session SceneSession) error {
	s.mu.Lock()
	s.store[key] = session
	s.mu.Unlock()
	return nil
}

// Delete removes the session stored under key.
func (s *MemorySessionStore) Delete(key string) error {
	s.mu.Lock()
	delete(s.store, key)
	s.mu.Unlock()
	return nil
}

// SceneResult describes how scene execution should proceed after a handler returns.
type SceneResult struct {
	Action SceneAction
	Next   string
}

// SceneAction controls how the bot updates scene state after a handler returns.
type SceneAction int

const (
	// SceneActionStay keeps the current scene and step active.
	SceneActionStay SceneAction = iota
	// SceneActionNext moves the session to another named step.
	SceneActionNext
	// SceneActionExit removes the current scene session.
	SceneActionExit
	// SceneActionPass lets normal bot routing continue after the scene handler.
	SceneActionPass
)

// SceneScope defines how scene sessions are keyed.
type SceneScope int

const (
	// SceneScopeUser shares a scene across all chats for one user.
	SceneScopeUser SceneScope = iota
	// SceneScopeChat shares a scene across all users in one chat.
	SceneScopeChat
	// SceneScopeUserChat isolates a scene per user-chat pair.
	SceneScopeUserChat
)

type sceneRuntime interface {
	findScene(name string) (*sceneMeta, bool)
	getSession(key string) (SceneSession, error)
	setSession(key string, session SceneSession) error
	deleteSession(key string) error
	findSceneSession(ctx *MessageContext) (string, SceneSession, error)
}

type sceneMeta struct {
	Name  string
	Scope SceneScope
	Entry string
	Steps map[string]struct{}
}
