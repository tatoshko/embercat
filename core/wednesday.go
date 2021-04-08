package core

import (
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "github.com/tatoshko/tbot/assets"
    "log"
    "time"
)

func handleWednesday(bot *tgbotapi.BotAPI, update tgbotapi.Update, value string)  {
    box := assets.GetBox()

    var pic string

    if time.Now().Weekday() == time.Wednesday {
        pic = "wednesday.jpg"
    } else {
        pic = "no-wednesday.jpg"
    }

    if b, err := box.Bytes(pic); err != nil {
        log.Println(err)
    } else {
        id := update.Message.Chat.ID
        msg := tgbotapi.NewPhotoUpload(id, tgbotapi.FileBytes{Name:  pic, Bytes: b})
        msg.UseExisting = true

        bot.Send(msg)
    }
}
