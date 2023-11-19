package handlerPic

import (
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Save(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    logger("TRY REPLY", fmt.Sprintf("%v", update.Message.ReplyToMessage))
}
