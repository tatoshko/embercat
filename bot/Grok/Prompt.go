package Grok

import (
    "embercat/huggingface"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "strings"
)

func Prompt(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    var err error
    var logger = getLogger("Prompt")

    var hfc *huggingface.HuggingFaceClient
    if hfc, err = huggingface.GetClient(); err != nil {
        logger(err.Error())
        return
    }

    text := strings.TrimPrefix(update.Message.Text, "уголек")
    text = strings.TrimPrefix(update.Message.Text, ",")
    text = strings.TrimPrefix(update.Message.Text, " ")

    var result string
    if result, err = hfc.Ask(text); err != nil {
        logger(err.Error())
        return
    }

    msg := tgbotapi.NewMessage(update.Message.Chat.ID, result)
    msg.ReplyToMessageID = update.Message.MessageID

    if _, err = bot.Send(msg); err != nil {
        logger("Message send error", err.Error())
    }
}
