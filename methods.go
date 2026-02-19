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

	for _, u := range updates {
		bot.SetUpdateOffset(u.UpdateID + 1)
		err = bot.GetQueue().Enqueue(&u)
		if err != nil {
			return nil, err
		}

		if bot.RequestLogger != nil {
			j, err := json.Marshal(u)
			if err != nil {
				bot.GetLogger().Error(err)
			}
			bot.RequestLogger.Debugf("UPDATE %s\n", j)
		}
	}
	return updates, err
}
