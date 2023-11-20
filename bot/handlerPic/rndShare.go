package handlerPic

import (
    "embercat/pgsql"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func RndShare(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
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

    msg := tgbotapi.NewPhotoShare(update.Message.Chat.ID, photoId)
    if _, err = bot.Send(msg); err != nil {
        logger(err.Error())
    }
}
