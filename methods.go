package laniakea

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"git.nix13.pw/scuroneko/laniakea/tgapi"
)

func (b *Bot) Updates() ([]tgapi.Update, error) {
	offset := b.GetUpdateOffset()
	params := tgapi.UpdateParams{
		Offset:         Ptr(offset),
		Timeout:        Ptr(30),
		AllowedUpdates: b.GetUpdateTypes(),
	}

	updates, err := b.api.GetUpdates(params)
	if err != nil {
		return nil, err
	}

	for _, u := range updates {
		b.SetUpdateOffset(u.UpdateID + 1)
		err = b.GetQueue().Enqueue(&u)
		if err != nil {
			return nil, err
		}

		if b.RequestLogger != nil {
			j, err := json.Marshal(u)
			if err != nil {
				b.Logger().Error(err)
			}
			b.RequestLogger.Debugf("UPDATE %s\n", j)
		}
	}
	return updates, err
}

func (b *Bot) GetFileByLink(link string) ([]byte, error) {
	u := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", b.token, link)
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return io.ReadAll(res.Body)
}
