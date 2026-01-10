package handlerPic

import (
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "io"
    "net/http"
)

func RndServer(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    var err error
    logger := getLogger("RND SERVER")
    chatID := update.Message.Chat.ID

    pic := NewPicRND()

    resp, err := http.Get(pic.URL())
    if err != nil {
        logger(err.Error())
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        logger(fmt.Sprintf("%d", resp.StatusCode))
        return
    }

    var b []byte
    if b, err = io.ReadAll(resp.Body); err != nil {
        logger(err.Error())
        return
    }

    msg := tgbotapi.NewPhoto(chatID, tgbotapi.FileBytes{Name: pic.URL(), Bytes: b})
    msg.Caption = pic.GetID()
    if _, err := bot.Send(msg); err != nil {
        logger(err.Error())
    }
}
