package handlerWednesday

import (
    "database/sql"
    "embercat/assets"
    "embercat/pgsql"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Check(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    var err error
    chatID := update.Message.Chat.ID
    logger := getLogger("CHECK")

    if ItIsWednesdayMyDudes() {
        pg := pgsql.GetClient()

        var row *sql.Row
        q := `SELECT photoid FROM frog ORDER BY random()`
        if row = pg.QueryRow(q); err != nil {
            logger(err.Error())
        }

        var frog string
        if err = row.Scan(&frog); err != nil {
            logger(err.Error())
            return
        }

        file := tgbotapi.FileID(frog)

        msg := tgbotapi.NewPhoto(chatID, file)
        if _, err := bot.Send(msg); err != nil {
            logger(err.Error())
        }
    } else {
        box := assets.GetBox()
        if b, err := box.Bytes(NO_WEDNESDAY); err != nil {
            logger(err.Error())
        } else {
            msg := tgbotapi.NewPhoto(chatID, tgbotapi.FileBytes{Name: NO_WEDNESDAY, Bytes: b})
            if _, err := bot.Send(msg); err != nil {
                logger(err.Error())
            }
        }
    }
}
