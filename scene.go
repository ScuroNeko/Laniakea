package laniakea

import (
	"encoding/json"
	"sync"
)

type SceneHandler[T any] func(ctx *SceneContext, db T) (SceneResult, error)
type Scene[T any] struct {
	Name       string
	Scope      SceneScope
	Entry      string // starting step
	PluginName string

	steps    map[string]SceneHandler[T]
	commands map[string]SceneHandler[T]
	message  SceneHandler[T]
}

func NewScene[T any](name string) *Scene[T] {
	return &Scene[T]{
		Name:     name,
		Scope:    SceneScopeUserChat,
		Entry:    "",
		steps:    make(map[string]SceneHandler[T]),
		commands: make(map[string]SceneHandler[T]),
		message:  nil,
	}
}
func (s *Scene[T]) SetScope(scope SceneScope) *Scene[T] {
	s.Scope = scope
	return s
}
func (s *Scene[T]) SetEntry(step string) *Scene[T] {
	s.Entry = step
	return s
}
func (s *Scene[T]) setPluginName(name string) *Scene[T] {
	s.PluginName = name
	return s
}

func (s *Scene[T]) OnStep(step string, handler SceneHandler[T]) *Scene[T] {
	s.steps[step] = handler
	return s
}
func (s *Scene[T]) OnCommand(cmd string, handler SceneHandler[T]) *Scene[T] {
	s.commands[cmd] = handler
	return s
}
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

type SceneSession struct {
	Scene string
	Step  string
	Data  []byte
}

// SetData sets the session data. It is thread-safe and can be used to store any arbitrary data as a byte slice.
func (s *SceneSession) SetData(data []byte) {
	s.Data = data
}

// GetData retrieves the session data. It is thread-safe and returns the data as a byte slice.
func (s *SceneSession) GetData() []byte {
	return s.Data
}
func (s *SceneSession) HasData() bool {
	return len(s.Data) > 0
}

// ClearData clears the session data. It is thread-safe and sets the data to nil.
func (s *SceneSession) ClearData() {
	s.Data = nil
}

// BindData binds the session data to the provided struct.
func (s *SceneSession) BindData(v any) error {
	if len(s.Data) == 0 {
		return nil // No data to bind, return nil error
	}
	return json.Unmarshal(s.Data, v)
}

// SaveData saves the provided struct as JSON in the session data.
func (s *SceneSession) SaveData(v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	s.Data = data
	return nil
}

type SessionStore interface {
	Get(key string) (SceneSession, error)
	Set(key string, session SceneSession) error
	Delete(key string) error
}

type MemorySessionStore struct {
	store map[string]SceneSession
	mu    sync.RWMutex
}

func NewMemorySessionStore() *MemorySessionStore {
	return &MemorySessionStore{
		store: make(map[string]SceneSession),
	}
}
func (s *MemorySessionStore) Get(key string) (SceneSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if session, ok := s.store[key]; ok {
		return session, nil
	}
	return SceneSession{}, nil
}
func (s *MemorySessionStore) Set(key string, session SceneSession) error {
	s.mu.Lock()
	s.store[key] = session
	s.mu.Unlock()
	return nil
}
func (s *MemorySessionStore) Delete(key string) error {
	s.mu.Lock()
	delete(s.store, key)
	s.mu.Unlock()
	return nil
}

type SceneResult struct {
	Action SceneAction
	Next   string
}

type SceneAction int

const (
	SceneActionStay SceneAction = iota
	SceneActionNext
	SceneActionExit
	SceneActionPass
)

type SceneScope int

const (
	SceneScopeUser SceneScope = iota
	SceneScopeChat
	SceneScopeUserChat
)

type sceneRuntime interface {
	FindScene(name string) (*sceneMeta, bool)
	GetSession(key string) (SceneSession, error)
	SetSession(key string, session SceneSession) error
	DeleteSession(key string) error
	BuildSceneKey(scope SceneScope, ctx *MsgContext) (string, bool)
	FindSceneSession(ctx *MsgContext) (string, SceneSession, error)
}

type sceneMeta struct {
	Name  string
	Scope SceneScope
	Entry string
	Steps map[string]struct{}
}
