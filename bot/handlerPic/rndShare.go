package handlerPic

import (
    "embercat/pgsql"
    tba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func RndShare(bot *tba.BotAPI, update tba.Update) {
    var err error
    logger := getLogger("LOAD")

    pg := pgsql.GetClient()
    q := `SELECT photoId from anime ORDER BY random()`

    row := pg.QueryRow(q)

    var photoId string
    if err = row.Scan(&photoId); err != nil {
        logger(err.Error())
        return
    }

    file := tba.FileID(photoId)

    msg := tba.NewPhoto(update.Message.Chat.ID, file)
    if _, err = bot.Send(msg); err != nil {
        logger(err.Error())
    }
}
