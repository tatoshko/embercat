package handlerPic

import (
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Save(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    reply := update.Message.ReplyToMessage

    if reply != nil && reply.Photo != nil {
        for _, pic := range *reply.Photo {
            logger("PSAVE", pic.FileID)
        }
    }
}
