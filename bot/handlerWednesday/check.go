package handlerWednesday

import (
    "database/sql"
    "embercat/assets"
    "embercat/pgsql"
    tba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Check(bot *tba.BotAPI, update tba.Update) {
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

        file := tba.FileID(frog)

        msg := tba.NewPhoto(chatID, file)
        if _, err := bot.Send(msg); err != nil {
            logger(err.Error())
        }
    } else {
        box := assets.GetBox()
        if b, err := box.Bytes(NO_WEDNESDAY); err != nil {
            logger(err.Error())
        } else {
            msg := tba.NewPhoto(chatID, tba.FileBytes{Name: NO_WEDNESDAY, Bytes: b})
            if _, err := bot.Send(msg); err != nil {
                logger(err.Error())
            }
        }
    }
}
