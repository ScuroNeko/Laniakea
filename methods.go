package laniakea

import (
	"encoding/json"

	"git.nix13.pw/scuroneko/laniakea/tgapi"
)

func (bot *Bot[T]) Updates() ([]tgapi.Update, error) {
	offset := bot.GetUpdateOffset()
	params := tgapi.UpdateParams{
		Offset:         Ptr(offset),
		Timeout:        Ptr(30),
		AllowedUpdates: bot.GetUpdateTypes(),
	}

	updates, err := bot.api.GetUpdates(params)
	if err != nil {
		return nil, err
	}

	if bot.RequestLogger != nil {
		for _, u := range updates {
			j, err := json.Marshal(u)
			if err != nil {
				bot.GetLogger().Error(err)
			}
			bot.RequestLogger.Debugf("UPDATE %s\n", j)
		}
	}
	if len(updates) > 0 {
		bot.SetUpdateOffset(updates[len(updates)-1].UpdateID + 1)
	}
	return updates, err
}
