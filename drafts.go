package laniakea

import (
	"math"
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
	return rand.Uint64N(math.MaxUint64)
}

type LinearDraftIdGenerator struct {
	draftIdGenerator
	lastId uint64
}

func (g *LinearDraftIdGenerator) Next() uint64 {
	return atomic.AddUint64(&g.lastId, 1)
}

type DraftProvider struct {
	api *tgapi.API

	chatID          int
	messageThreadID int
	parseMode       tgapi.ParseMode
	entities        []tgapi.MessageEntity

	drafts    map[uint64]*Draft
	generator draftIdGenerator
}
type Draft struct {
	api *tgapi.API

	chatID          int
	messageThreadID int
	parseMode       tgapi.ParseMode
	entities        []tgapi.MessageEntity

	ID      uint64
	Message string
}

func NewRandomDraftProvider(api *tgapi.API) *DraftProvider {
	return &DraftProvider{api: api, generator: &RandomDraftIdGenerator{}}
}
func NewLinearDraftProvider(api *tgapi.API, startValue uint64) *DraftProvider {
	return &DraftProvider{api: api, generator: &LinearDraftIdGenerator{lastId: startValue}}
}
func (d *DraftProvider) NewDraft() *Draft {
	id := d.generator.Next()
	draft := &Draft{d.api, d.chatID, d.messageThreadID, d.parseMode, d.entities, id, ""}
	d.drafts[id] = draft
	return draft
}

func (d *Draft) Push(draftId uint64, newText string) error {
	d.Message += newText
	params := tgapi.SendMessageDraftP{
		ChatID:    d.chatID,
		DraftID:   draftId,
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
