package laniakea

import (
	"math/rand/v2"
	"sync"
	"sync/atomic"

	"git.scuroneko.dev/scuroneko/laniakea/tgapi"
)

// Interface for generating unique draft IDs.
type draftIDGenerator interface {
	// Next returns the next unique draft ID.
	Next() uint64
}

// RandomDraftIDGenerator generates draft IDs using cryptographically secure random numbers.
// Suitable for distributed systems or when ID predictability is undesirable.
type RandomDraftIDGenerator struct{}

// Next returns a random 64-bit unsigned integer.
func (g *RandomDraftIDGenerator) Next() uint64 {
	return rand.Uint64()
}

// LinearDraftIDGenerator generates draft IDs using a monotonically increasing counter.
// Useful for debugging, persistence, or when drafts must be ordered.
type LinearDraftIDGenerator struct {
	lastID atomic.Uint64
}

// Next returns the next linear ID, atomically incremented.о
func (g *LinearDraftIDGenerator) Next() uint64 {
	return g.lastID.Add(1)
}

// DraftProvider manages a collection of Drafts and a shared draft ID generator.
//
// DraftProvider is safe for concurrent use.
type DraftProvider struct {
	mu        sync.RWMutex
	api       *tgapi.API
	drafts    map[uint64]*Draft
	generator draftIDGenerator
}

// NewRandomDraftProvider creates a new DraftProvider using random draft IDs.
//
// The provider will use cryptographically secure random numbers for draft IDs.
// All drafts created via this provider will have unpredictable, unique IDs.
func NewRandomDraftProvider(api *tgapi.API) *DraftProvider {
	return &DraftProvider{
		api: api, generator: &RandomDraftIDGenerator{},
		drafts: make(map[uint64]*Draft),
	}
}

// NewLinearDraftProvider creates a new DraftProvider using linear (incrementing) draft IDs.
//
// startValue is the initial value for the counter. Use 0 for fresh start, or a known
// value to resume from persisted state.
//
// This is useful when you need to store draft IDs externally (e.g., in a database)
// and want to reconstruct drafts after restart.
func NewLinearDraftProvider(api *tgapi.API, startValue uint64) *DraftProvider {
	g := &LinearDraftIDGenerator{}
	g.lastID.Store(startValue)
	return &DraftProvider{
		api:       api,
		generator: g,
		drafts:    make(map[uint64]*Draft),
	}
}

// GetDraft retrieves a draft by its ID.
//
// Returns the draft and true if found, or nil and false if not found.
func (p *DraftProvider) GetDraft(id uint64) (*Draft, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	draft, ok := p.drafts[id]
	return draft, ok
}

// FlushAll sends all pending drafts as final messages and clears them.
//
// If one or more drafts fail to send, FlushAll still attempts all drafts and
// returns the first encountered error.
//
// After successful flush, each draft is removed from the provider and cleared.
func (p *DraftProvider) FlushAll() error {
	p.mu.RLock()
	drafts := make([]*Draft, 0, len(p.drafts))
	for _, draft := range p.drafts {
		drafts = append(drafts, draft)
	}
	p.mu.RUnlock()

	var firstErr error
	for _, draft := range drafts {
		if err := draft.Flush(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// Draft represents a single message draft that can be edited and flushed.
//
// Drafts are safe to use from a single goroutine. Multiple goroutines must
// synchronize access manually.
//
// Drafts are automatically removed from the provider's map when Flush() succeeds.
type Draft struct {
	api      *tgapi.API
	provider *DraftProvider

	chatID          int64
	messageThreadID int
	parseMode       tgapi.ParseMode
	entities        []tgapi.MessageEntity

	ID      uint64
	Message string
}

// NewDraft creates a new draft with the provided parse mode.
//
// The caller must set a chat with SetChat before Push or Flush.
func (p *DraftProvider) NewDraft(parseMode tgapi.ParseMode) *Draft {
	id := p.generator.Next()
	draft := &Draft{
		api:       p.api,
		provider:  p,
		parseMode: parseMode,
		ID:        id,
		Message:   "",
	}
	p.mu.Lock()
	p.drafts[id] = draft
	p.mu.Unlock()
	return draft
}

// SetChat overrides the draft's target chat and message thread.
//
// This is useful for sending a draft to a different chat than the provider's default.
func (d *Draft) SetChat(chatID int64, messageThreadID int) *Draft {
	d.chatID = chatID
	d.messageThreadID = messageThreadID
	return d
}

// SetEntities replaces the draft's message entities.
//
// Entities are stored by reference. If you plan to mutate the slice later,
// pass a copy: `SetEntities(append([]tgapi.MessageEntity{}, myEntities...))`.
func (d *Draft) SetEntities(entities []tgapi.MessageEntity) *Draft {
	d.entities = entities
	return d
}

// Push appends text to the draft and attempts to update the server-side draft.
//
// Returns an error if the Telegram API rejects the update (e.g., due to network issues).
// The draft's Message field is always updated, even if the API call fails.
//
// Use this method to build the message incrementally.
func (d *Draft) Push(text string) error {
	return d.push(text)
}

// GetMessage returns the current content of the draft.
//
// Useful for inspection, logging, or validation before flushing.
func (d *Draft) GetMessage() string {
	return d.Message
}

// Clear resets the draft's message content to empty string.
//
// Does not affect server-side draft — use Flush() for that.
func (d *Draft) Clear() {
	d.Message = ""
}

// Delete removes the draft from its provider and clears its content.
//
// This is an internal method used by Flush(). You may call it manually if you
// want to cancel a draft without sending it.
func (d *Draft) Delete() {
	if d.provider != nil {
		d.provider.mu.Lock()
		delete(d.provider.drafts, d.ID)
		d.provider.mu.Unlock()
	}
	d.Clear()
}

// Flush sends the draft as a final message and clears it locally.
//
// If successful:
//   - The message is sent to Telegram.
//   - The draft's content is cleared.
//   - The draft is removed from the provider's map.
//
// If an error occurs:
//   - The message is NOT sent.
//   - The draft remains in the provider and retains its content.
//   - You can call Flush() again to retry.
//
// If the draft is empty, Flush() returns nil without calling the API.
func (d *Draft) Flush() error {
	if d.Message == "" {
		return nil
	}
	if d.chatID == 0 {
		return ErrDraftChatIDZero
	}
	if err := validateMessageText(d.Message); err != nil {
		return err
	}

	params := tgapi.SendMessage{
		ChatID:    d.chatID,
		ParseMode: d.parseMode,
		Entities:  d.entities,
		Text:      d.Message,
	}
	if d.messageThreadID > 0 {
		params.MessageThreadID = d.messageThreadID
	}

	_, err := d.api.SendMessage(params)
	if err == nil {
		d.Delete()
	}
	return err
}

// Internal helper for Push that updates the server-side draft.
func (d *Draft) push(text string) error {
	if d.chatID == 0 {
		return ErrDraftChatIDZero
	}
	d.Message += text
	if err := validateMessageText(d.Message); err != nil {
		return err
	}
	params := tgapi.SendMessageDraft{
		ChatID:    d.chatID,
		DraftID:   d.ID,
		Text:      d.Message,
		ParseMode: d.parseMode,
		Entities:  d.entities,
	}
	if d.messageThreadID > 0 {
		params.MessageThreadID = d.messageThreadID
	}
	_, err := d.api.SendMessageDraft(params)
	return err
}
