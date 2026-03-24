package handlerPic

import (
    "embercat/pgsql"
    tba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Save(bot *tba.BotAPI, update tba.Update) {
    logger := getLogger("SAVE")

    chatID := update.Message.Chat.ID
    reply := update.Message.ReplyToMessage

    if reply == nil || reply.Photo == nil {
        msg := tba.NewMessage(chatID, "Вызов этой команды возможен только в ответ на сообщение с картинкой")

        if _, err := bot.Send(msg); err != nil {
            logger(err.Error())
        }
        return
    }

    var photoID string
    for _, pic := range reply.Photo {
        photoID = pic.FileID
        break
    }

    pg := pgsql.GetClient()
    q := `insert into anime (photoid) values ($1)`
    if _, err := pg.Exec(q, photoID); err != nil {
        logger(err.Error())
    }
}
