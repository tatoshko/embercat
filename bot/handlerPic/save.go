package handlerPic

import (
    "embercat/pgsql"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Save(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    logger := getLogger("SAVE")

    chatID := update.Message.Chat.ID
    reply := update.Message.ReplyToMessage

    if reply == nil || reply.Photo == nil {
        msg := tgbotapi.NewMessage(chatID, "Вызов этой команды возможен только в ответ на сообщение с картинкой")

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
