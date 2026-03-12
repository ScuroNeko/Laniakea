// Package laniakea provides a safe, high-level interface for managing Telegram
// message drafts using the tgapi library. It allows creating, editing, and
// flushing drafts with automatic ID generation and optional bulk flushing.
//
// Drafts are designed to be ephemeral, mutable buffers that can be built up
// incrementally and then sent as final messages. The package ensures safe
// state management by copying entities and isolating draft contexts.
//
// Two draft ID generation strategies are supported:
//   - Random: Cryptographically secure random IDs (default). Ideal for distributed systems.
//   - Linear: Monotonically increasing IDs. Useful for persistence, debugging, or recovery.
//
// Example usage:
//
//	provider := laniakea.NewRandomDraftProvider(api)
//
//	draft := provider.NewDraft(tgapi.ParseModeMarkdown)
//	draft.SetChat(-1001234567890, 0)
//	draft.Push("*Hello*").Push(" **world**!")
//	err := draft.Flush() // Sends message and deletes draft
//	if err != nil {
//	    log.Printf("Failed to send draft: %v", err)
//	}
//
//	// Or flush all pending drafts at once:
//	err = provider.FlushAll() // Sends all drafts and clears them
//
// Note: Drafts are NOT thread-safe. Concurrent access requires external synchronization.
package laniakea

import (
	"math/rand/v2"
	"sync/atomic"

	"git.nix13.pw/scuroneko/laniakea/tgapi"
)

// draftIdGenerator defines an interface for generating unique draft IDs.
type draftIdGenerator interface {
	// Next returns the next unique draft ID.
	Next() uint64
}

// RandomDraftIdGenerator generates draft IDs using cryptographically secure random numbers.
// Suitable for distributed systems or when ID predictability is undesirable.
type RandomDraftIdGenerator struct{}

// Next returns a random 64-bit unsigned integer.
func (g *RandomDraftIdGenerator) Next() uint64 {
	return rand.Uint64()
}

// LinearDraftIdGenerator generates draft IDs using a monotonically increasing counter.
// Useful for debugging, persistence, or when drafts must be ordered.
type LinearDraftIdGenerator struct {
	lastId atomic.Uint64
}

// Next returns the next linear ID, atomically incremented.
func (g *LinearDraftIdGenerator) Next() uint64 {
	return g.lastId.Add(1)
}

// DraftProvider manages a collection of Drafts and provides methods to create and
// configure them. It holds shared configuration (chat, parse mode, entities) and
// a draft ID generator.
//
// DraftProvider is NOT thread-safe. Concurrent access from multiple goroutines
// requires external synchronization.
type DraftProvider struct {
	api       *tgapi.API
	drafts    map[uint64]*Draft
	generator draftIdGenerator

	// Internal defaults — not exposed directly to users.
	chatID          int64
	messageThreadID int
	parseMode       tgapi.ParseMode
	entities        []tgapi.MessageEntity
}

// NewRandomDraftProvider creates a new DraftProvider using random draft IDs.
//
// The provider will use cryptographically secure random numbers for draft IDs.
// All drafts created via this provider will have unpredictable, unique IDs.
func NewRandomDraftProvider(api *tgapi.API) *DraftProvider {
	return &DraftProvider{
		api: api, generator: &RandomDraftIdGenerator{},
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
	g := &LinearDraftIdGenerator{}
	g.lastId.Store(startValue)
	return &DraftProvider{
		api:       api,
		generator: g,
		drafts:    make(map[uint64]*Draft),
	}
}

// SetChat sets the target chat and optional message thread for all drafts created
// by this provider. Must be called before NewDraft().
//
// If not set, NewDraft() will create drafts with zero chatID, which will cause
// SendMessageDraft to fail. Use this method to avoid runtime errors.
func (p *DraftProvider) SetChat(chatID int64, messageThreadID int) *DraftProvider {
	p.chatID = chatID
	p.messageThreadID = messageThreadID
	return p
}

// SetParseMode sets the default parse mode for all new drafts.
// Overrides the parse mode passed to NewDraft() only if not specified there.
func (p *DraftProvider) SetParseMode(mode tgapi.ParseMode) *DraftProvider {
	p.parseMode = mode
	return p
}

// SetEntities sets the default message entities (e.g., bold, links, mentions)
// to be copied into every new draft.
//
// Entities are shallow-copied — if you mutate the slice later, it will affect
// future drafts. For safety, pass a copy if needed.
func (p *DraftProvider) SetEntities(entities []tgapi.MessageEntity) *DraftProvider {
	p.entities = entities
	return p
}

// GetDraft retrieves a draft by its ID.
//
// Returns the draft and true if found, or nil and false if not found.
func (p *DraftProvider) GetDraft(id uint64) (*Draft, bool) {
	draft, ok := p.drafts[id]
	return draft, ok
}

// FlushAll sends all pending drafts as final messages and clears them.
//
// If any draft fails to send, FlushAll returns the error immediately and
// leaves other drafts unflushed. This allows for retry logic or logging.
//
// After successful flush, each draft is removed from the provider and cleared.
func (p *DraftProvider) FlushAll() error {
	var lastErr error
	for _, draft := range p.drafts {
		if err := draft.Flush(); err != nil {
			lastErr = err
			break // Stop on first error to avoid partial state
		}
	}
	return lastErr
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
// The draft inherits the provider's chatID, messageThreadID, and entities.
// If parseMode is zero, the provider's default parseMode is used.
//
// Panics if chatID is zero — call SetChat() on the provider first.
func (p *DraftProvider) NewDraft(parseMode tgapi.ParseMode) *Draft {
	if p.chatID == 0 {
		panic("laniakea: DraftProvider.SetChat() must be called before NewDraft()")
	}

	id := p.generator.Next()
	draft := &Draft{
		api:             p.api,
		provider:        p,
		chatID:          p.chatID,
		messageThreadID: p.messageThreadID,
		parseMode:       parseMode,
		entities:        p.entities, // Shallow copy — caller must ensure immutability
		ID:              id,
		Message:         "",
	}
	p.drafts[id] = draft
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
// pass a copy: `SetEntities(append([]tgapi.MessageEntity{}, myEntities...))`
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
		delete(d.provider.drafts, d.ID)
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

	params := tgapi.SendMessageP{
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

// push is the internal helper for Push(). It updates the server draft via SendMessageDraft.
func (d *Draft) push(text string) error {
	d.Message += text
	params := tgapi.SendMessageDraftP{
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
