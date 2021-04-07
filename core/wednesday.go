package core

import (
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "github.com/tatoshko/tbot/assets"
    "log"
)

func handleWednesday(bot *tgbotapi.BotAPI, update tgbotapi.Update, value string)  {
    box := assets.GetBox()

    if b, err := box.Bytes("wednesday.jpg"); err != nil {
        log.Println(err)
    } else {
        id := update.Message.Chat.ID
        msg := tgbotapi.NewDocumentUpload(id, b)

        bot.Send(msg)
    }
}
