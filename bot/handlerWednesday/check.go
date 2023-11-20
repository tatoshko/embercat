package handlerWednesday

import (
    "database/sql"
    "embercat/assets"
    "embercat/pgsql"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "time"
)

func Check(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    var err error
    chatID := update.Message.Chat.ID
    logger := getLogger("CHECK")

    if time.Now().Weekday() == time.Wednesday {
        pg := pgsql.GetClient()
        if pg == nil {
            return
        }
        defer pg.Close()

        var row *sql.Row
        q := `SELECT * FROM frog ORDER BY random()`
        if row = pg.QueryRow(q); err != nil {
            logger(err.Error())
        }

        var frog string
        if err = row.Scan(&frog); err != nil {
            logger(err.Error())
            return
        }

        msg := tgbotapi.NewPhotoShare(chatID, frog)
        if _, err := bot.Send(msg); err != nil {
            logger(err.Error())
        }
    } else {
        box := assets.GetBox()
        if b, err := box.Bytes(NO_WEDNESDAY); err != nil {
            logger(err.Error())
        } else {
            msg := tgbotapi.NewPhotoUpload(chatID, tgbotapi.FileBytes{Name: NO_WEDNESDAY, Bytes: b})
            if _, err := bot.Send(msg); err != nil {
                logger(err.Error())
            }
        }
    }
}
