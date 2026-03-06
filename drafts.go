package laniakea

import (
	"math/rand/v2"
	"sync/atomic"

	"git.nix13.pw/scuroneko/laniakea/tgapi"
)

type draftIdGenerator interface {
	Next() uint64
}

type RandomDraftIdGenerator struct {
	draftIdGenerator
}

func (g *RandomDraftIdGenerator) Next() uint64 {
	return rand.Uint64()
}

type LinearDraftIdGenerator struct {
	draftIdGenerator
	lastId atomic.Uint64
}

func (g *LinearDraftIdGenerator) Next() uint64 {
	return g.lastId.Add(1)
}

type DraftProvider struct {
	api *tgapi.API

	chatID          int64
	messageThreadID int
	parseMode       tgapi.ParseMode
	entities        []tgapi.MessageEntity

	drafts    map[uint64]*Draft
	generator draftIdGenerator
}
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

func NewRandomDraftProvider(api *tgapi.API) *DraftProvider {
	return &DraftProvider{
		api: api, generator: &RandomDraftIdGenerator{},
		parseMode: tgapi.ParseMDV2,
		drafts:    make(map[uint64]*Draft),
	}
}
func NewLinearDraftProvider(api *tgapi.API, startValue uint64) *DraftProvider {
	g := &LinearDraftIdGenerator{}
	g.lastId.Store(startValue)
	return &DraftProvider{
		api:       api,
		generator: g,
		drafts:    make(map[uint64]*Draft),
	}
}
func (d *DraftProvider) NewDraft() *Draft {
	id := d.generator.Next()
	entitiesCopy := make([]tgapi.MessageEntity, 0)
	copy(entitiesCopy, d.entities)
	draft := &Draft{
		api:             d.api,
		provider:        d,
		chatID:          d.chatID,
		messageThreadID: d.messageThreadID,
		parseMode:       d.parseMode,
		entities:        entitiesCopy,
		ID:              id,
		Message:         "",
	}
	d.drafts[id] = draft
	return draft
}
func (d *Draft) push(text string, escapeMd bool) error {
	if escapeMd {
		d.Message += EscapeMarkdownV2(text)
	} else {
		d.Message += text
	}
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

func (d *Draft) Push(text string) error {
	return d.push(text, true)
}
func (d *Draft) PushMarkdown(text string) error {
	return d.push(text, false)
}

func (d *Draft) Clear() {
	d.Message = ""
}
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
		d.Clear()
		delete(d.provider.drafts, d.ID)
	}
	return err
}
