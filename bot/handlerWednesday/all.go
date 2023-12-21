package handlerWednesday

import (
    "database/sql"
    "embercat/pgsql"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func All(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    var err error
    logger := getLogger("ALL")

    pg := pgsql.GetClient()
    q := `select id, photoId from frog`

    var rows *sql.Rows
    if rows, err = pg.Query(q); err != nil {
        logger(err.Error())
        return
    }
    defer rows.Close()

    msg := tgbotapi.NewPhotoShare(update.Message.Chat.ID, "")
    for rows.Next() {
        var id, photoId string
        if err = rows.Scan(&id, &photoId); err != nil {
            logger(err.Error())
            return
        }

        msg.Caption = id
        msg.FileID = photoId

        if _, err = bot.Send(msg); err != nil {
            logger(err.Error())
        }
    }
}
