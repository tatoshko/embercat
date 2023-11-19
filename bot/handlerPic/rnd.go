package handlerPic

import (
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "io"
    "log"
    "net/http"
)

func Rnd(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    var err error
    chatID := update.Message.Chat.ID

    pic := NewPicRND()

    resp, err := http.Get(pic.URL())
    if err != nil {
        logger("HTTP.GET", err.Error())
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        logger("STATUS", string(resp.StatusCode))
        return
    }

    var b []byte
    if b, err = io.ReadAll(resp.Body); err != nil {
        logger("READ PIC", err.Error())
        return
    }

    msg := tgbotapi.NewPhotoUpload(chatID, tgbotapi.FileBytes{Name: pic.URL(), Bytes: b})
    msg.Caption = pic.GetID()

    if _, err := bot.Send(msg); err != nil {
        log.Printf("Pic.RND send message error: %s", err.Error())
    }
}
