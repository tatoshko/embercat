package core

import (
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "github.com/tatoshko/tbot/assets"
    "log"
    "time"
)

func handleWednesday(bot *tgbotapi.BotAPI, update tgbotapi.Update, value string)  {
    log.Printf("Weekday %v|%d", time.Now().Weekday(), time.Wednesday)

    if time.Now().Weekday() != time.Wednesday {
        return
    }


    box := assets.GetBox()

    if b, err := box.Bytes("wednesday.jpg"); err != nil {
        log.Println(err)
    } else {
        id := update.Message.Chat.ID
        msg := tgbotapi.NewPhotoUpload(id, tgbotapi.FileBytes{
            Name:  "wednesday.jpg",
            Bytes: b,
        })

        bot.Send(msg)
    }
}
