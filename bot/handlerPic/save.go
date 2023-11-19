package handlerPic

import (
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Save(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    reply := update.Message.ReplyToMessage

    if reply != nil {
        logger("PSAVE", fmt.Sprintf("%v", reply.Photo))
    }
}
