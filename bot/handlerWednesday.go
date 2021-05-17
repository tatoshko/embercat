package bot

import (
    tba "github.com/go-telegram-bot-api/telegram-bot-api"
    "embercat/assets"
    "log"
    "time"
)

func handleWednesday(bot *tba.BotAPI, update tba.Update)  {
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
        msg := tba.NewPhotoUpload(update.Message.Chat.ID, tba.FileBytes{Name: pic, Bytes: b})

        if _, err := bot.Send(msg); err != nil {
            log.Printf("Wednesday send error %s\n", err.Error())
        }
    }
}
