package bot

import (
    "embercat/assets"
    tba "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
    "time"
)

func handlerDay(bot *tba.BotAPI, update tba.Update) {
    box := assets.GetBox()

    var pic string

    switch time.Now().Weekday() {
    case time.Monday:
        pic = "mon.jpg"
    case time.Tuesday:
        pic = "tue.jpg"
    case time.Wednesday:
        pic = "wed.jpg"
    case time.Thursday:
        pic = "thu.jpg"
    case time.Friday:
        pic = "fri.jpg"
    case time.Saturday:
        pic = "sat.jpg"
    case time.Sunday:
        pic = "sun.jpg"
    }

    if b, err := box.Bytes("days/" + pic); err != nil {
        log.Println(err)
    } else {
        msg := tba.NewPhotoUpload(update.Message.Chat.ID, tba.FileBytes{Name: pic, Bytes: b})

        if _, err := bot.Send(msg); err != nil {
            log.Printf("Days send error %s\n", err.Error())
        }
    }
}
