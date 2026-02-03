package laniakea

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
)

type Uploader struct {
	bot *Bot
}

func NewUploader(bot *Bot) *Uploader {
	return &Uploader{bot: bot}
}
func (u *Uploader) UploadPhoto(chatId int, data []byte) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendPhoto", u.bot.token)

	buf := bytes.NewBuffer(nil)
	w := multipart.NewWriter(buf)
	defer w.Close()
	err := w.WriteField("chat_id", strconv.Itoa(chatId))
	if err != nil {
		return err
	}
	fw, err := w.CreateFormFile("photo", "photo.jpg")
	if err != nil {
		return err
	}
	_, err = fw.Write(data)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println("upload", string(body))
	return nil
}
